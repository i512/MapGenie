package fragments

import (
	"reflect"
)

type CastValueToPtrType struct {
	BaseMapStatement
	CastWith string
}

func NewCastValueToPtrType(base BaseMapStatement) *CastValueToPtrType {
	return &CastValueToPtrType{BaseMapStatement: base}
}

func (c *CastValueToPtrType) Generate(g *GenerationCtx) (string, error) {
	sourceTemplate :=
		`{{- if ne .CastWith "" }}
			{{ .OutField }} := {{ .CastWith }}(input.{{ .InField }})
			result.{{ .OutField }} = &{{ .OutField }}
		{{- else }}
			result.{{ .OutField }} = &input.{{ .InField }}
		{{- end }}`

	c.CastWith = c.CastExpression(c.In, c.In.Underlying(), g.NameResolver)

	return c.RunTemplate(c, sourceTemplate)
}

type ValueToPtrType struct {
	CastWith *Type
	Type     *Type
	Var      *Var
	BaseMapStatement
	BaseFrag
}

func NewValueToPtrType(base BaseMapStatement) *ValueToPtrType {
	f := &ValueToPtrType{BaseMapStatement: base}
	f.Type = &Type{Type: f.Out}
	f.Var = &Var{DesiredName: base.OutField}

	inUnderlying := f.In.Underlying()
	if !reflect.DeepEqual(base.In, inUnderlying) {
		f.CastWith = &Type{Type: inUnderlying}
	}
	return f
}

func (f *ValueToPtrType) Lines() []string {
	if f.CastWith == nil {
		return writer().s("&input.", f.InField).Lines()
	}

	return writer().
		s(f.Var.Name, " := ", f.CastWith.LocalName, "(input.", f.InField, ")").
		s("&", f.Var.Name).Lines()
}

func (f *ValueToPtrType) Deps(registry *DependencyRegistry) {
	registry.Type(f.CastWith)
	registry.Type(f.Type)
	registry.Var(f.Var)
}
