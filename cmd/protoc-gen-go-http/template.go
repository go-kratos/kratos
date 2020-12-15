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
func _HTTP_{{.ServiceType}}_{{.Name}}(srv interface{}, ctx context.Context, m http.Marshaler) ([]byte, error) {
	in := new({{.Request}})
	if err := m.Unmarshal(in{{.Body}}); err != nil {
		return nil, err
	}

{{ if ne (len .Params) 0 }}
	var(
		err error
		vars = m.PathParams()
	)
{{ end }}
{{ range .Params }}

	{{.ProtoName}}, ok := vars["{{.ProtoName}}"]
	if !ok {
		return nil, http.ErrInvalidArgument("missing parameter: {{.ProtoName}}")
	}
	in.{{.GoName}} = {{.ProtoName}}
{{- end }}

	reply, err := srv.({{.ServiceType}}Server).{{.Name}}(ctx, in)
	if err != nil {
		return nil, err
	}
	return m.Marshal(reply{{.ResponseBody}})
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
	Params  []pathParam
	Request string
	Reply   string
	// http_rule
	Path         string
	Method       string
	Body         string
	ResponseBody string
}

type pathParam struct {
	GoName    string
	ProtoName string
	Type      string
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
