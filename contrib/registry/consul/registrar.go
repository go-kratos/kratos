package consul

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/consul/api"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
)

// Option is consul registry option.
type RegistrarOption func(*kratosRegistrar)

// WithHealthCheck with registry health check option.
func WithHealthCheck(enable bool) RegistrarOption {
	return func(o *kratosRegistrar) {
		o.enableHealthCheck = enable
	}
}

// WithHeartbeat enable or disable heartbeat
func WithHeartbeat(enable bool) RegistrarOption {
	return func(o *kratosRegistrar) {
		o.heartbeat = enable
	}
}

// WithHealthCheckInterval with healthcheck interval in seconds.
func WithHealthCheckInterval(interval int) RegistrarOption {
	return func(o *kratosRegistrar) {
		o.healthcheckInterval = interval
	}
}

// WithDeregisterCriticalServiceAfter with deregister-critical-service-after in seconds.
func WithDeregisterCriticalServiceAfter(interval int) RegistrarOption {
	return func(o *kratosRegistrar) {
		o.deregisterCriticalServiceAfter = interval
	}
}

// WithServiceCheck with service checks
func WithServiceCheck(checks ...*api.AgentServiceCheck) RegistrarOption {
	return func(o *kratosRegistrar) {
		o.serviceChecks = checks
	}
}

type kratosRegistrar struct {
	// native consul client
	cli *api.Client

	// wheather enable health check
	enableHealthCheck bool

	// healthcheck time interval in seconds
	healthcheckInterval int

	// heartbeat enable heartbeat
	heartbeat bool

	// deregisterCriticalServiceAfter time interval in seconds
	deregisterCriticalServiceAfter int

	// serviceChecks  user custom checks
	serviceChecks api.AgentServiceChecks

	heartbeatDone chan bool

	// context for subroutines started by Registrar, if canceld, all subroutines will be closed
	runContext context.Context
}

func NewRegistrar(ctx context.Context, apiClient *api.Client, opts ...RegistrarOption) registry.Registrar {
	r := &kratosRegistrar{
		cli:                            apiClient,
		enableHealthCheck:              true,
		healthcheckInterval:            10,
		heartbeat:                      true,
		deregisterCriticalServiceAfter: 600,
		heartbeatDone:                  make(chan bool),
		runContext:                     ctx,
	}

	for _, o := range opts {
		o(r)
	}

	return r
}

func (r *kratosRegistrar) Register(_ context.Context, svc *registry.ServiceInstance) error {
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

	if svc.Version != "" {
		if svc.Metadata == nil {
			svc.Metadata = map[string]string{"version": svc.Version}
		} else {
			svc.Metadata["version"] = svc.Version
		}
	}

	asr := &api.AgentServiceRegistration{
		ID:              svc.ID,
		Name:            svc.Name,
		Meta:            svc.Metadata,
		TaggedAddresses: addresses,
	}
	if len(checkAddresses) > 0 {
		host, portRaw, _ := net.SplitHostPort(checkAddresses[0])
		port, _ := strconv.ParseInt(portRaw, 10, 32)
		asr.Address = host
		asr.Port = int(port)
	}
	if r.enableHealthCheck {
		for _, address := range checkAddresses {
			asr.Checks = append(asr.Checks, &api.AgentServiceCheck{
				Name:                           "TCP Connectivity to " + address,
				TCP:                            address,
				Interval:                       fmt.Sprintf("%ds", r.healthcheckInterval),
				DeregisterCriticalServiceAfter: fmt.Sprintf("%ds", r.deregisterCriticalServiceAfter),
				Timeout:                        "5s",
			})
		}
	}
	if r.heartbeat {
		asr.Checks = append(asr.Checks, &api.AgentServiceCheck{
			Name:                           "Heartbeat",
			CheckID:                        "heartbeat:" + svc.ID,
			Status:                         "passing",
			TTL:                            fmt.Sprintf("%ds", r.healthcheckInterval*2),
			DeregisterCriticalServiceAfter: fmt.Sprintf("%ds", r.deregisterCriticalServiceAfter),
		})
	}

	// custom checks if any
	asr.Checks = append(asr.Checks, r.serviceChecks...)

	err := r.cli.Agent().ServiceRegister(asr)
	if err != nil {
		return err
	}
	if r.heartbeat {
		go func() {
			ticker := time.NewTicker(time.Second * time.Duration(r.healthcheckInterval))
			defer ticker.Stop()
			for {
				select {
				case <-r.runContext.Done():
					// run context is canceled, exit now
					return
				case <-r.heartbeatDone:
					return
				case <-ticker.C:
					err = r.cli.Agent().UpdateTTL("heartbeat:"+svc.ID, "Service alive", "pass")
					if err != nil {
						log.Errorf("[Consul] update ttl heartbeat to consul failed! err=%v", err)
						// when the previous report fails, try to re register the service
						time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
						if err := r.cli.Agent().ServiceRegister(asr); err != nil {
							log.Errorf("[Consul] re registry service failed!, err=%v", err)
						} else {
							log.Warn("[Consul] re registry of service occurred success")
						}
					}
				}
			}
		}()
	}
	return nil
}

func (r *kratosRegistrar) Deregister(_ context.Context, svc *registry.ServiceInstance) error {
	r.heartbeatDone <- true
	return r.cli.Agent().ServiceDeregister(svc.ID)
}
