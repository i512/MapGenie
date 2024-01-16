package main

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
	"mapgenie/mapgen/edit"
	"mapgenie/mapgen/entities"
	"mapgenie/mapgen/typematch"
	"os"
	"regexp"
)

func main() {
	flag.Parse()
	patterns := flag.Args()
	if len(patterns) == 0 {
		panic("provide package patterns in arguments, for example: ./...")
	}

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
		panic(err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		panic("load errors")
	}

	// Print the names of the source files
	// for each package listed on the command line.
	for _, pkg := range pkgs {
		fmt.Println("analysis of pkg:", pkg.Name)

		fmt.Println(pkg.ID, pkg.GoFiles)
		for i, file := range pkg.Syntax {
			filePath := pkg.GoFiles[i]
			analyzePkgFile(pkg.Fset, file, filePath, pkg)
		}
	}
}

func analyzePkgFile(fset *token.FileSet, file *ast.File, filePath string, pkg *packages.Package) {
	toMapFuncRegex := regexp.MustCompile(`^(?:\w+) map this pls`)

	fileImports := edit.NewFileImports(file, pkg)

	astModified := false

	astutil.Apply(file, nil, func(cursor *astutil.Cursor) bool {
		funcDecl, ok := cursor.Node().(*ast.FuncDecl)
		if !ok {
			return true
		}

		if funcDecl.Doc == nil || !toMapFuncRegex.MatchString(funcDecl.Doc.Text()) {
			fmt.Println("ignored", funcDecl.Name)
			return true
		}

		fmt.Println("will map this func:", funcDecl.Name)

		tfs, err := getInputOutputTypes(funcDecl, pkg)
		if errors.Is(err, ErrFuncMismatchError) {
			fmt.Println("invalid function: ", err)
			return true
		}
		if err != nil {
			fmt.Println("failed to get func types: ", err)
			return true
		}

		mappable := typematch.MappableFields(tfs, fileImports)

		data := edit.MapTemplateData{
			FuncName: funcDecl.Name.Name,
			InType:   fileImports.ResolveTypeName(tfs.In.Named),
			InIsPtr:  tfs.In.IsPtr,
			OutType:  fileImports.ResolveTypeName(tfs.Out.Named),
			OutIsPtr: tfs.Out.IsPtr,
			Mappings: mappable,
		}

		newFuncDecl := edit.MapperFuncAst(fset, data)
		funcDecl.Body = newFuncDecl.Body
		funcDecl.Type.Params = newFuncDecl.Type.Params
		funcDecl.Type.Results = newFuncDecl.Type.Results

		astModified = true

		return true
	})

	if astModified {
		fileImports.WriteImportsToAst(fset, file)

		modifyFile(fset, file)
	}
}

func modifyFile(fset *token.FileSet, file *ast.File) {
	path := fset.Position(file.Pos()).Filename

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0755) // read existing permissions?
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = format.Node(f, fset, file)
	if err != nil {
		panic(err)
	}
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
		return sa, fmt.Errorf("%w: wrong number of arguments", ErrFuncMismatchError)
	}

	firstArg := tuple.At(0).Type()
	if ptr, ok := firstArg.(*types.Pointer); ok {
		sa.IsPtr = true
		firstArg = ptr.Elem()
	}
	named := firstArg.(*types.Named)

	structArg, ok := named.Underlying().(*types.Struct)
	if !ok {
		return sa, fmt.Errorf("%w: input argument is not a struct", ErrFuncMismatchError)
	}

	sa.Named = named
	sa.Struct = structArg

	return sa, nil
}
