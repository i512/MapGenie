package edit

import (
	"bytes"
	"mapgenie/mapgen/entities"
	"text/template"
)

type CastValueToPointerExpression struct {
	MapExpression
	Tm TemplateMapping
}

func (c *CastValueToPointerExpression) String(resolver entities.TypeNameResolver) string {
	sourceTemplate :=
		`{{- if and .OutPtr (ne .CastWith "")}}
			{{ .OutName }} := {{ .CastWith }}(input.{{ .InName }})
			result.{{ .OutName }} = &{{ .OutName }}
		{{- else }}
			result.{{ .OutName }} ={{" "}}
				{{- if ne .CastWith "" }}{{ .CastWith }}({{- end }}
				{{- if .OutPtr }}&{{- end }}
				{{- ""}}input.{{ .InName }}
				{{- if ne .CastWith "" }}){{- end }}
		{{- end }}`

	t := template.Must(template.New("map").Parse(sourceTemplate))
	buf := bytes.NewBuffer(nil)
	err := t.Execute(buf, c.Tm)
	if err != nil {
		panic(err)
	}

	return buf.String()
}
