package fragments

import (
	"context"
	"go/types"
)

type GenerationCtx struct {
	Ctx          context.Context
	NameResolver NameResolver
}

type NameResolver interface {
	ResolveTypeName(t types.Type) string
	ResolvePkgImport(pkgPath string) string
}
