package gen

import (
	"bytes"
	"go/types"
	"mapgenie/mapgen/entities"
	"reflect"
	"text/template"
)

type BaseMapStatement struct {
	In, Out           types.Type
	InField, OutField string
}

func (e BaseMapStatement) CastExpression(in, out types.Type, imports entities.TypeNameResolver) string {
	if reflect.DeepEqual(in, out) {
		return "" // same type, no cast needed
	}

	return imports.ResolveTypeName(out)
}

func (e BaseMapStatement) RunTemplate(exp any, temp string) (string, error) {
	t := template.Must(template.New("map_expression").Parse(temp))
	buf := bytes.NewBuffer(nil)

	err := t.Execute(buf, exp)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
