package p2c

import (
	"context"
	"math"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bilibili/kratos/pkg/conf/env"

	"github.com/bilibili/kratos/pkg/log"
	nmd "github.com/bilibili/kratos/pkg/net/metadata"
	wmd "github.com/bilibili/kratos/pkg/net/rpc/warden/internal/metadata"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
)

const (
	// The mean lifetime of `cost`, it reaches its half-life after Tau*ln(2).
	tau = int64(time.Millisecond * 600)
	// if statistic not collected,we add a big penalty to endpoint
	penalty = uint64(1000 * time.Millisecond * 250)

	forceGap = int64(time.Second * 3)
)

var _ base.PickerBuilder = &p2cPickerBuilder{}
var _ balancer.Picker = &p2cPicker{}

// Name is the name of pick of two random choices balancer.
const Name = "p2c"

// newBuilder creates a new weighted-roundrobin balancer builder.
func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, &p2cPickerBuilder{})
}

func init() {
	balancer.Register(newBuilder())
}

type subConn struct {
	// metadata
	conn balancer.SubConn
	addr resolver.Address
	meta wmd.MD

	//client statistic data
	lag      uint64
	success  uint64
	inflight int64
	// server statistic data
	svrCPU uint64

	//last collected timestamp
	stamp int64
	//last pick timestamp
	pick int64
	// request number in a period time
	reqs int64
}

func (sc *subConn) valid() bool {
	return sc.health() > 500 && atomic.LoadUint64(&sc.svrCPU) < 900
}

func (sc *subConn) health() uint64 {
	return atomic.LoadUint64(&sc.success)
}

func (sc *subConn) load() uint64 {
	lag := uint64(math.Sqrt(float64(atomic.LoadUint64(&sc.lag))) + 1)
	load := atomic.LoadUint64(&sc.svrCPU) * lag * uint64(atomic.LoadInt64(&sc.inflight))
	if load == 0 {
		// penalty是初始化没有数据时的惩罚值，默认为1e9 * 250
		load = penalty
	}
	return load
}

func (sc *subConn) cost() uint64 {
	load := atomic.LoadUint64(&sc.svrCPU) * atomic.LoadUint64(&sc.lag) * uint64(atomic.LoadInt64(&sc.inflight))
	if load == 0 {
		// penalty是初始化没有数据时的惩罚值，默认为1e9 * 250
		load = penalty
	}
	return load
}

// statistics is info for log
type statistic struct {
	addr     string
	score    float64
	cs       uint64
	lantency uint64
	cpu      uint64
	inflight int64
	reqs     int64
}

type p2cPickerBuilder struct{}

