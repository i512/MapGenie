package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
	"io"
	"os"
	"reflect"
	"regexp"
	"strings"
	"text/template"
	"unicode"
)

const mapTemplate = `func {{ .FuncName }}(input {{ .InputType }}) {{ .OutputType }} {
	var result {{ .OutputType }}

	{{- range .Mappings }}
	{{- if .InPtr }}
	if input.{{ .InName }} != nil {
		result.{{ .OutName }} = *input.{{ .InName }}
	}
	{{- else }}
	result.{{ .OutName }} = {{ if .OutPtr }}&{{ end }}input.{{ .InName }}
	{{- end }}
	{{- end }}

	return result
}
`

type TemplateMapping struct {
	InName  string
	OutName string
	OutPtr  bool
	InPtr   bool
}

type MapTemplateData struct {
	FuncName   string
	InputType  string
	InputVar   string
	OutputType string
	Mappings   []TemplateMapping
}

type TargetFuncSignature struct {
	InType, OutType     *types.Named
	InVar               string
	InStruct, OutStruct *types.Struct
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

type FileChange struct {
	pos, end int
	content  []byte
}

func analyzePkgFile(fset *token.FileSet, file *ast.File, filePath string, pkg *packages.Package) {
	toMapFuncRegex := regexp.MustCompile(`^(?:\w+) map this pls`)

	importMap := filePkgNameMap(file, pkg)

	var fileChange []FileChange

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
			FuncName:   funcDecl.Name.Name,
			InputType:  fileLocalName(tfs.InType, importMap),
			InputVar:   tfs.InVar,
			OutputType: fileLocalName(tfs.OutType, importMap),
			Mappings:   mappable,
		}

		t := template.Must(template.New("map").Parse(mapTemplate))
		funcSource := bytes.NewBuffer(nil)
		err = t.Execute(funcSource, data)
		if err != nil {
			panic(err)
		}

		fileChange = append(fileChange, FileChange{
			pos:     fset.Position(funcDecl.Pos()).Offset,
			end:     fset.Position(funcDecl.End()).Offset,
			content: funcSource.Bytes(),
		})

		return true
	})

	if len(fileChange) > 0 {
		modifyFile(filePath, fileChange)
	}
}

func modifyFile(filePath string, fileChange []FileChange) {
	f, err := os.OpenFile(filePath, os.O_RDWR, 0755) // read existing permissions?
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fileContent, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	for i := len(fileChange) - 1; i >= 0; i-- {
		change := fileChange[i]
		//fileContent = append(append(fileContent[:change.pos], change.content...), fileContent[change.end:]...)

		step1 := append(change.content, fileContent[change.end:]...)
		step2 := append(fileContent[:change.pos], step1...)

		fileContent = step2
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		panic(err)
	}

	_, err = f.Write(fileContent)
	if err != nil {
		panic(err)
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

func mappableFields(in, out map[string]types.Type) []TemplateMapping {
	list := make([]TemplateMapping, 0)

	for outMapField, outMapType := range out {
		inMapType, ok := in[outMapField]
		if !ok {
			fmt.Println("no matching field for ", outMapField)
			continue
		}

		if unicode.IsLower(([]rune(outMapField))[0]) {
			fmt.Println("output field is unexported ", outMapField)
			continue
		}

		if reflect.DeepEqual(inMapType, outMapType) {
			list = append(list, TemplateMapping{InName: outMapField, OutName: outMapField})
			continue
		}

		outPtr, ok := outMapType.(*types.Pointer)
		if ok && reflect.DeepEqual(inMapType, outPtr.Elem()) {
			list = append(list, TemplateMapping{InName: outMapField, OutName: outMapField, OutPtr: true})
			continue
		}

		inPtr, ok := inMapType.(*types.Pointer)
		if ok && reflect.DeepEqual(inPtr.Elem(), outMapType) {
			list = append(list, TemplateMapping{InName: outMapField, OutName: outMapField, InPtr: true})
			continue
		}

	}

	return list
}
