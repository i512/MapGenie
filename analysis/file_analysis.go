package analysis

import (
	"context"
	"errors"
	"fmt"
	"github.com/i512/mapgenie/entities"
	"github.com/i512/mapgenie/pkg/log"
	"github.com/i512/mapgenie/typematch"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
	"path/filepath"
	"regexp"
)

var targetFuncMagicComment = regexp.MustCompile(`^\w+ map this pls`)
var ErrFuncMismatchError = fmt.Errorf("function is not mappable")

type FileAnalysis struct {
	pkg  *packages.Package
	fset *token.FileSet
	ast  *ast.File

	providers []*types.Func
}

func NewFileAnalysis(pkg *packages.Package, fset *token.FileSet, ast *ast.File, providers []*types.Func) *FileAnalysis {
	return &FileAnalysis{
		pkg:  pkg,
		fset: fset,
		ast:  ast,

		providers: providers,
	}
}

func (file *FileAnalysis) name() string {
	return file.fset.File(file.ast.Pos()).Name()
}

func (file *FileAnalysis) FindTargets(ctx context.Context) *entities.TargetFile {
	ctx = log.Fold(ctx, "file: %s", filepath.Base(file.name()))

	funcs := make([]entities.TargetFunc, 0)

	astutil.Apply(file.ast, nil, func(cursor *astutil.Cursor) bool {
		funcDecl, ok := cursor.Node().(*ast.FuncDecl)
		if !ok || !file.isTargetFunc(funcDecl) {
			return true
		}

		f := file.targetFunc(ctx, funcDecl)
		if f != nil {
			funcs = append(funcs, *f)
		}

		return true
	})

	if len(funcs) == 0 {
		return nil
	}

	return &entities.TargetFile{
		Pkg:   file.pkg,
		Fset:  file.fset,
		Ast:   file.ast,
		Funcs: funcs,
	}
}

func (file *FileAnalysis) isTargetFunc(f *ast.FuncDecl) bool {
	return f.Doc != nil && targetFuncMagicComment.MatchString(f.Doc.Text())
}

func (file *FileAnalysis) targetFunc(ctx context.Context, f *ast.FuncDecl) *entities.TargetFunc {
	ctx = log.Fold(ctx, "mapper func: %s (line %d)", f.Name.String(), file.fset.Position(f.Pos()).Line)

	target, err := file.arguments(ctx, f)
	if errors.Is(err, ErrFuncMismatchError) {
		log.Errorf(ctx, "Invalid mapper: %s", err.Error())
		return nil
	}
	if err != nil {
		log.Errorf(ctx, "Failed to interpret func signature: %s", err.Error())
		return nil
	}

	target.Fragments = typematch.MappableFields(ctx, target, file.providers)

	return &target
}

func (file *FileAnalysis) arguments(ctx context.Context, f *ast.FuncDecl) (tf entities.TargetFunc, err error) {
	tf.FuncDecl = f

	funcType := file.pkg.Types.Scope().Lookup(f.Name.Name)
	signature := funcType.Type().(*types.Signature)
	tf.In, err = file.argument(signature.Params())
	if err != nil {
		return tf, fmt.Errorf("bad argument: %w", err)
	}

	tf.Out, err = file.argument(signature.Results())
	if err != nil {
		return tf, fmt.Errorf("bad return argument: %w", err)
	}

	return tf, nil
}

func (file *FileAnalysis) argument(tuple *types.Tuple) (arg entities.Argument, err error) {
	if tuple.Len() != 1 {
		return arg, fmt.Errorf("signature must have a single argument, have: %s. %w", tuple.String(), ErrFuncMismatchError)
	}

	firstArg := tuple.At(0).Type()
	if ptr, ok := firstArg.(*types.Pointer); ok {
		arg.IsPtr = true
		firstArg = ptr.Elem()
	}
	named, ok := firstArg.(*types.Named)
	if !ok {
		return arg, fmt.Errorf("`%s` is not a struct. %w", firstArg.String(), ErrFuncMismatchError)
	}

	structArg, ok := named.Underlying().(*types.Struct)
	if !ok {
		return arg, fmt.Errorf("`%s` is not a struct. %w", named.String(), ErrFuncMismatchError)
	}

	arg.Named = named
	arg.Struct = structArg
	arg.Local = arg.Named.Obj().Pkg().Path() == file.pkg.PkgPath

	return arg, nil
}
