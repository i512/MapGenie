package analysis

import (
	"context"
	"go/token"
	"golang.org/x/tools/go/packages"
	"mapgenie/entities"
	"mapgenie/pkg/log"
	"strings"
)

func FindTargetsInPackages(ctx context.Context, patterns ...string) []entities.TargetFile {
	ctx = log.Fold(ctx, "Analyze package paths (patterns): %s", strings.Join(patterns, ", "))
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
		return nil
	}

	if packages.PrintErrors(pkgs) > 0 {
		log.Errorf(ctx, "Encountered errors while parsing. The problem above ^^^")
		return nil
	}

	targetFiles := make([]entities.TargetFile, 0)

	for _, pkg := range pkgs {
		pkgCtx := log.Fold(ctx, "package: %s", pkg.Name)

		for _, file := range pkg.Syntax {
			file := NewFileAnalysis(pkg, fset, file)
			targetFile := file.FindTargets(pkgCtx)
			if targetFile == nil {
				continue
			}

			targetFiles = append(targetFiles, *targetFile)
		}
	}

	return targetFiles
}
