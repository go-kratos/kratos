package wrr

import (
	"context"
	"math"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bilibili/kratos/pkg/conf/env"
	"github.com/bilibili/kratos/pkg/log"
	nmd "github.com/bilibili/kratos/pkg/net/metadata"
	wmeta "github.com/bilibili/kratos/pkg/net/rpc/warden/internal/metadata"
	"github.com/bilibili/kratos/pkg/stat/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
)

var _ base.PickerBuilder = &wrrPickerBuilder{}
var _ balancer.Picker = &wrrPicker{}

// var dwrrFeature feature.Feature = "dwrr"

// Name is the name of round_robin balancer.
const Name = "wrr"

// newBuilder creates a new weighted-roundrobin balancer builder.
func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, &wrrPickerBuilder{})
}

func init() {
	//feature.DefaultGate.Add(map[feature.Feature]feature.Spec{
	//	dwrrFeature: {Default: false},
	//})

	balancer.Register(newBuilder())
}

type serverInfo struct {
	cpu     int64
	success uint64 // float64 bits
}

type subConn struct {
	conn balancer.SubConn
	addr resolver.Address
	meta wmeta.MD

	err     metric.RollingCounter
	latency metric.RollingGauge
	si      serverInfo
	// effective weight
	ewt int64
	// current weight
	cwt int64
	// last score
	score float64
}

func (c *subConn) errSummary() (err int64, req int64) {
	c.err.Reduce(func(iterator metric.Iterator) float64 {
		for iterator.Next() {
			bucket := iterator.Bucket()
			req += bucket.Count
			for _, p := range bucket.Points {
				err += int64(p)
			}
		}
		return 0
	})
	return
}

func (c *subConn) latencySummary() (latency float64, count int64) {
	c.latency.Reduce(func(iterator metric.Iterator) float64 {
		for iterator.Next() {
			bucket := iterator.Bucket()
			count += bucket.Count
			for _, p := range bucket.Points {
				latency += p
			}
		}
		return 0
	})
	return latency / float64(count), count
}

// statistics is info for log
type statistics struct {
	addr     string
	ewt      int64
	cs       float64
	ss       float64
	lantency float64
	cpu      float64
	req      int64
}

// Stats is grpc Interceptor for client to collect server stats
func Stats() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		var (
			trailer metadata.MD
			md      nmd.MD
			ok      bool
		)
		if md, ok = nmd.FromContext(ctx); !ok {
			md = nmd.MD{}
		} else {
			md = md.Copy()
		}
		ctx = nmd.NewContext(ctx, md)
		opts = append(opts, grpc.Trailer(&trailer))

		err = invoker(ctx, method, req, reply, cc, opts...)

		conn, ok := md["conn"].(*subConn)
		if !ok {
			return
		}
		if strs, ok := trailer[wmeta.CPUUsage]; ok {
			if cpu, err2 := strconv.ParseInt(strs[0], 10, 64); err2 == nil && cpu > 0 {
				atomic.StoreInt64(&conn.si.cpu, cpu)
			}
		}
		return
	}
}

type wrrPickerBuilder struct{}

func (*wrrPickerBuilder) Build(readySCs map[resolver.Address]balancer.SubConn) balancer.Picker {
	p := &wrrPicker{
		colors: make(map[string]*wrrPicker),
	}
	for addr, sc := range readySCs {
		meta, ok := addr.Metadata.(wmeta.MD)
		if !ok {
			meta = wmeta.MD{
				Weight: 10,
			}
		}
		subc := &subConn{
			conn: sc,
			addr: addr,

			meta:  meta,
			ewt:   int64(meta.Weight),
			score: -1,

			err: metric.NewRollingCounter(metric.RollingCounterOpts{
				Size:           10,
				BucketDuration: time.Millisecond * 100,
			}),
			latency: metric.NewRollingGauge(metric.RollingGaugeOpts{
				Size:           10,
				BucketDuration: time.Millisecond * 100,
			}),

			si: serverInfo{cpu: 500, success: math.Float64bits(1)},
		}
		if meta.Color == "" {
			p.subConns = append(p.subConns, subc)
			continue
		}
		// if color not empty, use color picker
		cp, ok := p.colors[meta.Color]
		if !ok {
			cp = &wrrPicker{}
			p.colors[meta.Color] = cp
		}
		cp.subConns = append(cp.subConns, subc)
	}
	return p
}

