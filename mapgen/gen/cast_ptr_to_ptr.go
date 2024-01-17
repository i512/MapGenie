package gen

import (
	"go/types"
)

type CastPtrToPtr struct {
	BaseExpression
	CastWith string
}

func NewCastPtrToPtr(base BaseExpression) *CastPtrToPtr {
	return &CastPtrToPtr{BaseExpression: base}
}

func (c *CastPtrToPtr) String(resolver *FileImports) string {
	sourceTemplate :=
		`if input.{{ .InField }} != nil {
			{{ .OutField }} := 
				 {{- if ne .CastWith "" }}{{ .CastWith }}({{- end }}
				 *input.{{ .InField }}
				 {{- if ne .CastWith "" }}){{- end }}

			result.{{ .OutField }} = &{{ .OutField }}
		}`

	c.CastWith = c.CastExpression(
		c.In.(*types.Pointer).Elem(),
		c.Out.(*types.Pointer).Elem(),
		resolver,
	)
	return c.RunTemplate(c, sourceTemplate)
}
