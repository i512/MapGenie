package fragments

type Fragment interface {
	Deps(registry *DependencyRegistry)
	ResVar() *Var

	Lines() []string
}

type BaseFrag struct{}

func (f BaseFrag) Deps(*DependencyRegistry) {}
func (f BaseFrag) ResVar() *Var {
	return nil
}

type DependencyRegister interface {
	Deps(registry DependencyRegistry)
}
