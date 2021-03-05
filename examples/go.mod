module github.com/go-kratos/examples

go 1.15

require (
	entgo.io/ent v0.6.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-kratos/kratos/v2 v2.0.0-alpha4
	github.com/go-redis/redis/extra/redisotel v0.2.0
	github.com/go-redis/redis/v8 v8.6.0
	github.com/go-sql-driver/mysql v1.5.1-0.20200311113236-681ffa848bae
	github.com/golang/protobuf v1.4.3
	github.com/google/wire v0.5.0
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/opentracing/opentracing-go v1.2.0
	go.opentelemetry.io/otel v0.17.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.17.0
	go.opentelemetry.io/otel/sdk v0.17.0
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20210228012217-479acdf4ea46 // indirect
	google.golang.org/genproto v0.0.0-20210226172003-ab064af71705
	google.golang.org/grpc v1.36.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/go-kratos/kratos/v2 => ./..
