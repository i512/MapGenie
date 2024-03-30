package typematch

import (
	"context"
	"github.com/i512/mapgenie/entities"
	"github.com/i512/mapgenie/gen/fragments"
	"github.com/i512/mapgenie/pkg/log"
	"go/token"
	"go/types"
	"reflect"
)

func MappableFields(ctx context.Context, tfs entities.TargetFunc, providers []*types.Func) map[string]entities.Fragment {
	in := tfs.In.FieldMap()

	list := make(map[string]entities.Fragment)

	for i := 0; i < tfs.Out.Struct.NumFields(); i++ {
		field := tfs.Out.Struct.Field(i)
		outFieldName := field.Name()
		outFieldType := field.Type()

		inFieldType, ok := in[outFieldName]
		if !ok {
			log.Debugf(ctx, "No matching field for: %s", outFieldName)
			continue
		}

		if !token.IsExported(outFieldName) && !(tfs.In.Local && tfs.Out.Local) {
			log.Debugf(ctx, "Output field is unexported: %s", outFieldName)
			continue
		}

		mapping, ok := createMapping(outFieldName, inFieldType, outFieldType, providers)
		if !ok {
			continue
		}

		list[outFieldName] = mapping
	}

	return list
}

func createMapping(fieldName string, in, out types.Type, providers []*types.Func) (entities.Fragment, bool) {
	base := fragments.BaseMapStatement{
		In:       in,
		Out:      out,
		InField:  fieldName,
		OutField: fieldName,
	}

	if typesAreCastable(in, out) {
		return fragments.NewCast(base), true
	}

	if provider := FindProvider(in, out, providers); provider != nil {
		return fragments.NewProviderCall(base, provider), true
	}

	outPtr, ok := out.(*types.Pointer)
	if ok && typesAreCastable(in, outPtr.Elem()) {
		return fragments.NewValueToPtr(base), true
	}

	outPtr, ok = out.Underlying().(*types.Pointer)
	if ok && typesAreCastable(in, outPtr.Elem()) {
		return fragments.NewValueToPointerType(base), true
	}

	inPtr, ok := in.(*types.Pointer)
	if ok && typesAreCastable(inPtr.Elem(), out) {
		return fragments.NewPtrToValue(base), true
	}

	inPtr, ok = in.Underlying().(*types.Pointer)
	if ok && typesAreCastable(inPtr.Elem(), out) {
		return fragments.NewPtrToValue(base), true
	}

	if inPtr != nil && outPtr != nil && typesAreCastable(inPtr.Elem(), outPtr.Elem()) {
		return fragments.NewPtrToPtr(base), true
	}

	if typeIsIntegerOrFloat(in) && isString(out) {
		return fragments.NewNumberToString(base), true
	}

	if isString(in) && typeIsIntegerOrFloat(out) {
		return fragments.NewNumberFromString(base), true
	}

	if isTime(in) && isBasic(out, types.Int, types.Int64) {
		return fragments.NewTimeToNumber(base), true
	}

	if isBasic(in, types.Int, types.Int64) && isTime(out) {
		return fragments.NewNumberToTime(base), true
	}

	if isTime(in) && isString(out) {
		return fragments.NewTimeToString(base), true
	}

	if isString(in) && isTime(out) {
		return fragments.NewStringToTime(base), true
	}

	if CheckImplementsStringer(in) && isString(out) {
		return fragments.NewStringerToString(base), true
	}

	return nil, false
}

func isBasic(t types.Type, kinds ...types.BasicKind) bool {
	basic, ok := t.(*types.Basic)
	if !ok {
		return false
	}

	for _, kind := range kinds {
		if basic.Kind() == kind {
			return true
		}
	}

	return false
}

func typesAreCastable(in, out types.Type) bool {
	in = getUnderlying(in)
	out = getUnderlying(out)

	same := reflect.DeepEqual(in, out)
	numbers := typeIsIntegerOrFloat(in) && typeIsIntegerOrFloat(out)
	bytes := typeIsStringOrByteSlice(in) && typeIsStringOrByteSlice(out)
	return same || numbers || bytes
}

func typeIsIntegerOrFloat(t types.Type) bool {
	basic, isBasic := t.(*types.Basic)
	return isBasic && (basic.Info()&types.IsInteger > 0 || basic.Info()&types.IsFloat > 0)
}

func typeIsStringOrByteSlice(t types.Type) bool {
	return isString(t) || typeIsByteSlice(t)
}

func isString(t types.Type) bool {
	return isBasic(t, types.String)
}

func typeIsByteSlice(t types.Type) bool {
	slice, ok := t.(*types.Slice)
	if !ok {
		return false
	}

	basic, ok := slice.Elem().(*types.Basic)
	return ok && basic.Kind()&types.Byte > 0
}

func getUnderlying(t types.Type) types.Type {
	for t != t.Underlying() {
		t = t.Underlying()
	}
	return t
}

func isTime(t types.Type) bool {
	return t.String() == "time.Time"
}
