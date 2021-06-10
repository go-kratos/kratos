package http

import (
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/transport"
)

var (
	_ transport.Transporter = &Transport{}
)

type Transport struct {
	endpoint      string
	serviceMethod string
	metadata      metadata.Metadata
}

func (tr *Transport) Kind() string {
	return transport.KindHTTP
}

func (tr *Transport) Endpoint() string {
	return tr.endpoint
}

func (tr *Transport) ServiceMethod() string {
	return tr.serviceMethod
}

func (tr *Transport) SetServiceMethod(serviceMethod string) {
	tr.serviceMethod = serviceMethod
}

func (tr *Transport) Metadata() metadata.Metadata {
	return tr.metadata
}

func (tr *Transport) WithMetadata(md metadata.Metadata) {
	if tr.metadata == nil {
		tr.metadata = md
		return
	}
	for k, v := range md {
		tr.metadata.Set(k, v)
	}
}

func (tr *Transport) Clone() transport.Transporter {
	return &Transport{
		endpoint:      tr.endpoint,
		serviceMethod: tr.serviceMethod,
		metadata:      tr.metadata.Clone(),
	}
}
