package fragments

import (
	"github.com/i512/mapgenie/entities"
	"go/types"
)

type ProviderCall struct {
	BaseMapStatement
	BaseFrag
	provider *types.Func
	Var      *entities.Var
}

func NewProviderCall(base BaseMapStatement, provider *types.Func) *ProviderCall {
	f := &ProviderCall{
		BaseMapStatement: base,
		provider:         provider,
	}

	if f.returnsErr() {
		f.Var = &entities.Var{DesiredName: f.OutField}
	}

	return f
}

func (f *ProviderCall) returnsErr() bool {
	return f.provider.Type().(*types.Signature).Results().Len() != 1

}

func (f *ProviderCall) Body() entities.Writer {
	if f.Var == nil {
		return writer()
	}

	return writer().Lnf("%s, _ := %s(input.%s)", f.Var.Name, f.provider.Name(), f.InField)

}

func (f *ProviderCall) Result() entities.Writer {
	if f.Var == nil {
		return writer().Lnf("%s(input.%s)", f.provider.Name(), f.InField)
	}

	return writer().Ln(f.Var.Name)
}

func (f *ProviderCall) Deps(registry entities.DepReg) {
	registry.Var(f.Var)
}
