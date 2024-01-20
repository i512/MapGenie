package gen

type CastValueToPtrType struct {
	BaseMapStatement
	CastWith string
}

func NewCastValueToPtrType(base BaseMapStatement) *CastValueToPtrType {
	return &CastValueToPtrType{BaseMapStatement: base}
}

func (c *CastValueToPtrType) String(resolver *FileImports) string {
	sourceTemplate :=
		`{{- if ne .CastWith "" }}
			{{ .OutField }} := {{ .CastWith }}(input.{{ .InField }})
			result.{{ .OutField }} = &{{ .OutField }}
		{{- else }}
			result.{{ .OutField }} = &input.{{ .InField }}
		{{- end }}`

	c.CastWith = c.CastExpression(c.In, c.In.Underlying(), resolver)

	return c.RunTemplate(c, sourceTemplate)
}
