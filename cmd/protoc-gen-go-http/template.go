package main

import (
	"bytes"
	"strings"
	"text/template"
)

var httpTemplate = `
type {{.ServiceType}}Handler interface {
{{range .MethodSets}}
	{{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
{{end}}
}

func New{{.ServiceType}}Handler(srv {{.ServiceType}}Handler, opts ...http1.HandleOption) http.Handler {
	r := mux.NewRouter()
	{{range .Methods}}
	r.Handle("{{.Path}}", http1.NewHandler(srv.{{.Name}}, opts...)).Methods("{{.Method}}")
	{{end}}
	return r
}

type {{.ServiceType}}HttpClient interface {
{{range .MethodSets}}
	{{.Name}}(ctx context.Context, req *{{.Request}}, opts ...http1.CallOption) (rsp *{{.Reply}}, err error) 
{{end}}
}
	
type {{.ServiceType}}HttpClientImpl struct{
	cc *http1.Client
}
	
func New{{.ServiceType}}HttpClient (client *http1.Client) {{.ServiceType}}HttpClient {
	return &{{.ServiceType}}HttpClientImpl{client}
}
	
{{$svrType := .ServiceType}}
{{$svrName := .ServiceName}}
{{range .MethodSets}}
func (c *{{$svrType}}HttpClientImpl) {{.Name}}(ctx context.Context, in *{{.Request}}, opts ...http1.CallOption) (out *{{.Reply}}, err error) {
	path := "{{.Path}}"
	method := "{{.Method}}"
	body := "{{.Body}}"
	
	out = &{{.Reply}}{}
	err = c.cc.Invoke(ctx, path, in, out, http1.Body(body), http1.Method(method))
	if err != nil {
		return
	}
	return 
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
