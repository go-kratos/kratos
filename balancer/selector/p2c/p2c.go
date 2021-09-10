package p2c

import (
	"context"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/go-kratos/kratos/v2/balancer"
)

const (
	forcePick = time.Second * 3
	// Name is balancer name
	Name = "p2c"
)

var _ balancer.Selector = &Selector{}

// New p2c
func New() balancer.Selector {
	p := &Selector{
		r: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	return p
}

// Selector is p2c selector
type Selector struct {
	r  *rand.Rand
	lk int64
}

// choose two distinct nodes
func (s *Selector) prePick(nodes []balancer.Node) (nodeA balancer.Node, nodeB balancer.Node) {
	a := s.r.Intn(len(nodes))
	b := s.r.Intn(len(nodes) - 1)
	if b >= a {
		b = b + 1
	}
	nodeA, nodeB = nodes[a], nodes[b]
	return
}

func (s *Selector) Select(ctx context.Context, nodes []balancer.Node) (balancer.Node, balancer.Done, error) {
	if len(nodes) == 0 {
		return nil, nil, balancer.ErrNoAvailable
	} else if len(nodes) == 1 {
		done := nodes[0].Pick()
		return nodes[0], done, nil
	}

	var pc, upc balancer.Node
	var done func(ctx context.Context, di balancer.DoneInfo)
	nodeA, nodeB := s.prePick(nodes)
	// meta.Weight为服务发布者在discovery中设置的权重
	if nodeB.Weight() > nodeA.Weight() {
		pc, upc = nodeB, nodeA
	} else {
		pc, upc = nodeA, nodeB
	}
	// 如果选中的节点，在forceGap期间内从来没有被选中一次，则强制选一次
	// 利用强制的机会，来触发成功率、延迟的更新
	if upc.PickElapsed() > forcePick {
		if atomic.CompareAndSwapInt64(&s.lk, 0, 1) {
			pc = upc
			done = pc.Pick()
			atomic.StoreInt64(&s.lk, 0)
		} else {
			done = pc.Pick()
		}
	} else {
		done = pc.Pick()
	}

	return pc, done, nil
}

/*

func (p *P2cPicker) PrintStats() {
	if len(p.subConns) == 0 {
		return
	}
	stats := make([]statistic, 0, len(p.subConns))
	var serverName string
	var reqs int64
	var now = time.Now().UnixNano()
	for _, conn := range p.subConns {
		var stat statistic
		stat.addr = conn.node.Endpoints[0]
		stat.cs = atomic.LoadUint64(&conn.success)
		stat.inflight = atomic.LoadInt64(&conn.inflight)
		stat.lantency = time.Duration(atomic.LoadInt64(&conn.lag))
		stat.reqs = atomic.SwapInt64(&conn.reqs, 0)
		stat.load = conn.load(now)
		stat.predict = time.Duration(atomic.LoadInt64(&conn.predict))
		stats = append(stats, stat)
		if serverName == "" {
			serverName = conn.node.Name
		}
		reqs += stat.reqs
	}
	if reqs > 10 {
		//log.DefaultLog.Debugf("p2c %s : %+v", serverName, stats)
	}
}

// statistics is info for log
type statistic struct {
	addr     string
	score    float64
	cs       uint64
	lantency time.Duration
	load     uint64
	inflight int64
	reqs     int64
	predict  time.Duration
}
*/
