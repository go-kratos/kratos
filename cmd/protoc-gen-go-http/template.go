package main

import (
	"bytes"
	"html/template"
	"strings"
)

var httpTemplate = `
type {{.ServiceType}}HTTPServer interface {
{{ range .Methods }}
	{{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
{{- end }}
}

func Register{{.ServiceType}}HTTPServer(s http.ServiceRegistrar, srv {{.ServiceType}}HTTPServer) {
	s.RegisterService(&_HTTP_{{.ServiceType}}_serviceDesc, srv)
}

{{ range .Methods }}
func _HTTP_{{.ServiceType}}_{{.Name}}(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new({{.Request}})
	if err := dec(in); err != nil {
		return nil, err
	}
	return srv.({{.ServiceType}}Server).{{.Name}}(ctx, in)
}
{{- end }}

var _HTTP_{{.ServiceType}}_serviceDesc = http.ServiceDesc{
	ServiceName: "{{.ServiceName}}",
	HandlerType: (*{{.ServiceType}}HTTPServer)(nil),
	Methods: []http.MethodDesc{
{{ range .Methods }}
		{
			Path:    "{{.Path}}",
			Method:  "{{.Method}}",
			Body:    "{{.Body}}",
			ResponseBody: "{{.ResponseBody}}",
			Handler: _HTTP_{{.ServiceType}}_{{.Name}},
		},
{{- end }}
	},
	Metadata: "{{.Metadata}}",
}
`

type serviceDesc struct {
	ServiceType string // Greeter
	ServiceName string // helloworld.Greeter
	Metadata    string // api/helloworld/helloworld.proto
	Methods     []*methodDesc
}

type methodDesc struct {
	// service
	ServiceType string // Greeter
	// method
	Name    string
	Request string
	Reply   string
	// http_rule
	Path         string
	Method       string
	Body         string
	ResponseBody string
}

func (s *serviceDesc) execute() string {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("http").Parse(strings.TrimSpace(httpTemplate))
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, s); err != nil {
		panic(err)
	}
	return string(buf.Bytes())
}
