package main

import (
	"bytes"
	"html/template"
	"strings"
)

var httpTemplate = `
type {{.ServiceType}}HTTPServer interface {
{{range .MethodSets}}
	{{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
{{end}}
}
func Register{{.ServiceType}}HTTPServer(s http1.ServiceRegistrar, srv {{.ServiceType}}HTTPServer) {
	s.RegisterService(&_HTTP_{{.ServiceType}}_serviceDesc, srv)
}
{{range .Methods}}
func _HTTP_{{$.ServiceType}}_{{.Name}}_{{.Num}}(srv interface{}, ctx context.Context, req *http.Request, dec func(interface{}) error) (interface{}, error) {
	var in {{.Request}}
{{if eq .Body ""}}
	if err := http1.BindForm(req, &in); err != nil {
		return nil, err
	}
{{else if eq .Body ".*"}}
	if err := dec(&in); err != nil {
		return nil, err
	}
{{else}}
	if err := dec(in{{.Body}}); err != nil {
		return nil, err
	}
{{end}}
{{if ne (len .Vars) 0}}
	if err := http1.BindVars(req, &in); err != nil {
		return nil, err
	}
{{end}}
	out, err := srv.({{$.ServiceType}}Server).{{.Name}}(ctx, &in)
	if err != nil {
		return nil, err
	}
	return out{{.ResponseBody}}, nil
}
{{end}}
var _HTTP_{{.ServiceType}}_serviceDesc = http1.ServiceDesc{
	ServiceName: "{{.ServiceName}}",
	Methods: []http1.MethodDesc{
{{range .Methods}}
		{
			Path:    "{{.Path}}",
			Method:  "{{.Method}}",
			Handler: _HTTP_{{$.ServiceType}}_{{.Name}}_{{.Num}},
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
	MethodSets  map[string]*methodDesc
}

type methodDesc struct {
	// method
	Name    string
	Num     int
	Vars    []string
	Forms   []string
	Request string
	Reply   string
	// http_rule
	Path         string
	Method       string
	Body         string
	ResponseBody string
}

func (s *serviceDesc) execute() string {
	s.MethodSets = make(map[string]*methodDesc)
	for _, m := range s.Methods {
		s.MethodSets[m.Name] = m
	}
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
