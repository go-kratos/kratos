package naming

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"sort"

	"github.com/bilibili/kratos/pkg/conf/env"
	"github.com/bilibili/kratos/pkg/log"

	"github.com/dgryski/go-farm"
)

// BuildOptions build options.
type BuildOptions struct {
	Filter     func(map[string][]*Instance) map[string][]*Instance
	Subset     func([]*Instance, int) []*Instance
	SubsetSize int
	ClientZone string
	Scheduler  func(*InstancesInfo) []*Instance
}

// BuildOpt build option interface.
type BuildOpt interface {
	Apply(*BuildOptions)
}

type funcOpt struct {
	f func(*BuildOptions)
}

func (f *funcOpt) Apply(opt *BuildOptions) {
	f.f(opt)
}

// Filter filter option.
func Filter(schema string, clusters map[string]struct{}) BuildOpt {
	return &funcOpt{f: func(opt *BuildOptions) {
		opt.Filter = func(inss map[string][]*Instance) map[string][]*Instance {
			newInss := make(map[string][]*Instance)
			for zone := range inss {
				var instances []*Instance
				for _, ins := range inss[zone] {
					//如果r.clusters的长度大于0说明需要进行集群选择
					if len(clusters) > 0 {
						if _, ok := clusters[ins.Metadata[MetaCluster]]; !ok {
							continue
						}
					}
					var addr string
					for _, a := range ins.Addrs {
						u, err := url.Parse(a)
						if err == nil && u.Scheme == schema {
							addr = u.Host
						}
					}
					if addr == "" {
						fmt.Fprintf(os.Stderr, "resolver: app(%s,%s) no valid grpc address(%v) found!", ins.AppID, ins.Hostname, ins.Addrs)
						log.Warn("resolver: invalid rpc address(%s,%s,%v) found!", ins.AppID, ins.Hostname, ins.Addrs)
						continue
					}
					instances = append(instances, ins)
				}
				newInss[zone] = instances
			}
			return newInss
		}
	}}
}

func defulatSubset(inss []*Instance, size int) []*Instance {
	backends := inss
	if len(backends) <= int(size) {
		return backends
	}
	clientID := env.Hostname
	sort.Slice(backends, func(i, j int) bool {
		return backends[i].Hostname < backends[j].Hostname
	})
	count := len(backends) / size
	// hash得到ID
	id := farm.Fingerprint64([]byte(clientID))
	// 获得rand轮数
	round := int64(id / uint64(count))

	s := rand.NewSource(round)
	ra := rand.New(s)
	//  根据source洗牌
	ra.Shuffle(len(backends), func(i, j int) {
		backends[i], backends[j] = backends[j], backends[i]
	})
	start := (id % uint64(count)) * uint64(size)
	return backends[int(start) : int(start)+int(size)]
}

// Subset Subset option.
func Subset(defaultSize int) BuildOpt {
	return &funcOpt{f: func(opt *BuildOptions) {
		opt.SubsetSize = defaultSize
		opt.Subset = defulatSubset
	}}
}

// ScheduleNode ScheduleNode option.
func ScheduleNode(clientZone string) BuildOpt {
	return &funcOpt{f: func(opt *BuildOptions) {
		opt.ClientZone = clientZone
		opt.Scheduler = func(app *InstancesInfo) (instances []*Instance) {
			type Zone struct {
				inss   []*Instance
				weight int64
				name   string
				score  float64
			}
			var zones []*Zone

			if app.Scheduler != nil {
				si, err := json.Marshal(app.Scheduler)
				if err == nil {
					log.Info("schedule info: %s", string(si))
				}
				if strategy, ok := app.Scheduler.Clients[clientZone]; ok {
					var min *Zone
					for name, zone := range strategy.Zones {
						inss := app.Instances[name]
						if len(inss) == 0 {
							continue
						}
						z := &Zone{
							inss:   inss,
							weight: zone.Weight,
							name:   name,
							score:  float64(len(inss)) / float64(zone.Weight),
						}
						if min == nil || z.score < min.score {
							min = z
						}
						zones = append(zones, z)
					}
					if opt.SubsetSize != 0 && len(min.inss) > opt.SubsetSize {
						min.score = float64(opt.SubsetSize) / float64(min.weight)
					}
					for _, z := range zones {
						nums := int(min.score * float64(z.weight))
						if nums == 0 {
							nums = 1
						}
						if nums < len(z.inss) {
							if opt.Subset != nil {
								z.inss = opt.Subset(z.inss, nums)
							} else {
								z.inss = defulatSubset(z.inss, nums)
							}
						}
					}
				}
			}
			for _, zone := range zones {
				for _, ins := range zone.inss {
					instances = append(instances, ins)
				}
			}
			//如果没有拿到节点，则选择直接获取
			if len(instances) == 0 {
				instances = app.Instances[clientZone]
				if len(instances) == 0 {
					for _, value := range app.Instances {
						instances = append(instances, value...)
					}
				}
			}
			return
		}
	}}
}
