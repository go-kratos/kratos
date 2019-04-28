package warden

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"go-common/library/ecode"
	"go-common/library/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// recovery is a server interceptor that recovers from any panics.
func (s *Server) recovery() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, args *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if rerr := recover(); rerr != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				rs := runtime.Stack(buf, false)
				if rs > size {
					rs = size
				}
				buf = buf[:rs]
				pl := fmt.Sprintf("grpc server panic: %v\n%v\n%s\n", req, rerr, buf)
				fmt.Fprintf(os.Stderr, pl)
				log.Error(pl)
				err = status.Errorf(codes.Unknown, ecode.ServerErr.Error())
			}
		}()
		resp, err = handler(ctx, req)
		return
	}
}

// recovery return a client interceptor  that recovers from any panics.
func (c *Client) recovery() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		defer func() {
			if rerr := recover(); rerr != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				rs := runtime.Stack(buf, false)
				if rs > size {
					rs = size
				}
				buf = buf[:rs]
				pl := fmt.Sprintf("grpc client panic: %v\n%v\n%v\n%s\n", req, reply, rerr, buf)
				fmt.Fprintf(os.Stderr, pl)
				log.Error(pl)
				err = ecode.ServerErr
			}
		}()
		err = invoker(ctx, method, req, reply, cc, opts...)
		return
	}
}
