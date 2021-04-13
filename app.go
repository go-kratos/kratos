package kratos

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

// App is an application components lifecycle manager
type App struct {
	opts     options
	instance *registry.ServiceInstance
	log      *log.Helper
}

// New create an application lifecycle manager.
func New(opts ...Option) *App {
	options := options{
		logger: log.DefaultLogger,
	}
	if id, err := uuid.NewUUID(); err == nil {
		options.id = id.String()
	}
	for _, o := range opts {
		o(&options)
	}
	return &App{
		opts:     options,
		instance: buildInstance(options),
		log:      log.NewHelper("app", options.logger),
	}
}

// Run executes all OnStart hooks registered with the application's Lifecycle.
func (a *App) Run(ctx context.Context) error {
	a.log.Infow(
		"service_id", a.opts.id,
		"service_name", a.opts.name,
		"version", a.opts.version,
	)
	g, ctx := errgroup.WithContext(ctx)
	for _, srv := range a.opts.servers {
		srv := srv
		g.Go(func() error {
			<-ctx.Done() // wait for stop signal
			a.Stop()
			return srv.Stop()
		})
		g.Go(func() error {
			return srv.Start()
		})
	}
	if a.opts.registrar != nil {
		if err := a.opts.registrar.Register(ctx, a.instance); err != nil {
			return err
		}
	}
	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

// Stop gracefully stops the application.
func (a *App) Stop() error {
	if a.opts.registrar != nil {
		if err := a.opts.registrar.Deregister(context.Background(), a.instance); err != nil {
			return err
		}
	}
	return nil
}

func buildInstance(o options) *registry.ServiceInstance {
	if len(o.endpoints) == 0 {
		for _, srv := range o.servers {
			if e, err := srv.Endpoint(); err == nil {
				o.endpoints = append(o.endpoints, e)
			}
		}
	}
	return &registry.ServiceInstance{
		ID:        o.id,
		Name:      o.name,
		Version:   o.version,
		Metadata:  o.metadata,
		Endpoints: o.endpoints,
	}
}

// OnInterrupt returns a context that is canceled on SIGINT, SIGQUIT and SIGTERM.
func OnInterrupt() (context.Context, func()) {
	return Wrap(context.Background(), syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
}

// Wrap returns a new context wrapping the provided context and
// canceling it on the provided signals.
func Wrap(ctx context.Context, signals ...os.Signal) (context.Context, func()) {
	ctx, closer := context.WithCancel(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, signals...)

	go func() {
		select {
		case <-c:
			closer()
		case <-ctx.Done():
		}
	}()

	return ctx, closer
}
