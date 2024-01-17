package gen

type Cast struct {
	BaseExpression
	CastWith string
}

func NewCast(base BaseExpression) *Cast {
	return &Cast{BaseExpression: base}
}

func (c *Cast) String(resolver *FileImports) string {
	c.CastWith = c.CastExpression(c.In, c.Out, resolver)
	sourceTemplate :=
		`result.{{ .OutField }} = 
			 {{- if ne .CastWith "" }}{{ .CastWith }}({{- end }}
			 input.{{ .InField }}
			 {{- if ne .CastWith "" }}){{- end }}`

	return c.RunTemplate(c, sourceTemplate)
}
