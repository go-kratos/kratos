package rpcx

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/internal/endpoint"
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

// clientOptions is RPCx Client
type clientOptions struct {
	endpoint   string
	tlsConf    *tls.Config
	timeout    time.Duration
	discovery  registry.Discovery
	middleware []middleware.Middleware
	rpcXOpts   client.Option
	selectMode client.SelectMode
	failMode   client.FailMode
}

// Dial returns a RPCx connection.
func Dial(ctx context.Context, opts ...ClientOption) (client.XClient, error) {
	return dial(ctx, false, opts...)
}

// DialInsecure returns an insecure RPCx connection.
func DialInsecure(ctx context.Context, opts ...ClientOption) (client.XClient, error) {
	return dial(ctx, true, opts...)
}

func dial(ctx context.Context, insecure bool, opts ...ClientOption) (client.XClient, error) {
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
	var d client.ServiceDiscovery
	if options.discovery != nil {
		var KVPair []*client.KVPair
		if options.discovery != nil {
			service, err := options.discovery.GetService(ctx, options.endpoint)
			if err != nil {
				panic(err)
			}
			for _, instance := range service {
				endpoint, err := endpoint.ParseEndpoint(instance.Endpoints, "rpcx", !insecure)
				if err != nil {
					//r.log.Errorf("[resolver] Failed to parse discovery endpoint: %v", err)
					continue
				}
				if endpoint == "" {
					continue
				}
				meta, _ := json.Marshal(instance.Metadata)
				KVPair = append(KVPair, &client.KVPair{
					Key:   endpoint,
					Value: string(meta),
				})
			}
		}
		fmt.Println(&KVPair[0])
		d, _ = client.NewMultipleServersDiscovery(KVPair)
	} else {
		d, _ = client.NewPeer2PeerDiscovery(options.endpoint, "")
	}
	xclient := client.NewXClient("Greeter", options.failMode, options.selectMode, d, options.rpcXOpts)
	return xclient, nil
}
