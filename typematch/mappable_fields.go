package typematch

import (
	"context"
	"go/token"
	"go/types"
	"mapgenie/entities"
	"mapgenie/gen"
	"mapgenie/gen/fragments"
	"mapgenie/pkg/log"
	"reflect"
)

func MappableFields(ctx context.Context, tfs entities.TargetFunc) []entities.Statement {
	in := tfs.In.FieldMap()

	list := make([]entities.Statement, 0)

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

		mapping, ok := createMapping(outFieldName, inFieldType, outFieldType)
		if !ok {
			continue
		}

		list = append(list, mapping)
	}

	return list
}

func createMapping(fieldName string, in, out types.Type) (gen.MapStatement, bool) {
	base := fragments.BaseMapStatement{
		In:       in,
		Out:      out,
		InField:  fieldName,
		OutField: fieldName,
	}

	if typesAreCastable(in, out) {
		return fragments.NewCast(base), true
	}

	outPtr, ok := out.(*types.Pointer)
	if ok && typesAreCastable(in, outPtr.Elem()) {
		return fragments.NewCastValueToPtr(base), true
	}

	outPtr, ok = out.Underlying().(*types.Pointer)
	if ok && typesAreCastable(in, outPtr.Elem()) {
		return fragments.NewCastValueToPtrType(base), true
	}

	inPtr, ok := in.(*types.Pointer)
	if ok && typesAreCastable(inPtr.Elem(), out) {
		return fragments.NewCastPtrToValue(base), true
	}

	inPtr, ok = in.Underlying().(*types.Pointer)
	if ok && typesAreCastable(inPtr.Elem(), out) {
		return fragments.NewCastPtrToValue(base), true
	}

	if inPtr != nil && outPtr != nil && typesAreCastable(inPtr.Elem(), outPtr.Elem()) {
		return fragments.NewCastPtrToPtr(base), true
	}

	if typeIsIntegerOrFloat(in) && isString(out) {
		return fragments.NewAssignNumberToString(base), true
	}

	if isString(in) && typeIsIntegerOrFloat(out) {
		return fragments.NewParseNumberFromString(base), true
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
		return fragments.NewParseTimeFromString(base), true
	}

	if CheckImplementsStringer(in) && isString(out) {
		return fragments.NewAssignStringer(base), true
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
