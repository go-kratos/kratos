package ewma

import (
	"context"
	"math"
	"net"
	"sync/atomic"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/selector"
)

const (
	// The mean lifetime of `cost`, it reaches its half-life after Tau*ln(2).
	tau = int64(time.Millisecond * 600)
	// if statistic not collected,we add a big lag penalty to endpoint
	penalty = uint64(time.Microsecond * 100)
)

var (
	_ selector.WeightedNode        = (*Node)(nil)
	_ selector.WeightedNodeBuilder = (*Builder)(nil)
)

// Node is endpoint instance
type Node struct {
	selector.Node

	// client statistic data
	lag       int64
	success   uint64
	inflight  int64
	inflights [200]int64
	// last collected timestamp
	stamp int64
	// request number in a period time
	reqs int64
	// last lastPick timestamp
	lastPick int64

	errHandler   func(err error) (isErr bool)
	cachedWeight *atomic.Value
}

type nodeWeight struct {
	value    float64
	updateAt int64
}

// Builder is ewma node builder.
type Builder struct {
	ErrHandler func(err error) (isErr bool)
}

// Build create a weighted node.
func (b *Builder) Build(n selector.Node) selector.WeightedNode {
	s := &Node{
		Node:         n,
		lag:          0,
		success:      1000,
		inflight:     1,
		errHandler:   b.ErrHandler,
		cachedWeight: &atomic.Value{},
	}
	return s
}

func (n *Node) health() uint64 {
	return atomic.LoadUint64(&n.success)
}

func (n *Node) load() (load uint64) {
	now := time.Now().UnixNano()
	avgLag := atomic.LoadInt64(&n.lag)
	predict := n.predict(avgLag, now)

	if avgLag == 0 {
		// penalty is the penalty value when there is no data when the node is just started.
		load = penalty * uint64(atomic.LoadInt64(&n.inflight))
		return
	}
	if predict > avgLag {
		avgLag = predict
	}
	// add 5ms to eliminate the latency gap between different zones
	avgLag += int64(time.Millisecond * 5)
	avgLag = int64(math.Sqrt(float64(avgLag)))
	load = uint64(avgLag) * uint64(atomic.LoadInt64(&n.inflight))
	return load
}

func (n *Node) predict(avgLag int64, now int64) (predict int64) {
	var (
		total    int64
		slowNum  int
		totalNum int
	)
	for i := range n.inflights {
		start := atomic.LoadInt64(&n.inflights[i])
		if start != 0 {
			totalNum++
			lag := now - start
			if lag > avgLag {
				slowNum++
				total += lag
			}
		}
	}
	if slowNum >= (totalNum/2 + 1) {
		predict = total / int64(slowNum)
	}
	return
}

// Pick pick a node.
func (n *Node) Pick() selector.DoneFunc {
	start := time.Now().UnixNano()
	atomic.StoreInt64(&n.lastPick, start)
	atomic.AddInt64(&n.inflight, 1)
	reqs := atomic.AddInt64(&n.reqs, 1)
	slot := reqs % 200
	swapped := atomic.CompareAndSwapInt64(&n.inflights[slot], 0, start)
	return func(_ context.Context, di selector.DoneInfo) {
		if swapped {
			atomic.CompareAndSwapInt64(&n.inflights[slot], start, 0)
		}
		atomic.AddInt64(&n.inflight, -1)

		now := time.Now().UnixNano()
		// get moving average ratio w
		stamp := atomic.SwapInt64(&n.stamp, now)
		td := now - stamp
		if td < 0 {
			td = 0
		}
		w := math.Exp(float64(-td) / float64(tau))

		lag := now - start
		if lag < 0 {
			lag = 0
		}
		oldLag := atomic.LoadInt64(&n.lag)
		if oldLag == 0 {
			w = 0.0
		}
		lag = int64(float64(oldLag)*w + float64(lag)*(1.0-w))
		atomic.StoreInt64(&n.lag, lag)

		success := uint64(1000) // error value ,if error set 1
		if di.Err != nil {
			if n.errHandler != nil && n.errHandler(di.Err) {
				success = 0
			}
			var netErr net.Error
			if errors.Is(context.DeadlineExceeded, di.Err) || errors.Is(context.Canceled, di.Err) ||
				errors.IsServiceUnavailable(di.Err) || errors.IsGatewayTimeout(di.Err) || errors.As(di.Err, &netErr) {
				success = 0
			}
		}
		oldSuc := atomic.LoadUint64(&n.success)
		success = uint64(float64(oldSuc)*w + float64(success)*(1.0-w))
		atomic.StoreUint64(&n.success, success)
	}
}

// Weight is node effective weight.
func (n *Node) Weight() (weight float64) {
	w, ok := n.cachedWeight.Load().(*nodeWeight)
	now := time.Now().UnixNano()
	if !ok || time.Duration(now-w.updateAt) > (time.Millisecond*5) {
		health := n.health()
		load := n.load()
		weight = float64(health*uint64(time.Microsecond)*10) / float64(load)
		n.cachedWeight.Store(&nodeWeight{
			value:    weight,
			updateAt: now,
		})
	} else {
		weight = w.value
	}
	return
}

func (n *Node) PickElapsed() time.Duration {
	return time.Duration(time.Now().UnixNano() - atomic.LoadInt64(&n.lastPick))
}

func (n *Node) Raw() selector.Node {
	return n.Node
}
