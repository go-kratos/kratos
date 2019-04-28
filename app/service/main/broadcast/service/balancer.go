package service

import (
	"fmt"
	"math"
	"sort"
	"sync"

	"go-common/app/service/main/broadcast/model"
	"go-common/library/log"
)

const (
	_minWeight = 1
	_maxWeight = 1 << 20
	_maxNodes  = 5
)

// haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// Distance function returns the distance (in meters) between two points of
//     a given longitude and latitude relatively accurately (using a spherical
//     approximation of the Earth) through the Haversin Distance Formula for
//     great arc distance on a sphere with accuracy for small distances
//
// point coordinates are supplied in degrees and converted into rad. in the func
//
// distance returned is METERS!!!!!!
// http://en.wikipedia.org/wiki/Haversine_formula
func distance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180
	r = 6378100 // Earth radius in METERS
	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)
	return 2 * r * math.Asin(math.Sqrt(h))
}

type weightedNode struct {
	region        string
	hostname      string
	addrs         []string
	addrsV6       []string
	fixedWeight   int64
	currentWeight int64
	currentConns  int64
	updated       int64
	lat           float64
	lng           float64
	overseas      bool
}

func (w *weightedNode) String() string {
	return fmt.Sprintf("region:%s fixedWeight:%d, currentWeight:%d, currentConns:%d", w.region, w.fixedWeight, w.currentWeight, w.currentConns)
}

func (w *weightedNode) chosen() {
	w.currentConns++
}

func (w *weightedNode) reset() {
	w.currentWeight = 0
}

func (w *weightedNode) calcuateWeightByGeo(totalWeight, totalConns, currentConns int64, totalDist, currentDist float64) {
	fixedWeight := float64(w.fixedWeight)
	if totalConns > 0 && w.lat != 0 && w.lng != 0 {
		var (
			connRatio, geoRatio float64
			weightRatio         = fixedWeight / float64(totalWeight)
		)
		if totalConns > 0 {
			connRatio = float64(currentConns) / float64(totalConns) * 0.4
		}
		if totalDist > 0 {
			geoRatio = float64(currentDist) / float64(totalDist) * 0.6
		}
		diff := weightRatio - connRatio - geoRatio
		multiple := diff * float64(float64(totalConns)+totalDist)
		floor := math.Floor(multiple)
		if floor-multiple >= -0.5 {
			w.currentWeight = int64(fixedWeight + floor)
		} else {
			w.currentWeight = int64(fixedWeight + math.Ceil(multiple))
		}
		if diff < 0 {
			// we always return the max from minWeight and calculated Current weight
			if _minWeight > w.currentWeight {
				w.currentWeight = _minWeight
			}
		} else {
			// we always return the min from maxWeight and calculated Current weight
			if _maxWeight < w.currentWeight {
				w.currentWeight = _maxWeight
			}
		}
	} else {
		w.reset()
	}
}

func (w *weightedNode) calcuateWeightByConns(totalWeight, totalConns, currentConns int64, gainWeight float64) {
	fixedWeight := float64(w.fixedWeight) * gainWeight
	totalWeight += int64(fixedWeight) - w.fixedWeight
	if totalConns > 0 {
		var (
			connRatio   float64
			weightRatio = fixedWeight / float64(totalWeight)
		)
		if totalConns > 0 {
			connRatio = float64(currentConns) / float64(totalConns)
		}
		diff := weightRatio - connRatio
		multiple := diff * float64(totalConns)
		floor := math.Floor(multiple)
		if floor-multiple >= -0.5 {
			w.currentWeight = int64(fixedWeight + floor)
		} else {
			w.currentWeight = int64(fixedWeight + math.Ceil(multiple))
		}
		if diff < 0 {
			// we always return the max from minWeight and calculated Current weight
			if _minWeight > w.currentWeight {
				w.currentWeight = _minWeight
			}
		} else {
			// we always return the min from maxWeight and calculated Current weight
			if _maxWeight < w.currentWeight {
				w.currentWeight = _maxWeight
			}
		}
	} else {
		w.reset()
	}
}

// LoadBalancer load balancer.
type LoadBalancer struct {
	totalConns  int64
	totalWeight int64
	nodes       map[string]*weightedNode
	nodesMutex  sync.Mutex
}

// NewLoadBalancer new a load balancer.
func NewLoadBalancer() *LoadBalancer {
	lb := &LoadBalancer{
		nodes: make(map[string]*weightedNode),
	}
	return lb
}

// Size return node size.
func (lb *LoadBalancer) Size() int {
	return len(lb.nodes)
}

