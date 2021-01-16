package registry

import "context"

// Registry is registry interface.
type Registry interface {
	Register(ctx context.Context, svc Service) error
	Deregister(ctx context.Context, svc Service) error
}

// Discovery is service discovery interface.
type Discovery interface {
	Service(ctx context.Context, name string) ([]Service, error)
	Resolve(name string) Watcher
}

// Watcher is service watcher.
type Watcher interface {
	Watch(ctx context.Context) ([]Service, error)
	Close()
}

// Service is service interface.
type Service interface {
	ID() string
	Name() string
	Version() string
	Metadata() map[string]string
	Host() string
	Endpoints() []string
}
