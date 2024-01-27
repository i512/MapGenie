package typematch

import "go/types"

type receiver interface {
	NumMethods() int
	Method(i int) *types.Func
}

func CheckImplementsStringer(t types.Type) bool {
	named, ok := t.(*types.Named)
	if ok && hasStringerMethod(named) {
		return true
	}

	iface, ok := t.Underlying().(*types.Interface)
	return ok && hasStringerMethod(iface)
}

func hasStringerMethod(r receiver) bool {
	for i := 0; i < r.NumMethods(); i++ {
		if isStringerMethod(r.Method(i)) {
			return true
		}
	}

	return false
}

func isStringerMethod(iface interface{}) bool {
	f, ok := iface.(*types.Func)
	if !ok {
		return false
	}

	if f.Name() != "String" {
		return false
	}

	signature, ok := f.Type().(*types.Signature)
	if !ok {
		return false
	}

	if signature.Params() != nil || signature.Results().Len() != 1 {
		return false
	}

	return isString(signature.Results().At(0).Type())
}
