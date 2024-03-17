package fragments

import "go/types"

type ParseNumberFromString struct {
	BaseMapStatement
	StrconvName string
	ConvFunc    string
	FuncArgs    string
}

func NewParseNumberFromString(base BaseMapStatement) *ParseNumberFromString {
	return &ParseNumberFromString{BaseMapStatement: base}
}

func (c *ParseNumberFromString) Generate(g *GenerationCtx) (string, error) {
	c.StrconvName = g.NameResolver.ResolvePkgImport("strconv")
	c.ConvFunc, c.FuncArgs = c.funcAndArgs()

	sourceTemplate := `
{{ .OutField }}, err := 
{{ if ne .FuncArgs "" }}
	{{ .StrconvName }}.{{ .ConvFunc }}(input.{{ .InField }}, {{ .FuncArgs }})
{{ else }}
	{{ .StrconvName }}.{{ .ConvFunc }}(input.{{ .InField }})
{{ end }}

if err == nil {
	{{ if ne .CastWith "" }}
		result.{{ .OutField }} = {{ .CastWith }}({{ .OutField }})
	{{ else }}
		result.{{ .OutField }} = {{ .OutField }}
	{{ end }}
}
`

	return c.RunTemplate(c, sourceTemplate)
}

func (c *ParseNumberFromString) CastWith() string {
	kind := c.Out.(*types.Basic).Kind()
	if !(kind == types.Int || kind == types.Int64 || kind == types.Uint64 || kind == types.Float64) {
		return c.Out.String()
	}

	return ""
}

func (c *ParseNumberFromString) funcAndArgs() (string, string) {
	basic := c.Out.(*types.Basic)

	if basic.Info()&types.IsInteger > 0 {
		if basic.Kind() == types.Int {
			return "Atoi", ""
		}
		if basic.Info()&types.IsUnsigned > 0 {
			return "ParseUint", "10, " + c.IntBits(basic)
		} else {
			return "ParseInt", "10, " + c.IntBits(basic)
		}
	}

	if basic.Kind() == types.Float64 {
		return "ParseFloat", "64"
	}

	if basic.Kind() == types.Float32 {
		return "ParseFloat", "32"
	}

	return "unknown", "unknown"
}

func (c *ParseNumberFromString) IntBits(basic *types.Basic) string {
	switch basic.Kind() {
	case types.Int8, types.Uint8:
		return "8"
	case types.Int16, types.Uint16:
		return "16"
	case types.Int32, types.Uint32:
		return "32"
	case types.Int, types.Uint, types.Int64, types.Uint64:
		return "64"
	default:
		return ""
	}
}

type NumberFromString struct {
	StrConv *Pkg
	Var     *Var
	BaseMapStatement
	BaseFrag
}

func NewNumberFromString(base BaseMapStatement) *NumberFromString {
	f := &NumberFromString{
		BaseMapStatement: base,
		StrConv:          &Pkg{Path: "strconv"},
		Var:              &Var{DesiredName: base.OutField},
	}

	return f
}

func (f *NumberFromString) Lines() []string {
	fun, args := f.funcAndArgs()
	w := writer()
	if args == "" {
		w.s(f.Var.Name, ", _ := ", f.StrConv.LocalName, ".", fun, "(input.", f.InField, ")")
	} else {
		w.s(f.Var.Name, ", _ := ", f.StrConv.LocalName, ".", fun, "(input.", f.InField, ", ", args, ")")
	}

	castWith := f.CastWith()
	if castWith == "" {
		w.s(f.Var.Name)
	} else {
		w.s(castWith, "(", f.Var.Name, ")")
	}
	return w.Lines()
}

func (f *NumberFromString) Deps(registry *DependencyRegistry) {
	registry.Pkg(f.StrConv)
	registry.Var(f.Var)
}

func (c *NumberFromString) CastWith() string {
	kind := c.Out.(*types.Basic).Kind()
	if !(kind == types.Int || kind == types.Int64 || kind == types.Uint64 || kind == types.Float64) {
		return c.Out.String()
	}

	return ""
}

func (c *NumberFromString) funcAndArgs() (string, string) {
	basic := c.Out.(*types.Basic)

	if basic.Info()&types.IsInteger > 0 {
		if basic.Kind() == types.Int {
			return "Atoi", ""
		}
		if basic.Info()&types.IsUnsigned > 0 {
			return "ParseUint", "10, " + c.IntBits(basic)
		} else {
			return "ParseInt", "10, " + c.IntBits(basic)
		}
	}

	if basic.Kind() == types.Float64 {
		return "ParseFloat", "64"
	}

	if basic.Kind() == types.Float32 {
		return "ParseFloat", "32"
	}

	return "unknown", "unknown"
}

func (c *NumberFromString) IntBits(basic *types.Basic) string {
	switch basic.Kind() {
	case types.Int8, types.Uint8:
		return "8"
	case types.Int16, types.Uint16:
		return "16"
	case types.Int32, types.Uint32:
		return "32"
	case types.Int, types.Uint, types.Int64, types.Uint64:
		return "64"
	default:
		return ""
	}
}
