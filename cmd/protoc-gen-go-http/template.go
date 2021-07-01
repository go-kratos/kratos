package main

import (
	"bytes"
	"strings"
	"text/template"
)

var httpTemplate = `
{{$svrType := .ServiceType}}
{{$svrName := .ServiceName}}
type {{.ServiceType}}HTTPServer interface {
{{- range .MethodSets}}
	{{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
{{- end}}
}

func Register{{.ServiceType}}HTTPServer(s *http.Server, srv {{.ServiceType}}HTTPServer) {
	r := s.Route("/")
	{{- range .Methods}}
	r.{{.Method}}("{{.Path}}", _{{$svrType}}_{{.Name}}{{.Num}}_HTTP_Handler(srv))
	{{- end}}
}

{{range .Methods}}
func _{{$svrType}}_{{.Name}}{{.Num}}_HTTP_Handler(srv {{$svrType}}HTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in {{.Request}}
		{{- if .HasBody}}
		if err := ctx.Bind(&in{{.Body}}); err != nil {
			return err
		}
		{{- else}}
		if err := ctx.BindQuery(&in{{.Body}}); err != nil {
			return err
		}
		{{- end}}
		{{- if .HasVars}}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		{{- end}}
		http.SetOperation(ctx,"/{{$svrName}}/{{.Name}}")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.{{.Name}}(ctx, req.(*{{.Request}}))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*{{.Reply}})
		return ctx.Result(200, reply{{.ResponseBody}})
	}
}
{{end}}

type {{.ServiceType}}HTTPClient interface {
{{- range .MethodSets}}
	{{.Name}}(ctx context.Context, req *{{.Request}}, opts ...http.CallOption) (rsp *{{.Reply}}, err error) 
{{- end}}
}
	
type {{.ServiceType}}HTTPClientImpl struct{
	cc *http.Client
}
	
func New{{.ServiceType}}HTTPClient (client *http.Client) {{.ServiceType}}HTTPClient {
	return &{{.ServiceType}}HTTPClientImpl{client}
}

{{range .MethodSets}}
func (c *{{$svrType}}HTTPClientImpl) {{.Name}}(ctx context.Context, in *{{.Request}}, opts ...http.CallOption) (*{{.Reply}}, error) {
	var out {{.Reply}}
	pattern := "{{.Path}}"
	path := binding.EncodeURL(pattern, in, {{not .HasBody}})
	opts = append(opts, http.Operation("/{{$svrName}}/{{.Name}}"))
	opts = append(opts, http.PathTemplate(pattern))
	{{if .HasBody -}}
	err := c.cc.Invoke(ctx, "{{.Method}}", path, in{{.Body}}, &out{{.ResponseBody}}, opts...)
	{{else -}} 
	err := c.cc.Invoke(ctx, "{{.Method}}", path, nil, &out{{.ResponseBody}}, opts...)
	{{end -}}
	if err != nil {
		return nil, err
	}
	return &out, err
}
{{end}}
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
	Request string
	Reply   string
	// http_rule
	Path         string
	Method       string
	HasVars      bool
	HasBody      bool
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
	return strings.Trim(string(buf.Bytes()), "\r\n")
}
