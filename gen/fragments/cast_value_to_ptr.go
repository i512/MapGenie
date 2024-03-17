package fragments

import (
	"go/types"
	"reflect"
)

type CastValueToPtr struct {
	BaseMapStatement
	CastWith string
}

func NewCastValueToPtr(base BaseMapStatement) *CastValueToPtr {
	return &CastValueToPtr{BaseMapStatement: base}
}

func (c *CastValueToPtr) Generate(g *GenerationCtx) (string, error) {
	sourceTemplate :=
		`{{- if ne .CastWith "" }}
			{{ .OutField }} := {{ .CastWith }}(input.{{ .InField }})
			result.{{ .OutField }} = &{{ .OutField }}
		{{- else }}
			result.{{ .OutField }} = &input.{{ .InField }}
		{{- end }}`

	c.CastWith = c.CastExpression(c.In, c.Out.(*types.Pointer).Elem(), g.NameResolver)

	return c.RunTemplate(c, sourceTemplate)
}

type ValueToPtr struct {
	CastWith *Type
	Var      *Var
	BaseMapStatement
	BaseFrag
}

func NewValueToPtr(base BaseMapStatement) *ValueToPtr {
	f := &ValueToPtr{BaseMapStatement: base}

	outElemType := base.Out.(*types.Pointer).Elem()
	if !reflect.DeepEqual(base.In, outElemType) {
		f.CastWith = &Type{Type: outElemType}
		f.Var = &Var{DesiredName: base.OutField}
	}
	return f
}

func (f *ValueToPtr) Lines() []string {
	if f.CastWith == nil {
		return writer().s("&input.", f.InField).Lines()
	}

	return writer().
		s(f.Var.Name, " := ", f.CastWith.LocalName, "(input.", f.InField, ")").
		s("&", f.Var.Name).Lines()
}

func (f *ValueToPtr) Deps(registry *DependencyRegistry) {
	registry.Type(f.CastWith)
	registry.Var(f.Var)
}
