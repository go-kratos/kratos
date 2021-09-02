package rpcx

import (
	"crypto/tls"
	"github.com/go-kratos/kratos/v2/transport/rpcx/resolver/discovery"
	"github.com/smallnest/rpcx/client"
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/registry"
)

// ClientOption is RPCx client option.
type ClientOption func(o *clientOptions)

// WithEndpoint with client endpoint.
func WithEndpoint(endpoint string) ClientOption {
	return func(o *clientOptions) {
		o.endpoint = endpoint
	}
}

// WithTimeout with client timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.timeout = timeout
	}
}

// WithMiddleware with client middleware.
func WithMiddleware(m ...middleware.Middleware) ClientOption {
	return func(o *clientOptions) {
		o.middleware = m
	}
}

// WithDiscovery with client discovery.
func WithDiscovery(d registry.Discovery) ClientOption {
	return func(o *clientOptions) {
		o.discovery = d
	}
}

// WithTLSConfig with TLS config.
func WithTLSConfig(c *tls.Config) ClientOption {
	return func(o *clientOptions) {
		o.tlsConf = c
	}
}

// WithOption with RPCx options.
func WithOption(opt client.Option) ClientOption {
	return func(o *clientOptions) {
		o.rpcXOpts = opt
	}
}

// WithSelectMode with RPCx options.
func WithSelectMode(s client.SelectMode) ClientOption {
	return func(o *clientOptions) {
		o.selectMode = s
	}
}

// WithFailMode with RPCx options.
func WithFailMode(f client.FailMode) ClientOption {
	return func(o *clientOptions) {
		o.failMode = f
	}
}

// WithServerPath with RPCx serverPath.
func WithServerPath(servicePath string) ClientOption {
	return func(o *clientOptions) {
		o.servicePath = servicePath
	}
}

// clientOptions is RPCx Client
type clientOptions struct {
	endpoint    string
	tlsConf     *tls.Config
	timeout     time.Duration
	discovery   registry.Discovery
	middleware  []middleware.Middleware
	rpcXOpts    client.Option
	selectMode  client.SelectMode
	failMode    client.FailMode
	servicePath string
}

// Dial returns a RPCx connection.
func Dial(opts ...ClientOption) (client.XClient, error) {
	return dial(false, opts...)
}

// DialInsecure returns an insecure RPCx connection.
func DialInsecure(opts ...ClientOption) (client.XClient, error) {
	return dial(true, opts...)
}

func dial(insecure bool, opts ...ClientOption) (client.XClient, error) {
	options := clientOptions{
		timeout:    2000 * time.Millisecond,
		rpcXOpts:   client.DefaultOption,
		selectMode: client.RandomSelect,
		failMode:   client.Failfast,
	}
	for _, o := range opts {
		o(&options)
	}
	if options.tlsConf != nil {
		options.rpcXOpts.TLSConfig = options.tlsConf
	}
	var xClient client.XClient
	if options.discovery != nil {
		cc, _ := client.NewMultipleServersDiscovery([]*client.KVPair{})
		b := discovery.NewBuilder(options.discovery, discovery.WithInsecure(insecure))
		err := b.Build(options.endpoint, cc)
		if err != nil {
			return nil, err
		}
		xClient = buildXClient(cc, options)
		return xClient, nil
	}
	cc, _ := client.NewPeer2PeerDiscovery(options.endpoint, "")
	xClient = buildXClient(cc, options)
	return xClient, nil
}

func buildXClient(cc client.ServiceDiscovery, options clientOptions) client.XClient {
	return client.NewXClient(options.servicePath, options.failMode, options.selectMode, cc, options.rpcXOpts)
}
