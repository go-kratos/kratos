module github.com/go-kratos/kratos/middleware/tracing/v2

go 1.15

require (
	github.com/go-kratos/kratos/v2 v2.0.5
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/otel v1.0.0-RC2
	go.opentelemetry.io/otel/sdk v1.0.0-RC2
	go.opentelemetry.io/otel/trace v1.0.0-RC2
	google.golang.org/grpc v1.39.1
	google.golang.org/protobuf v1.27.1
)

replace github.com/go-kratos/kratos/v2 => ../../../kratos