func (lb *LoadBalancer) weightedByGeo(lat, lng float64) (nodes []*weightedNode) {
	var totalDist float64
	for _, n := range lb.nodes {
		if n.lat != 0 && n.lng != 0 {
			totalDist += distance(n.lat, n.lng, lat, lng)
		}
	}
	for _, n := range lb.nodes {
		if n.lat != 0 && n.lng != 0 {
			n.calcuateWeightByGeo(lb.totalWeight, lb.totalConns, n.currentConns, totalDist, distance(n.lat, n.lng, lat, lng))
			nodes = append(nodes, n)
		}
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].currentWeight > nodes[j].currentWeight
	})
	if len(nodes) > 0 {
		nodes[0].chosen()
		lb.totalConns++
	}
	return
}

func (lb *LoadBalancer) weightedByConns(region string, regionWeight float64) (nodes []*weightedNode) {
	for _, n := range lb.nodes {
		var gainWeight = 1.0
		if n.region == region {
			gainWeight *= regionWeight
		}
		n.calcuateWeightByConns(lb.totalWeight, lb.totalConns, n.currentConns, gainWeight)
		nodes = append(nodes, n)
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].currentWeight > nodes[j].currentWeight
	})
	if len(nodes) > 0 {
		nodes[0].chosen()
		lb.totalConns++
	}
	return
}

// NodeDetails return nodes.
func (lb *LoadBalancer) NodeDetails(region string, regionWeight float64) (res []map[string]interface{}) {
	lb.nodesMutex.Lock()
	nodes := lb.weightedByConns(region, regionWeight)
	res = make([]map[string]interface{}, 0, len(nodes))
	for _, n := range nodes {
		r := map[string]interface{}{
			"region":        n.region,
			"hostname":      n.hostname,
			"fixedWeight":   n.fixedWeight,
			"currentWeight": n.currentWeight,
			"currentConns":  n.currentConns,
			"addrs":         n.addrs,
			"updated":       n.updated,
		}
		res = append(res, r)
	}
	lb.nodesMutex.Unlock()
	return
}

// NodeAddrsByRegion return node addrs.
func (lb *LoadBalancer) NodeAddrsByRegion(region, domain string, regionWeight float64, overseas, ipv6 bool) (domains, addrs []string) {
	lb.nodesMutex.Lock()
	nodes := lb.weightedByConns(region, regionWeight)
	lb.nodesMutex.Unlock()
	for _, n := range nodes {
		if n.overseas && !overseas {
			continue
		}
		if !ipv6 && len(n.addrs) > 0 {
			// ipv4
			domains = append(domains, n.hostname+domain)
			addrs = append(addrs, n.addrs...)
		} else if ipv6 && len(n.addrsV6) > 0 {
			// ipv6
			domains = append(domains, n.hostname+domain)
			addrs = append(addrs, n.addrsV6...)
		}
		if len(addrs) == _maxNodes {
			break
		}
	}
	return
}

// NodeAddrsByGeo return node addrs by geo.
func (lb *LoadBalancer) NodeAddrsByGeo(domain string, lat, lng float64, overseas, ipv6 bool) (domains, addrs []string) {
	lb.nodesMutex.Lock()
	nodes := lb.weightedByGeo(lat, lng)
	lb.nodesMutex.Unlock()
	for _, n := range nodes {
		if n.overseas && !overseas {
			continue
		}
		if !ipv6 && len(n.addrs) > 0 {
			// ipv4
			domains = append(domains, n.hostname+domain)
			addrs = append(addrs, n.addrs...)
		} else if ipv6 && len(n.addrsV6) > 0 {
			// ipv6
			domains = append(domains, n.hostname+domain)
			addrs = append(addrs, n.addrsV6...)
		}
		if len(addrs) == _maxNodes {
			break
		}
	}
	return
}

// Update update server nodes.
func (lb *LoadBalancer) Update(srvs []*model.ServerInfo) {
	var (
		totalConns  int64
		totalWeight int64
		nodes       = make(map[string]*weightedNode, len(srvs))
	)
	if len(srvs) == 0 || float32(len(srvs))/float32(len(lb.nodes)) < 0.5 {
		log.Error("load balancer update src:%d target:%d less than half", len(lb.nodes), len(srvs))
		return
	}
	lb.nodesMutex.Lock()
	for _, s := range srvs {
		if old, ok := lb.nodes[s.Server]; ok && old.updated == s.Updated {
			nodes[s.Server] = old
			totalConns += old.currentConns
			totalWeight += old.fixedWeight
		} else {
			node := &weightedNode{
				region:       s.Region,
				hostname:     s.Server,
				fixedWeight:  int64(s.Weight),
				currentConns: int64(s.ConnCount),
				addrs:        s.IPAddrs,
				addrsV6:      s.IPAddrsV6,
				updated:      s.Updated,
				lat:          s.Latitude,
				lng:          s.Longitude,
				overseas:     s.Overseas,
			}
			nodes[s.Server] = node
			totalConns += int64(s.ConnCount)
			totalWeight += int64(s.Weight)
		}
	}
	lb.nodes = nodes
	lb.totalConns = totalConns
	lb.totalWeight = totalWeight
	lb.nodesMutex.Unlock()
}
