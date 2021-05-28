package kratos

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

// App is an application components lifecycle manager
type App struct {
	opts     options
	ctx      context.Context
	cancel   func()
	instance *registry.ServiceInstance
	log      *log.Helper
}

// New create an application lifecycle manager.
func New(opts ...Option) *App {
	options := options{
		ctx:    context.Background(),
		logger: log.DefaultLogger,
		sigs:   []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
	}
	if id, err := uuid.NewUUID(); err == nil {
		options.id = id.String()
	}
	for _, o := range opts {
		o(&options)
	}
	ctx, cancel := context.WithCancel(options.ctx)
	return &App{
		ctx:    ctx,
		cancel: cancel,
		opts:   options,
		log:    log.NewHelper(options.logger),
	}
}

// Run executes all OnStart hooks registered with the application's Lifecycle.
func (a *App) Run() error {
	a.log.Infow(
		"service_id", a.opts.id,
		"service_name", a.opts.name,
		"version", a.opts.version,
	)
	instance, err := buildInstance(a.opts)
	if err != nil {
		return err
	}
	eg, ctx := errgroup.WithContext(a.ctx)
	wg := sync.WaitGroup{}
	for _, srv := range a.opts.servers {
		srv := srv
		eg.Go(func() error {
			<-ctx.Done() // wait for stop signal
			return srv.Stop()
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			return srv.Start()
		})
	}
	wg.Wait()
	if a.opts.registrar != nil {
		if err := a.opts.registrar.Register(a.opts.ctx, instance); err != nil {
			return err
		}
		a.instance = instance
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				a.Stop()
			}
		}
	})
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

// Stop gracefully stops the application.
func (a *App) Stop() error {
	if a.opts.registrar != nil && a.instance != nil {
		if err := a.opts.registrar.Deregister(a.opts.ctx, a.instance); err != nil {
			return err
		}
	}
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}

func buildInstance(o options) (*registry.ServiceInstance, error) {
	if len(o.endpoints) == 0 {
		for _, srv := range o.servers {
			if r, ok := srv.(transport.Endpointer); ok {
				e, err := r.Endpoint()
				if err != nil {
					return nil, err
				}
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
	}, nil
}
