package kratos

import (
	"context"
	"net/url"
	"os"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"
)

// Option is an application option.
type Option func(o *options)

// options is an application options.
type options struct {
	id        string
	name      string
	version   string
	metadata  map[string]string
	endpoints []*url.URL

	ctx  context.Context
	sigs []os.Signal

	logger           log.Logger
	registrar        registry.Registrar
	registrarTimeout time.Duration
	stopTimeout      time.Duration
	servers          []transport.Server

	// Before and After funcs
	beforeStart []func(context.Context) error
	beforeStop  []func(context.Context) error
	afterStart  []func(context.Context) error
	afterStop   []func(context.Context) error
}

// ID with service id.
func ID(id string) Option {
	return func(o *options) { o.id = id }
}

// Name with service name.
func Name(name string) Option {
	return func(o *options) { o.name = name }
}

// Version with service version.
func Version(version string) Option {
	return func(o *options) { o.version = version }
}

// Metadata with service metadata.
func Metadata(md map[string]string) Option {
	return func(o *options) { o.metadata = md }
}

// Endpoint with service endpoint.
func Endpoint(endpoints ...*url.URL) Option {
	return func(o *options) { o.endpoints = endpoints }
}

// Context with service context.
func Context(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// Logger with service logger.
func Logger(logger log.Logger) Option {
	return func(o *options) { o.logger = logger }
}

// Server with transport servers.
func Server(srv ...transport.Server) Option {
	return func(o *options) { o.servers = srv }
}

// Signal with exit signals.
func Signal(sigs ...os.Signal) Option {
	return func(o *options) { o.sigs = sigs }
}

// Registrar with service registry.
func Registrar(r registry.Registrar) Option {
	return func(o *options) { o.registrar = r }
}

// RegistrarTimeout with registrar timeout.
func RegistrarTimeout(t time.Duration) Option {
	return func(o *options) { o.registrarTimeout = t }
}

// StopTimeout with app stop timeout.
func StopTimeout(t time.Duration) Option {
	return func(o *options) { o.stopTimeout = t }
}

// Before and Afters

// BeforeStart run funcs before app starts
func BeforeStart(fn func(context.Context) error) Option {
	return func(o *options) {
		o.beforeStart = append(o.beforeStart, fn)
	}
}

// BeforeStop run funcs before app stops
func BeforeStop(fn func(context.Context) error) Option {
	return func(o *options) {
		o.beforeStop = append(o.beforeStop, fn)
	}
}

// AfterStart run funcs after app starts
func AfterStart(fn func(context.Context) error) Option {
	return func(o *options) {
		o.afterStart = append(o.afterStart, fn)
	}
}

// AfterStop run funcs after app stops
func AfterStop(fn func(context.Context) error) Option {
	return func(o *options) {
		o.afterStop = append(o.afterStop, fn)
	}
}
