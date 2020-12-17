package http

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

// ServerInfo is HTTP server infomation.
type ServerInfo struct {
	Request  *http.Request
	Response http.ResponseWriter
}

type serverKey struct{}

// NewContext returns a new Context that carries value.
func NewContext(ctx context.Context, s ServerInfo) context.Context {
	return context.WithValue(ctx, serverKey{}, s)
}

// FromContext returns the Transport value stored in ctx, if any.
func FromContext(ctx context.Context) (s ServerInfo, ok bool) {
	s, ok = ctx.Value(serverKey{}).(ServerInfo)
	return
}

// PathParams returns path parameters.
func PathParams(req *http.Request) map[string]string {
	return mux.Vars(req)
}
