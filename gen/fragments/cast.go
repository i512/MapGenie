package fragments

import "reflect"

type Cast struct {
	BaseMapStatement
	CastWith string
}

func NewCast(base BaseMapStatement) *Cast {
	return &Cast{BaseMapStatement: base}
}

func (c *Cast) Generate(g *GenerationCtx) (string, error) {
	c.CastWith = c.CastExpression(c.In, c.Out, g.NameResolver)
	sourceTemplate :=
		`result.{{ .OutField }} = 
			 {{- if ne .CastWith "" }}{{ .CastWith }}({{- end }}
			 input.{{ .InField }}
			 {{- if ne .CastWith "" }}){{- end }}`

	return c.RunTemplate(c, sourceTemplate)
}

type CastField struct {
	BaseFrag
	BaseMapStatement
	outType *Type
}

func NewCastField(base BaseMapStatement) *CastField {
	f := &CastField{
		BaseMapStatement: base,
	}

	if !reflect.DeepEqual(base.In, base.Out) {
		f.outType = &Type{Type: base.Out}
	}

	return f
}

func (f *CastField) Lines() []string {
	castWith := ""
	if f.outType != nil {
		castWith = f.outType.LocalName
	}
	return NewCastF(&OneLinerFrag{Line: "input." + f.InField}, castWith).Lines()
}

func (f *CastField) TypeSet(s TypeSet) {
	s[f.outType] = struct{}{}
}
