package gen

type Cast struct {
	BaseMapStatement
	CastWith string
}

func NewCast(base BaseMapStatement) *Cast {
	return &Cast{BaseMapStatement: base}
}

func (c *Cast) Generate(resolver *FileImports) (string, error) {
	c.CastWith = c.CastExpression(c.In, c.Out, resolver)
	sourceTemplate :=
		`result.{{ .OutField }} = 
			 {{- if ne .CastWith "" }}{{ .CastWith }}({{- end }}
			 input.{{ .InField }}
			 {{- if ne .CastWith "" }}){{- end }}`

	return c.RunTemplate(c, sourceTemplate)
}
