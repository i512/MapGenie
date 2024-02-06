package gen

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
	"mapgenie/pkg/log"
	"strings"
)

type FileImports struct {
	pkg        *packages.Package
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

	return &FileImports{importMap: imports, pkg: pkg}
}

func (f *FileImports) ResolveTypeName(t types.Type) string {
	val, ok := t.(*types.Pointer)
	if ok {
		return "*" + f.ResolveTypeName(val.Elem())
	}

	globalName := t.String()
	parts := strings.Split(globalName, ".")
	if len(parts) == 1 {
		return globalName
	}
	if len(parts) != 2 {
		log.Fatalf(context.Background(), "cannot parse object name: %s", globalName)
	}

	pkgPath, objName := parts[0], parts[1]

	name := f.ResolvePkgImport(pkgPath)
	if name == "" {
		return objName
	}

	return name + "." + objName
}

func (f *FileImports) ResolvePkgImport(pkgPath string) string {
	if name, ok := f.importMap[pkgPath]; ok {
		return name
	}

	f.newImports = append(f.newImports, pkgPath)
	name := f.pkgDefaultAlias(pkgPath)
	alias := f.resolvePkgNameCollision(name, pkgPath)
	f.importMap[pkgPath] = alias

	return alias
}

func (f *FileImports) pkgDefaultAlias(pkgPath string) string {
	parts := strings.Split(pkgPath, "/")
	return parts[len(parts)-1]
}

func (f *FileImports) resolvePkgNameCollision(name string, pkgPath string) string {
	if f.checkNameAvailable(name, pkgPath) {
		return name
	}

	for i := 2; i < 1000; i++ {
		alias := fmt.Sprintf("%s%d", name, i)
		if f.checkNameAvailable(alias, pkgPath) {
			return alias
		}
	}

	log.Fatalf(context.Background(), "could not find an available name to use for: %s", name)
	return ""
}

func (f *FileImports) checkNameAvailable(name string, pkgPath string) bool {
	for path, usedImportName := range f.importMap {
		if usedImportName == name {
			return path == pkgPath
		}
	}

	return f.pkg.Types.Scope().Lookup(name) == nil
}

func (f *FileImports) WriteImportsToAst(fset *token.FileSet, file *ast.File) {
	for _, pkgPath := range f.newImports {
		alias := f.importMap[pkgPath]
		name := ""
		if alias != f.pkgDefaultAlias(pkgPath) {
			name = alias
		}
		astutil.AddNamedImport(fset, file, name, pkgPath)
	}
}
