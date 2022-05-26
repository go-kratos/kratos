package registry

import (
	"context"

	"github.com/go-kratos/kratos/v2/registry"
	"golang.org/x/sync/errgroup"
)

// Group Aggregate a set of registrars into a group.
func Group(r ...registry.Registrar) registry.Registrar {
	return &registrarGroup{registrars: r}
}

type registrarGroup struct {
	registrars []registry.Registrar
}

// Register the registration.
func (g *registrarGroup) Register(ctx context.Context, service *registry.ServiceInstance) error {
	eg := &errgroup.Group{}
	for _, reg := range g.registrars {
		r := reg
		eg.Go(func() error {
			return r.Register(ctx, service)
		})
	}
	return eg.Wait()
}

// Deregister the registration.
func (g *registrarGroup) Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	eg := &errgroup.Group{}
	for _, reg := range g.registrars {
		r := reg
		eg.Go(func() error {
			return r.Deregister(ctx, service)
		})
	}
	return eg.Wait()
}
