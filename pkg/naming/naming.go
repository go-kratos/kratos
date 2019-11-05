package naming

import (
	"context"
)

// metadata common key
const (
	MetaWeight  = "weight"
	MetaCluster = "cluster"
	MetaZone    = "zone"
	MetaColor   = "color"
)

// Instance represents a server the client connects to.
type Instance struct {
	// Region is region.
	Region string `json:"region"`
	// Zone is IDC.
	Zone string `json:"zone"`
	// Env prod/pre、uat/fat1
	Env string `json:"env"`
	// AppID is mapping servicetree appid.
	AppID string `json:"appid"`
	// Hostname is hostname from docker.
	Hostname string `json:"hostname"`
	// Addrs is the address of app instance
	// format: scheme://host
	Addrs []string `json:"addrs"`
	// Version is publishing version.
	Version string `json:"version"`
	// LastTs is instance latest updated timestamp
	LastTs int64 `json:"latest_timestamp"`
	// Metadata is the information associated with Addr, which may be used
	// to make load balancing decision.
	Metadata map[string]string `json:"metadata"`
	// Status instance status, eg: 1UP 2Waiting
	Status int64 `json:"status"`
}

// Resolver resolve naming service
type Resolver interface {
	Fetch(context.Context) (*InstancesInfo, bool)
	Watch() <-chan struct{}
	Close() error
}

// Registry Register an instance and renew automatically.
type Registry interface {
	Register(ctx context.Context, ins *Instance) (cancel context.CancelFunc, err error)
	Close() error
}

// Builder resolver builder.
type Builder interface {
	Build(id string, options ...BuildOpt) Resolver
	Scheme() string
}

// InstancesInfo instance info.
type InstancesInfo struct {
	Instances map[string][]*Instance `json:"instances"`
	LastTs    int64                  `json:"latest_timestamp"`
	Scheduler *Scheduler             `json:"scheduler"`
}

// Scheduler scheduler.
type Scheduler struct {
	Clients map[string]*ZoneStrategy `json:"clients"`
}

// ZoneStrategy is the scheduling strategy of all zones
type ZoneStrategy struct {
	Zones map[string]*Strategy `json:"zones"`
}

// Strategy is zone scheduling strategy.
type Strategy struct {
	Weight int64 `json:"weight"`
}
