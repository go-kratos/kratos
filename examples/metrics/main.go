package main

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	prom "github.com/go-kratos/prometheus/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"

	"github.com/go-kratos/kratos/examples/helloworld/helloworld"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/prometheus/client_golang/prometheus"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = "metrics"
	// Version is the version of the compiled software.
	Version = "v1.0.0"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	helloworld.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	if in.Name == "error" {
		return nil, errors.BadRequest("custom_error", fmt.Sprintf("invalid argument %s", in.Name))
	}
	if in.Name == "panic" {
		panic("server panic")
	}
	return &helloworld.HelloReply{Message: fmt.Sprintf("Hello %+v", in.Name)}, nil
}

func main() {
	s := &server{}

	grpcHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "grpc_server",
		Name:      "request_duration_millisecond",
		Help:      "server requests duration(ms).",
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	}, []string{"kind", "operation"})

	grpcCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "kratos_client",
		Name:      "api_requests_total",
		Help:      "The total number of processed requests",
	}, []string{"kind", "operation", "code", "reason"})

	httpHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "http_server",
		Name:      "request_duration_millisecond",
		Help:      "server requests duration(ms).",
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	}, []string{"kind", "operation"})

	httpCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "http_server",
		Name:      "api_requests_total",
		Help:      "The total number of processed requests",
	}, []string{"kind", "operation", "code", "reason"})

	prometheus.MustRegister(grpcCounter, grpcHistogram, httpCounter, httpHistogram)
	grpcSrv := grpc.NewServer(
		grpc.Address(":9000"),
		grpc.Middleware(
			recovery.Recovery(),
			metrics.Server(
				metrics.WithSeconds(
					prom.NewHistogram(grpcHistogram),
				),
				metrics.WithRequests(prom.NewCounter(grpcCounter)),
			),
		),
	)
	httpSrv := http.NewServer(
		http.Address(":8000"),
		http.Middleware(
			recovery.Recovery(),
			metrics.Server(
				metrics.WithSeconds(
					prom.NewHistogram(httpHistogram),
				),
				metrics.WithRequests(prom.NewCounter(httpCounter)),
			),
		),
	)

	httpSrv.Handle("/metrics", promhttp.Handler())

	helloworld.RegisterGreeterServer(grpcSrv, s)
	helloworld.RegisterGreeterHTTPServer(httpSrv, s)

	app := kratos.New(
		kratos.Name(Name),
		kratos.Server(
			httpSrv,
			grpcSrv,
		),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
