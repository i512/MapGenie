package analysis

import (
	"context"
	"github.com/i512/mapgenie/entities"
	"github.com/i512/mapgenie/pkg/log"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
	"regexp"
	"strings"
)

func FindTargetsInPackages(ctx context.Context, patterns ...string) []entities.TargetFile {
	ctx = log.Fold(ctx, "Analyze package paths (patterns): %s", strings.Join(patterns, ", "))
	fset := token.NewFileSet()

	cfg := packages.Config{
		Fset: fset,
		Mode: packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedTypesInfo |
			packages.NeedImports |
			packages.NeedName |
			packages.NeedFiles,
	}

	pkgs, err := packages.Load(&cfg, patterns...)
	if err != nil {
		log.Errorf(ctx, "Failed to load packages: %s", err.Error())
		return nil
	}

	if packages.PrintErrors(pkgs) > 0 {
		log.Errorf(ctx, "Encountered errors while parsing. The problem above ^^^")
		return nil
	}

	targetFiles := make([]entities.TargetFile, 0)

	for _, pkg := range pkgs {
		pkgCtx := log.Fold(ctx, "package: %s", pkg.Name)

		providers := findProviders(pkg)

		for _, file := range pkg.Syntax {
			file := NewFileAnalysis(pkg, fset, file, providers)
			targetFile := file.FindTargets(pkgCtx)
			if targetFile == nil {
				continue
			}

			targetFiles = append(targetFiles, *targetFile)
		}
	}

	return targetFiles
}

func findProviders(pkg *packages.Package) []*types.Func {
	funcs := make([]*types.Func, 0)

	for _, file := range pkg.Syntax {
		astutil.Apply(file, nil, func(cursor *astutil.Cursor) bool {
			funcDecl, ok := cursor.Node().(*ast.FuncDecl)
			if !ok || !isProviderFunc(funcDecl) {
				return true
			}

			funcType, ok := pkg.Types.Scope().Lookup(funcDecl.Name.Name).(*types.Func)
			if !ok {
				return true
			}

			funcs = append(funcs, funcType)

			return true
		})
	}

	return funcs
}

var providerFuncMagicComment = regexp.MustCompile(`^\w+ magic provider`)

func isProviderFunc(f *ast.FuncDecl) bool {
	return f.Doc != nil && providerFuncMagicComment.MatchString(f.Doc.Text())
}
