package gen

import "go/types"

type TimeToNumber struct {
	BaseMapStatement
	CastWith string
}

func NewTimeToNumber(base BaseMapStatement) *TimeToNumber {
	return &TimeToNumber{BaseMapStatement: base}
}

func (c *TimeToNumber) String(resolver *FileImports) string {
	if b, ok := c.Out.(*types.Basic); ok && b.Kind() != types.Int64 {
		c.CastWith = c.Out.String()
	}

	var sourceTemplate string
	if c.CastWith == "" {
		sourceTemplate = `result.{{ .OutField }} = input.{{ .InField }}.Unix()`
	} else {
		sourceTemplate = `result.{{ .OutField }} = {{ .CastWith }}(input.{{ .InField }}.Unix())`
	}

	return c.RunTemplate(c, sourceTemplate)
}
