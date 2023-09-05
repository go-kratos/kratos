package consul

import (
	"context"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/hashicorp/consul/api"
)

// A mix up of implmentations of registry.Registrar and registry.Discovery
// Still use this name for backward compatibility
// @deprecated use NewDiscovery and NewRegistrar instead
type Registry struct {
	*kratosRegistrar
	*kratosDiscovery
	runContext context.Context
}

// @deprecated
type Option interface{}

// @deprecated
type Datacenter string

// @deprecated
const (
	SingleDatacenter Datacenter = "SINGLE"
	MultiDatacenter  Datacenter = "MULTI"
)

// @deprecated useless now, remove this option in your code
func WithDatacenter(dc Datacenter) DiscoveryOption {
	return func(o *kratosDiscovery) {
	}
}

// New creates consul registry
// @deprecated
func New(apiClient *api.Client, opts ...Option) *Registry {
	r := &Registry{
		// old interface have no way to specify context, just use background
		runContext: context.Background(),
	}

	regOpts := []RegistrarOption{}
	disOpts := []DiscoveryOption{
		WithOldGetServiceBehavior(),
	}

	for _, o := range opts {
		switch v := o.(type) {
		case RegistrarOption:
			regOpts = append(regOpts, v)
		case DiscoveryOption:
			disOpts = append(disOpts, v)
		default:
			panic("invalid option")
		}
	}

	r.kratosDiscovery = NewDiscovery(r.runContext, apiClient, disOpts...).(*kratosDiscovery)
	r.kratosRegistrar = NewRegistrar(r.runContext, apiClient, regOpts...).(*kratosRegistrar)
	return r
}

func (r *Registry) ListServices() (allServices map[string][]*registry.ServiceInstance, err error) {
	return r.kratosDiscovery.ListServices(r.kratosDiscovery.ctx)
}
