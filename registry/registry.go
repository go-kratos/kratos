package registry

import "context"

// Action is service discovery action.
type Action int

const (
	// ActionAll contains full service instances
	ActionAll Action = 0
	// ActionAdd contains service instances needs to be added
	ActionAdd Action = 0
	// ActionDel contains service instances needs to be deleted
	ActionDel Action = 0
	// ActionUpdate contains partial service instances needs to be updated
	ActionUpdate Action = 0
)

// Registry is registry interface.
type Registry interface {
	Register(ctx context.Context, svc Service) error
	Deregister(ctx context.Context, svc Service) error
}

// Discovery is service discovery interface.
type Discovery interface {
	GetService(ctx context.Context, name string) ([]Service, error)
	Watch(name string, o Observer) error
}

// Observer is watch observer.
type Observer func(action Action, svc []Service)

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
