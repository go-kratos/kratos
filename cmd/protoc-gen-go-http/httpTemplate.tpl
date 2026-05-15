{{$svrType := .ServiceType}}
{{$svrName := .ServiceName}}

{{- range .MethodSets}}
const Operation{{$svrType}}{{.OriginalName}} = "/{{$svrName}}/{{.OriginalName}}"
{{- end}}

type {{.ServiceType}}HTTPServer interface {
{{- range .MethodSets}}
	{{- if ne .Comment ""}}
	{{.Comment}}
	{{- end}}
	{{- if .ClientStreaming}}
	{{.Name}}({{$svrType}}_{{.Name}}Server) error
	{{- else if .ServerStreaming}}
	{{.Name}}(*{{.Request}}, {{$svrType}}_{{.Name}}Server) error
	{{- else}}
	{{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
	{{- end}}
{{- end}}
}

func Register{{.ServiceType}}HTTPServer(s *http.Server, srv {{.ServiceType}}HTTPServer) {
	r := s.Route("/")
	{{- range .Methods}}
	{{- if .ClientStreaming}}
	r.Handle("GET", "{{.Path}}", _{{$svrType}}_{{.Name}}{{.Num}}_HTTP_Handler(srv))
	{{- else}}
	r.Handle("{{.Method}}", "{{.Path}}", _{{$svrType}}_{{.Name}}{{.Num}}_HTTP_Handler(srv))
	{{- end}}
	{{- end}}
}

{{range .MethodSets}}
{{- if or .ClientStreaming .ServerStreaming}}
type {{$svrType}}_{{.Name}}HTTPServer struct {
	http.ServerStream
}

{{- if .ServerStreaming}}
func (x *{{$svrType}}_{{.Name}}HTTPServer) Send(m *{{.Reply}}) error {
	return x.ServerStream.Send(m)
}
{{- end}}

{{- if .ClientStreaming}}
func (x *{{$svrType}}_{{.Name}}HTTPServer) Recv() (*{{.Request}}, error) {
	m := new({{.Request}})
	if err := x.ServerStream.Recv(m); err != nil {
		return nil, err
	}
	return m, nil
}
{{- end}}

{{- if and .ClientStreaming (not .ServerStreaming)}}
func (x *{{$svrType}}_{{.Name}}HTTPServer) SendAndClose(m *{{.Reply}}) error {
	return x.ServerStream.SendAndClose(m)
}
{{- end}}
{{- end}}
{{end}}

{{range .Methods}}
func _{{$svrType}}_{{.Name}}{{.Num}}_HTTP_Handler(srv {{$svrType}}HTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		{{- if .ClientStreaming}}
		stream, err := http.NewWebSocketServerStream(ctx)
		if err != nil {
			return err
		}
		http.SetOperation(ctx,Operation{{$svrType}}{{.OriginalName}})
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			stream.SetContext(ctx)
			return nil, srv.{{.Name}}(&{{$svrType}}_{{.Name}}HTTPServer{ServerStream: stream})
		})
		_, err = h(ctx, nil)
		return stream.Close(err)
		{{- else if .ServerStreaming}}
		var in {{.Request}}
		{{- if .HasBody}}
		if err := ctx.Bind(&in{{.Body}}); err != nil {
			return err
		}
		{{- end}}
		{{- if not .HasBody}}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		{{- else if ne .BodyField "*"}}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		{{- end}}
		{{- if .HasVars}}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		{{- end}}
		stream := http.NewServerSentEventServerStream(ctx)
		http.SetOperation(ctx,Operation{{$svrType}}{{.OriginalName}})
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			stream.SetContext(ctx)
			return nil, srv.{{.Name}}(req.(*{{.Request}}), &{{$svrType}}_{{.Name}}HTTPServer{ServerStream: stream})
		})
		_, err := h(ctx, &in)
		return stream.Close(err)
		{{- else}}
		var in {{.Request}}
		{{- if .HasBody}}
		if err := ctx.Bind(&in{{.Body}}); err != nil {
			return err
		}
		{{- end}}
		{{- if not .HasBody}}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		{{- else if ne .BodyField "*"}}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		{{- end}}
		{{- if .HasVars}}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		{{- end}}
		http.SetOperation(ctx,Operation{{$svrType}}{{.OriginalName}})
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.{{.Name}}(ctx, req.(*{{.Request}}))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*{{.Reply}})
		return ctx.Result(200, reply{{.ResponseBody}})
		{{- end}}
	}
}
{{end}}

