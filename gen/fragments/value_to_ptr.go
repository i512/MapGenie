package fragments

import (
	"github.com/i512/mapgenie/entities"
	"go/types"
	"reflect"
)

type ValueToPtr struct {
	CastWith *entities.Type
	Var      *entities.Var
	BaseMapStatement
	BaseFrag
}

func NewValueToPtr(base BaseMapStatement) *ValueToPtr {
	f := &ValueToPtr{BaseMapStatement: base}

	outElemType := base.Out.(*types.Pointer).Elem()
	if !reflect.DeepEqual(base.In, outElemType) {
		f.CastWith = &entities.Type{Type: outElemType}
		f.Var = &entities.Var{DesiredName: base.OutField}
	}
	return f
}

func (f *ValueToPtr) Body() entities.Writer {
	if f.CastWith == nil {
		return writer()
	}

	return writer().Ln(f.Var.Name, " := ", f.CastWith.LocalName, "(input.", f.InField, ")")
}

func (f *ValueToPtr) Result() entities.Writer {
	if f.CastWith == nil {
		return writer().Ln("&input.", f.InField)
	}

	return writer().Ln("&", f.Var.Name)
}

func (f *ValueToPtr) Deps(registry entities.DepReg) {
	registry.Type(f.CastWith)
	registry.Var(f.Var)
}
