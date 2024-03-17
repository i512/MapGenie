package fragments

import (
	"go/types"
)

type BaseMapStatement struct {
	In, Out           types.Type
	InField, OutField string
}
