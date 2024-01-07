package edit

import (
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
	"strings"
)

type FileImports struct {
	importMap  map[string]string
	newImports []string
}

func NewFileImports(file *ast.File, pkg *packages.Package) *FileImports {
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

	return &FileImports{importMap: imports}
}

func (f *FileImports) Resolve(t types.Type) string {
	globalName := t.String()
	parts := strings.Split(globalName, ".")
	if len(parts) == 1 {
		return globalName
	}
	if len(parts) != 2 {
		panic("could not detect obj name")
	}

	pkgPath, objName := parts[0], parts[1]

	name, ok := f.importMap[pkgPath]
	if !ok {
		name = f.AppendImport(pkgPath)
	}

	if name == "" {
		return objName
	}

	return name + "." + objName
}

func (f *FileImports) AppendImport(pkgPath string) string {
	if name, ok := f.importMap[pkgPath]; ok {
		return name
	}

	f.newImports = append(f.newImports, pkgPath)
	parts := strings.Split(pkgPath, "/")
	name := parts[len(parts)-1]
	f.importMap[pkgPath] = name

	return name
}

func (f *FileImports) WriteImportsToAst(fset *token.FileSet, file *ast.File) {
	for _, pkgPath := range f.newImports {
		astutil.AddImport(fset, file, pkgPath)
	}
}
