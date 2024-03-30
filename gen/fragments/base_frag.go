package fragments

import "github.com/i512/mapgenie/entities"

type BaseFrag struct{}

func (f BaseFrag) Deps(entities.DepReg) {}
func (f BaseFrag) ResVar() *entities.Var {
	return nil
}

func (f BaseFrag) Body() entities.Writer {
	return writer()
}
