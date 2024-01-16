package edit

import (
	"bytes"
	"mapgenie/mapgen/entities"
	"text/template"
)

type CastPointerToPointerExpression struct {
	MapExpression
	Tm TemplateMapping
}

func (c *CastPointerToPointerExpression) String(resolver entities.TypeNameResolver) string {
	sourceTemplate :=
		`if input.{{ .InName }} != nil {
			{{ .OutName }} := 
				 {{- if ne .CastWith "" }}{{ .CastWith }}({{- end }}
				 *input.{{ .InName }}
				 {{- if ne .CastWith "" }}){{- end }}

			result.{{ .OutName }} = &{{ .OutName }}
		}`

	t := template.Must(template.New("map").Parse(sourceTemplate))
	buf := bytes.NewBuffer(nil)
	err := t.Execute(buf, c.Tm)
	if err != nil {
		panic(err)
	}

	return buf.String()
}
