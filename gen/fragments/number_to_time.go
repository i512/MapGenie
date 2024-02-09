package fragments

import "go/types"

type NumberToTime struct {
	BaseMapStatement
	CastWith string
	TimeName string
}

func NewNumberToTime(base BaseMapStatement) *NumberToTime {
	return &NumberToTime{BaseMapStatement: base}
}

func (c *NumberToTime) Generate(g *GenerationCtx) (string, error) {
	c.TimeName = g.NameResolver.ResolvePkgImport("time")

	if b, ok := c.In.(*types.Basic); ok && b.Kind() != types.Int64 {
		c.CastWith = "int64"
	}

	var sourceTemplate string
	if c.CastWith == "" {
		sourceTemplate = `result.{{ .OutField }} = time.Unix(input.{{ .InField }}, 0).UTC()`
	} else {
		sourceTemplate = `result.{{ .OutField }} = time.Unix({{ .CastWith }}(input.{{ .InField }}), 0).UTC()`
	}

	return c.RunTemplate(c, sourceTemplate)
}
