package discovery

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/SeeMusic/kratos/v2/registry"
)

var (
	ErrDuplication = errors.New("register failed: instance duplicated: ")
	ErrServerError = errors.New("server error")
)

const (
	// Discovery server resource uri
	_registerURL = "http://%s/discovery/register"
	//_setURL      = "http://%s/discovery/set"
	_cancelURL = "http://%s/discovery/cancel"
	_renewURL  = "http://%s/discovery/renew"
	_pollURL   = "http://%s/discovery/polls"

	// Discovery server error codes
	_codeOK          = 0
	_codeNotFound    = -404
	_codeNotModified = -304
	//_SERVER_ERROR = -500

	// _registerGap is the gap to renew instance registration.
	_registerGap    = 30 * time.Second
	_statusUP       = "1"
	_discoveryAppID = "infra.discovery"
)

// Config Discovery configures.
type Config struct {
	Nodes  []string
	Region string
	Zone   string
	Env    string
	Host   string
}

func fixConfig(c *Config) error {
	if c.Host == "" {
		c.Host, _ = os.Hostname()
	}
	if len(c.Nodes) == 0 || c.Region == "" || c.Zone == "" || c.Env == "" || c.Host == "" {
		return fmt.Errorf(
			"invalid Discovery config nodes:%+v region:%s zone:%s deployEnv:%s host:%s",
			c.Nodes,
			c.Region,
			c.Zone,
			c.Env,
			c.Host,
		)
	}
	return nil
}

// discoveryInstance represents a server the client connects to.
type discoveryInstance struct {
	Region   string   `json:"region"`           // Region is region.
	Zone     string   `json:"zone"`             // Zone is IDC.
	Env      string   `json:"env"`              // Env prod/pre/uat/fat1
	AppID    string   `json:"appid"`            // AppID is mapping service-tree appId.
	Hostname string   `json:"hostname"`         // Hostname is hostname from docker
	Addrs    []string `json:"addrs"`            // Addrs is the address of app instance format: scheme://host
	Version  string   `json:"version"`          // Version is publishing version.
	LastTs   int64    `json:"latest_timestamp"` // LastTs is instance latest updated timestamp
	// Metadata is the information associated with Addr, which may be used to make load balancing decision.
	Metadata map[string]string `json:"metadata"`
	Status   int64             `json:"status"` // Status instance status, eg: 1UP 2Waiting
}

const _reservedInstanceIDKey = "kratos.v2.serviceinstance.id"

// fromServerInstance convert registry.ServiceInstance into discoveryInstance
func fromServerInstance(ins *registry.ServiceInstance, config *Config) *discoveryInstance {
	if ins == nil {
		return nil
	}

	metadata := ins.Metadata
	if ins.Metadata == nil {
		metadata = make(map[string]string, 8)
	}
	metadata[_reservedInstanceIDKey] = ins.ID

	return &discoveryInstance{
		Region:   config.Region,
		Zone:     config.Zone,
		Env:      config.Env,
		AppID:    ins.Name,
		Hostname: config.Host,
		Addrs:    ins.Endpoints,
		Version:  ins.Version,
		LastTs:   time.Now().Unix(),
		Metadata: metadata,
		Status:   1,
	}
}

// toServiceInstance convert discoveryInstance into registry.ServiceInstance
func toServiceInstance(ins *discoveryInstance) *registry.ServiceInstance {
	if ins == nil {
		return nil
	}

	md := map[string]string{
		"region":   ins.Region,
		"zone":     ins.Zone,
		"lastTs":   strconv.Itoa(int(ins.LastTs)),
		"env":      ins.Env,
		"hostname": ins.Hostname,
	}

	if len(ins.Metadata) != 0 {
		for k, v := range ins.Metadata {
			md[k] = v
		}
	}

	return &registry.ServiceInstance{
		ID:        ins.Metadata[_reservedInstanceIDKey],
		Name:      ins.AppID,
		Version:   ins.Version,
		Metadata:  md,
		Endpoints: ins.Addrs,
	}
}

// disInstancesInfo instance info.
type disInstancesInfo struct {
	Instances map[string][]*discoveryInstance `json:"instances"`
	LastTs    int64                           `json:"latest_timestamp"`
	Scheduler *scheduler                      `json:"scheduler"`
}

// scheduler scheduler.
type scheduler struct {
	Clients map[string]*zoneStrategy `json:"clients"`
}

// zoneStrategy is the scheduling strategy of all zones
type zoneStrategy struct {
	Zones map[string]*strategy `json:"zones"`
}

// strategy is zone scheduling strategy.
type strategy struct {
	Weight int64 `json:"weight"`
}

const (
	_paramKeyRegion   = "region"
	_paramKeyZone     = "zone"
	_paramKeyEnv      = "env"
	_paramKeyHostname = "hostname"
	_paramKeyAppID    = "appid"
	_paramKeyAddrs    = "addrs"
	_paramKeyVersion  = "version"
	_paramKeyStatus   = "status"
	_paramKeyMetadata = "metadata"
)

func newParams(c *Config) url.Values {
	p := make(url.Values, 8)
	if c == nil {
		return p
	}

	p.Set(_paramKeyRegion, c.Region)
	p.Set(_paramKeyZone, c.Zone)
	p.Set(_paramKeyEnv, c.Env)
	p.Set(_paramKeyHostname, c.Host)
	return p
}

type discoveryCommonResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type discoveryPollsResp struct {
	Code int                          `json:"code"`
	Data map[string]*disInstancesInfo `json:"data"`
}
