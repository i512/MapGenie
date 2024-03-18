package fragments

type Writer interface {
	Ln(...string) Writer
	Merge(Writer) Writer
	Indent(func(Writer)) Writer
	String() string
	Lines() []string
	PreLastLine() Writer
	LastLine() string
}

type Fragment interface {
	Deps(registry *DependencyRegistry)
	ResVar() *Var

	Lines() Writer
}

type BaseFrag struct{}

func (f BaseFrag) Deps(*DependencyRegistry) {}
func (f BaseFrag) ResVar() *Var {
	return nil
}

type DependencyRegister interface {
	Deps(registry DependencyRegistry)
}
