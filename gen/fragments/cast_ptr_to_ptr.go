package fragments

import (
	"go/types"
	"reflect"
)

type CastPtrToPtr struct {
	BaseMapStatement
	CastWith string
}

func NewCastPtrToPtr(base BaseMapStatement) *CastPtrToPtr {
	return &CastPtrToPtr{BaseMapStatement: base}
}

func (c *CastPtrToPtr) Generate(g *GenerationCtx) (string, error) {
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
		g.NameResolver,
	)
	return c.RunTemplate(c, sourceTemplate)
}

type PtrToPtr struct {
	CastWith *Type
	Var      *Var
	BaseMapStatement
	BaseFrag
}

func NewPtrToPtr(base BaseMapStatement) *PtrToPtr {
	f := &PtrToPtr{BaseMapStatement: base}

	outElemType := base.Out.(*types.Pointer).Elem()
	if !reflect.DeepEqual(base.In, outElemType) {
		f.CastWith = &Type{Type: outElemType}
		f.Var = &Var{DesiredName: base.OutField}
	}
	return f
}

func (f *PtrToPtr) Lines() []string {
	w := writer()
	w.s("var ", f.Var.Name, " ", f.CastWith.LocalName)
	w.s("if input.", f.InField, " != nil {").ident(func(w *LineWriter) {
		w.s(f.Var.Name, " = ").
			a((&CastFrag{Fragment: &OneLinerFrag{Line: "*input." + f.InField}, CastWith: f.CastWith.LocalName}).Lines())
	}).s("}")
	w.s("&", f.Var.Name)

	return w.Lines()
}

func (f *PtrToPtr) Deps(registry *DependencyRegistry) {
	registry.Type(f.CastWith)
	registry.Var(f.Var)
}
