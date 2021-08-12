package grpc

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/go-kratos/kratos/v2/transport"
)

func TestTransport_Kind(t *testing.T) {
	o := &Transport{}
	assert.Equal(t, transport.KindGRPC, o.Kind())
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

func TestTransport_RequestHeader(t *testing.T) {
	v := headerCarrier{}
	v.Set("a", "1")
	o := &Transport{reqHeader: v}
	assert.Equal(t, "1", o.RequestHeader().Get("a"))
	assert.Equal(t, "", o.RequestHeader().Get("notfound"))

}

func TestTransport_ReplyHeader(t *testing.T) {
	v := headerCarrier{}
	v.Set("a", "1")
	o := &Transport{replyHeader: v}
	assert.Equal(t, "1", o.ReplyHeader().Get("a"))
}

func TestHeaderCarrier_Keys(t *testing.T) {
	v := headerCarrier{}
	v.Set("abb", "1")
	v.Set("bcc", "2")
	assert.ElementsMatch(t, []string{"abb", "bcc"}, v.Keys())
}
