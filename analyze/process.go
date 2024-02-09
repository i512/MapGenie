package analyze

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
	"mapgenie/entities"
	"mapgenie/gen"
	"mapgenie/pkg/log"
	"mapgenie/typematch"
	"os"
	"regexp"
)

func ProcessPackages(ctx context.Context, patterns ...string) {
	fset := token.NewFileSet()

	cfg := packages.Config{
		Fset: fset,
		Mode: packages.NeedTypes |
			packages.NeedTypesSizes |
			packages.NeedSyntax |
			packages.NeedTypesInfo |
			packages.NeedImports |
			packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles,
	}

	pkgs, err := packages.Load(&cfg, patterns...)
	if err != nil {
		log.Errorf(ctx, "Failed to load packages: %s", err.Error())
		return
	}
	if packages.PrintErrors(pkgs) > 0 {
		log.Errorf(ctx, "Encountered errors while parsing")
		return
	}

	for _, pkg := range pkgs {
		pkgCtx := log.Fold(ctx, "Analysing package: %s", pkg.Name)

		for _, file := range pkg.Syntax {
			analyzePkgFile(pkgCtx, pkg.Fset, file, pkg)
		}
	}
}

var targetFuncComment = regexp.MustCompile(`^\w+ map this pls`)

func analyzePkgFile(ctx context.Context, fset *token.FileSet, file *ast.File, pkg *packages.Package) {
	ctx = log.Fold(ctx, "Analysing file: %s", fset.File(file.Pos()).Name())

	fileImports := gen.NewFileImports(file, pkg)

	astModified := false

	astutil.Apply(file, nil, func(cursor *astutil.Cursor) bool {
		funcDecl, ok := cursor.Node().(*ast.FuncDecl)
		if !ok {
			return true
		}

		if funcDecl.Doc == nil || !targetFuncComment.MatchString(funcDecl.Doc.Text()) {
			return true
		}

		processMapper(ctx, pkg, fset, funcDecl, fileImports)

		astModified = true

		return true
	})

	if astModified {
		fileImports.WriteImportsToAst(fset, file)

		err := modifyFile(fset, file)
		if err != nil {
			log.Errorf(ctx, "Failed to change file: %s", err.Error())
		}
	}
}

func processMapper(
	ctx context.Context,
	pkg *packages.Package,
	fset *token.FileSet,
	funcDecl *ast.FuncDecl,
	imports *gen.FileImports,
) (astModified bool) {
	ctx = log.Fold(ctx, "Mapper found: %s", funcDecl.Name)

	tfs, err := getInputOutputTypes(funcDecl, pkg)
	if errors.Is(err, ErrFuncMismatchError) {
		log.Errorf(ctx, "Invalid mapper: %s", err.Error())
		return false
	}
	if err != nil {
		log.Errorf(ctx, "Failed to interpret func signature: %s", err.Error())
		return false
	}

	mappable := typematch.MappableFields(ctx, tfs)

	data := gen.MapTemplateData{
		FuncName: funcDecl.Name.Name,
		InType:   imports.ResolveTypeName(tfs.In.Named),
		InIsPtr:  tfs.In.IsPtr,
		OutType:  imports.ResolveTypeName(tfs.Out.Named),
		OutIsPtr: tfs.Out.IsPtr,
		Maps:     mappable,
		Resolver: imports,
	}

	newFuncDecl, err := gen.MapperFuncAst(ctx, fset, data)
	if err != nil {
		log.Errorf(ctx, "Generation failed: %s", err.Error())
		return false
	}

	funcDecl.Body = newFuncDecl.Body
	funcDecl.Type.Params = newFuncDecl.Type.Params
	funcDecl.Type.Results = newFuncDecl.Type.Results

	return true
}

func modifyFile(fset *token.FileSet, file *ast.File) error {
	path := fset.Position(file.Pos()).Filename

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0) // read existing permissions?
	if err != nil {
		return err
	}
	defer f.Close()

	return format.Node(f, fset, file)
}

var ErrFuncMismatchError = fmt.Errorf("function is not mappable")

func getInputOutputTypes(f *ast.FuncDecl, pkg *packages.Package) (tfs entities.TargetFuncSignature, err error) {
	funcType := pkg.Types.Scope().Lookup(f.Name.Name)

	signature := funcType.Type().(*types.Signature)
	tfs.In, err = structArgFromTuple(signature.Params())
	if err != nil {
		return tfs, fmt.Errorf("bad argument: %w", err)
	}

	tfs.Out, err = structArgFromTuple(signature.Results())
	if err != nil {
		return tfs, fmt.Errorf("bad return argument: %w", err)
	}

	tfs.In.Local = tfs.In.Named.Obj().Pkg().Path() == pkg.PkgPath
	tfs.Out.Local = tfs.Out.Named.Obj().Pkg().Path() == pkg.PkgPath

	return tfs, nil
}

func structArgFromTuple(tuple *types.Tuple) (sa entities.Argument, err error) {
	if tuple.Len() != 1 {
		return sa, fmt.Errorf("signature must have a single argument, have: %s. %w", tuple.String(), ErrFuncMismatchError)
	}

	firstArg := tuple.At(0).Type()
	if ptr, ok := firstArg.(*types.Pointer); ok {
		sa.IsPtr = true
		firstArg = ptr.Elem()
	}
	named, ok := firstArg.(*types.Named)
	if !ok {
		return sa, fmt.Errorf("`%s` is not a struct. %w", firstArg.String(), ErrFuncMismatchError)
	}

	structArg, ok := named.Underlying().(*types.Struct)
	if !ok {
		return sa, fmt.Errorf("`%s` is not a struct. %w", named.String(), ErrFuncMismatchError)
	}

	sa.Named = named
	sa.Struct = structArg

	return sa, nil
}
