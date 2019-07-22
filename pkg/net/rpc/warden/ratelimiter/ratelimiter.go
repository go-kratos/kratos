package ratelimiter

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/bilibili/kratos/pkg/log"
	limit "github.com/bilibili/kratos/pkg/ratelimit"
	"github.com/bilibili/kratos/pkg/ratelimit/bbr"
	"github.com/bilibili/kratos/pkg/stat/metric"
	"google.golang.org/grpc"
)

var (
	_metricServerBBR = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: "grpc_server",
		Subsystem: "",
		Name:      "bbr_total",
		Help:      "grpc server bbr total.",
		Labels:    []string{"url"},
	})
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
			_metricServerBBR.Inc(uri)
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
