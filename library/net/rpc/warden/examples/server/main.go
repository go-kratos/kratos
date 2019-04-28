package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/library/ecode"
	epb "go-common/library/ecode/pb"
	"go-common/library/log"
	"go-common/library/net/rpc/warden"
	pb "go-common/library/net/rpc/warden/proto/testproto"
	xtime "go-common/library/time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
)

type helloServer struct {
	addr string
}

func (s *helloServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	if in.Name == "err_detail_test" {
		any, _ := ptypes.MarshalAny(&pb.HelloReply{Success: true, Message: "this is test detail"})
		err := epb.From(ecode.AccessDenied)
		err.ErrDetail = any
		return nil, err
	}
	return &pb.HelloReply{Message: fmt.Sprintf("hello %s from %s", in.Name, s.addr)}, nil
}

func (s *helloServer) StreamHello(ss pb.Greeter_StreamHelloServer) error {
	for i := 0; i < 3; i++ {
		in, err := ss.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		ret := &pb.HelloReply{Message: "Hello " + in.Name, Success: true}
		err = ss.Send(ret)
		if err != nil {
			return err
		}
	}
	return nil
}

func runServer(addr string) *warden.Server {
	server := warden.NewServer(&warden.ServerConfig{
		//服务端每个请求的默认超时时间
		Timeout: xtime.Duration(time.Second),
	})
	server.Use(middleware())
	pb.RegisterGreeterServer(server.Server(), &helloServer{addr: addr})
	go func() {
		err := server.Run(addr)
		if err != nil {
			panic("run server failed!" + err.Error())
		}
	}()
	return server
}

func main() {
	log.Init(&log.Config{Stdout: true})
	server := runServer("0.0.0.0:8080")
	signalHandler(server)
}

//类似于中间件
func middleware() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		//记录调用方法
		log.Info("method:%s", info.FullMethod)
		//call chain
		resp, err = handler(ctx, req)
		return
	}
}

func signalHandler(s *warden.Server) {
	var (
		ch = make(chan os.Signal, 1)
	)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("get a signal %s, stop the consume process", si.String())
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()
			//gracefully shutdown with timeout
			s.Shutdown(ctx)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
