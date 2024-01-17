package gen

import (
	"bytes"
	"go/types"
	"mapgenie/mapgen/entities"
	"reflect"
	"text/template"
)

type BaseExpression struct {
	In, Out           types.Type
	InField, OutField string
}

func (e BaseExpression) CastExpression(in, out types.Type, imports entities.TypeNameResolver) string {
	if reflect.DeepEqual(in, out) {
		return "" // same type, no cast needed
	}

	return imports.ResolveTypeName(out)
}

func (e BaseExpression) RunTemplate(exp any, temp string) string {
	t := template.Must(template.New("map_expression").Parse(temp))
	buf := bytes.NewBuffer(nil)

	err := t.Execute(buf, exp)
	if err != nil {
		panic(err)
	}

	return buf.String()
}
