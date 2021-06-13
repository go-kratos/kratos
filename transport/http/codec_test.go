package http

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	nethttp "net/http"
	"testing"
)

func TestDefaultRequestDecoder(t *testing.T) {
	req1 := &nethttp.Request{
		Header: make(nethttp.Header),
		Body:   ioutil.NopCloser(bytes.NewBufferString("{\"a\":\"1\", \"b\": 2}")),
	}
	req1.Header.Set("Content-Type", "application/json")

	v1 := &struct {
		A string `json:"a"`
		B int64  `json:"b"`
	}{}
	err1 := DefaultRequestDecoder(req1, &v1)
	assert.Nil(t, err1)
	assert.Equal(t, "1", v1.A)
	assert.Equal(t, int64(2), v1.B)
}


func TestCodecForRequest(t *testing.T) {
	req1 := &nethttp.Request{
		Header: make(nethttp.Header),
		Body:   ioutil.NopCloser(bytes.NewBufferString("{\"a\":\"1\", \"b\": 2}")),
	}
	req1.Header.Set("Content-Type", "application/xml")

	c, ok := CodecForRequest(req1,"Content-Type")
	assert.True(t, ok)
	assert.Equal(t, "xml", c.Name())

	req2 := &nethttp.Request{
		Header: make(nethttp.Header),
		Body:   ioutil.NopCloser(bytes.NewBufferString("{\"a\":\"1\", \"b\": 2}")),
	}
	req2.Header.Set("Content-Type", "blablablabla")

	c, ok = CodecForRequest(req2,"Content-Type")
	assert.False(t, ok)
	assert.Equal(t, "json", c.Name())
}

