package http

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

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
	ctx        context.Context
	timeout    time.Duration
	userAgent  string
	transport  http.RoundTripper
	middleware middleware.Middleware
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
		ctx:       ctx,
		timeout:   500 * time.Millisecond,
		transport: http.DefaultTransport,
	}
	for _, o := range opts {
		o(options)
	}
	return &baseTransport{
		middleware: options.middleware,
		userAgent:  options.userAgent,
		timeout:    options.timeout,
		base:       options.transport,
	}, nil
}

type baseTransport struct {
	userAgent  string
	timeout    time.Duration
	base       http.RoundTripper
	middleware middleware.Middleware
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
		return t.base.RoundTrip(in.(*http.Request))
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
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		se := &errors.StatusError{Code: 2}
		if err := decodeResponse(res, se); err != nil {
			return err
		}
		return se
	}
	return decodeResponse(res, target)
}

func decodeResponse(res *http.Response, target interface{}) error {
	subtype := contentSubtype(res.Header.Get(contentTypeHeader))
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
