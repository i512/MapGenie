package fragments

import (
	"go/types"
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
