package main

import (
	"context"
	"net/http"
	"time"

	"go-common/library/net/rpc/warden"
	pb "go-common/library/net/rpc/warden/proto/testproto"
	xtime "go-common/library/time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	config = &warden.ServerConfig{Timeout: xtime.Duration(time.Second)}
)

func main() {
	newServer()
}

type hello struct {
}

func (s *hello) SayHello(c context.Context, in *pb.HelloRequest) (out *pb.HelloReply, err error) {
	out = new(pb.HelloReply)
	out.Message = in.Name
	return
}

func (s *hello) StreamHello(ss pb.Greeter_StreamHelloServer) error {
	return nil
}
func newServer() {
	server := warden.NewServer(config)
	pb.RegisterGreeterServer(server.Server(), &hello{})
	go func() {
		http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
			h := promhttp.Handler()
			h.ServeHTTP(w, r)
		})
		http.ListenAndServe("0.0.0.0:9998", nil)
	}()
	err := server.Run(":9999")
	if err != nil {
		return
	}

}
