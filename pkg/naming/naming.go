package naming

import (
	"context"
	"strconv"
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
	Build(id string) Resolver
	Scheme() string
}

// InstancesInfo instance info.
type InstancesInfo struct {
	Instances map[string][]*Instance `json:"instances"`
	LastTs    int64                  `json:"latest_timestamp"`
	Scheduler []Zone                 `json:"scheduler"`
}

// Zone zone scheduler info.
type Zone struct {
	Src string           `json:"src"`
	Dst map[string]int64 `json:"dst"`
}

// UseScheduler use scheduler info on instances.
// if instancesInfo contains scheduler info about zone,
// return releated zone's instances weighted by scheduler.
// if not,only zone instances be returned.
func (insInf *InstancesInfo) UseScheduler(zone string) (inss []*Instance) {
	var scheduler struct {
		zone    []string
		weights []int64
	}
	var oriWeights []int64
	for _, sch := range insInf.Scheduler {
		if sch.Src == zone {
			for zone, schWeight := range sch.Dst {
				if zins, ok := insInf.Instances[zone]; ok {
					var totalWeight int64
					for _, ins := range zins {
						var weight int64
						if weight, _ = strconv.ParseInt(ins.Metadata[MetaWeight], 10, 64); weight <= 0 {
							weight = 10
						}
						totalWeight += weight
					}
					oriWeights = append(oriWeights, totalWeight)
					inss = append(inss, zins...)
				}
				scheduler.weights = append(scheduler.weights, schWeight)
				scheduler.zone = append(scheduler.zone, zone)
			}
		}
	}
	if len(inss) == 0 {
		var ok bool
		if inss, ok = insInf.Instances[zone]; ok {
			return
		}
		for _, v := range insInf.Instances {
			inss = append(inss, v...)
		}
		return
	}
	var comMulti int64 = 1
	for _, weigth := range oriWeights {
		comMulti *= weigth
	}
	var fixWeight = make(map[string]int64, len(scheduler.weights))
	for i, zone := range scheduler.zone {
		fixWeight[zone] = scheduler.weights[i] * comMulti / oriWeights[i]
	}
	for _, ins := range inss {
		var weight int64
		if weight, _ = strconv.ParseInt(ins.Metadata[MetaWeight], 10, 64); weight <= 0 {
			weight = 10
		}
		if fix, ok := fixWeight[ins.Zone]; ok {
			weight = weight * fix
		}
		ins.Metadata[MetaWeight] = strconv.FormatInt(weight, 10)
	}
	return
}
