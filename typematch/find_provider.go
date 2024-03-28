package typematch

import (
	"go/types"
	"reflect"
)

func FindProvider(in, out types.Type, providers []*types.Func) *types.Func {
	for _, pf := range providers {
		signature := pf.Type().(*types.Signature)
		nParams := signature.Params().Len()
		nResults := signature.Results().Len()
		if nParams > 1 || nResults == 0 || nResults > 2 {
			continue
		}

		if !reflect.DeepEqual(signature.Params().At(0).Type(), in) || !reflect.DeepEqual(signature.Results().At(0).Type(), out) {
			continue
		}

		if nResults == 2 && signature.Results().At(1).Type().String() != "error" {
			continue
		}

		return pf
	}

	return nil
}
