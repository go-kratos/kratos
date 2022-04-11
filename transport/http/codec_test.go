package http

import (
	"bytes"
	"io"
	nethttp "net/http"
	"reflect"
	"testing"

	"github.com/go-kratos/kratos/v2/errors"
)

func TestDefaultRequestDecoder(t *testing.T) {
	req1 := &nethttp.Request{
		Header: make(nethttp.Header),
		Body:   io.NopCloser(bytes.NewBufferString("{\"a\":\"1\", \"b\": 2}")),
	}
	req1.Header.Set("Content-Type", "application/json")

	v1 := &struct {
		A string `json:"a"`
		B int64  `json:"b"`
	}{}
	err1 := DefaultRequestDecoder(req1, &v1)
	if err1 != nil {
		t.Errorf("expected no error, got %v", err1)
	}
	if !reflect.DeepEqual("1", v1.A) {
		t.Errorf("expected %v, got %v", "1", v1.A)
	}
	if !reflect.DeepEqual(int64(2), v1.B) {
		t.Errorf("expected %v, got %v", 2, v1.B)
	}
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
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual("application/json", w.Header().Get("Content-Type")) {
		t.Errorf("expected %v, got %v", "application/json", w.Header().Get("Content-Type"))
	}
	if !reflect.DeepEqual(200, w.StatusCode) {
		t.Errorf("expected %v, got %v", 200, w.StatusCode)
	}
	if w.Data == nil {
		t.Errorf("expected not nil, got %v", w.Data)
	}
}

func TestDefaultResponseEncoderWithError(t *testing.T) {
	w := &mockResponseWriter{header: make(nethttp.Header)}
	req := &nethttp.Request{
		Header: make(nethttp.Header),
	}
	req.Header.Set("Content-Type", "application/json")

	se := errors.New(511, "", "")
	DefaultErrorEncoder(w, req, se)
	if !reflect.DeepEqual("application/json", w.Header().Get("Content-Type")) {
		t.Errorf("expected %v, got %v", "application/json", w.Header().Get("Content-Type"))
	}
	if !reflect.DeepEqual(511, w.StatusCode) {
		t.Errorf("expected %v, got %v", 511, w.StatusCode)
	}
	if w.Data == nil {
		t.Errorf("expected not nil, got %v", w.Data)
	}
}

func TestCodecForRequest(t *testing.T) {
	req1 := &nethttp.Request{
		Header: make(nethttp.Header),
		Body:   io.NopCloser(bytes.NewBufferString("<xml></xml>")),
	}
	req1.Header.Set("Content-Type", "application/xml")

	c, ok := CodecForRequest(req1, "Content-Type")
	if !ok {
		t.Errorf("expected true, got %v", ok)
	}
	if !reflect.DeepEqual("xml", c.Name()) {
		t.Errorf("expected %v, got %v", "xml", c.Name())
	}

	req2 := &nethttp.Request{
		Header: make(nethttp.Header),
		Body:   io.NopCloser(bytes.NewBufferString("{\"a\":\"1\", \"b\": 2}")),
	}
	req2.Header.Set("Content-Type", "blablablabla")

	c, ok = CodecForRequest(req2, "Content-Type")
	if ok {
		t.Errorf("expected false, got %v", ok)
	}
	if !reflect.DeepEqual("json", c.Name()) {
		t.Errorf("expected %v, got %v", "json", c.Name())
	}
}
