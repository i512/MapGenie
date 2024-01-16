package typematch

import (
	"fmt"
	"go/token"
	"go/types"
	"mapgenie/mapgen/edit"
	"mapgenie/mapgen/entities"
	"reflect"
)

func MappableFields(tfs entities.TargetFuncSignature, imports *edit.FileImports) []edit.TemplateMapping {
	in := tfs.In.FieldMap()

	list := make([]edit.TemplateMapping, 0)

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

		mapping, ok := createMapping(outFieldName, inFieldType, outFieldType, imports)
		if !ok {
			continue
		}

		list = append(list, mapping)
	}

	return list
}

func createMapping(fieldName string, in, out types.Type, imports *edit.FileImports) (edit.TemplateMapping, bool) {
	mapping := edit.TemplateMapping{
		InName:  fieldName,
		OutName: fieldName,
	}

	if typesAreCastable(in, out) {
		mapping.CastWith = castWith(in, out, imports)
		return mapping, true
	}

	outPtr, ok := out.(*types.Pointer)
	if ok && typesAreCastable(in, outPtr.Elem()) {
		mapping.OutPtr = true
		mapping.CastWith = castWith(in, outPtr.Elem(), imports)
		return mapping, true
	}

	inPtr, ok := in.(*types.Pointer)
	if ok && typesAreCastable(inPtr.Elem(), out) {
		mapping.InPtr = true
		mapping.CastWith = castWith(inPtr.Elem(), out, imports)
		return mapping, true
	}

	if inPtr != nil && outPtr != nil && typesAreCastable(inPtr.Elem(), outPtr.Elem()) {
		mapping.InPtr = true
		mapping.OutPtr = true
		mapping.CastWith = castWith(inPtr.Elem(), outPtr.Elem(), imports)
		return mapping, true
	}

	if hasUnderlying(in) {
		mapping, ok := createMapping(fieldName, in.Underlying(), out, imports)
		if ok {
			return mapping, true
		}
	}

	return edit.TemplateMapping{}, false
}

func typesAreCastable(in, out types.Type) bool {
	in = getUnderlying(in)
	out = getUnderlying(out)

	same := reflect.DeepEqual(in, out)
	numbers := typeIsIntegerOrFloat(in) && typeIsIntegerOrFloat(out)
	bytes := typeIsStringOrByteSlice(in) && typeIsStringOrByteSlice(out)
	return same || numbers || bytes
}

func castWith(in, out types.Type, imports *edit.FileImports) string {
	if reflect.DeepEqual(in, out) {
		return "" // same type, no cast needed
	}

	return imports.ResolveTypeName(out)
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

func typeIsUnderlying(base, derived types.Type) bool {
	return reflect.DeepEqual(base, derived.Underlying())
}

func hasUnderlying(t types.Type) bool {
	return t != t.Underlying()
}

func getUnderlying(t types.Type) types.Type {
	for t != t.Underlying() {
		t = t.Underlying()
	}
	return t
}
