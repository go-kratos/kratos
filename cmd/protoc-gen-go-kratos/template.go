package main

import (
	"bytes"
	"strings"
	"text/template"
)

var clientTemplate = `


//New{{ .ServiceName }}HTTPClient 
func New{{ .ServiceName }}HTTPKratosClient(opts ...http.ClientOption) ({{ .ServiceName }}HTTPClient, error) {
	opts = append(opts, http.WithEndpoint({{ .Endpoint }}))
	client, err := http.NewClient(context.Background(),
		opts...,
	)
	if err != nil {
		return nil, err
	}
	return New{{ .ServiceName }}HTTPClient(client), nil
}

//New{{ .ServiceName }}GRPCClient 
func New{{ .ServiceName }}GRPCKratosClient(opts ...grpc.ClientOption) ({{ .ServiceName }}Client, error) {
	opts = append(opts, grpc.WithEndpoint({{ .Endpoint }}))
	conn, err := grpc.DialInsecure(context.Background(), opts...)
	if err != nil {
		return nil, err
	}
	return New{{ .ServiceName }}Client(conn), nil
}



`

type ClientInfo struct {
	ServiceName string // proto service
	Endpoint    string // default_host
}

type ClientTemplate struct {
	ClientInfoList []ClientInfo
}

//NewClientTemplate new client template
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

//NewClientTemplate new client template
//func NewClientTemplate(serviceName, endpoint string) *ClientTemplate {
//	return &ClientTemplate{
//		ServiceName: serviceName,
//		Endpoint:    endpoint,
//	}
//}

//Parse create content
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
