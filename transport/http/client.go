package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/internal/httputil"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http/balancer"
	"github.com/go-kratos/kratos/v2/transport/http/balancer/random"
)

// DecodeErrorFunc is decode error func.
type DecodeErrorFunc func(ctx context.Context, res *http.Response) error

// EncodeRequestFunc is request encode func.
type EncodeRequestFunc func(ctx context.Context, in interface{}) (contentType string, body []byte, err error)

// DecodeResponseFunc is response decode func.
type DecodeResponseFunc func(ctx context.Context, res *http.Response, out interface{}) error

// ClientOption is HTTP client option.
type ClientOption func(*clientOptions)

// Client is an HTTP transport client.
type clientOptions struct {
	ctx          context.Context
	timeout      time.Duration
	endpoint     string
	userAgent    string
	encoder      EncodeRequestFunc
	decoder      DecodeResponseFunc
	errorDecoder DecodeErrorFunc
	transport    http.RoundTripper
	balancer     balancer.Balancer
	discovery    registry.Discovery
	middleware   middleware.Middleware
}

// WithTransport with client transport.
func WithTransport(trans http.RoundTripper) ClientOption {
	return func(o *clientOptions) {
		o.transport = trans
	}
}

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

// WithMiddleware with client middleware.
func WithMiddleware(m ...middleware.Middleware) ClientOption {
	return func(o *clientOptions) {
		o.middleware = middleware.Chain(m...)
	}
}

// WithEndpoint with client addr.
func WithEndpoint(endpoint string) ClientOption {
	return func(o *clientOptions) {
		o.endpoint = endpoint
	}
}

// WithRequestEncoder with client request encoder.
func WithRequestEncoder(encoder EncodeRequestFunc) ClientOption {
	return func(o *clientOptions) {
		o.encoder = encoder
	}
}

// WithResponseDecoder with client response decoder.
func WithResponseDecoder(decoder DecodeResponseFunc) ClientOption {
	return func(o *clientOptions) {
		o.decoder = decoder
	}
}

// WithErrorDecoder with client error decoder.
func WithErrorDecoder(errorDecoder DecodeErrorFunc) ClientOption {
	return func(o *clientOptions) {
		o.errorDecoder = errorDecoder
	}
}

// WithDiscovery with client discovery.
func WithDiscovery(d registry.Discovery) ClientOption {
	return func(o *clientOptions) {
		o.discovery = d
	}
}

// WithBalancer with client balancer.
// Experimental
// Notice: This type is EXPERIMENTAL and may be changed or removed in a later release.
func WithBalancer(b balancer.Balancer) ClientOption {
	return func(o *clientOptions) {
		o.balancer = b
	}
}

// Client is an HTTP client.
type Client struct {
	opts   clientOptions
	target *Target
	r      *resolver
	cc     *http.Client
}

// NewClient returns an HTTP client.
func NewClient(ctx context.Context, opts ...ClientOption) (*Client, error) {
	options := clientOptions{
		ctx:          ctx,
		timeout:      500 * time.Millisecond,
		encoder:      DefaultRequestEncoder,
		decoder:      DefaultResponseDecoder,
		errorDecoder: DefaultErrorDecoder,
		transport:    http.DefaultTransport,
		balancer:     random.New(),
	}
	for _, o := range opts {
		o(&options)
	}
	target, err := parseTarget(options.endpoint)
	if err != nil {
		return nil, err
	}
	var r *resolver
	if options.discovery != nil {
		if target.Scheme == "discovery" {
			if r, err = newResolver(ctx, options.discovery, target); err != nil {
				return nil, fmt.Errorf("[http client] new resolver failed!err: %v", options.endpoint)
			}
		} else {
			return nil, fmt.Errorf("[http client] invalid endpoint format: %v", options.endpoint)
		}
	}
	return &Client{
		opts:   options,
		target: target,
		r:      r,
		cc: &http.Client{
			Timeout:   options.timeout,
			Transport: options.transport,
		},
	}, nil
}

