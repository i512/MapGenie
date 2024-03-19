package fragments

import (
	"mapgenie/entities"
)

type MapFunc struct {
	tf        entities.TargetFunc
	structLit *StructLit
	in, out   *entities.Type
	BaseFrag
}

func NewMapFunc(tf entities.TargetFunc) *MapFunc {
	f := &MapFunc{
		tf:  tf,
		in:  &entities.Type{Type: tf.In.Named},
		out: &entities.Type{Type: tf.Out.Named},
	}
	f.structLit = f.structFrag()

	return f
}

func (f *MapFunc) Result() entities.Writer {
	inPtr := f.tf.In.IsPtr
	outPtr := f.tf.Out.IsPtr
	w := writer()
	w.Lnf("func %s(input %s%s) %s%s {", f.tf.Name(), opt(inPtr, "*"), f.in.LocalName, opt(outPtr, "*"), f.out.LocalName).Indent(func(w entities.Writer) {
		if f.tf.In.IsPtr {
			w.Ln("if input == nil {").Indent(func(w entities.Writer) {
				w.Lnf("return %s%s{}", opt(outPtr, "&"), f.out.LocalName)
			}).Ln("}")
		}

		w.Merge(f.structLit.Body())
		w.Ln("return ", f.structLit.Result().String())
	}).Ln("}")

	return w
}

func (f *MapFunc) structFrag() *StructLit {
	assigns := make([]StructAssign, 0)
	for outField, Fragment := range f.tf.Fragments {
		assigns = append(assigns, StructAssign{OutField: outField, Fragment: Fragment})
	}

	return NewStructLit(f.tf.Out.Named, f.tf.Out.IsPtr, assigns)
}

func (f *MapFunc) Deps(registry entities.DepReg) {
	registry.Register(f.structLit)
	registry.Type(f.in, f.out)
}

func opt(b bool, s string) string {
	if b {
		return s
	}

	return ""
}
