package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/http/binding"
)

// SupportPackageIsVersion1 These constants should not be referenced from any other code.
const SupportPackageIsVersion1 = true

var _ StatusCoder = (*errors.StatusError)(nil)

// DecodeRequestFunc is decode request func.
type DecodeRequestFunc func(*http.Request, interface{}) error

// EncodeResponseFunc is encode response func.
type EncodeResponseFunc func(http.ResponseWriter, interface{}) error

// EncodeErrorFunc is encode error func.
type EncodeErrorFunc func(http.ResponseWriter, error)

// StatusCoder is returns the HTTPStatus code.
type StatusCoder interface {
	HTTPStatus() int
}

// HandleOption is handle option.
type HandleOption func(*HandleOptions)

// HandleOptions is handle options.
type HandleOptions struct {
	Decode     DecodeRequestFunc
	Encode     EncodeResponseFunc
	Error      EncodeErrorFunc
	Middleware middleware.Middleware
}

// DefaultHandleOptions returns a default handle options.
func DefaultHandleOptions() HandleOptions {
	return HandleOptions{
		Decode: DecodeRequest,
		Encode: EncodeResponse,
		Error:  EncodeError,
	}
}

// RequestDecoder with request decoder.
func RequestDecoder(dec DecodeRequestFunc) HandleOption {
	return func(o *HandleOptions) {
		o.Decode = dec
	}
}

// ResponseEncoder with response encoder.
func ResponseEncoder(en EncodeResponseFunc) HandleOption {
	return func(o *HandleOptions) {
		o.Encode = en
	}
}

// ErrorEncoder with error encoder.
func ErrorEncoder(en EncodeErrorFunc) HandleOption {
	return func(o *HandleOptions) {
		o.Error = en
	}
}

// Middleware with middleware option.
func Middleware(m middleware.Middleware) HandleOption {
	return func(o *HandleOptions) {
		o.Middleware = m
	}
}

// DecodeRequest decodes the request body to object.
func DecodeRequest(req *http.Request, v interface{}) error {
	switch stripContentType(req.Header.Get("content-type")) {
	case "application/json":
		return json.NewDecoder(req.Body).Decode(v)
	}
	return binding.BindForm(req, v)
}

// EncodeResponse encodes the object to the HTTP response.
func EncodeResponse(w http.ResponseWriter, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

// EncodeError encodes the erorr to the HTTP response.
func EncodeError(w http.ResponseWriter, err error) {
	if c, ok := err.(StatusCoder); ok {
		w.WriteHeader(c.HTTPStatus())
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(err)
}

func stripContentType(contentType string) string {
	i := strings.Index(contentType, ";")
	if i != -1 {
		contentType = contentType[:i]
	}
	return contentType
}