func (*p2cPickerBuilder) Build(readySCs map[resolver.Address]balancer.SubConn) balancer.Picker {
	p := &p2cPicker{
		colors: make(map[string]*p2cPicker),
		r:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	for addr, sc := range readySCs {
		meta, ok := addr.Metadata.(wmd.MD)
		if !ok {
			meta = wmd.MD{
				Weight: 10,
			}
		}
		subc := &subConn{
			conn: sc,
			addr: addr,
			meta: meta,

			svrCPU:   500,
			lag:      0,
			success:  1000,
			inflight: 1,
		}
		if meta.Color == "" {
			p.subConns = append(p.subConns, subc)
			continue
		}
		// if color not empty, use color picker
		cp, ok := p.colors[meta.Color]
		if !ok {
			cp = &p2cPicker{r: rand.New(rand.NewSource(time.Now().UnixNano()))}
			p.colors[meta.Color] = cp
		}
		cp.subConns = append(cp.subConns, subc)
	}
	return p
}

type p2cPicker struct {
	// subConns is the snapshot of the weighted-roundrobin balancer when this picker was
	// created. The slice is immutable. Each Get() will do a round robin
	// selection from it and return the selected SubConn.
	subConns []*subConn
	colors   map[string]*p2cPicker
	logTs    int64
	r        *rand.Rand
	lk       sync.Mutex
}

func (p *p2cPicker) Pick(ctx context.Context, opts balancer.PickOptions) (balancer.SubConn, func(balancer.DoneInfo), error) {
	// FIXME refactor to unify the color logic
	color := nmd.String(ctx, nmd.Color)
	if color == "" && env.Color != "" {
		color = env.Color
	}
	if color != "" {
		if cp, ok := p.colors[color]; ok {
			return cp.pick(ctx, opts)
		}
	}
	return p.pick(ctx, opts)
}

// choose two distinct nodes
func (p *p2cPicker) prePick() (nodeA *subConn, nodeB *subConn) {
	for i := 0; i < 3; i++ {
		p.lk.Lock()
		a := p.r.Intn(len(p.subConns))
		b := p.r.Intn(len(p.subConns) - 1)
		p.lk.Unlock()
		if b >= a {
			b = b + 1
		}
		nodeA, nodeB = p.subConns[a], p.subConns[b]
		if nodeA.valid() || nodeB.valid() {
			break
		}
	}
	return
}

func (p *p2cPicker) pick(ctx context.Context, opts balancer.PickOptions) (balancer.SubConn, func(balancer.DoneInfo), error) {
	var pc, upc *subConn
	start := time.Now().UnixNano()

	if len(p.subConns) <= 0 {
		return nil, nil, balancer.ErrNoSubConnAvailable
	} else if len(p.subConns) == 1 {
		pc = p.subConns[0]
	} else {
		nodeA, nodeB := p.prePick()
		// meta.Weight为服务发布者在disocvery中设置的权重
		if nodeA.load()*nodeB.health()*nodeB.meta.Weight > nodeB.load()*nodeA.health()*nodeA.meta.Weight {
			pc, upc = nodeB, nodeA
		} else {
			pc, upc = nodeA, nodeB
		}
		// 如果选中的节点，在forceGap期间内没有被选中一次，那么强制一次
		// 利用强制的机会，来触发成功率、延迟的衰减
		// 原子锁conn.pick保证并发安全，放行一次
		pick := atomic.LoadInt64(&upc.pick)
		if start-pick > forceGap && atomic.CompareAndSwapInt64(&upc.pick, pick, start) {
			pc = upc
		}
	}

	// 节点未发生切换才更新pick时间
	if pc != upc {
		atomic.StoreInt64(&pc.pick, start)
	}
	atomic.AddInt64(&pc.inflight, 1)
	atomic.AddInt64(&pc.reqs, 1)
	return pc.conn, func(di balancer.DoneInfo) {
		atomic.AddInt64(&pc.inflight, -1)
		now := time.Now().UnixNano()
		// get moving average ratio w
		stamp := atomic.SwapInt64(&pc.stamp, now)
		td := now - stamp
		if td < 0 {
			td = 0
		}
		w := math.Exp(float64(-td) / float64(tau))

		lag := now - start
		if lag < 0 {
			lag = 0
		}
		oldLag := atomic.LoadUint64(&pc.lag)
		if oldLag == 0 {
			w = 0.0
		}
		lag = int64(float64(oldLag)*w + float64(lag)*(1.0-w))
		atomic.StoreUint64(&pc.lag, uint64(lag))

		success := uint64(1000) // error value ,if error set 1
		if di.Err != nil {
			if st, ok := status.FromError(di.Err); ok {
				// only counter the local grpc error, ignore any business error
				if st.Code() != codes.Unknown && st.Code() != codes.OK {
					success = 0
				}
			}
		}
		oldSuc := atomic.LoadUint64(&pc.success)
		success = uint64(float64(oldSuc)*w + float64(success)*(1.0-w))
		atomic.StoreUint64(&pc.success, success)

		trailer := di.Trailer
		if strs, ok := trailer[wmd.CPUUsage]; ok {
			if cpu, err2 := strconv.ParseUint(strs[0], 10, 64); err2 == nil && cpu > 0 {
				atomic.StoreUint64(&pc.svrCPU, cpu)
			}
		}

		logTs := atomic.LoadInt64(&p.logTs)
		if now-logTs > int64(time.Second*3) {
			if atomic.CompareAndSwapInt64(&p.logTs, logTs, now) {
				p.printStats()
			}
		}
	}, nil
}

func (p *p2cPicker) printStats() {
	if len(p.subConns) <= 0 {
		return
	}
	stats := make([]statistic, 0, len(p.subConns))
	for _, conn := range p.subConns {
		var stat statistic
		stat.addr = conn.addr.Addr
		stat.cpu = atomic.LoadUint64(&conn.svrCPU)
		stat.cs = atomic.LoadUint64(&conn.success)
		stat.inflight = atomic.LoadInt64(&conn.inflight)
		stat.lantency = atomic.LoadUint64(&conn.lag)
		stat.reqs = atomic.SwapInt64(&conn.reqs, 0)
		load := conn.load()
		if load != 0 {
			stat.score = float64(stat.cs*conn.meta.Weight*1e8) / float64(load)
		}
		stats = append(stats, stat)
	}
	log.Info("p2c %s : %+v", p.subConns[0].addr.ServerName, stats)
	//fmt.Printf("%+v\n", stats)
}
