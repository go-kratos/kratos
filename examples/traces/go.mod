module github.com/go-kratos/kratos/examples/traces

go 1.15

require (
	github.com/go-kratos/kratos/v2 v2.0.0-beta4
	github.com/gorilla/mux v1.8.0
	go.opentelemetry.io/otel v0.20.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.20.0
	go.opentelemetry.io/otel/sdk v0.20.0
	go.opentelemetry.io/otel/trace v0.20.0
	google.golang.org/genproto v0.0.0-20210518161634-ec7691c0a37d
	google.golang.org/grpc v1.37.1
	google.golang.org/protobuf v1.26.0
)
