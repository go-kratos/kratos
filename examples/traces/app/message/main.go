package main

import (
	"context"
	"os"

	v1 "github.com/SeeMusic/kratos/examples/traces/api/message"
	"github.com/SeeMusic/kratos/v2"
	"github.com/SeeMusic/kratos/v2/log"
	"github.com/SeeMusic/kratos/v2/middleware"
	"github.com/SeeMusic/kratos/v2/middleware/logging"
	"github.com/SeeMusic/kratos/v2/middleware/recovery"
	"github.com/SeeMusic/kratos/v2/middleware/tracing"
	"github.com/SeeMusic/kratos/v2/transport/grpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = "message"
	// Version is the version of the compiled software.
	// Version = "v1.0.0"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	v1.UnimplementedMessageServiceServer
}

// set trace provider
func setTracerProvider(url string) error {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return err
	}
	tp := tracesdk.NewTracerProvider(
		// Set the sampling rate based on the parent span to 100%
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(Name),
			attribute.String("env", "dev"),
		)),
	)
	otel.SetTracerProvider(tp)
	return nil
}

func (s *server) GetUserMessage(ctx context.Context, request *v1.GetUserMessageRequest) (*v1.GetUserMessageReply, error) {
	msgs := &v1.GetUserMessageReply{}
	for i := 0; i < int(request.Count); i++ {
		msgs.Messages = append(msgs.Messages, &v1.Message{Content: "Teletubbies say hello."})
	}
	return msgs, nil
}

func main() {
	logger := log.NewStdLogger(os.Stdout)
	logger = log.With(logger, "trace_id", tracing.TraceID())
	logger = log.With(logger, "span_id", tracing.SpanID())
	log := log.NewHelper(logger)

	url := "http://jaeger:14268/api/traces"
	if os.Getenv("jaeger_url") != "" {
		url = os.Getenv("jaeger_url")
	}
	err := setTracerProvider(url)
	if err != nil {
		log.Error(err)
	}

	s := &server{}
	// grpc server
	grpcSrv := grpc.NewServer(
		grpc.Address(":9000"),
		grpc.Middleware(
			middleware.Chain(
				recovery.Recovery(),
				tracing.Server(),
				logging.Server(logger),
			),
		))
	v1.RegisterMessageServiceServer(grpcSrv, s)

	app := kratos.New(
		kratos.Name(Name),
		kratos.Server(
			grpcSrv,
		),
	)

	if err := app.Run(); err != nil {
		log.Error(err)
	}
}
