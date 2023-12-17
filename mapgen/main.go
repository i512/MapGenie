package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
	"os"
	"regexp"
	"text/template"
)

const mapTemplate = `
func f(input {{ .InputType }}) {{ .OutputType }} {
	result := {{ .OutputType }}{}

	{{ range .Fields }}
	result.{{ .Name }} = input.{{  .Name }}
	{{ end }}

	return result
}
`

type Field struct {
	Name string
}

type MapTemplateData struct {
	InputType  string
	OutputType string
	Fields     []Field
}

type TargetFuncSignature struct {
	InType, OutType     types.Object
	InStruct, OutStruct *types.Struct
}

func main() {
	//Many tools pass their command-line arguments (after any flags)
	//uninterpreted to packages.Load so that it can interpret them
	//according to the conventions of the underlying build system.
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

	pkgs, err := packages.Load(&cfg, "./...")
	if err != nil {
		fmt.Fprintf(os.Stderr, "load: %v\n", err)
		os.Exit(1)
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	// Print the names of the source files
	// for each package listed on the command line.
	for _, pkg := range pkgs {
		fmt.Println("analysis of pkg:", pkg.Name)

		fmt.Println(pkg.ID, pkg.GoFiles)
		for _, file := range pkg.Syntax {
			analyzePkgFile(pkg.Fset, file, pkg)
		}
	}
}

func analyzePkgFile(fset *token.FileSet, file *ast.File, pkg *packages.Package) {
	regex := regexp.MustCompile(`^(?:\w+) map this pls`)

	astutil.Apply(file, nil, func(cursor *astutil.Cursor) bool {
		funcDecl, ok := cursor.Node().(*ast.FuncDecl)
		if !ok {
			return true
		}

		if funcDecl.Doc == nil || !regex.MatchString(funcDecl.Doc.Text()) {
			fmt.Println("ignored", funcDecl.Name)
			return true
		}

		fmt.Println("will map this func:", funcDecl.Name)

		tfs, err := getInputOutputTypes(funcDecl, pkg)
		if errors.Is(err, ErrFuncMismatchError) {
			fmt.Println("invalid function: ", err)
			return true
		}

		inMap := getStructFieldMap(tfs.InStruct, pkg)
		outMap := getStructFieldMap(tfs.OutStruct, pkg)
		mappable := mappableFields(inMap, outMap)

		data := MapTemplateData{
			InputType:  tfs.InType.Name(),
			OutputType: tfs.OutType.Name(),
			Fields:     mappable,
		}

		t := template.Must(template.New("map").Parse(mapTemplate))
		err = t.Execute(os.Stdout, data)
		fmt.Println(err)

		fmt.Println("mappable: ", mappable)

		return true
	})
}

var ErrFuncMismatchError = fmt.Errorf("function is not mappable")

func getInputOutputTypes(f *ast.FuncDecl, pkg *packages.Package) (tfs TargetFuncSignature, err error) {
	paramList := f.Type.Params.List
	if len(paramList) != 1 {
		return tfs, fmt.Errorf("%w: wrong number of arguments", ErrFuncMismatchError)
	}

	inputIdent, ok := paramList[0].Type.(*ast.Ident)
	if !ok {
		return tfs, fmt.Errorf("failed to get type ident")
	}

	inputType := pkg.Types.Scope().Lookup(inputIdent.Name)

	structArg, ok := inputType.Type().Underlying().(*types.Struct)

	if !ok {
		return tfs, fmt.Errorf("input argument is not a struct")
	}

	results := f.Type.Results
	if results == nil || len(results.List) != 1 {
		return tfs, fmt.Errorf("%w: wrong number of return arguments", ErrFuncMismatchError)
	}

	resultIdent, ok := results.List[0].Type.(*ast.Ident)
	if !ok {
		return tfs, fmt.Errorf("failed to get result ident")
	}

	resultType := pkg.Types.Scope().Lookup(resultIdent.Name)
	resultStruct, ok := resultType.Type().Underlying().(*types.Struct)
	if !ok {
		return tfs, fmt.Errorf("%w: result type is not a struct")
	}

	tfs = TargetFuncSignature{
		InType:    inputType,
		OutType:   resultType,
		InStruct:  structArg,
		OutStruct: resultStruct,
	}

	return tfs, nil
}

func getStructFieldMap(s *types.Struct, pkg *packages.Package) map[string]types.Type {
	result := map[string]types.Type{}

	for i := 0; i < s.NumFields(); i++ {
		f := s.Field(i)
		result[f.Name()] = f.Type()

		if bt, ok := f.Type().Underlying().(*types.Basic); ok {
			if (bt.Info() & types.IsInteger) == types.IsInteger {
				fmt.Println("found integer", f.Name())
			}
		}
	}

	return result
}

func mappableFields(in, out map[string]types.Type) []Field {
	list := make([]Field, 0)

	for outMapField, outMapType := range out {
		inMapType, ok := in[outMapField]
		if !ok {
			fmt.Println("no matching field for ", outMapField)
			continue
		}

		if inMapType != outMapType {
			fmt.Println("different field types ", outMapField)
			continue
		}

		list = append(list, Field{Name: outMapField})
	}

	return list
}
