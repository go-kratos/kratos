package main

import (
	"bytes"
	"text/template"
)

var errorsTemplate = `
{{ range .Errors }}

func Is{{.CamelValue}}(err error) bool {
	e := errors.FromError(err)
	return e.Reason == {{.Name}}_{{.Value}}.String() && e.Code == {{.HTTPCode}} 
}

func Error{{.CamelValue}}(format string, args ...interface{}) error {
	 return errors.New({{.HttpCode}}, {{.Name}}_{{.Value}}.String(), fmt.Sprintf(format, args...))

func {{.CamelValue}}(format string, args ...interface{}) error {
	var message string = "{{.Message}}"
	if format != "" {
		message = fmt.Sprintf(format, args...)
	}
	return errors.New({{.HttpCode}}, {{.Name}}_{{.Value}}.String(), message)
}

{{- end }}
`

type errorInfo struct {
	Name       string
	Message    string
	Value      string
	HttpCode   int
	CamelValue string
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
	return buf.String()
}