type {{.ServiceType}}HTTPClient interface {
{{- range .MethodSets}}
	{{- if ne .Comment ""}}
	{{.Comment}}
	{{- end}}
	{{- if .ClientStreaming}}
	{{.Name}}(ctx context.Context, opts ...http.CallOption) ({{$svrType}}_{{.Name}}Client, error)
	{{- else if .ServerStreaming}}
	{{.Name}}(ctx context.Context, req *{{.Request}}, opts ...http.CallOption) ({{$svrType}}_{{.Name}}Client, error)
	{{- else}}
	{{.Name}}(ctx context.Context, req *{{.Request}}, opts ...http.CallOption) (rsp *{{.Reply}}, err error)
	{{- end}}
{{- end}}
}

type {{.ServiceType}}HTTPClientImpl struct{
	cc *http.Client
}

func New{{.ServiceType}}HTTPClient (client *http.Client) {{.ServiceType}}HTTPClient {
	return &{{.ServiceType}}HTTPClientImpl{client}
}

{{range .MethodSets}}
{{- if or .ClientStreaming .ServerStreaming}}
type {{$svrType}}_{{.Name}}HTTPClient struct {
	http.ClientStream
	{{- if .ClientStreaming}}
	ctx context.Context
	cc *http.Client
	pattern string
	opts []http.CallOption
	{{- end}}
}

{{- if .ClientStreaming}}
func (x *{{$svrType}}_{{.Name}}HTTPClient) open(m *{{.Request}}) error {
	if x.ClientStream != nil {
		return nil
	}
	{{- if .BodyHTTPBody}}
	opts := append([]http.CallOption{
		http.ContentType(http.BodyContentType(m)),
	}, x.opts...)
	{{- else}}
	opts := x.opts
	{{- end}}
	{{- if .HasBody}}
		{{- if or (eq .BodyField "*") (eq .BodyField "")}}
	path := http.BuildPath(x.pattern, m)
		{{- else}}
	path := http.BuildPath(x.pattern, m, http.WithQueryParams(), http.WithOmitFields("{{.BodyQueryName}}"))
		{{- end}}
	{{- else}}
	path := http.BuildPath(x.pattern, m, http.WithQueryParams())
	{{- end}}
	stream, err := x.cc.WebSocket(x.ctx, path, opts...)
	if err != nil {
		return err
	}
	x.ClientStream = stream
	return nil
}

func (x *{{$svrType}}_{{.Name}}HTTPClient) CloseSend() error {
	if err := x.open(nil); err != nil {
		return err
	}
	return x.ClientStream.CloseSend()
}

func (x *{{$svrType}}_{{.Name}}HTTPClient) Send(m *{{.Request}}) error {
	if err := x.open(m); err != nil {
		return err
	}
	return x.ClientStream.Send(m)
}
{{- end}}

{{- if .ServerStreaming}}
func (x *{{$svrType}}_{{.Name}}HTTPClient) Recv() (*{{.Reply}}, error) {
	{{- if .ClientStreaming}}
	if err := x.open(nil); err != nil {
		return nil, err
	}
	{{- end}}
	m := new({{.Reply}})
	if err := x.ClientStream.Recv(m); err != nil {
		return nil, err
	}
	return m, nil
}
{{- end}}

{{- if and .ClientStreaming (not .ServerStreaming)}}
func (x *{{$svrType}}_{{.Name}}HTTPClient) CloseAndRecv() (*{{.Reply}}, error) {
	if err := x.open(nil); err != nil {
		return nil, err
	}
	m := new({{.Reply}})
	if err := x.ClientStream.CloseAndRecv(m); err != nil {
		return nil, err
	}
	return m, nil
}
{{- end}}
{{- end}}
{{end}}

