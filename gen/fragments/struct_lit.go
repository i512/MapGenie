package fragments

import (
	"github.com/i512/mapgenie/entities"
	"go/types"
)

type StructAssign struct {
	OutField string
	Fragment entities.Fragment
}

type StructLit struct {
	OutType *entities.Type
	OutPtr  bool
	Assigns []StructAssign
	BaseFrag
}

func NewStructLit(outType types.Type, ptr bool, assigns []StructAssign) *StructLit {
	return &StructLit{
		OutType: &entities.Type{Type: outType},
		Assigns: assigns,
		OutPtr:  ptr,
	}
}

func (f *StructLit) Body() entities.Writer {
	body := writer()

	for _, assign := range f.Assigns {
		body.Merge(assign.Fragment.Body())
	}

	return body
}

func (f *StructLit) Result() entities.Writer {
	w := writer()
	w.Lnf("%s%s{", opt(f.OutPtr, "&"), f.OutType.LocalName).Indent(func(w entities.Writer) {
		for _, assign := range f.Assigns {
			w.Ln(assign.OutField, ": ", assign.Fragment.Result().String(), ",")
		}
	}).Ln("}")

	return w
}

func (f *StructLit) Deps(r entities.DepReg) {
	r.Type(f.OutType)
	for _, assign := range f.Assigns {
		r.Register(assign.Fragment)
	}
}
