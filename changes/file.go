package changes

import (
	"context"
	"go/ast"
	"go/format"
	"go/token"
	"mapgenie/entities"
	"mapgenie/gen"
	"mapgenie/pkg/log"
	"os"
)

func ApplyFilesChanges(ctx context.Context, files []entities.TargetFile) {
	for _, file := range files {
		ApplyFile(ctx, file)
	}
}

func ApplyFile(ctx context.Context, file entities.TargetFile) {
	imports := gen.NewFileImports(file.Ast, file.Pkg)

	modified := false

	for _, tf := range file.Funcs {
		fileModified := processMapper(ctx, file.Fset, tf.FuncDecl, tf, imports)
		modified = modified || fileModified
	}

	if modified {
		imports.WriteImportsToAst(file.Fset, file.Ast)
		modifyFile(file.Fset, file.Ast)
	}
}

func processMapper(
	ctx context.Context,
	fset *token.FileSet,
	funcDecl *ast.FuncDecl,
	tf entities.TargetFunc,
	imports *gen.FileImports,
) (astModified bool) {
	data := gen.MapTemplateData{
		FuncName: funcDecl.Name.Name,
		InType:   imports.ResolveTypeName(tf.In.Named), // TODO: move type resolution to gen
		InIsPtr:  tf.In.IsPtr,
		OutType:  imports.ResolveTypeName(tf.Out.Named),
		OutIsPtr: tf.Out.IsPtr,
		Maps:     tf.Statements,
		Resolver: imports,
	}

	newFuncDecl, err := gen.MapperFuncAst(ctx, fset, data)
	if err != nil {
		log.Errorf(ctx, "Generation failed: %s", err.Error())
		return false
	}

	funcDecl.Body = newFuncDecl.Body
	funcDecl.Type.Params = newFuncDecl.Type.Params
	funcDecl.Type.Results = newFuncDecl.Type.Results

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
