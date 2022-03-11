package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"

	pb "github.com/SeeMusic/kratos/examples/helloworld/helloworld"
	"github.com/SeeMusic/kratos/v2/transport/grpc"
	"github.com/SeeMusic/kratos/v2/transport/http"
)

func main() {
	// Load CA certificate pem file.
	b, err := os.ReadFile("../cert/ca.crt")
	if err != nil {
		panic(err)
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		panic(err)
	}
	tlsConf := &tls.Config{ServerName: "www.kratos.com", RootCAs: cp}
	callHTTP(tlsConf)
	callGRPC(tlsConf)
}

func callHTTP(tlsConf *tls.Config) {
	conn, err := http.NewClient(
		context.Background(),
		http.WithEndpoint("https://127.0.0.1:8000"),
		http.WithTLSConfig(tlsConf),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pb.NewGreeterHTTPClient(conn)
	reply, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[http] SayHello %s\n", reply.Message)
}

func callGRPC(tlsConf *tls.Config) {
	conn, err := grpc.Dial(
		context.Background(),
		grpc.WithEndpoint("127.0.0.1:9000"),
		grpc.WithTLSConfig(tlsConf),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)
	reply, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[grpc] SayHello %+v\n", reply)
}
