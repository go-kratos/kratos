package warden

import (
	"context"
	"strconv"
	"time"

	"go-common/library/ecode"
	nmd "go-common/library/net/metadata"
	"go-common/library/stat/summary"
	"go-common/library/stat/sys/cpu"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (s *Server) stats() grpc.UnaryServerInterceptor {
	errstat := summary.New(time.Second*3, 10)

	return func(ctx context.Context, req interface{}, args *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		c := int64(0)
		if resp, err = handler(ctx, req); err != nil {
			if ecode.ServerErr.Equal(err) || ecode.ServiceUnavailable.Equal(err) ||
				ecode.Deadline.Equal(err) || ecode.LimitExceed.Equal(err) {
				c = 1
			}
		}
		errstat.Add(c)
		errors, requests := errstat.Value()
		kv := []string{
			nmd.Errors, strconv.FormatInt(errors, 10),
			nmd.Requests, strconv.FormatInt(requests, 10),
		}
		var cpustat cpu.Stat
		cpu.ReadStat(&cpustat)
		if cpustat.Usage != 0 {
			kv = append(kv, nmd.CPUUsage, strconv.FormatInt(int64(cpustat.Usage), 10))
		}
		trailer := metadata.Pairs(kv...)
		grpc.SetTrailer(ctx, trailer)
		return
	}
}
