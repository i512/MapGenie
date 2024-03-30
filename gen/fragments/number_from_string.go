package fragments

import (
	"github.com/i512/mapgenie/entities"
	"go/types"
)

type NumberFromString struct {
	StrConv *entities.Pkg
	Var     *entities.Var
	BaseMapStatement
	BaseFrag
}

func NewNumberFromString(base BaseMapStatement) *NumberFromString {
	f := &NumberFromString{
		BaseMapStatement: base,
		StrConv:          &entities.Pkg{Path: "strconv"},
		Var:              &entities.Var{DesiredName: base.OutField},
	}

	return f
}

func (f *NumberFromString) Body() entities.Writer {
	fun, args := f.funcAndArgs()
	w := writer()
	if args == "" {
		w.Ln(f.Var.Name, ", _ := ", f.StrConv.LocalName, ".", fun, "(input.", f.InField, ")")
	} else {
		w.Ln(f.Var.Name, ", _ := ", f.StrConv.LocalName, ".", fun, "(input.", f.InField, ", ", args, ")")
	}

	return w
}

func (f *NumberFromString) Result() entities.Writer {
	w := writer()

	castWith := f.CastWith()
	if castWith == "" {
		return w.Ln(f.Var.Name)
	}

	return w.Ln(castWith, "(", f.Var.Name, ")")
}

func (f *NumberFromString) Deps(registry entities.DepReg) {
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
