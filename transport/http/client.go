package http

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-kratos/kratos/v2/encoding"
	xhttp "github.com/go-kratos/kratos/v2/internal/http"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

// DecodeErrorFunc is decode error func.
type DecodeErrorFunc func(ctx context.Context, w *http.Response) error

// ClientOption is HTTP client option.
type ClientOption func(*clientOptions)

// WithTimeout with client request timeout.
func WithTimeout(d time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.timeout = d
	}
}

// WithUserAgent with client user agent.
func WithUserAgent(ua string) ClientOption {
	return func(o *clientOptions) {
		o.userAgent = ua
	}
}

// WithTransport with client transport.
func WithTransport(trans http.RoundTripper) ClientOption {
	return func(o *clientOptions) {
		o.transport = trans
	}
}

// WithMiddleware with client middleware.
func WithMiddleware(m middleware.Middleware) ClientOption {
	return func(o *clientOptions) {
		o.middleware = m
	}
}

// Client is a HTTP transport client.
type clientOptions struct {
	ctx          context.Context
	timeout      time.Duration
	userAgent    string
	transport    http.RoundTripper
	errorDecoder DecodeErrorFunc
	middleware   middleware.Middleware
}

// NewClient returns an HTTP client.
func NewClient(ctx context.Context, opts ...ClientOption) (*http.Client, error) {
	trans, err := NewTransport(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &http.Client{Transport: trans}, nil
}

// NewTransport creates an http.RoundTripper.
func NewTransport(ctx context.Context, opts ...ClientOption) (http.RoundTripper, error) {
	options := &clientOptions{
		ctx:          ctx,
		timeout:      500 * time.Millisecond,
		transport:    http.DefaultTransport,
		errorDecoder: CheckResponse,
	}
	for _, o := range opts {
		o(options)
	}
	return &baseTransport{
		errorDecoder: options.errorDecoder,
		middleware:   options.middleware,
		userAgent:    options.userAgent,
		timeout:      options.timeout,
		base:         options.transport,
	}, nil
}

type baseTransport struct {
	userAgent    string
	timeout      time.Duration
	base         http.RoundTripper
	errorDecoder DecodeErrorFunc
	middleware   middleware.Middleware
}

func (t *baseTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.userAgent != "" && req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", t.userAgent)
	}
	ctx := transport.NewContext(req.Context(), transport.Transport{Kind: transport.KindHTTP})
	ctx = NewClientContext(ctx, ClientInfo{Request: req})
	ctx, cancel := context.WithTimeout(ctx, t.timeout)
	defer cancel()

	h := func(ctx context.Context, in interface{}) (interface{}, error) {
		res, err := t.base.RoundTrip(in.(*http.Request))
		if err != nil {
			return nil, err
		}
		if err := t.errorDecoder(ctx, res); err != nil {
			return nil, err
		}
		return res, nil
	}
	if t.middleware != nil {
		h = t.middleware(h)
	}
	res, err := h(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*http.Response), nil
}

// Do send an HTTP request and decodes the body of response into target.
// returns an error (of type *Error) if the response status code is not 2xx.
func Do(client *http.Client, req *http.Request, target interface{}) error {
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	subtype := xhttp.ContentSubtype(res.Header.Get(xhttp.HeaderContentType))
	codec := encoding.GetCodec(subtype)
	if codec == nil {
		codec = encoding.GetCodec("json")
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return codec.Unmarshal(data, target)
}

// CheckResponse returns an error (of type *Error) if the response
// status code is not 2xx.
func CheckResponse(ctx context.Context, res *http.Response) error {
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return nil
	}
	defer res.Body.Close()
	if data, err := ioutil.ReadAll(res.Body); err == nil {
		st := new(spb.Status)
		if err = protojson.Unmarshal(data, st); err == nil {
			return status.ErrorProto(st)
		}
	}
	return status.Error(xhttp.GRPCCodeFromStatus(res.StatusCode), res.Status)
}
