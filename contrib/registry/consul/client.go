package consul

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

	"github.com/hashicorp/consul/api"
)

type ClusterMode string

const (
	Single        ClusterMode = "Single"
	Peering       ClusterMode = "Cluster Peering"
	WanFederation ClusterMode = "WAN Federation"
)

// Client is consul client config
type Client struct {
	consul *api.Client
	ctx    context.Context
	cancel context.CancelFunc

	// resolve service entry endpoints
	resolver ServiceResolver
	// healthcheck time interval in seconds
	healthcheckInterval int
	// heartbeat enable heartbeat
	heartbeat bool
	// deregisterCriticalServiceAfter time interval in seconds
	deregisterCriticalServiceAfter int
	// serviceChecks  user custom checks
	serviceChecks api.AgentServiceChecks

	// multiClusterMode is the consul cluster mode
	multiClusterMode ClusterMode

	// allow re-registration of services
	allowReRegistration bool

	// clusters specify the cluster to be used, if not set, obtain all currently associated clusters
	clusters []string

	lock                sync.RWMutex
	deregisteredService map[string]struct{}
}

func defaultResolver(_ context.Context, entries []*api.ServiceEntry) []*registry.ServiceInstance {
	services := make([]*registry.ServiceInstance, 0, len(entries))
	for _, entry := range entries {
		var version string
		for _, tag := range entry.Service.Tags {
			ss := strings.SplitN(tag, "=", 2)
			if len(ss) == 2 && ss[0] == "version" {
				version = ss[1]
			}
		}
		endpoints := make([]string, 0)
		for scheme, addr := range entry.Service.TaggedAddresses {
			if scheme == "lan_ipv4" || scheme == "wan_ipv4" || scheme == "lan_ipv6" || scheme == "wan_ipv6" {
				continue
			}
			endpoints = append(endpoints, addr.Address)
		}
		if len(endpoints) == 0 && entry.Service.Address != "" && entry.Service.Port != 0 {
			endpoints = append(endpoints, fmt.Sprintf("http://%s:%d", entry.Service.Address, entry.Service.Port))
		}
		services = append(services, &registry.ServiceInstance{
			ID:        entry.Service.ID,
			Name:      entry.Service.Service,
			Metadata:  entry.Service.Meta,
			Version:   version,
			Endpoints: endpoints,
		})
	}

	return services
}

// ServiceResolver is used to resolve service endpoints
type ServiceResolver func(ctx context.Context, entries []*api.ServiceEntry) []*registry.ServiceInstance

// service get services from consul
func (c *Client) service(ctx context.Context, service string, passingOnly bool, opts *api.QueryOptions) ([]*registry.ServiceInstance, uint64, error) {
	entries, meta, err := c.consul.Health().Service(service, "", passingOnly, opts)
	if err != nil {
		return nil, 0, err
	}

	for _, entry := range entries {
		if entry.Service.Meta == nil {
			entry.Service.Meta = make(map[string]string, 1)
		}
		if _, ok := entry.Service.Meta["cluster"]; !ok {
			entry.Service.Meta["cluster"] = entry.Node.Datacenter
		}
	}

	return c.resolver(ctx, entries), meta.LastIndex, nil
}

// Register register service instance to consul
func (c *Client) Register(ctx context.Context, svc *registry.ServiceInstance, enableHealthCheck bool) error {
	addresses := make(map[string]api.ServiceAddress, len(svc.Endpoints))
	checkAddresses := make([]string, 0, len(svc.Endpoints))
	for _, endpoint := range svc.Endpoints {
		raw, err := url.Parse(endpoint)
		if err != nil {
			return err
		}
		addr := raw.Hostname()
		port, _ := strconv.ParseUint(raw.Port(), 10, 16)

		checkAddresses = append(checkAddresses, net.JoinHostPort(addr, strconv.FormatUint(port, 10)))
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
				Interval:                       fmt.Sprintf("%ds", c.healthcheckInterval),
				DeregisterCriticalServiceAfter: fmt.Sprintf("%ds", c.deregisterCriticalServiceAfter),
				Timeout:                        "5s",
			})
		}
		// custom checks
		asr.Checks = append(asr.Checks, c.serviceChecks...)
	}
	if c.heartbeat {
		asr.Checks = append(asr.Checks, &api.AgentServiceCheck{
			CheckID:                        "service:" + svc.ID,
			TTL:                            fmt.Sprintf("%ds", c.healthcheckInterval*2),
			DeregisterCriticalServiceAfter: fmt.Sprintf("%ds", c.deregisterCriticalServiceAfter),
		})
	}

	err := c.consul.Agent().ServiceRegister(asr)
	if err != nil {
		return err
	}
	if c.heartbeat {
		go func() {
			time.Sleep(time.Second)
			err = c.consul.Agent().UpdateTTL("service:"+svc.ID, "pass", "pass")
			if err != nil {
				log.Errorf("[Consul]update ttl heartbeat to consul failed!err:=%v", err)
			}
			ticker := time.NewTicker(time.Second * time.Duration(c.healthcheckInterval))
			defer ticker.Stop()

			for range ticker.C {
				// ensure that unregistered services will not be re-registered by mistake
				c.lock.RLock()
				if _, ok := c.deregisteredService[svc.ID]; ok {
					c.lock.RUnlock()
					return
				}

				err = c.consul.Agent().UpdateTTL("service:"+svc.ID, "pass", "pass")
				if err != nil {
					log.Errorf("[Consul] update ttl heartbeat to consul failed! err=%v", err)
					if !c.allowReRegistration {
						c.lock.RUnlock()
						continue
					}
					// when the previous report fails, try to re register the service
					time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
					if err := c.consul.Agent().ServiceRegister(asr); err != nil {
						log.Errorf("[Consul] re registry service failed!, err=%v", err)
					} else {
						log.Warn("[Consul] re registry of service occurred success")
					}
					c.lock.RUnlock()
				}
			}
		}()
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	// the service may need to be re registered, so the record needs to be deleted
	delete(c.deregisteredService, svc.ID)
	return nil
}

// Deregister service by service ID
func (c *Client) Deregister(_ context.Context, serviceID string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.deregisteredService[serviceID] = struct{}{}
	return c.consul.Agent().ServiceDeregister(serviceID)
}
