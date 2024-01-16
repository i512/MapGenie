package entities

import "go/types"

type TargetFuncSignature struct {
	In, Out Argument
}

type Argument struct {
	Named  *types.Named
	Struct *types.Struct
	IsPtr  bool
	Local  bool
}

func (s Argument) FieldMap() map[string]types.Type {
	result := map[string]types.Type{}

	for i := 0; i < s.Struct.NumFields(); i++ {
		f := s.Struct.Field(i)
		result[f.Name()] = f.Type()
	}

	return result
}
