package typematch

import (
	"fmt"
	"go/token"
	"go/types"
	"mapgenie/mapgen/entities"
	"mapgenie/mapgen/gen"
	"reflect"
)

func MappableFields(tfs entities.TargetFuncSignature) []gen.MapExpression {
	in := tfs.In.FieldMap()

	list := make([]gen.MapExpression, 0)

	for i := 0; i < tfs.Out.Struct.NumFields(); i++ {
		field := tfs.Out.Struct.Field(i)
		outFieldName := field.Name()
		outFieldType := field.Type()

		inFieldType, ok := in[outFieldName]
		if !ok {
			fmt.Println("no matching field for ", outFieldName)
			continue
		}

		if !token.IsExported(outFieldName) && !(tfs.In.Local && tfs.Out.Local) {
			fmt.Println("output field is unexported ", outFieldName)
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

func createMapping(fieldName string, in, out types.Type) (gen.MapExpression, bool) {
	base := gen.BaseMapStatement{
		In:       in,
		Out:      out,
		InField:  fieldName,
		OutField: fieldName,
	}

	if typesAreCastable(in, out) {
		return gen.NewCast(base), true
	}

	outPtr, ok := out.(*types.Pointer)
	if ok && typesAreCastable(in, outPtr.Elem()) {
		return gen.NewCastValueToPtr(base), true
	}

	outPtr, ok = out.Underlying().(*types.Pointer)
	if ok && typesAreCastable(in, outPtr.Elem()) {
		return gen.NewCastValueToPtrType(base), true
	}

	inPtr, ok := in.(*types.Pointer)
	if ok && typesAreCastable(inPtr.Elem(), out) {
		return gen.NewCastPtrToValue(base), true
	}

	inPtr, ok = in.Underlying().(*types.Pointer)
	if ok && typesAreCastable(inPtr.Elem(), out) {
		return gen.NewCastPtrToValue(base), true
	}

	if inPtr != nil && outPtr != nil && typesAreCastable(inPtr.Elem(), outPtr.Elem()) {
		return gen.NewCastPtrToPtr(base), true
	}

	return nil, false
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
	if basic, ok := t.(*types.Basic); ok && basic.Kind()&types.String > 0 {
		return true
	}

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
