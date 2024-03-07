package gen

import "mapgenie/gen/fragments"

type LocalVariables struct {
	fileImports *FileImports
}

func NewLocalVariables() *LocalVariables {
	return &LocalVariables{}
}

func (v *LocalVariables) AssignNames(fragments.VarSet) {

}