type wrrPicker struct {
	// subConns is the snapshot of the weighted-roundrobin balancer when this picker was
	// created. The slice is immutable. Each Get() will do a round robin
	// selection from it and return the selected SubConn.
	subConns []*subConn
	colors   map[string]*wrrPicker
	updateAt int64

	mu sync.Mutex
}

func (p *wrrPicker) Pick(ctx context.Context, opts balancer.PickOptions) (balancer.SubConn, func(balancer.DoneInfo), error) {
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

func (p *wrrPicker) pick(ctx context.Context, opts balancer.PickOptions) (balancer.SubConn, func(balancer.DoneInfo), error) {
	var (
		conn        *subConn
		totalWeight int64
	)
	if len(p.subConns) <= 0 {
		return nil, nil, balancer.ErrNoSubConnAvailable
	}
	p.mu.Lock()
	// nginx wrr load balancing algorithm: http://blog.csdn.net/zhangskd/article/details/50194069
	for _, sc := range p.subConns {
		totalWeight += sc.ewt
		sc.cwt += sc.ewt
		if conn == nil || conn.cwt < sc.cwt {
			conn = sc
		}
	}
	conn.cwt -= totalWeight
	p.mu.Unlock()
	start := time.Now()
	if cmd, ok := nmd.FromContext(ctx); ok {
		cmd["conn"] = conn
	}
	//if !feature.DefaultGate.Enabled(dwrrFeature) {
	//	return conn.conn, nil, nil
	//}
	return conn.conn, func(di balancer.DoneInfo) {
		ev := int64(0) // error value ,if error set 1
		if di.Err != nil {
			if st, ok := status.FromError(di.Err); ok {
				// only counter the local grpc error, ignore any business error
				if st.Code() != codes.Unknown && st.Code() != codes.OK {
					ev = 1
				}
			}
		}
		conn.err.Add(ev)

		now := time.Now()
		conn.latency.Add(now.Sub(start).Nanoseconds() / 1e5)
		u := atomic.LoadInt64(&p.updateAt)
		if now.UnixNano()-u < int64(time.Second) {
			return
		}
		if !atomic.CompareAndSwapInt64(&p.updateAt, u, now.UnixNano()) {
			return
		}
		var (
			stats = make([]statistics, len(p.subConns))
			count int
			total float64
		)
		for i, conn := range p.subConns {
			cpu := float64(atomic.LoadInt64(&conn.si.cpu))
			ss := math.Float64frombits(atomic.LoadUint64(&conn.si.success))
			errc, req := conn.errSummary()
			lagv, lagc := conn.latencySummary()

			if req > 0 && lagc > 0 && lagv > 0 {
				// client-side success ratio
				cs := 1 - (float64(errc) / float64(req))
				if cs <= 0 {
					cs = 0.1
				} else if cs <= 0.2 && req <= 5 {
					cs = 0.2
				}
				conn.score = math.Sqrt((cs * ss * ss * 1e9) / (lagv * cpu))
				stats[i] = statistics{cs: cs, ss: ss, lantency: lagv, cpu: cpu, req: req}
			}
			stats[i].addr = conn.addr.Addr

			if conn.score > 0 {
				total += conn.score
				count++
			}
		}
		// count must be greater than 1,otherwise will lead ewt to 0
		if count < 2 {
			return
		}
		avgscore := total / float64(count)
		p.mu.Lock()
		for i, conn := range p.subConns {
			if conn.score <= 0 {
				conn.score = avgscore
			}
			conn.ewt = int64(conn.score * float64(conn.meta.Weight))
			stats[i].ewt = conn.ewt
		}
		p.mu.Unlock()
		log.Info("warden wrr(%s): %+v", conn.addr.ServerName, stats)
	}, nil

}
