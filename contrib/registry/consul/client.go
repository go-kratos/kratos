package consul

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/hashicorp/consul/api"
)

// Client is consul client config
type Client struct {
	cli    *api.Client
	ctx    context.Context
	cancel context.CancelFunc
}

// NewClient creates consul client
func NewClient(cli *api.Client) *Client {
	c := &Client{cli: cli}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	return c
}

// Service get services from consul
func (d *Client) Service(ctx context.Context, service string, index uint64, passingOnly bool) ([]*registry.ServiceInstance, uint64, error) {
	opts := &api.QueryOptions{
		WaitIndex: index,
		WaitTime:  time.Second * 55,
	}
	opts = opts.WithContext(ctx)
	entries, meta, err := d.cli.Health().Service(service, "", passingOnly, opts)
	if err != nil {
		return nil, 0, err
	}

	services := make([]*registry.ServiceInstance, 0)

	for _, entry := range entries {
		var version string
		var scheme string
		for _, tag := range entry.Service.Tags {
			strs := strings.SplitN(tag, "=", 2)
			if len(strs) == 2 && strs[0] == "version" {
				version = strs[1]
			}
			if len(strs) == 2 && strs[0] == "scheme" {
				scheme = strs[1]
			}
		}
		endpoint := fmt.Sprintf("%s://%s:%d", scheme, entry.Service.Address, entry.Service.Port)

		services = append(services, &registry.ServiceInstance{
			ID:        entry.Service.ID,
			Name:      entry.Service.Service,
			Metadata:  entry.Service.Meta,
			Version:   version,
			Endpoints: []string{endpoint},
		})
	}
	return services, meta.LastIndex, nil
}

// Register register service instacen to consul
func (d *Client) Register(ctx context.Context, svc *registry.ServiceInstance, enableHealthCheck bool) error {
	// register grpc or http service with different srvId
	// and save scheme in tag.
	for _, endpoint := range svc.Endpoints {

		raw, err := url.Parse(endpoint)
		if err != nil {
			return err
		}
		addr := raw.Hostname()
		port, _ := strconv.ParseUint(raw.Port(), 10, 16)

		srvId := svc.ID + "-" + raw.Scheme
		asr := &api.AgentServiceRegistration{
			ID:      srvId,
			Name:    svc.Name,
			Meta:    svc.Metadata,
			Tags:    []string{fmt.Sprintf("version=%s", svc.Version), fmt.Sprintf("scheme=%s", raw.Scheme)},
			Address: addr,
			Port:    int(port),
		}

		if enableHealthCheck {
			asr.Checks = append(asr.Checks, &api.AgentServiceCheck{
				TCP:                            fmt.Sprintf("%s:%d", addr, port),
				Interval:                       "20s",
				DeregisterCriticalServiceAfter: "70s",
			})
		}

		err = d.cli.Agent().ServiceRegister(asr)
		if err != nil {
			return err
		}
		go func() {
			ticker := time.NewTicker(time.Second * 20)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					_ = d.cli.Agent().UpdateTTL("service:"+srvId, "pass", "pass")
				case <-d.ctx.Done():
					return
				}
			}
		}()
	}
	return nil
}

// Deregister deregister service by service ID
func (d *Client) Deregister(ctx context.Context, serviceID string) error {
	d.cancel()
	return d.cli.Agent().ServiceDeregister(serviceID)
}
