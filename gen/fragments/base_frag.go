package fragments

import "mapgenie/entities"

type BaseFrag struct{}

func (f BaseFrag) Deps(entities.DepReg) {}
func (f BaseFrag) ResVar() *entities.Var {
	return nil
}

func (f BaseFrag) Body() entities.Writer {
	return writer()
}
