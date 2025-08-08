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

func Error{{ .CamelValue }}WithContext(ctx context.Context, args ...interface{}) *errors.Error {
     if len(args) == 0 {
          return errors.NewWithContext(ctx, {{ .HTTPCode }}, {{ .Name }}_{{ .Value }}.String(), "")
     }
	 return errors.NewWithContext(ctx, {{ .HTTPCode }}, {{ .Name }}_{{ .Value }}.String(), args[0])
}

{{- end }}
