package kratos

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"

	"github.com/google/uuid"
)

// AppInfo is application context value.
type AppInfo interface {
	ID() string
	Name() string
	Version() string
	Metadata() map[string]string
	Endpoint() []string
}

// App is an application components lifecycle manager.
type App struct {
	opts     options
	ctx      context.Context
	cancel   func()
	instance *registry.ServiceInstance

	started []transport.Server
}

// New create an application lifecycle manager.
func New(opts ...Option) *App {
	o := options{
		ctx:              context.Background(),
		logger:           log.NewHelper(log.DefaultLogger),
		sigs:             []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		registrarTimeout: 10 * time.Second,
	}
	if id, err := uuid.NewUUID(); err == nil {
		o.id = id.String()
	}
	for _, opt := range opts {
		opt(&o)
	}
	ctx, cancel := context.WithCancel(o.ctx)
	return &App{
		ctx:     ctx,
		cancel:  cancel,
		opts:    o,
		started: make([]transport.Server, 0, len(o.servers)),
	}
}

// ID returns app instance id.
func (a *App) ID() string { return a.opts.id }

// Name returns service name.
func (a *App) Name() string { return a.opts.name }

// Version returns app version.
func (a *App) Version() string { return a.opts.version }

// Metadata returns service metadata.
func (a *App) Metadata() map[string]string { return a.opts.metadata }

// Endpoint returns endpoints.
func (a *App) Endpoint() []string {
	if a.instance == nil {
		return []string{}
	}
	return a.instance.Endpoints
}

// Run executes all OnStart hooks registered with the application's Lifecycle.
func (a *App) Run() error {
	defer func() {
		e := a.Stop()
		if e != nil {
			a.opts.logger.Errorf("[kratos]failed to stop app: %v", e)
		}
	}()

	ctx := NewContext(a.ctx, a)
	endpoints := make([]string, 0, len(a.opts.servers)+len(a.opts.endpoints))
	for _, e := range a.opts.endpoints {
		endpoints = append(endpoints, e.String())
	}
	for _, srv := range a.opts.servers {
		err := srv.Start(ctx)
		if err != nil {
			return err
		}
		if r, ok := srv.(transport.Endpointer); ok && len(a.opts.endpoints) == 0 {
			e, err := r.Endpoint()
			if err != nil {
				return err
			}
			endpoints = append(endpoints, e.String())
		}
		a.started = append(a.started, srv)
	}

	instance, err := a.buildInstance(endpoints)
	if err != nil {
		return err
	}
	if a.opts.registrar != nil {
		rctx, rcancel := context.WithTimeout(a.opts.ctx, a.opts.registrarTimeout)
		defer rcancel()
		if err := a.opts.registrar.Register(rctx, instance); err != nil {
			return err
		}
		a.instance = instance
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)
	select {
	case <-ctx.Done():
		if !errors.Is(ctx.Err(), context.Canceled) {
			return ctx.Err()
		}
	case <-c:
	}
	return nil
}

// Stop gracefully stops the application.
func (a *App) Stop() error {
	if a.opts.registrar != nil && a.instance != nil {
		ctx, cancel := context.WithTimeout(a.opts.ctx, a.opts.registrarTimeout)
		defer cancel()
		if err := a.opts.registrar.Deregister(ctx, a.instance); err != nil {
			return err
		}
	}
	a.cancel()
	for _, srv := range a.started {
		if err := srv.Stop(context.Background()); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) buildInstance(endpoints []string) (*registry.ServiceInstance, error) {
	return &registry.ServiceInstance{
		ID:        a.opts.id,
		Name:      a.opts.name,
		Version:   a.opts.version,
		Metadata:  a.opts.metadata,
		Endpoints: endpoints,
	}, nil
}

type appKey struct{}

// NewContext returns a new Context that carries value.
func NewContext(ctx context.Context, s AppInfo) context.Context {
	return context.WithValue(ctx, appKey{}, s)
}

// FromContext returns the Transport value stored in ctx, if any.
func FromContext(ctx context.Context) (s AppInfo, ok bool) {
	s, ok = ctx.Value(appKey{}).(AppInfo)
	return
}
