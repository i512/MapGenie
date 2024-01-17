package gen

type CastPtrToValue struct {
	BaseExpression
	CastWith string
}

func NewCastPtrToValue(base BaseExpression) *CastPtrToValue {
	return &CastPtrToValue{BaseExpression: base}
}

func (c *CastPtrToValue) String(resolver *FileImports) string {
	sourceTemplate :=
		`if input.{{ .InField }} != nil {
			result.{{ .OutField }} ={{" "}}
				 {{- if ne .CastWith "" }}{{ .CastWith }}({{- end }}
				 *input.{{ .InField }}
				 {{- if ne .CastWith "" }}){{- end }}
		}`

	c.CastWith = resolver.ResolveTypeName(c.Out)
	return c.RunTemplate(c, sourceTemplate)
}
