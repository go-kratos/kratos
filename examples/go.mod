module github.com/go-kratos/kratos/examples

go 1.16

require (
	entgo.io/ent v0.8.0
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.1103 // indirect
	github.com/armon/go-metrics v0.3.8 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.6.1
	github.com/fatih/color v1.11.0 // indirect
	github.com/gin-gonic/gin v1.7.2
	github.com/go-errors/errors v1.4.0 // indirect
	github.com/go-kratos/consul v0.0.0-20210425141546-e047a9f6ec87
	github.com/go-kratos/etcd v0.0.0-20210423155933-752c35c0d203
	github.com/go-kratos/kratos/v2 v2.0.0-beta4.0.20210510133946-86193a2a5ff6
	github.com/go-kratos/nacos v0.0.0-20210423031012-5dee2cc4aea1
	github.com/go-playground/validator/v10 v10.6.1 // indirect
	github.com/go-redis/redis/extra/redisotel v0.3.0
	github.com/go-redis/redis/v8 v8.8.3
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/google/wire v0.5.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/hashicorp/consul/api v1.8.1
	github.com/hashicorp/go-hclog v0.16.1 // indirect
	github.com/hashicorp/go-sockaddr v1.0.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/nacos-group/nacos-sdk-go v1.0.7
	github.com/ugorji/go v1.2.6 // indirect
	go.etcd.io/etcd/client/v3 v3.5.0-beta.3
	go.etcd.io/etcd/pkg/v3 v3.5.0-beta.3 // indirect
	go.opentelemetry.io/otel v0.20.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.20.0
	go.opentelemetry.io/otel/sdk v0.20.0
	go.opentelemetry.io/otel/trace v0.20.0
	go.uber.org/multierr v1.7.0 // indirect
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/net v0.0.0-20210521195947-fe42d452be8f // indirect
	golang.org/x/sys v0.0.0-20210521203332-0cec03c779c1 // indirect
	golang.org/x/tools v0.1.1 // indirect
	google.golang.org/genproto v0.0.0-20210521181308-5ccab8a35a9a
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/go-kratos/kratos/v2 => ../
