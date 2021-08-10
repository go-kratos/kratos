package auth

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type headerCarrier http.Header

func (hc headerCarrier) Get(key string) string { return http.Header(hc).Get(key) }

func (hc headerCarrier) Set(key string, value string) { http.Header(hc).Set(key, value) }

// Keys lists the keys stored in this carrier.
func (hc headerCarrier) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range http.Header(hc) {
		keys = append(keys, k)
	}
	return keys
}

func newTokenHeader(headerKey, token string) *headerCarrier {
	header := &headerCarrier{}
	header.Set(headerKey, token)
	return header
}

type Transport struct {
	kind      transport.Kind
	endpoint  string
	operation string
	reqHeader transport.Header
}

func (tr *Transport) Kind() transport.Kind {
	return tr.kind
}
func (tr *Transport) Endpoint() string {
	return tr.endpoint
}
func (tr *Transport) Operation() string {
	return tr.operation
}
func (tr *Transport) RequestHeader() transport.Header {
	return tr.reqHeader
}
func (tr *Transport) ReplyHeader() transport.Header {
	return nil
}

type testParser struct {
}

func (t testParser) ParseToken(token string) (interface{}, error) {
	if token == "" {
		return nil, fmt.Errorf("can not find token")
	}
	return token, nil
}

type tokeProvider struct {
	token string
}

func (t tokeProvider) GetToken() string {
	return t.token
}

func TestDefaultAuthHeaderKey(t *testing.T) {
	token := "testToken"
	tests := []struct {
		name      string
		ctx       context.Context
		headerKey string
	}{
		{
			name:      "with headerKey",
			ctx:       transport.NewServerContext(context.Background(), &Transport{reqHeader: newTokenHeader("token", token)}),
			headerKey: "token",
		},
		{
			name:      "without headerKey",
			ctx:       transport.NewServerContext(context.Background(), &Transport{reqHeader: newTokenHeader("Authorization", token)}),
			headerKey: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			next := func(ctx context.Context, req interface{}) (interface{}, error) {
				t.Log(req)
				return "reply", nil
			}
			server := Server(testParser{}, WithAuthHeaderKey(test.headerKey))(next)
			_, err := server(test.ctx, "hhhh")
			assert.Nil(t, err)
		})
	}
}
