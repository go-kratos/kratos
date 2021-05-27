package http

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/internal/httputil"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http/binding"
)

// SupportPackageIsVersion1 These constants should not be referenced from any other code.
const SupportPackageIsVersion1 = true

// DecodeRequestFunc is decode request func.
type DecodeRequestFunc func(*http.Request, interface{}) error

// EncodeResponseFunc is encode response func.
type EncodeResponseFunc func(http.ResponseWriter, *http.Request, interface{}) error

// EncodeErrorFunc is encode error func.
type EncodeErrorFunc func(http.ResponseWriter, *http.Request, error)

// HandleOption is handle option.
type HandleOption func(*HandleOptions)

// HandleOptions is handle options.
// Deprecated: use Handler instead.
type HandleOptions struct {
	Decode     DecodeRequestFunc
	Encode     EncodeResponseFunc
	Error      EncodeErrorFunc
	Middleware middleware.Middleware
}

// DefaultHandleOptions returns a default handle options.
// Deprecated: use NewHandler instead.
func DefaultHandleOptions() HandleOptions {
	return HandleOptions{
		Decode:     defaultRequestDecoder,
		Encode:     defaultResponseEncoder,
		Error:      defaultErrorEncoder,
		Middleware: recovery.Recovery(),
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
func Middleware(m ...middleware.Middleware) HandleOption {
	return func(o *HandleOptions) {
		o.Middleware = middleware.Chain(m...)
	}
}

// Handler is handle options.
type Handler struct {
	method reflect.Value
	in     reflect.Type
	out    reflect.Type
	opts   HandleOptions
}

// NewHandler new a HTTP handler.
func NewHandler(handler interface{}, opts ...HandleOption) http.Handler {
	if err := validateHandler(handler); err != nil {
		panic(err)
	}
	typ := reflect.TypeOf(handler)
	h := &Handler{
		method: reflect.ValueOf(handler),
		in:     typ.In(1).Elem(),
		out:    typ.Out(0).Elem(),
		opts:   DefaultHandleOptions(),
	}
	for _, o := range opts {
		o(&h.opts)
	}
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	in := reflect.New(h.in).Interface()
	if err := h.opts.Decode(req, in); err != nil {
		h.opts.Error(w, req, err)
		return
	}
	invoke := func(ctx context.Context, in interface{}) (interface{}, error) {
		ret := h.method.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(in),
		})
		if ret[1].IsNil() {
			return ret[0].Interface(), nil
		}
		return nil, ret[1].Interface().(error)
	}
	if h.opts.Middleware != nil {
		invoke = h.opts.Middleware(invoke)
	}
	out, err := invoke(req.Context(), in)
	if err != nil {
		h.opts.Error(w, req, err)
		return
	}
	if err := h.opts.Encode(w, req, out); err != nil {
		h.opts.Error(w, req, err)
	}
}

func validateHandler(handler interface{}) error {
	typ := reflect.TypeOf(handler)
	if typ.NumIn() != 2 || typ.NumOut() != 2 {
		return fmt.Errorf("invalid types, in: %d out: %d", typ.NumIn(), typ.NumOut())
	}
	if typ.In(1).Kind() != reflect.Ptr || typ.Out(0).Kind() != reflect.Ptr {
		return fmt.Errorf("invalid types is not a pointer")
	}
	if !typ.In(0).Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
		return fmt.Errorf("input does not implement the context")
	}
	if !typ.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return fmt.Errorf("input does not implement the error")
	}
	return nil
}

// defaultRequestDecoder decodes the request body to object.
func defaultRequestDecoder(req *http.Request, v interface{}) error {
	subtype := httputil.ContentSubtype(req.Header.Get("Content-Type"))
	if codec := encoding.GetCodec(subtype); codec != nil {
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return errors.BadRequest("CODEC", err.Error())
		}
		if err := codec.Unmarshal(data, v); err != nil {
			return errors.BadRequest("CODEC", err.Error())
		}
	} else {
		if err := binding.BindForm(req, v); err != nil {
			return errors.BadRequest("CODEC", err.Error())
		}
	}
	return nil
}

// defaultResponseEncoder encodes the object to the HTTP response.
func defaultResponseEncoder(w http.ResponseWriter, r *http.Request, v interface{}) error {
	codec := CodecForRequest(r)
	data, err := codec.Marshal(v)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", httputil.ContentType(codec.Name()))
	if sc, ok := v.(interface {
		StatusCode() int
	}); ok {
		w.WriteHeader(sc.StatusCode())
	}
	_, _ = w.Write(data)
	return nil
}

// defaultErrorEncoder encodes the error to the HTTP response.
func defaultErrorEncoder(w http.ResponseWriter, r *http.Request, se error) {
	codec := CodecForRequest(r)
	body, err := codec.Marshal(se)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", httputil.ContentType(codec.Name()))
	if sc, ok := se.(interface {
		StatusCode() int
	}); ok {
		w.WriteHeader(sc.StatusCode())
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(body)
}

// CodecForRequest get encoding.Codec via http.Request
func CodecForRequest(r *http.Request) encoding.Codec {
	var codec encoding.Codec
	for _, accept := range r.Header["Accept"] {
		if codec = encoding.GetCodec(httputil.ContentSubtype(accept)); codec != nil {
			break
		}
	}
	if codec == nil {
		codec = encoding.GetCodec("json")
	}
	return codec
}
