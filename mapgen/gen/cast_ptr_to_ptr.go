package gen

import (
	"go/types"
)

type CastPtrToPtr struct {
	BaseMapStatement
	CastWith string
}

func NewCastPtrToPtr(base BaseMapStatement) *CastPtrToPtr {
	return &CastPtrToPtr{BaseMapStatement: base}
}

func (c *CastPtrToPtr) Generate(resolver *FileImports) (string, error) {
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
