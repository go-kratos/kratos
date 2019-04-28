package rpc

import (
	"context"
	"reflect"
	"sync/atomic"

	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

const (
	_clientsPool = 3
)

// balancer interface.
type balancer interface {
	Boardcast(context.Context, string, interface{}, interface{}) error
	Call(context.Context, string, interface{}, interface{}) error
}

// wrr get avaliable rpc client by wrr strategy.
type wrr struct {
	pool   []*Client
	weight int64
	server int64
	idx    int64
}

// Boardcast broad cast to all Client.
// NOTE: reply must be ptr.
func (r *wrr) Boardcast(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) (err error) {
	if r.weight == 0 {
		log.Error("wrr get() error weight:%d server:%d idx:%d", len(r.pool), r.server, r.idx)
		return ErrNoClient
	}
	rtp := reflect.TypeOf(reply).Elem()
	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i < int(r.server); i++ {
		j := i
		g.Go(func() error {
			nrp := reflect.New(rtp).Interface()
			return r.pool[j].Call(ctx, serviceMethod, args, nrp)
		})
	}
	return g.Wait()
}

func (r *wrr) Call(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) (err error) {
	if r.weight == 0 {
		log.Error("wrr get() error weight:%d server:%d idx:%d", len(r.pool), r.server, r.idx)
		return ErrNoClient
	}
	v := atomic.AddInt64(&r.idx, 1)
	for i := int64(0); i < r.server; i++ {
		cli := r.pool[int((v+i)%r.weight)]
		if err = cli.Call(ctx, serviceMethod, args, reply); err != ErrNoClient {
			return
		}
	}
	return ErrNoClient
}

type key interface {
	Key() int64
}

type sharding struct {
	pool   []*Client
	weight int64
	server int64
	idx    int64
}

// Boardcast broad cast to all clients.
// NOTE: reply must be ptr.
func (r *sharding) Boardcast(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) (err error) {
	if r.weight == 0 {
		log.Error("wrr get() error weight:%d server:%d idx:%d", len(r.pool), r.server, r.idx)
		return ErrNoClient
	}
	rtp := reflect.TypeOf(reply).Elem()
	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i < int(r.server); i++ {
		j := i
		g.Go(func() error {
			nrp := reflect.New(rtp).Interface()
			return r.pool[j].Call(ctx, serviceMethod, args, nrp)
		})
	}
	return g.Wait()
}

func (r *sharding) Call(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) (err error) {
	if r.weight == 0 {
		log.Error("wrr get() error weight:%d server:%d idx:%d", len(r.pool), r.server, r.idx)
		return ErrNoClient
	}
	if k, ok := args.(key); ok {
		if err = r.pool[int(k.Key()%r.server)].Call(ctx, serviceMethod, args, reply); err != ErrNoClient {
			return
		}
	}
	return ErrNoClient
}
