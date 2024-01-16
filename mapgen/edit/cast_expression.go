package edit

import (
	"bytes"
	"mapgenie/mapgen/entities"
	"text/template"
)

type CastExpression struct {
	MapExpression
	Tm TemplateMapping
}

func (c *CastExpression) String(resolver entities.TypeNameResolver) string {
	sourceTemplate :=
		`result.{{ .OutName }} = 
			 {{- if ne .CastWith "" }}{{ .CastWith }}({{- end }}
			 input.{{ .InName }}
			 {{- if ne .CastWith "" }}){{- end }}`

	t := template.Must(template.New("map").Parse(sourceTemplate))
	buf := bytes.NewBuffer(nil)
	err := t.Execute(buf, c.Tm)
	if err != nil {
		panic(err)
	}

	return buf.String()
}
