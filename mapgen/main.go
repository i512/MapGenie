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
)

const mapTemplate = `func {{ .FuncName }}(input {{ .InTypeArg }}) {{ .OutTypeArg }} {
	var result {{ .OutType }}

	{{- if .InIsPtr }}
	if input == nil {
		return {{ if .OutIsPtr }}&{{ end }}result
	}
	{{ end }}

	{{- range .Mappings }}
	{{- if .InPtr }}
	if input.{{ .InName }} != nil {
		result.{{ .OutName }} = *input.{{ .InName }}
	}
	{{- else }}
	result.{{ .OutName }} = {{ if .OutPtr }}&{{ end }}input.{{ .InName }}
	{{- end }}
	{{- end }}

	return {{ if .OutIsPtr }}&{{ end }}result
}
`

type TemplateMapping struct {
	InName  string
	OutName string
	OutPtr  bool
	InPtr   bool
}

type MapTemplateData struct {
	FuncName string
	InType   string
	InIsPtr  bool
	InputVar string
	OutType  string
	OutIsPtr bool
	Mappings []TemplateMapping
}

func (d MapTemplateData) InTypeArg() string {
	if d.InIsPtr {
		return "*" + d.InType
	}

	return d.InType
}

func (d MapTemplateData) OutTypeArg() string {
	if d.OutIsPtr {
		return "*" + d.OutType
	}

	return d.OutType
}

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

		mappable := mappableFields(tfs)

		data := MapTemplateData{
			FuncName: funcDecl.Name.Name,
			InType:   fileLocalName(tfs.In.Named, importMap),
			InIsPtr:  tfs.In.IsPtr,
			OutType:  fileLocalName(tfs.Out.Named, importMap),
			OutIsPtr: tfs.Out.IsPtr,
			Mappings: mappable,
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

func mappableFields(tfs TargetFuncSignature) []TemplateMapping {
	in := tfs.In.FieldMap()
	out := tfs.Out.FieldMap()

	list := make([]TemplateMapping, 0)

	for outMapField, outMapType := range out {
		inMapType, ok := in[outMapField]
		if !ok {
			fmt.Println("no matching field for ", outMapField)
			continue
		}

		if !token.IsExported(outMapField) && !(tfs.In.Local && tfs.Out.Local) {
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
