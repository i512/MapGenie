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
	"os"
	"reflect"
	"regexp"
)

type TargetFuncSignature struct {
	In, Out StructArg
}

type StructArg struct {
	Named  *types.Named
	Struct *types.Struct
	IsPtr  bool
	Local  bool
}

func (s StructArg) FieldMap() map[string]types.Type {
	result := map[string]types.Type{}

	for i := 0; i < s.Struct.NumFields(); i++ {
		f := s.Struct.Field(i)
		result[f.Name()] = f.Type()
	}

	return result
}

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

		mappable := mappableFields(tfs, fileImports)

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

func getInputOutputTypes(f *ast.FuncDecl, pkg *packages.Package) (tfs TargetFuncSignature, err error) {
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

func structArgFromTuple(tuple *types.Tuple) (sa StructArg, err error) {
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

func mappableFields(tfs TargetFuncSignature, imports *edit.FileImports) []edit.TemplateMapping {
	in := tfs.In.FieldMap()

	list := make([]edit.TemplateMapping, 0)

	for i := 0; i < tfs.Out.Struct.NumFields(); i++ {
		field := tfs.Out.Struct.Field(i)
		outFieldName := field.Name()
		outFieldType := field.Type()

		inFieldType, ok := in[outFieldName]
		if !ok {
			fmt.Println("no matching field for ", outFieldName)
			continue
		}

		if !token.IsExported(outFieldName) && !(tfs.In.Local && tfs.Out.Local) {
			fmt.Println("output field is unexported ", outFieldName)
			continue
		}

		if typesAreCastable(inFieldType, outFieldType) {
			mapping := edit.TemplateMapping{
				InName:   outFieldName,
				OutName:  outFieldName,
				CastWith: castWith(inFieldType, outFieldType, imports),
			}
			list = append(list, mapping)
			continue
		}

		outPtr, ok := outFieldType.(*types.Pointer)
		if ok && typesAreCastable(inFieldType, outPtr.Elem()) {
			list = append(list, edit.TemplateMapping{
				InName:   outFieldName,
				OutName:  outFieldName,
				OutPtr:   true,
				CastWith: castWith(inFieldType, outPtr.Elem(), imports),
			})
			continue
		}

		inPtr, ok := inFieldType.(*types.Pointer)
		if ok && typesAreCastable(inPtr.Elem(), outFieldType) {
			list = append(list, edit.TemplateMapping{
				InName:   outFieldName,
				OutName:  outFieldName,
				InPtr:    true,
				CastWith: castWith(inPtr.Elem(), outFieldType, imports),
			})
			continue
		}

		if inPtr != nil && outPtr != nil && typesAreCastable(inPtr.Elem(), outPtr.Elem()) {
			list = append(list, edit.TemplateMapping{
				InName:   outFieldName,
				OutName:  outFieldName,
				InPtr:    true,
				OutPtr:   true,
				CastWith: castWith(inPtr.Elem(), outPtr.Elem(), imports),
			})
			continue
		}
	}

	return list
}

func typesAreCastable(in, out types.Type) bool {
	same := reflect.DeepEqual(in, out)
	numbers := typeIsIntegerOrFloat(in) && typeIsIntegerOrFloat(out)
	stringLike := typeIsStringOrByteSlice(in) && typeIsStringOrByteSlice(out)
	derivedType := typeIsUnderlying(in, out) || typeIsUnderlying(out, in)
	return same || numbers || stringLike || derivedType
}

func castWith(in, out types.Type, imports *edit.FileImports) string {
	if reflect.DeepEqual(in, out) {
		return "" // same type, no cast needed
	}

	return imports.ResolveTypeName(out)
}

func typeIsIntegerOrFloat(t types.Type) bool {
	basic, isBasic := t.(*types.Basic)
	return isBasic && (basic.Info()&types.IsInteger > 0 || basic.Info()&types.IsFloat > 0)
}

func typeIsStringOrByteSlice(t types.Type) bool {
	if basic, ok := t.(*types.Basic); ok && basic.Kind()&types.String > 0 {
		return true
	}

	slice, ok := t.(*types.Slice)
	if !ok {
		return false
	}

	basic, ok := slice.Elem().(*types.Basic)
	return ok && basic.Kind()&types.Byte > 0
}

func typeIsUnderlying(base, derived types.Type) bool {
	return reflect.DeepEqual(base, derived.Underlying())
}
