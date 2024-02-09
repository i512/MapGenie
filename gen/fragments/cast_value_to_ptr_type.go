package fragments

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
