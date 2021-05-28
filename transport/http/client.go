package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
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

// Client is an HTTP client.
type Client struct {
	cc *http.Client
	r  *resolver
	b  balancer.Balancer

	scheme       string
	target       Target
	userAgent    string
	middleware   middleware.Middleware
	encoder      EncodeRequestFunc
	decoder      DecodeResponseFunc
	errorDecoder DecodeErrorFunc
	discovery    registry.Discovery
}

// DecodeErrorFunc is decode error func.
type DecodeErrorFunc func(ctx context.Context, res *http.Response) error

// EncodeRequestFunc is request encode func.
type EncodeRequestFunc func(ctx context.Context, in interface{}) (contentType string, body []byte, err error)

// DecodeResponseFunc is response decode func.
type DecodeResponseFunc func(ctx context.Context, res *http.Response, out interface{}) error

// ClientOption is HTTP client option.
type ClientOption func(*clientOptions)

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

// WithScheme with client schema.
func WithScheme(scheme string) ClientOption {
	return func(o *clientOptions) {
		o.scheme = scheme
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

// Client is an HTTP transport client.
type clientOptions struct {
	ctx          context.Context
	transport    http.RoundTripper
	middleware   middleware.Middleware
	timeout      time.Duration
	scheme       string
	endpoint     string
	userAgent    string
	encoder      EncodeRequestFunc
	decoder      DecodeResponseFunc
	errorDecoder DecodeErrorFunc
	discovery    registry.Discovery
	balancer     balancer.Balancer
}

// NewClient returns an HTTP client.
func NewClient(ctx context.Context, opts ...ClientOption) (*Client, error) {
	options := &clientOptions{
		ctx:          ctx,
		scheme:       "http",
		timeout:      1 * time.Second,
		encoder:      DefaultRequestEncoder,
		decoder:      DefaultResponseDecoder,
		errorDecoder: DefaultErrorDecoder,
		transport:    http.DefaultTransport,
		balancer:     random.New(),
	}
	for _, o := range opts {
		o(options)
	}
	target := Target{
		Scheme:   options.scheme,
		Endpoint: options.endpoint,
	}
	var r *resolver
	if options.endpoint != "" && options.discovery != nil {
		u, err := url.Parse(options.endpoint)
		if err != nil {
			u, err = url.Parse("http://" + options.endpoint)
			if err != nil {
				return nil, fmt.Errorf("[http client] invalid endpoint format: %v", options.endpoint)
			}
		}
		if u.Scheme == "discovery" && len(u.Path) > 1 {
			target = Target{
				Scheme:    u.Scheme,
				Authority: u.Host,
				Endpoint:  u.Path[1:],
			}
			r, err = newResolver(ctx, options.scheme, options.discovery, target)
			if err != nil {
				return nil, fmt.Errorf("[http client] new resolver failed!err: %v", options.endpoint)
			}
		} else {
			return nil, fmt.Errorf("[http client] invalid endpoint format: %v", options.endpoint)
		}
	}

	return &Client{
		cc:           &http.Client{Timeout: options.timeout, Transport: options.transport},
		r:            r,
		encoder:      options.encoder,
		decoder:      options.decoder,
		errorDecoder: options.errorDecoder,
		middleware:   options.middleware,
		userAgent:    options.userAgent,
		target:       target,
		scheme:       options.scheme,
		discovery:    options.discovery,
		b:            options.balancer,
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
		contentType, body, err = client.encoder(ctx, args)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(body)
	}
	url := fmt.Sprintf("%s://%s%s", client.scheme, client.target.Endpoint, path)
	req, err := http.NewRequest(c.method, url, reqBody)
	if err != nil {
		return err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if client.userAgent != "" {
		req.Header.Set("User-Agent", client.userAgent)
	}

	ctx = transport.NewContext(ctx, transport.Transport{Kind: transport.KindHTTP})
	ctx = NewClientContext(ctx, ClientInfo{
		PathPattern: c.pathPattern,
		Request:     req,
	})

	return client.invoke(ctx, req, args, reply, c)
}

func (client *Client) invoke(ctx context.Context, req *http.Request, args interface{}, reply interface{}, c callInfo) error {
	h := func(ctx context.Context, in interface{}) (interface{}, error) {
		var done func(context.Context, balancer.DoneInfo)
		if client.r != nil {
			nodes := client.r.fetch(ctx)
			if len(nodes) == 0 {
				return nil, errors.ServiceUnavailable("NODE_NOT_FOUND", "fetch error")
			}
			var node *registry.ServiceInstance
			var err error
			node, done, err = client.b.Pick(ctx, c.pathPattern, nodes)
			if err != nil {
				return nil, errors.ServiceUnavailable("NODE_NOT_FOUND", err.Error())
			}
			req = req.Clone(ctx)
			addr, err := parseEndpoint(client.scheme, node.Endpoints)
			if err != nil {
				return nil, errors.ServiceUnavailable("NODE_NOT_FOUND", err.Error())
			}
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
		if err := client.decoder(ctx, res, reply); err != nil {
			return nil, err
		}
		return reply, nil
	}
	if client.middleware != nil {
		h = client.middleware(h)
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
	if err := client.errorDecoder(ctx, resp); err != nil {
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
