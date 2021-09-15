package discovery

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/go-kratos/kratos/v2/registry"
)

var (
	ErrDuplication = errors.New("register failed: instance duplicated: ")
	ErrServerError = errors.New("server error")
)

const (
	// discovery server resource uri
	_registerURL = "http://%s/discovery/register"
	_setURL      = "http://%s/discovery/set"
	_cancelURL   = "http://%s/discovery/cancel"
	_renewURL    = "http://%s/discovery/renew"
	_pollURL     = "http://%s/discovery/polls"

	// Discovery server error codes
	_OK           = 0
	_NOT_FOUND    = -404
	_NOT_MODIFIED = -304
	_SERVER_ERROR = -500

	// _registerGap is the gap to renew instance registration.
	_registerGap    = 30 * time.Second
	_statusUP       = "1"
	_discoveryAppID = "infra.discovery"
)

// deploy env.
const (
	DeployEnvDev  = "dev"
	DeployEnvFat  = "fat"
	DeployEnvUat  = "uat"
	DeployEnvPre  = "pre"
	DeployEnvProd = "prod"
)

// env default value.
const (
	// env
	_region    = "region01"
	_zone      = "zone01"
	_deployEnv = "dev"
)

// env configuration.
var (
	// EnvRegion available region where app at.
	EnvRegion string
	// EnvZone available zone where app at.
	EnvZone string
	// EnvHostname machine hostname.
	EnvHostname string
	// EnvDeployEnv deploy env where app at.
	EnvDeployEnv string
	// EnvAppID is global unique application id, register by service tree.
	// such as main.arch.disocvery.
	EnvAppID string
	// EnvColor is the identification of different experimental group in one caster cluster.
	EnvColor string
	// EnvDiscoveryNodes is seed nodes.
	EnvDiscoveryNodes string
)

func init() {
	var err error
	EnvHostname = os.Getenv("HOSTNAME")
	if EnvHostname == "" {
		EnvHostname, err = os.Hostname()
		if err != nil {
			EnvHostname = strconv.Itoa(int(time.Now().UnixNano()))
		}
	}
	addFlag(flag.CommandLine)
}

func addFlag(fs *flag.FlagSet) {
	// env
	fs.StringVar(&EnvRegion, "region", defaultString("REGION", _region), "available region. or use REGION env variable, value: sh etc.")
	fs.StringVar(&EnvZone, "zone", defaultString("ZONE", _zone), "available zone. or use ZONE env variable, value: sh001/sh002 etc.")
	fs.StringVar(&EnvAppID, "appid", os.Getenv("APP_ID"), "appid is global unique application id, register by service tree. or use APP_ID env variable.")
	fs.StringVar(&EnvDeployEnv, "deploy.env", defaultString("DEPLOY_ENV", _deployEnv), "deploy env. or use DEPLOY_ENV env variable, value: dev/fat1/uat/pre/prod etc.")
	fs.StringVar(&EnvColor, "deploy.color", os.Getenv("DEPLOY_COLOR"), "deploy.color is the identification of different experimental group.")
	fs.StringVar(&EnvDiscoveryNodes, "discovery.nodes", os.Getenv("DISCOVERY_NODES"), "discovery.nodes is seed nodes. value: 127.0.0.1:7171,127.0.0.2:7171 etc.")
}

func defaultString(env, value string) string {
	v := os.Getenv(env)
	if v == "" {
		return value
	}
	return v
}

// Config discovery configures.
type Config struct {
	Nodes  []string
	Region string
	Zone   string
	Env    string
	Host   string
}

func fixConfig(c *Config) error {
	if len(c.Nodes) == 0 && EnvDiscoveryNodes != "" {
		c.Nodes = strings.Split(EnvDiscoveryNodes, ",")
	}
	if c.Region == "" {
		c.Region = EnvRegion
	}
	if c.Zone == "" {
		c.Zone = EnvZone
	}
	if c.Env == "" {
		c.Env = EnvDeployEnv
	}
	if c.Host == "" {
		c.Host = EnvHostname
	}
	if len(c.Nodes) == 0 || c.Region == "" || c.Zone == "" || c.Env == "" || c.Host == "" {
		return fmt.Errorf(
			"invalid discovery config nodes:%+v region:%s zone:%s deployEnv:%s host:%s",
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

func fromServerInstance(ins *registry.ServiceInstance) *discoveryInstance {
	if ins == nil {
		return nil
	}

	var (
		region   string
		zone     string
		env      string
		hostname string
		lastTs   int64
	)

	var metadata = ins.Metadata
	if ins.Metadata == nil {
		metadata = make(map[string]string, 8)
	}
	if v := metadata["region"]; v != "" {
		region = v
	}
	if v := metadata["zone"]; v != "" {
		zone = v
	}
	if v := metadata["env"]; v != "" {
		env = v
	}
	if v := metadata["hostname"]; v != "" {
		hostname = v
	}
	if v := metadata["lastTs"]; v != "" {
		lastTs, _ = strconv.ParseInt(v, 10, 64)
	}
	metadata["reserved.id"] = ins.ID

	return &discoveryInstance{
		Region:   region,
		Zone:     zone,
		Env:      env,
		AppID:    ins.Name,
		Hostname: hostname,
		Addrs:    ins.Endpoints,
		Version:  ins.Version,
		LastTs:   lastTs,
		Metadata: metadata,
		Status:   0, // TODO(@yeqown)
	}
}

func toServiceInstance(ins *discoveryInstance) *registry.ServiceInstance {
	if ins == nil {
		return nil
	}

	return &registry.ServiceInstance{
		ID:      ins.Metadata["reserved.id"],
		Name:    ins.AppID,
		Version: ins.Version,
		Metadata: map[string]string{
			"region":   ins.Region,
			"zone":     ins.Region,
			"lastTs":   strconv.Itoa(int(ins.LastTs)),
			"env":      ins.Env,
			"hostname": ins.Hostname,
		},
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

func newParams(c *Config) url.Values {
	p := make(url.Values, 8)
	if c == nil {
		return p
	}

	p.Set("region", c.Region)
	p.Set("zone", c.Zone)
	p.Set("env", c.Env)
	p.Set("hostname", c.Host)
	return p
}
