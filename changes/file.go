package changes

import (
	"context"
	"github.com/i512/mapgenie/entities"
	"github.com/i512/mapgenie/gen"
	"github.com/i512/mapgenie/pkg/log"
	"go/ast"
	"go/format"
	"go/token"
	"os"
)

func ApplyFilesChanges(ctx context.Context, files []entities.TargetFile) {
	ctx = log.Fold(ctx, "Writing changes to %d files", len(files))

	for _, file := range files {
		ApplyFile(ctx, file)
	}
}

func ApplyFile(ctx context.Context, file entities.TargetFile) {
	ctx = log.Fold(ctx, "file: %s", file.Name())
	defer log.Debugf(ctx, "written")

	imports := gen.NewFileImports(file.Ast, file.Pkg)

	modified := false

	for _, tf := range file.Funcs {
		fileModified := processMapper(ctx, file.Fset, tf, imports)
		modified = modified || fileModified
	}

	if !modified {
		log.Warnf(ctx, "No modifications applied")
		return
	}

	imports.WriteImportsToAst(file.Fset, file.Ast)
	err := modifyFile(file.Fset, file.Ast)
	if err != nil {
		log.Errorf(ctx, "Failed to generate: %s", err.Error())
	}
}

func processMapper(
	ctx context.Context,
	fset *token.FileSet,
	tf entities.TargetFunc,
	imports *gen.FileImports,
) (astModified bool) {
	ctx = log.Fold(ctx, "Processing func: %s", tf.Name())

	newFuncDecl, err := gen.FuncAst(ctx, tf, fset, imports)
	if err != nil {
		log.Errorf(ctx, "Generation failed: %s", err.Error())
		return false
	}

	f := tf.FuncDecl
	f.Body = newFuncDecl.Body
	f.Type.Params = newFuncDecl.Type.Params
	f.Type.Results = newFuncDecl.Type.Results

	return true
}

func modifyFile(fset *token.FileSet, file *ast.File) error {
	path := fset.Position(file.Pos()).Filename

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0) // read existing permissions?
	if err != nil {
		return err
	}
	defer f.Close()

	return format.Node(f, fset, file)
}
