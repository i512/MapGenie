package fragments

import "go/types"

type TimeToNumber struct {
	BaseMapStatement
	CastWith string
}

func NewTimeToNumber(base BaseMapStatement) *TimeToNumber {
	return &TimeToNumber{BaseMapStatement: base}
}

func (c *TimeToNumber) Generate(_ *GenerationCtx) (string, error) {
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

type TimeToNumber2 struct {
	BaseMapStatement
	BaseFrag
}

func NewTimeToNumber2(base BaseMapStatement) *TimeToNumber2 {
	f := &TimeToNumber2{
		BaseMapStatement: base,
	}

	return f
}

func (f *TimeToNumber2) Lines() []string {
	w := writer()

	if b, ok := f.Out.(*types.Basic); ok && b.Kind() != types.Int64 {
		castWith := f.Out.String()
		w.s(castWith, "(input.", f.InField, ".Unix())")
	} else {
		w.s("input.", f.InField, ".Unix()")
	}

	return w.Lines()
}
