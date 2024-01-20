package gen

type Cast struct {
	BaseMapStatement
	CastWith string
}

func NewCast(base BaseMapStatement) *Cast {
	return &Cast{BaseMapStatement: base}
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
