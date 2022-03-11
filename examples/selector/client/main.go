package main

import (
	"context"
	"log"
	"time"

	"github.com/SeeMusic/kratos/contrib/registry/consul/v2"
	"github.com/SeeMusic/kratos/examples/helloworld/helloworld"
	"github.com/SeeMusic/kratos/v2/middleware/recovery"
	"github.com/SeeMusic/kratos/v2/selector/filter"
	"github.com/SeeMusic/kratos/v2/selector/p2c"
	"github.com/SeeMusic/kratos/v2/selector/wrr"
	"github.com/SeeMusic/kratos/v2/transport/grpc"
	"github.com/SeeMusic/kratos/v2/transport/http"
	"github.com/hashicorp/consul/api"
)

func main() {
	consulCli, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}
	r := consul.New(consulCli)

	// new grpc client
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///helloworld"),
		grpc.WithDiscovery(r),
		// 由于gRPC框架的限制只能使用全局balancer+filter的方式来实现selector
		// 这里使用weighted round robin算法的balancer+静态version=1.0.0的Filter
		grpc.WithBalancerName(wrr.Name),
		grpc.WithFilter(
			filter.Version("1.0.0"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	gClient := helloworld.NewGreeterClient(conn)

	// new http client
	hConn, err := http.NewClient(
		context.Background(),
		http.WithMiddleware(
			recovery.Recovery(),
		),
		http.WithEndpoint("discovery:///helloworld"),
		http.WithDiscovery(r),
		// 这里使用p2c算法的balancer+静态version=2.0.0的Filter组成一个selector
		http.WithSelector(
			p2c.New(p2c.WithFilter(filter.Version("2.0.0"))),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer hConn.Close()
	hClient := helloworld.NewGreeterHTTPClient(hConn)

	for {
		time.Sleep(time.Second)
		callGRPC(gClient)
		callHTTP(hClient)
	}
}

func callGRPC(client helloworld.GreeterClient) {
	reply, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[grpc] SayHello %+v\n", reply)
}

func callHTTP(client helloworld.GreeterHTTPClient) {
	reply, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[http] SayHello %s\n", reply.Message)
}
