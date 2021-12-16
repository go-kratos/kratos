package consul

import (
	"context"
	"fmt"
	"net"
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

	// resolve service entry endpoints
	resolver Resolver
}

// NewClient creates consul client
func NewClient(cli *api.Client) *Client {
	c := &Client{cli: cli, resolver: defaultResolver}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	return c
}

func defaultResolver(_ context.Context, entry *api.ServiceEntry) []string {
	var endpoints []string
	for scheme, addr := range entry.Service.TaggedAddresses {
		if scheme == "lan_ipv4" || scheme == "wan_ipv4" || scheme == "lan_ipv6" || scheme == "wan_ipv6" {
			continue
		}
		endpoints = append(endpoints, addr.Address)
	}
	return endpoints
}

// Resolver is a wrapper for endpoints
type Resolver func(ctx context.Context, entry *api.ServiceEntry) []string

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
		for _, tag := range entry.Service.Tags {
			ss := strings.SplitN(tag, "=", 2)
			if len(ss) == 2 && ss[0] == "version" {
				version = ss[1]
			}
		}
		services = append(services, &registry.ServiceInstance{
			ID:        entry.Service.ID,
			Name:      entry.Service.Service,
			Metadata:  entry.Service.Meta,
			Version:   version,
			Endpoints: d.resolver(ctx, entry),
		})
	}
	return services, meta.LastIndex, nil
}

// Register register service instance to consul
func (d *Client) Register(_ context.Context, svc *registry.ServiceInstance, enableHealthCheck bool) error {
	addresses := make(map[string]api.ServiceAddress)
	checkAddresses := make([]string, 0, len(svc.Endpoints))
	for _, endpoint := range svc.Endpoints {
		raw, err := url.Parse(endpoint)
		if err != nil {
			return err
		}
		addr := raw.Hostname()
		port, _ := strconv.ParseUint(raw.Port(), 10, 16)
		checkAddresses = append(checkAddresses, fmt.Sprintf("%s:%d", addr, port))
		addresses[raw.Scheme] = api.ServiceAddress{Address: endpoint, Port: int(port)}
	}
	asr := &api.AgentServiceRegistration{
		ID:              svc.ID,
		Name:            svc.Name,
		Meta:            svc.Metadata,
		Tags:            []string{fmt.Sprintf("version=%s", svc.Version)},
		TaggedAddresses: addresses,
	}
	if len(checkAddresses) > 0 {
		host, portRaw, _ := net.SplitHostPort(checkAddresses[0])
		port, _ := strconv.ParseInt(portRaw, 10, 32)
		asr.Address = host
		asr.Port = int(port)
	}
	if enableHealthCheck {
		for _, address := range checkAddresses {
			asr.Checks = append(asr.Checks, &api.AgentServiceCheck{
				TCP:                            address,
				Interval:                       "20s",
				DeregisterCriticalServiceAfter: "70s",
			})
		}
	}
	err := d.cli.Agent().ServiceRegister(asr)
	if err != nil {
		return err
	}
	go func() {
		ticker := time.NewTicker(time.Second * 20)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				_ = d.cli.Agent().UpdateTTL("service:"+svc.ID, "pass", "pass")
			case <-d.ctx.Done():
				return
			}
		}
	}()
	return nil
}

// Deregister deregister service by service ID
func (d *Client) Deregister(ctx context.Context, serviceID string) error {
	d.cancel()
	return d.cli.Agent().ServiceDeregister(serviceID)
}
