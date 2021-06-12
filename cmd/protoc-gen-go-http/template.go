package main

import (
	"bytes"
	"strings"
	"text/template"
)

var httpTemplate = `
{{$svrType := .ServiceType}}
{{$svrName := .ServiceName}}
type {{.ServiceType}}Handler interface {
{{range .MethodSets}}
	{{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
{{end}}
}

func Register{{.ServiceType}}HTTPServer(s *http.Server, srv {{.ServiceType}}Handler) {
	r := s.Route("/")
	{{range .Methods}}
	r.{{.Method}}("{{.Path}}", func(ctx http.Context) error {
		var in {{.Request}}
		if err := ctx.Bind(&in{{.Body}}); err != nil {
			return err
		}
		{{if ne (len .Vars) 0}}
		if err := binding.BindVars(ctx.Vars(), &in); err != nil {
			return err
		}
		{{end}}
		transport.SetOperation(ctx,"/{{$svrName}}/{{.Name}}")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.{{.Name}}(ctx, req.(*{{.Request}}))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*{{.Reply}})
		return ctx.Result(200, reply{{.ResponseBody}})
	})
	{{end}}
}

type {{.ServiceType}}HTTPClient interface {
{{range .MethodSets}}
	{{.Name}}(ctx context.Context, req *{{.Request}}, opts ...http.CallOption) (rsp *{{.Reply}}, err error) 
{{end}}
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
	path := binding.EncodeVars("{{.Path}}", in, {{.IsQuery}})
	opts = append(opts, http.Operation("/{{$svrName}}/{{.Name}}"))
	{{if .HasBody }}
	err := c.cc.Invoke(ctx, "{{.Method}}", path, in{{.Body}}, &out{{.ResponseBody}}, opts...)
	{{else}} 
	err := c.cc.Invoke(ctx, "{{.Method}}", path, nil, &out{{.ResponseBody}}, opts...)
	{{end}}
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
	Vars    []string
	Request string
	Reply   string
	// http_rule
	Path         string
	Method       string
	HasBody      bool
	Body         string
	ResponseBody string
	IsQuery      bool
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
