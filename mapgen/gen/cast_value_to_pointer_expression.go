package gen

import (
	"go/types"
)

type CastValueToPtr struct {
	BaseExpression
	CastWith string
}

func NewCastValueToPtr(base BaseExpression) *CastValueToPtr {
	return &CastValueToPtr{BaseExpression: base}
}

func (c *CastValueToPtr) String(resolver *FileImports) string {
	sourceTemplate :=
		`{{- if ne .CastWith "" }}
			{{ .OutField }} := {{ .CastWith }}(input.{{ .InField }})
			result.{{ .OutField }} = &{{ .OutField }}
		{{- else }}
			result.{{ .OutField }} = &input.{{ .InField }}
		{{- end }}`

	c.CastWith = c.CastExpression(c.In, c.Out.(*types.Pointer).Elem(), resolver)

	return c.RunTemplate(c, sourceTemplate)
}
