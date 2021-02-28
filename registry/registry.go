package registry

import "context"

// Registrar is service registrar.
type Registrar interface {
	// Register the registration.
	Register(ctx context.Context, service *ServiceInstance) error
	// Deregister the registration.
	Deregister(ctx context.Context, service *ServiceInstance) error
}

// Discoverer is service discovery.
type Discoverer interface {
	// Service return the service instances in memory according to the service name.
	Fetch(ctx context.Context, serviceName string) ([]*ServiceInstance, error)
	// Watch creates a watcher according to the service name.
	Watch(ctx context.Context, serviceName string) (Watcher, error)
}

// Watcher is service watcher.
type Watcher interface {
	// Watch returns services in the following two cases:
	// 1.the first time to watch and the service instance list is not empty.
	// 2.any service instance changes found.
	// if the above two conditions are not met, it will block until context deadline exceeded or canceled
	Next() ([]*ServiceInstance, error)
	// Close close the watcher.
	Close() error
}

// ServiceInstance is an instance of a service in a discovery system.
type ServiceInstance struct {
	// ID is the unique instance ID as registered.
	ID string
	// Name is the service name as registered.
	Name string
	// Version is the version of the compiled.
	Version string
	// Metadata is the kv pair metadata associated with the service instance.
	Metadata map[string]string
	// Endpoints is endpoint addresses of the service instance.
	// schema:
	//   http://127.0.0.1:8000?isSecure=false
	//   grpc://127.0.0.1:9000?isSecure=false
	Endpoints []string
}
