package main

import (
	"context"
	pb "github.com/go-kratos/kratos/examples/traces/api/message"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
	"os"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = "message"
	// Version is the version of the compiled software.
	Version = "v1.0.0"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedMessageServiceServer
	tracer trace.TracerProvider
}

// Get trace provider
func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.NewRawExporter(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.ServiceNameKey.String(pb.MessageService_ServiceDesc.ServiceName),
			attribute.String("environment", "development"),
			attribute.Int64("ID", 1),
		)),
	)
	return tp, nil
}

func (s *server) GetUserMessage(ctx context.Context, request *pb.GetUserMessageRequest) (*pb.GetUserMessageReply, error) {
	msgs := &pb.GetUserMessageReply{}
	for i := 0; i < int(request.Count); i++ {
		msgs.Messages = append(msgs.Messages,&pb.Message{Content: "Teletubbies say hello."})
	}
	return msgs, nil
}

func main() {
	logger := log.NewStdLogger(os.Stdout)

	log := log.NewHelper(logger)

	tp, err := tracerProvider("http://jaeger:14268/api/traces")
	if err != nil {
		log.Error(err)
	}

	s := &server{tracer: tp}
	// grpc server
	grpcSrv := grpc.NewServer(
		grpc.Address(":9000"),
		grpc.Middleware(
			middleware.Chain(
				recovery.Recovery(),
				// Configuring tracing Middleware
				tracing.Server(
					tracing.WithTracerProvider(tp),
					tracing.WithPropagators(
						propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{}),
					),
				),
				logging.Server(logger),
			),
		))
	pb.RegisterMessageServiceServer(grpcSrv, s)

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
