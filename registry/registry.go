package registry

import "context"

// Registry is registry interface.
type Registry interface {
	Register(Service) error
	Deregister(Service) error
	GetService(string) ([]Service, error)
	ListServices() ([]*Service, error)
	Watch(ctx context.Context, name string) (chan Event, error)
}

// Service is service interface.
type Service interface {
	ID() string
	Name() string
	Version() string
	Metadata() map[string]string
	Endpoints() []*Endpoint
}

// Endpoint is endpoint interface.
type Endpoint interface {
	Scheme() string
	Host() string
	Port() int
	IsSecure() bool
}

// Event is watch event.
type Event interface {
	Type() string
	Service() Service
}
