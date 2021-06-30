package main

import (
	"context"
	"os"
	"time"

	messagev1 "github.com/go-kratos/kratos/examples/traces/api/message"
	v1 "github.com/go-kratos/kratos/examples/traces/api/user"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
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
	v1.UnimplementedUserServer
	tracer trace.TracerProvider
}

// Get trace provider
func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(v1.User_ServiceDesc.ServiceName),
			attribute.String("environment", "development"),
			attribute.Int64("ID", 1),
		)),
	)
	return tp, nil
}

func (s *server) GetMyMessages(ctx context.Context, in *v1.GetMyMessagesRequest) (*v1.GetMyMessagesReply, error) {
	// create grpc conn
	conn, err := grpc.DialInsecure(ctx,
		grpc.WithEndpoint("127.0.0.1:9000"),
		grpc.WithMiddleware(
			recovery.Recovery(),
			tracing.Client(
				tracing.WithTracerProvider(s.tracer),
			),
		),
		grpc.WithTimeout(2*time.Second),
	)
	if err != nil {
		return nil, err
	}
	msg := messagev1.NewMessageServiceClient(conn)
	// Method of calling blog service
	reply, err := msg.GetUserMessage(ctx, &messagev1.GetUserMessageRequest{Id: 123, Count: in.Count})
	if err != nil {
		return nil, err
	}
	res := &v1.GetMyMessagesReply{}
	for _, v := range reply.Messages {
		res.Messages = append(res.Messages, &v1.Message{Content: v.Content})
	}
	return res, nil
}

func main() {
	logger := log.NewStdLogger(os.Stdout)
	logger = log.With(logger, "trace_id", log.TraceID())
	logger = log.With(logger, "span_id", log.SpanID())
	log := log.NewHelper(logger)

	tp, err := tracerProvider("http://jaeger:14268/api/traces")
	if err != nil {
		log.Error(err)
	}

	httpSrv := http.NewServer(
		http.Address(":8000"),
		http.Middleware(
			recovery.Recovery(),
			// Configuring tracing middleware
			tracing.Server(
				tracing.WithTracerProvider(tp),
			),
			logging.Server(logger),
		),
	)
	s := &server{tracer: tp}
	v1.RegisterUserHTTPServer(httpSrv, s)

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
