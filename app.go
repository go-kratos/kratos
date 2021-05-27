package kratos

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

// Service is an instance of a service in a discovery system.
type Service struct {
	id        string
	name      string
	version   string
	metadata  map[string]string
	endpoints []string
}

// ID is service id
func (s *Service) ID() string {
	return s.id
}

// Name is service name
func (s *Service) Name() string {
	return s.name
}

// Version is service Version
func (s *Service) Version() string {
	return s.version
}

// Metadata is service Metadata
func (s *Service) Metadata() map[string]string {
	return s.metadata
}

// Endpoints is service Endpoints
func (s *Service) Endpoints() []string {
	return s.endpoints
}

// App is an application components lifecycle manager
type App struct {
	opts     options
	ctx      context.Context
	cancel   func()
	instance *Service
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
		opts:     options,
		ctx:      ctx,
		cancel:   cancel,
		instance: buildInstance(options),
		log:      log.NewHelper(options.logger),
	}
}

// Run executes all OnStart hooks registered with the application's Lifecycle.
func (a *App) Run() error {
	a.log.Infow(
		"service_id", a.opts.id,
		"service_name", a.opts.name,
		"version", a.opts.version,
	)
	g, ctx := errgroup.WithContext(a.ctx)
	for _, srv := range a.opts.servers {
		srv := srv
		g.Go(func() error {
			<-ctx.Done() // wait for stop signal
			return srv.Stop()
		})
		g.Go(func() error {
			return srv.Start()
		})
	}
	if err := a.waitForReady(ctx); err != nil {
		return err
	}
	if a.opts.registrar != nil {
		if err := a.opts.registrar.Register(a.opts.ctx, a.instance); err != nil {
			return err
		}
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				a.Stop()
			}
		}
	})
	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

// Stop gracefully stops the application.
func (a *App) Stop() error {
	if a.opts.registrar != nil {
		if err := a.opts.registrar.Deregister(a.opts.ctx, a.instance); err != nil {
			return err
		}
	}
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}

func (a *App) waitForReady(ctx context.Context) error {
retry:
	for _, srv := range a.opts.servers {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		e, err := srv.Endpoint()
		if err != nil {
			return err
		}
		if strings.HasSuffix(e, ":0") {
			time.Sleep(time.Millisecond * 100)
			goto retry
		}
	}
	a.instance = buildInstance(a.opts)
	return nil
}

func buildInstance(o options) *Service {
	if len(o.endpoints) == 0 {
		for _, srv := range o.servers {
			if e, err := srv.Endpoint(); err == nil {
				o.endpoints = append(o.endpoints, e)
			}
		}
	}
	return &Service{
		id:        o.id,
		name:      o.name,
		version:   o.version,
		metadata:  o.metadata,
		endpoints: o.endpoints,
	}
}
