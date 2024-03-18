package fragments

type StructAssign struct {
	OutField string
	Fragment Fragment
}

type StructLit struct {
	OutType string
	OutPtr  bool
	Assigns []StructAssign
}

func NewStructLit(OutType string, ptr bool, assigns []StructAssign) *StructLit {
	return &StructLit{
		OutType: OutType,
		Assigns: assigns,
		OutPtr:  ptr,
	}
}

func (f *StructLit) Lines() []string {
	preLines := writer()
	litLines := writer()

	for _, assign := range f.Assigns {
		fragLines := assign.Fragment.Lines()
		if assign.Fragment.ResVar() != nil {
			preLines.W(fragLines)
			litLines.S(assign.OutField, ": ", assign.Fragment.ResVar().Name, ",")
		} else {
			preLines.W(fragLines.PreLastLine())
			litLines.S(assign.OutField, ": ", fragLines.LastLine(), ",")
		}
	}

	takePtr := ""
	if f.OutPtr {
		takePtr = "&"
	}
	return preLines.S("return ", takePtr, f.OutType, "{").W(litLines).S("}").Lines()
}

func (f *StructLit) Deps(r *DependencyRegistry) {
	for _, assign := range f.Assigns {
		r.Register(assign.Fragment)
	}
}

func (f *StructLit) ResVars() []*Var {
	return nil
}
