package http

import (
	"context"
	"net/http"

	"github.com/go-kratos/kratos/v2/transport"
)

var (
	_ transport.Transporter = &Transport{}
)

// Transport is an HTTP transport.
type Transport struct {
	endpoint     string
	operation    string
	header       headerCarrier
	request      *http.Request
	pathTemplate string
}

// Kind returns the transport kind.
func (tr *Transport) Kind() transport.Kind {
	return transport.KindHTTP
}

// Endpoint returns the transport endpoint.
func (tr *Transport) Endpoint() string {
	return tr.endpoint
}

// Operation returns the transport operation.
func (tr *Transport) Operation() string {
	return tr.operation
}

// Header returns the transport header.
func (tr *Transport) Header() transport.Header {
	return tr.header
}

// Request returns the transport request.
func (tr *Transport) Request() *http.Request {
	return tr.request
}

// PathTemplate returns the http path template.
func (tr *Transport) PathTemplate() string {
	return tr.pathTemplate
}

// SetOperation sets the transport operation.
func SetOperation(ctx context.Context, op string) {
	if tr, ok := transport.FromServerContext(ctx); ok {
		if tr, ok := tr.(*Transport); ok {
			tr.operation = op
		}
	}
}

type headerCarrier http.Header

// Get returns the value associated with the passed key.
func (hc headerCarrier) Get(key string) string {
	return http.Header(hc).Get(key)
}

// Set stores the key-value pair.
func (hc headerCarrier) Set(key string, value string) {
	http.Header(hc).Set(key, value)
}

// Keys lists the keys stored in this carrier.
func (hc headerCarrier) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range http.Header(hc) {
		keys = append(keys, k)
	}
	return keys
}
