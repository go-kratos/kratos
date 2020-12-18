package main

import (
	"bytes"
	"html/template"
	"strings"
)

var httpTemplate = `
type {{.ServiceType}}HTTPServer interface {
{{range .Methods}}
	{{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
{{end}}
}

func Register{{.ServiceType}}HTTPServer(s http1.ServiceRegistrar, srv {{.ServiceType}}HTTPServer) {
	s.RegisterService(&_HTTP_{{.ServiceType}}_serviceDesc, srv)
}

{{range .Methods}}
func _HTTP_{{.ServiceType}}_{{.Name}}(srv interface{}, ctx context.Context, dec func(interface{}) error, req *http.Request) (interface{}, error) {
	var in {{.Request}}
{{if eq .Body ".*"}}
	if err := dec(&in); err != nil {
		return nil, err
	}
{{else if ne .Body ""}}
	if err := dec(&in{{.Body}}); err != nil {
		return nil, err
	}
{{end}}
{{if ne (len .Params) 0}}
	var (
		ok bool
		err error
		value string
		params = http1.PathParams(req)
	)
{{end }}
{{range .Params}}
	if value, ok = params["{{.ProtoName}}"]; !ok {
		return nil, errors.InvalidArgument("Errors_InvalidArgument", "Missing parameter: {{.ProtoName}}")
	}
	if in.{{.GoName}}, err = http1.{{.Kind}}(value); err != nil {
		return nil, errors.InvalidArgument("Errors_InvalidArgument", "Failed to parse {{.ProtoName}}: %s error = %v", value, err)
	}
{{end}}
	out, err := srv.({{.ServiceType}}Server).{{.Name}}(ctx, &in)
	if err != nil {
		return nil, err
	}
	return out{{.ResponseBody}}, nil
}
{{end}}

var _HTTP_{{.ServiceType}}_serviceDesc = http1.ServiceDesc{
	ServiceName: "{{.ServiceName}}",
	HandlerType: (*{{.ServiceType}}HTTPServer)(nil),
	Methods: []http1.MethodDesc{
{{range .Methods}}
		{
			Path:    "{{.Path}}",
			Method:  "{{.Method}}",
			Handler: _HTTP_{{.ServiceType}}_{{.Name}},
		},
{{end}}
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
	Kind      string
	GoName    string
	ProtoName string
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
