package main

import (
	"bytes"
	"text/template"
)

var errorsTemplate = `const (
{{ range .Errors }}
	{{.Name}}_{{.Value}} = "{{.Name}}_{{.Value}}"
{{- end }}
)

{{ range .Errors }}

func Is{{.Value}}(err error) bool {
	return errors.Reason(err).Reason == {{.Name}}_{{.Value}}
}
{{- end }}
`

// Error is a enum error.
type Error struct {
	Name  string
	Value string
}
type errorWrapper struct {
	Errors []*Error
}

func rendering(e errorWrapper) string {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("errors").Parse(errorsTemplate)
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, e); err != nil {
		panic(err)
	}
	return string(buf.Bytes())
}
