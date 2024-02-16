package analysis

import (
	"context"
	"go/token"
	"golang.org/x/tools/go/packages"
	"mapgenie/entities"
	"mapgenie/pkg/log"
)

func FindTargetsInPackages(ctx context.Context, patterns ...string) []entities.TargetFile {
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
		log.Errorf(ctx, "Encountered errors while parsing")
		return nil
	}

	targetFiles := make([]entities.TargetFile, 0)

	for _, pkg := range pkgs {
		pkgCtx := log.Fold(ctx, "Analysing package: %s", pkg.Name)

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
