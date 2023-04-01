{{ range .Errors }}

{{ if .HasComment }}{{ .Comment }}{{ end -}}
func Is{{.CamelValue}}(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == {{ .Name }}_{{ .Value }}.String() && e.Code == {{ .HTTPCode }}
}

{{ if .HasComment }}{{ .Comment }}{{ end -}}
func Error{{ .CamelValue }}(format string, args ...interface{}) *errors.Error {
	 return errors.New({{ .HTTPCode }}, {{ .Name }}_{{ .Value }}.String(), fmt.Sprintf(format, args...))
}

{{ if .HasDetails }}

{{ if .Details.HasFormat }}

{{ if .HasComment }}{{ .Comment }}{{ end -}}
func Error{{ .CamelValue }}WithFormat(args ...interface{}) *errors.Error {
	 return errors.New({{ .HTTPCode }}, {{ .Name }}_{{ .Value }}.String(), fmt.Sprintf("{{ .Details.Format }}", args...))
}

{{- end }}

{{- end }}

{{- end }}