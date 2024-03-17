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

type NumberToTime2 struct {
	Time *Pkg
	BaseMapStatement
	BaseFrag
}

func NewNumberToTime2(base BaseMapStatement) *NumberToTime2 {
	f := &NumberToTime2{
		BaseMapStatement: base,
		Time:             &Pkg{Path: "time"},
	}

	return f
}

func (f *NumberToTime2) Lines() []string {
	w := writer()

	if b, ok := f.In.(*types.Basic); ok && b.Kind() != types.Int64 {
		w.s(f.Time.LocalName, ".Unix(int64(input.", f.InField, "), 0).UTC()")
	} else {
		w.s(f.Time.LocalName, ".Unix(input.", f.InField, ", 0).UTC()")
	}

	return w.Lines()
}

func (f *NumberToTime2) Deps(registry *DependencyRegistry) {
	registry.Pkg(f.Time)
}
