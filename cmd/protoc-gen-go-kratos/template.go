package main

import (
	"bytes"
	"strings"
	"text/template"
)

var clientTemplate = `{{range .ClientInfoList }}
//New{{ .ServiceName }}GRPCClient create grpc client for kratos
func New{{ .ServiceName }}GRPCClient(ctx context.Context,opts ...grpc.ClientOption) (cli {{ .ServiceName }}Client, err error) {
	targetOpts := make([]grpc.ClientOption, 0, len(opts)+2)
	targetOpts = append(targetOpts,
		grpc.WithBalancerName(wrr.Name),
		{{- if .Endpoint }}
		grpc.WithEndpoint("{{ .Endpoint }}"),
		{{- end }}
	)
	conn, err := grpc.DialInsecure(ctx, append(targetOpts, opts...)...)
	if err != nil {
		return nil, err
	}
	return New{{ .ServiceName }}Client(conn), nil
}

//New{{ .ServiceName }}KratosHTTPClientV2 create http client for kratos
func New{{ .ServiceName }}KratosHTTPClientV2(ctx context.Context, opts ...http.ClientOption) (cli {{ .ServiceName }}HTTPClient, err error) {
	{{- if .Endpoint }}
	targetOpts := make([]http.ClientOption, 0, len(opts)+2)
	targetOpts = append(targetOpts, http.WithEndpoint("{{ .Endpoint }}"))
	client, err := http.NewClient(ctx, append(targetOpts, opts...)...)
	{{- else }}
	client, err := http.NewClient(ctx, opts...)
	{{- end }}
	if err != nil {
		return nil, err
	}
	return New{{ .ServiceName }}HTTPClient(client), nil
}
{{- end }}
`

type ClientInfo struct {
	ServiceName string // proto service
	Endpoint    string // default_host
}

type ClientTemplate struct {
	ClientInfoList []ClientInfo
}

// NewClientTemplate new client template
func NewClientTemplate() *ClientTemplate {
	return &ClientTemplate{
		ClientInfoList: make([]ClientInfo, 0, 5),
	}
}

func (receiver *ClientTemplate) AppendClientInfo(serviceName, endpoint string) {
	receiver.ClientInfoList = append(receiver.ClientInfoList, ClientInfo{
		ServiceName: serviceName,
		Endpoint:    endpoint,
	})
}

// Parse create content
func (receiver *ClientTemplate) execute() string {
	parser, err := template.New("clientTemplate").Parse(clientTemplate)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	if err := parser.Execute(buf, receiver); err != nil {
		panic(err)
	}
	return strings.Trim(buf.String(), "\r\n")
}
