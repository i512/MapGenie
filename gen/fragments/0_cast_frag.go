package fragments

type CastFrag struct {
	Fragment Fragment
	CastWith string
}

func NewCastF(fragment Fragment, castWith string) *CastFrag {
	return &CastFrag{
		Fragment: fragment,
		CastWith: castWith,
	}
}

func (f *CastFrag) Lines() []string {
	if f.Fragment.ResVar() != nil {
		if f.CastWith == "" {
			return f.Fragment.Lines()
		}

		return append(f.Fragment.Lines(), f.CastWith+"("+f.Fragment.ResVar().Name+")")
	}

	if f.CastWith == "" {
		return f.Fragment.Lines()
	}

	lines := f.Fragment.Lines()
	return append(lines[:len(lines)-1], f.CastWith+"("+lines[len(lines)-1]+")")
}
