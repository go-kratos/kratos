package http

import (
	"bytes"
	"io/ioutil"
	nethttp "net/http"
	"testing"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/stretchr/testify/assert"
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

type mockResponseWriter struct {
	StatusCode int
	Data       []byte
	header     nethttp.Header
}

func (w *mockResponseWriter) Header() nethttp.Header {
	return w.header
}

func (w *mockResponseWriter) Write(b []byte) (int, error) {
	w.Data = b
	return len(b), nil
}

func (w *mockResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}

type dataWithStatusCode struct {
	A string `json:"a"`
	B int64  `json:"b"`
}

func TestDefaultResponseEncoder(t *testing.T) {
	w := &mockResponseWriter{StatusCode: 200, header: make(nethttp.Header)}
	req1 := &nethttp.Request{
		Header: make(nethttp.Header),
	}
	req1.Header.Set("Content-Type", "application/json")

	v1 := &dataWithStatusCode{A: "1", B: 2}
	err := DefaultResponseEncoder(w, req1, v1)
	assert.Nil(t, err)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, 200, w.StatusCode)
	assert.NotNil(t, w.Data)
}

func TestDefaultResponseEncoderWithError(t *testing.T) {
	w := &mockResponseWriter{header: make(nethttp.Header)}
	req := &nethttp.Request{
		Header: make(nethttp.Header),
	}
	req.Header.Set("Content-Type", "application/json")

	se := &errors.Error{Code: 511}
	DefaultErrorEncoder(w, req, se)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, 511, w.StatusCode)
	assert.NotNil(t, w.Data)
}

func TestCodecForRequest(t *testing.T) {
	req1 := &nethttp.Request{
		Header: make(nethttp.Header),
		Body:   ioutil.NopCloser(bytes.NewBufferString("<xml></xml>")),
	}
	req1.Header.Set("Content-Type", "application/xml")

	c, ok := CodecForRequest(req1, "Content-Type")
	assert.True(t, ok)
	assert.Equal(t, "xml", c.Name())

	req2 := &nethttp.Request{
		Header: make(nethttp.Header),
		Body:   ioutil.NopCloser(bytes.NewBufferString("{\"a\":\"1\", \"b\": 2}")),
	}
	req2.Header.Set("Content-Type", "blablablabla")

	c, ok = CodecForRequest(req2, "Content-Type")
	assert.False(t, ok)
	assert.Equal(t, "json", c.Name())
}
