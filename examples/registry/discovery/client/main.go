package main

import (
	"context"
	"log"
	"time"

	srcgrpc "google.golang.org/grpc"

	"github.com/SeeMusic/kratos/contrib/registry/discovery/v2"
	"github.com/SeeMusic/kratos/examples/helloworld/helloworld"
	"github.com/SeeMusic/kratos/v2/transport/grpc"
)

func main() {
	r := discovery.New(&discovery.Config{
		Nodes:  []string{"0.0.0.0:7171"},
		Env:    "dev",
		Region: "sh1",
		Zone:   "zone1",
		Host:   "localhost",
	}, nil)

	connGRPC, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///helloworld"),
		grpc.WithDiscovery(r),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer connGRPC.Close()

	//connHTTP, err := http.NewClient(
	//	context.Background(),
	//	http.WithEndpoint("discovery:///helloworld"),
	//	http.WithDiscovery(r),
	//	http.WithBlock(),
	//)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer connHTTP.Close()

	for {
		callGRPC(connGRPC)
		time.Sleep(time.Second)
	}
}

func callGRPC(conn *srcgrpc.ClientConn) {
	client := helloworld.NewGreeterClient(conn)
	reply, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[grpc] SayHello %+v\n", reply)
}
