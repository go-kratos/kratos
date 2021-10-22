package http

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-kratos/kratos/v2/transport"
	"github.com/stretchr/testify/assert"
)

func TestTransport_Kind(t *testing.T) {
	o := &Transport{}
	assert.Equal(t, transport.KindHTTP, o.Kind())
}

func TestTransport_Endpoint(t *testing.T) {
	v := "hello"
	o := &Transport{endpoint: v}
	assert.Equal(t, v, o.Endpoint())
}

func TestTransport_Operation(t *testing.T) {
	v := "hello"
	o := &Transport{operation: v}
	assert.Equal(t, v, o.Operation())
}

func TestTransport_Request(t *testing.T) {
	v := &http.Request{}
	o := &Transport{request: v}
	assert.Same(t, v, o.Request())
}

func TestTransport_RequestHeader(t *testing.T) {
	v := headerCarrier{}
	v.Set("a", "1")
	o := &Transport{reqHeader: v}
	assert.Equal(t, "1", o.RequestHeader().Get("a"))
}

func TestTransport_ReplyHeader(t *testing.T) {
	v := headerCarrier{}
	v.Set("a", "1")
	o := &Transport{replyHeader: v}
	assert.Equal(t, "1", o.ReplyHeader().Get("a"))
}

func TestTransport_PathTemplate(t *testing.T) {
	v := "template"
	o := &Transport{pathTemplate: v}
	assert.Equal(t, v, o.PathTemplate())
}

func TestHeaderCarrier_Keys(t *testing.T) {
	v := headerCarrier{}
	v.Set("abb", "1")
	v.Set("bcc", "2")
	assert.ElementsMatch(t, []string{"Abb", "Bcc"}, v.Keys())
}

func TestSetOperation(t *testing.T) {
	tr := &Transport{}
	ctx := transport.NewServerContext(context.Background(), tr)
	SetOperation(ctx, "kratos")
	assert.Equal(t, tr.operation, "kratos")
}
