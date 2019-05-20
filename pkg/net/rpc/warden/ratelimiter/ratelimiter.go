package ratelimiter

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/bilibili/kratos/pkg/log"
	limit "github.com/bilibili/kratos/pkg/ratelimit"
	"github.com/bilibili/kratos/pkg/ratelimit/bbr"
	"github.com/bilibili/kratos/pkg/stat/prom"

	"google.golang.org/grpc"
)

const (
	_statName = "go_grpc_bbr"
)

var (
	stats = prom.New().WithState("go_grpc_bbr", []string{"url"})
)

// RateLimiter bbr middleware.
type RateLimiter struct {
	group   *bbr.Group
	logTime int64
}

// New return a ratelimit middleware.
func New(conf *bbr.Config) (s *RateLimiter) {
	return &RateLimiter{
		group:   bbr.NewGroup(conf),
		logTime: time.Now().UnixNano(),
	}
}

func (b *RateLimiter) printStats(fullMethod string, limiter limit.Limiter) {
	now := time.Now().UnixNano()
	if now-atomic.LoadInt64(&b.logTime) > int64(time.Second*3) {
		atomic.StoreInt64(&b.logTime, now)
		log.Info("grpc.bbr path:%s stat:%+v", fullMethod, limiter.(*bbr.BBR).Stat())
	}
}

// Limit is a server interceptor that detects and rejects overloaded traffic.
func (b *RateLimiter) Limit() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, args *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		uri := args.FullMethod
		limiter := b.group.Get(uri)
		done, err := limiter.Allow(ctx)
		if err != nil {
			stats.Incr(_statName, uri)
			return
		}
		defer func() {
			done(limit.DoneInfo{Op: limit.Success})
			b.printStats(uri, limiter)
		}()
		resp, err = handler(ctx, req)
		return
	}
}
