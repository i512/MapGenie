package fragments

type StructAssign struct {
	OutField string
	Fragment Fragment
}

type StructLitFrag struct {
	OutType string
	OutPtr  bool
	Assigns []StructAssign
}

func NewStructLit(OutType string, ptr bool, assigns []StructAssign) *StructLitFrag {
	return &StructLitFrag{
		OutType: OutType,
		Assigns: assigns,
		OutPtr:  ptr,
	}
}

func (f *StructLitFrag) Lines() []string {
	preLines := writer()
	litLines := writer()

	for _, assign := range f.Assigns {
		fragLines := assign.Fragment.Lines()
		if assign.Fragment.ResVar() != nil {
			preLines.a(fragLines)
			litLines.s(assign.OutField, ": ", assign.Fragment.ResVar().Name, ",")
		} else {
			preLines.a(fragLines[:len(fragLines)-1])
			litLines.s(assign.OutField, ": ", fragLines[len(fragLines)-1], ",")
		}
	}

	takePtr := ""
	if f.OutPtr {
		takePtr = "&"
	}
	return preLines.s("return ", takePtr, f.OutType, "{").w(litLines).s("}").Lines()
}

func (f *StructLitFrag) VarSet(set VarSet) {
	for _, assign := range f.Assigns {
		assign.Fragment.VarSet(set)
	}
}

func (f *StructLitFrag) TypeSet(set TypeSet) {
	for _, assign := range f.Assigns {
		assign.Fragment.TypeSet(set)
	}
}

func (f *StructLitFrag) ResVars() []*Var {
	return nil
}
