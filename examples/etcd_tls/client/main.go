package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/go-kratos/kratos/v2/transport/http"
	"io/ioutil"
	"log"
	"time"

	"github.com/go-kratos/etcd/registry"
	"github.com/go-kratos/kratos/examples/helloworld/helloworld"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {

	b, err := ioutil.ReadFile("../cert/server.crt")
	if err != nil {
		panic(err)
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		panic(err)
	}
	tlsConf := &tls.Config{ServerName: "www.kratos.com", RootCAs: cp}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		panic(err)
	}
	r := registry.New(cli)
	for {
		callGRPC(r, tlsConf)
		callHTTP(r, tlsConf)
		time.Sleep(time.Second)
	}
}

func callGRPC(r *registry.Registry, tlsConf *tls.Config) {
	conn, err := grpc.Dial(
		context.Background(),
		grpc.WithEndpoint("discovery:///helloworld"),
		grpc.WithDiscovery(r),
		grpc.WithTLSConfig(tlsConf),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := helloworld.NewGreeterClient(conn)
	reply, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[grpc] SayHello %+v\n", reply)
}

func callHTTP(r *registry.Registry, tlsConf *tls.Config) {
	conn, err := http.NewClient(
		context.Background(),
		http.WithEndpoint("discovery:///helloworld"),
		http.WithDiscovery(r),
		http.WithBlock(),
		http.WithTLSConfig(tlsConf),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := helloworld.NewGreeterHTTPClient(conn)
	reply, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[http] SayHello %+v\n", reply)
}