{{range .MethodSets}}
	{{- if ne .Comment ""}}
	{{.Comment}}
	{{- end}}
{{- if .ClientStreaming}}
func (c *{{$svrType}}HTTPClientImpl) {{.Name}}(ctx context.Context, opts ...http.CallOption) ({{$svrType}}_{{.Name}}Client, error) {
	pattern := "{{.PathTemplate}}"
	opts = append([]http.CallOption{
		http.Accept("application/protojson"),
		{{- if not .BodyHTTPBody}}
		http.ContentType("application/protojson"),
		{{- end}}
		http.Operation(Operation{{$svrType}}{{.OriginalName}}),
		http.PathTemplate(pattern),
	}, opts...)
	return &{{$svrType}}_{{.Name}}HTTPClient{ctx: ctx, cc: c.cc, pattern: pattern, opts: opts}, nil
}
{{- else if .ServerStreaming}}
func (c *{{$svrType}}HTTPClientImpl) {{.Name}}(ctx context.Context, in *{{.Request}}, opts ...http.CallOption) ({{$svrType}}_{{.Name}}Client, error) {
	pattern := "{{.PathTemplate}}"
	{{- if .HasBody}}
		{{- if or (eq .BodyField "*") (eq .BodyField "")}}
	path := http.BuildPath(pattern, in)
		{{- else}}
	path := http.BuildPath(pattern, in, http.WithQueryParams(), http.WithOmitFields("{{.BodyQueryName}}"))
		{{- end}}
	opts = append([]http.CallOption{
		http.Accept("text/event-stream"),
		{{- if .BodyHTTPBody}}
		http.ContentType(http.BodyContentType(in{{.Body}})),
		{{- else}}
		http.ContentType("application/protojson"),
		{{- end}}
		http.Operation(Operation{{$svrType}}{{.OriginalName}}),
		http.PathTemplate(pattern),
	}, opts...)
	stream, err := c.cc.ServerSentEvent(ctx, "{{.Method}}", path, in{{.Body}}, opts...)
	{{- else}}
	path := http.BuildPath(pattern, in, http.WithQueryParams())
	opts = append([]http.CallOption{
		http.Accept("text/event-stream"),
		http.ContentType("application/protojson"),
		http.Operation(Operation{{$svrType}}{{.OriginalName}}),
		http.PathTemplate(pattern),
	}, opts...)
	stream, err := c.cc.ServerSentEvent(ctx, "{{.Method}}", path, nil, opts...)
	{{- end}}
	if err != nil {
		return nil, err
	}
	return &{{$svrType}}_{{.Name}}HTTPClient{ClientStream: stream}, nil
}
{{- else}}
func (c *{{$svrType}}HTTPClientImpl) {{.Name}}(ctx context.Context, in *{{.Request}}, opts ...http.CallOption) (*{{.Reply}}, error) {
	var out {{.Reply}}
	pattern := "{{.PathTemplate}}"
	{{- if .HasBody}}
		{{- if or (eq .BodyField "*") (eq .BodyField "")}}
	path := http.BuildPath(pattern, in)
		{{- else}}
	path := http.BuildPath(pattern, in, http.WithQueryParams(), http.WithOmitFields("{{.BodyQueryName}}"))
		{{- end}}
	opts = append([]http.CallOption{
		http.Accept("application/protojson"),
		{{- if .BodyHTTPBody}}
		http.ContentType(http.BodyContentType(in{{.Body}})),
		{{- else}}
		http.ContentType("application/protojson"),
		{{- end}}
		http.Operation(Operation{{$svrType}}{{.OriginalName}}),
		http.PathTemplate(pattern),
	}, opts...)
	{{- else}}
	path := http.BuildPath(pattern, in, http.WithQueryParams())
	opts = append([]http.CallOption{
		http.Accept("application/protojson"),
		http.Operation(Operation{{$svrType}}{{.OriginalName}}),
		http.PathTemplate(pattern),
	}, opts...)
	{{- end}}
	{{if .HasBody -}}
	err := c.cc.Invoke(ctx, "{{.Method}}", path, in{{.Body}}, &out{{.ResponseBody}}, opts...)
	{{else -}}
	err := c.cc.Invoke(ctx, "{{.Method}}", path, nil, &out{{.ResponseBody}}, opts...)
	{{end -}}
	if err != nil {
		return nil, err
	}
	return &out, nil
}
{{- end}}
{{end}}
