package fragments

import (
	"github.com/i512/mapgenie/entities"
	"reflect"
)

type ValueToPointerType struct {
	CastWith *entities.Type
	Type     *entities.Type
	Var      *entities.Var
	BaseMapStatement
	BaseFrag
}

func NewValueToPointerType(base BaseMapStatement) *ValueToPointerType {
	f := &ValueToPointerType{BaseMapStatement: base}
	f.Type = &entities.Type{Type: f.Out}
	f.Var = &entities.Var{DesiredName: base.OutField}

	inUnderlying := f.In.Underlying()
	if !reflect.DeepEqual(base.In, inUnderlying) {
		f.CastWith = &entities.Type{Type: inUnderlying}
	}
	return f
}

func (f *ValueToPointerType) Body() entities.Writer {
	if f.CastWith == nil {
		return writer()
	}

	return writer().Ln(f.Var.Name, " := ", f.CastWith.LocalName, "(input.", f.InField, ")")
}

func (f *ValueToPointerType) Result() entities.Writer {
	if f.CastWith == nil {
		return writer().Ln("&input.", f.InField)
	}

	return writer().Ln("&", f.Var.Name)
}

func (f *ValueToPointerType) Deps(registry entities.DepReg) {
	registry.Type(f.CastWith)
	registry.Type(f.Type)
	registry.Var(f.Var)
}
