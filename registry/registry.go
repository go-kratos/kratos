package registry

import "context"

// Registry is registry interface.
type Registry interface {
	Register(ctx context.Context, svc Service) error
	Deregister(ctx context.Context, svc Service) error
	Services(ctx context.Context, name string) ([]Service, error)
	Watch(name string, o Observer) error
}

// Observer is watch observer.
type Observer func(action string, svc Service)

// Service is service interface.
type Service interface {
	ID() string
	Name() string
	Version() string
	Metadata() map[string]string
	Endpoints() []Endpoint
}

// Endpoint is endpoint interface.
type Endpoint interface {
	Scheme() string
	Host() string
	Port() int
	IsSecure() bool
}
