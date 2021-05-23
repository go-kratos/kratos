module github.com/go-kratos/kratos/examples

go 1.16

require (
	entgo.io/ent v0.6.0
	github.com/envoyproxy/protoc-gen-validate v0.6.1
	github.com/gin-gonic/gin v1.6.3
	github.com/go-kratos/consul v0.0.0-20210425141546-e047a9f6ec87
	github.com/go-kratos/etcd v0.0.0-20210423155933-752c35c0d203
	github.com/go-kratos/kratos/v2 v2.0.0-beta4.0.20210510133946-86193a2a5ff6
	github.com/go-kratos/nacos v0.0.0-20210423031012-5dee2cc4aea1
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/go-redis/redis/extra/redisotel v0.3.0
	github.com/go-redis/redis/v8 v8.7.1
	github.com/go-sql-driver/mysql v1.5.1-0.20200311113236-681ffa848bae
	github.com/golang/protobuf v1.5.2
	github.com/google/wire v0.5.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/hashicorp/consul/api v1.8.1
	github.com/hashicorp/go-sockaddr v1.0.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/nacos-group/nacos-sdk-go v1.0.7
	github.com/ugorji/go v1.2.3 // indirect
	go.etcd.io/etcd/client/v3 v3.5.0-alpha.0
	go.opentelemetry.io/otel v0.20.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.20.0
	go.opentelemetry.io/otel/sdk v0.20.0
	go.opentelemetry.io/otel/trace v0.20.0
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83 // indirect
	golang.org/x/net v0.0.0-20210521195947-fe42d452be8f // indirect
	golang.org/x/sys v0.0.0-20210521203332-0cec03c779c1 // indirect
	golang.org/x/tools v0.1.1 // indirect
	google.golang.org/genproto v0.0.0-20210309190941-1aeedc14537d
	google.golang.org/grpc v1.36.0
	google.golang.org/protobuf v1.26.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/go-kratos/kratos/v2 => ../
