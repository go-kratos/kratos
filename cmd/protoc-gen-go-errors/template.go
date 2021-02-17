package main

import (
	"bytes"
	"text/template"
)

var errorsTemplate = `const (
{{ range .Errors }}
	Errors_{{.Value}} = "{{.Name}}_{{.Value}}"
{{- end }}
)

{{ range .Errors }}

func Is{{.Value}}(err error) bool {
	return errors.Reason(err) == Errors_{{.Value}}
}
{{- end }}
`

type errorInfo struct {
	Name  string
	Value string
}
type errorWrapper struct {
	Errors []*errorInfo
}

func (e *errorWrapper) execute() string {
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
