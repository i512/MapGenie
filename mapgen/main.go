package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
	"os"
	"reflect"
	"regexp"
	"strings"
	"text/template"
)

const mapTemplate = `
func (input {{ .InputType }}) {{ .OutputType }} {
	var result {{ .OutputType }}

	{{ range .Fields }}result.{{ .Name }} = {{ $.InputVar }}.{{ .Name }}
	{{ end }}

	return result
}
`

type Field struct {
	Name string
}

type MapTemplateData struct {
	InputType  string
	InputVar   string
	OutputType string
	Fields     []Field
}

type TargetFuncSignature struct {
	InType, OutType     *types.Named
	InVar               string
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
		for i, file := range pkg.Syntax {
			filePath := pkg.GoFiles[i]
			analyzePkgFile(pkg.Fset, file, filePath, pkg)
		}
	}
}

func analyzePkgFile(fset *token.FileSet, file *ast.File, filePath string, pkg *packages.Package) {
	toMapFuncRegex := regexp.MustCompile(`^(?:\w+) map this pls`)

	astModified := false

	importMap := filePkgNameMap(file, pkg)

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

		inMap := getStructFieldMap(tfs.InStruct, pkg)
		outMap := getStructFieldMap(tfs.OutStruct, pkg)
		mappable := mappableFields(inMap, outMap)

		data := MapTemplateData{
			InputType:  fileLocalName(tfs.InType, importMap),
			InputVar:   tfs.InVar,
			OutputType: fileLocalName(tfs.OutType, importMap),
			Fields:     mappable,
		}

		t := template.Must(template.New("map").Parse(mapTemplate))
		funcSource := bytes.NewBuffer(nil)
		err = t.Execute(funcSource, data)
		if err != nil {
			panic(err)
		}

		fmt.Println(funcSource.String())
		fmt.Println("mappable: ", mappable)

		replacementFile := pkg.Name + "_" + funcDecl.Name.Name + "_replacement.go"
		fmt.Println(replacementFile)
		funcAst, err := parser.ParseExprFrom(fset, replacementFile, funcSource.String(), 0)
		if err != nil {
			panic(err)
		}

		funcDecl.Body = funcAst.(*ast.FuncLit).Body

		astModified = true

		return true
	})

	if astModified {
		f, err := os.OpenFile(filePath, os.O_RDWR|os.O_TRUNC, 0755)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		err = printer.Fprint(f, fset, file)
		if err != nil {
			panic(err)
		}
	}
}

func fileLocalName(t *types.Named, importMap map[string]string) string {
	globalName := t.String()
	parts := strings.Split(globalName, ".")
	if len(parts) != 2 {
		panic("could not detect obj name")
	}

	pkgName, objName := parts[0], parts[1]

	name, ok := importMap[pkgName]
	if !ok {
		panic("failed to obtain type's import name")
	}

	if name == "" {
		return objName
	}

	return name + "." + objName
}

func filePkgNameMap(file *ast.File, pkg *packages.Package) map[string]string {
	imports := make(map[string]string)

	imports[pkg.PkgPath] = ""

	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		pathWithoutV2 := strings.TrimSuffix(path, "/v2")
		parts := strings.Split(pathWithoutV2, "/")

		localName := parts[len(parts)-1]
		if imp.Name != nil {
			localName = imp.Name.Name
		}

		imports[path] = localName
	}

	return imports
}

var ErrFuncMismatchError = fmt.Errorf("function is not mappable")

func getInputOutputTypes(f *ast.FuncDecl, pkg *packages.Package) (tfs TargetFuncSignature, err error) {
	funcType := pkg.Types.Scope().Lookup(f.Name.Name)

	signature := funcType.Type().(*types.Signature)
	if signature.Params().Len() != 1 {
		return tfs, fmt.Errorf("%w: wrong number of arguments", ErrFuncMismatchError)
	}

	inputType := signature.Params().At(0).Type()

	structArg, ok := inputType.Underlying().(*types.Struct)
	if !ok {
		return tfs, fmt.Errorf("input argument is not a struct")
	}

	if signature.Results().Len() != 1 {
		return tfs, fmt.Errorf("%w: wrong number of return arguments", ErrFuncMismatchError)
	}

	resultType := signature.Results().At(0).Type()
	resultStruct, ok := resultType.Underlying().(*types.Struct)
	if !ok {
		return tfs, fmt.Errorf("%w: result type is not a struct")
	}

	tfs = TargetFuncSignature{
		InType:    inputType.(*types.Named),
		InVar:     signature.Params().At(0).Name(),
		OutType:   resultType.(*types.Named),
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

		if !reflect.DeepEqual(inMapType, outMapType) {
			fmt.Println("different field types ", outMapField)
			continue
		}

		list = append(list, Field{Name: outMapField})
	}

	return list
}