// Invoke makes an rpc call procedure for remote service.
func (client *Client) Invoke(ctx context.Context, path string, args interface{}, reply interface{}, opts ...CallOption) error {
	var (
		reqBody     io.Reader
		contentType string
	)
	c := defaultCallInfo()
	for _, o := range opts {
		if err := o.before(&c); err != nil {
			return err
		}
	}
	if args != nil {
		var (
			body []byte
			err  error
		)
		contentType, body, err = client.opts.encoder(ctx, args)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(body)
	}
	url := fmt.Sprintf("%s://%s%s", client.target.Scheme, client.target.Authority, path)
	req, err := http.NewRequest(c.method, url, reqBody)
	if err != nil {
		return err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if client.opts.userAgent != "" {
		req.Header.Set("User-Agent", client.opts.userAgent)
	}
	ctx = transport.NewContext(ctx, transport.Transport{Kind: transport.KindHTTP, Endpoint: client.opts.endpoint})
	ctx = NewClientContext(ctx, ClientInfo{PathPattern: c.pathPattern, Request: req})
	return client.invoke(ctx, req, args, reply, c)
}

func (client *Client) invoke(ctx context.Context, req *http.Request, args interface{}, reply interface{}, c callInfo) error {
	h := func(ctx context.Context, in interface{}) (interface{}, error) {
		var done func(context.Context, balancer.DoneInfo)
		if client.r != nil {
			var (
				err   error
				node  *registry.ServiceInstance
				nodes = client.r.fetch(ctx)
			)
			if node, done, err = client.opts.balancer.Pick(ctx, c.pathPattern, nodes); err != nil {
				return nil, errors.ServiceUnavailable("NODE_NOT_FOUND", err.Error())
			}
			scheme, addr, err := parseEndpoint(node.Endpoints)
			if err != nil {
				return nil, errors.ServiceUnavailable("NODE_NOT_FOUND", err.Error())
			}
			req = req.Clone(ctx)
			req.URL.Scheme = scheme
			req.URL.Host = addr
		}
		res, err := client.do(ctx, req, c)
		if done != nil {
			done(ctx, balancer.DoneInfo{Err: err})
		}
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		if err := client.opts.decoder(ctx, res, reply); err != nil {
			return nil, err
		}
		return reply, nil
	}
	if client.opts.middleware != nil {
		h = client.opts.middleware(h)
	}
	_, err := h(ctx, args)
	return err
}

// Do send an HTTP request and decodes the body of response into target.
// returns an error (of type *Error) if the response status code is not 2xx.
func (client *Client) Do(req *http.Request, opts ...CallOption) (*http.Response, error) {
	c := defaultCallInfo()
	for _, o := range opts {
		if err := o.before(&c); err != nil {
			return nil, err
		}
	}
	return client.do(req.Context(), req, c)
}

func (client *Client) do(ctx context.Context, req *http.Request, c callInfo) (*http.Response, error) {
	resp, err := client.cc.Do(req)
	if err != nil {
		return nil, err
	}
	if err := client.opts.errorDecoder(ctx, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// DefaultRequestEncoder is an HTTP request encoder.
func DefaultRequestEncoder(ctx context.Context, in interface{}) (string, []byte, error) {
	body, err := encoding.GetCodec("json").Marshal(in)
	if err != nil {
		return "", nil, err
	}
	return "application/json", body, err
}

// DefaultResponseDecoder is an HTTP response decoder.
func DefaultResponseDecoder(ctx context.Context, res *http.Response, v interface{}) error {
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return CodecForResponse(res).Unmarshal(data, v)
}

// DefaultErrorDecoder is an HTTP error decoder.
func DefaultErrorDecoder(ctx context.Context, res *http.Response) error {
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return nil
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err == nil {
		e := new(errors.Error)
		if err = CodecForResponse(res).Unmarshal(data, e); err == nil {
			return e
		}
	}
	return errors.Errorf(res.StatusCode, errors.UnknownReason, err.Error())
}

// CodecForResponse get encoding.Codec via http.Response
func CodecForResponse(r *http.Response) encoding.Codec {
	codec := encoding.GetCodec(httputil.ContentSubtype("Content-Type"))
	if codec != nil {
		return codec
	}
	return encoding.GetCodec("json")
}
