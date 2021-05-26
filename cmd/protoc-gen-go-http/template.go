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
	h := http1.DefaultHandleOptions()
	for _, o := range opts {
		o(&h)
	}
	r := mux.NewRouter()
	{{range .Methods}}
	r.HandleFunc("{{.Path}}", func(w http.ResponseWriter, r *http.Request) {
		var in {{.Request}}
		if err := h.Decode(r, &in{{.Body}}); err != nil {
			h.Error(w, r, err)
			return
		}
		{{if ne (len .Vars) 0}}
		if err := binding.BindVars(mux.Vars(r), &in); err != nil {
			h.Error(w, r, err)
			return
		}
		{{end}}
		next := func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.{{.Name}}(ctx, req.(*{{.Request}}))
		}
		if h.Middleware != nil {
			next = h.Middleware(next)
		}
		out, err := next(r.Context(), &in)
		if err != nil {
			h.Error(w, r, err)
			return
		}
		reply := out.(*{{.Reply}})
		if err := h.Encode(w, r, reply{{.ResponseBody}}); err != nil {
			h.Error(w, r, err)
		}
	}).Methods("{{.Method}}")
	{{end}}
	return r
}

type {{.ServiceType}}HTTPClient interface {
{{range .MethodSets}}
	{{.Name}}(ctx context.Context, req *{{.Request}}, opts ...http1.CallOption) (rsp *{{.Reply}}, err error) 
{{end}}
}
	
type {{.ServiceType}}HTTPClientImpl struct{
	cc *http1.Client
}
	
func New{{.ServiceType}}HTTPClient (client *http1.Client) {{.ServiceType}}HTTPClient {
	return &{{.ServiceType}}HTTPClientImpl{client}
}
	
{{$svrType := .ServiceType}}
{{$svrName := .ServiceName}}
{{range .MethodSets}}
func (c *{{$svrType}}HTTPClientImpl) {{.Name}}(ctx context.Context, in *{{.Request}}, opts ...http1.CallOption) (out *{{.Reply}}, err error) {
	path := binding.EncodePath("{{.Method}}", "{{.Path}}", in)
	out = &{{.Reply}}{}
	{{if .HasBody }}
	err = c.cc.Invoke(ctx, path, in{{.Body}}, &out{{.ResponseBody}}, http1.Method("{{.Method}}"), http1.PathPattern("{{.Path}}"))
	{{else}} 
	err = c.cc.Invoke(ctx, path, nil, &out{{.ResponseBody}}, http1.Method("{{.Method}}"), http1.PathPattern("{{.Path}}"))
	{{end}}
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
	Request string
	Reply   string
	// http_rule
	Path         string
	Method       string
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
	return string(buf.Bytes())
}
