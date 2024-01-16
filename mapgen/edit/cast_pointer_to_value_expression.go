package edit

import (
	"bytes"
	"mapgenie/mapgen/entities"
	"text/template"
)

type CastPointerToValueExpression struct {
	MapExpression
	Tm TemplateMapping
}

func (c *CastPointerToValueExpression) String(resolver entities.TypeNameResolver) string {
	sourceTemplate :=
		`if input.{{ .InName }} != nil {
			result.{{ .OutName }} ={{" "}}
				 {{- if ne .CastWith "" }}{{ .CastWith }}({{- end }}
				 *input.{{ .InName }}
				 {{- if ne .CastWith "" }}){{- end }}
		}`

	t := template.Must(template.New("map").Parse(sourceTemplate))
	buf := bytes.NewBuffer(nil)
	err := t.Execute(buf, c.Tm)
	if err != nil {
		panic(err)
	}

	return buf.String()
}
