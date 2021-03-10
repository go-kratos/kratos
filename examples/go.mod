module github.com/go-kratos/examples

go 1.15

require (
	entgo.io/ent v0.6.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-kratos/consul v0.0.0-20210228130349-387ba759cd99
	github.com/go-kratos/kratos/v2 v2.0.0-alpha4
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/go-redis/redis/extra/redisotel v0.2.0
	github.com/go-redis/redis/v8 v8.6.0
	github.com/go-sql-driver/mysql v1.5.1-0.20200311113236-681ffa848bae
	github.com/golang/protobuf v1.4.3
	github.com/google/wire v0.5.0
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/consul/api v1.8.1
	github.com/hashicorp/go-sockaddr v1.0.2 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/ugorji/go v1.2.3 // indirect
	go.opentelemetry.io/otel v0.17.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.17.0
	go.opentelemetry.io/otel/sdk v0.17.0
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	google.golang.org/genproto v0.0.0-20210226172003-ab064af71705
	google.golang.org/grpc v1.36.0
	google.golang.org/protobuf v1.25.1-0.20200805231151-a709e31e5d12
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/go-kratos/kratos/v2 => ./..
