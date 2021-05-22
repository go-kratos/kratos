package main

import (
	"context"
	message "github.com/go-kratos/kratos/examples/traces/api/message"
	pb "github.com/go-kratos/kratos/examples/traces/api/user"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
	"os"
	"time"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = "user"
	// Version is the version of the compiled software.
	Version = "v1.0.0"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedUserServer
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
			semconv.ServiceNameKey.String(pb.User_ServiceDesc.ServiceName),
			attribute.String("environment", "development"),
			attribute.Int64("ID", 1),
		)),
	)
	return tp, nil
}

func (s *server) GetMyMessages(ctx context.Context, in *pb.GetMyMessagesRequest) (*pb.GetMyMessagesReply, error) {
	// create grpc conn
	conn, err := grpc.DialInsecure(ctx,
		grpc.WithEndpoint("127.0.0.1:9000"),
		grpc.WithMiddleware(middleware.Chain(
			tracing.Client(
				tracing.WithTracerProvider(s.tracer),
				tracing.WithPropagators(
					propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{}),
				),
			),
			recovery.Recovery())),
		grpc.WithTimeout(2*time.Second),
	)
	if err != nil {
		return nil, err
	}
	msg := message.NewMessageServiceClient(conn)
	// Method of calling blog service
	reply, err := msg.GetUserMessage(ctx, &message.GetUserMessageRequest{Id: 123, Count: in.Count})
	if err != nil {
		return nil, err
	}
	res := &pb.GetMyMessagesReply{}
	for _, v := range reply.Messages {
		res.Messages = append(res.Messages,&pb.Message{Content: v.Content})
	}
	return res, nil
}

func main() {
	logger := log.NewStdLogger(os.Stdout)

	log := log.NewHelper(logger)

	tp, err := tracerProvider("http://jaeger:14268/api/traces")
	if err != nil {
		log.Error(err)
	}

	s := &server{tracer: tp}

	httpSrv := http.NewServer(http.Address(":8000"))
	httpSrv.HandlePrefix("/", pb.NewUserHandler(s,
		http.Middleware(
			middleware.Chain(
				recovery.Recovery(),
				// Configuring tracing middleware
				tracing.Server(
					tracing.WithTracerProvider(tp),
					tracing.WithPropagators(
						propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{}),
					),
				),
				logging.Server(logger),
			),
		)),
	)

	app := kratos.New(
		kratos.Name(Name),
		kratos.Server(
			httpSrv,
		),
	)

	if err := app.Run(); err != nil {
		log.Error(err)
	}
}
