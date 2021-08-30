module github.com/go-kratos/kratos/examples

go 1.16

require (
	entgo.io/ent v0.9.0
	github.com/BurntSushi/toml v0.3.1
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.1110 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.6.1
	github.com/gin-gonic/gin v1.7.3
	github.com/go-errors/errors v1.4.0 // indirect
	github.com/go-kratos/consul v0.1.4
	github.com/go-kratos/etcd v0.1.3
	github.com/go-kratos/gin v0.1.0
	github.com/go-kratos/kratos/middleware/logging/v2 v2.0.0-00010101000000-000000000000
	github.com/go-kratos/kratos/middleware/metadata/v2 v2.0.0-00010101000000-000000000000
	github.com/go-kratos/kratos/middleware/metrics/v2 v2.0.0-00010101000000-000000000000
	github.com/go-kratos/kratos/middleware/recovery/v2 v2.0.0-00010101000000-000000000000
	github.com/go-kratos/kratos/middleware/tracing/v2 v2.0.0-00010101000000-000000000000
	github.com/go-kratos/kratos/middleware/validate/v2 v2.0.0-00010101000000-000000000000
	github.com/go-kratos/kratos/registry/consul/v2 v2.0.0-00010101000000-000000000000
	github.com/go-kratos/kratos/registry/etcd/v2 v2.0.0-00010101000000-000000000000
	github.com/go-kratos/kratos/registry/nacos/v2 v2.0.0-00010101000000-000000000000
	github.com/go-kratos/kratos/registry/zookeeper/v2 v2.0.0-00010101000000-000000000000
	github.com/go-kratos/kratos/v2 v2.0.5
	github.com/go-kratos/prometheus v0.0.0-20210522055322-137e29e7cf47
	github.com/go-kratos/swagger-api v1.0.0
	github.com/go-playground/validator/v10 v10.6.1 // indirect
	github.com/go-redis/redis/extra/redisotel v0.3.0
	github.com/go-redis/redis/v8 v8.11.2
	github.com/go-sql-driver/mysql v1.6.0
	github.com/google/wire v0.5.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.5.0
	github.com/hashicorp/consul/api v1.9.1
	github.com/hashicorp/go-sockaddr v1.0.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/labstack/echo/v4 v4.5.0
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/nacos-group/nacos-sdk-go v1.0.8
	github.com/nicksnyder/go-i18n/v2 v2.1.2
	github.com/prometheus/client_golang v1.11.0
	github.com/segmentio/kafka-go v0.4.17
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/ugorji/go v1.2.6 // indirect
	go.etcd.io/etcd/client/v3 v3.5.0
	go.opentelemetry.io/otel v1.0.0-RC2
	go.opentelemetry.io/otel/exporters/jaeger v1.0.0-RC2
	go.opentelemetry.io/otel/sdk v1.0.0-RC2
	go.opentelemetry.io/otel/trace v1.0.0-RC2
	go.uber.org/zap v1.19.0
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e // indirect
	golang.org/x/text v0.3.6
	google.golang.org/genproto v0.0.0-20210805201207-89edb61ffb67
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
)

replace (
	github.com/go-kratos/kratos/middleware/logging/v2 => ../middleware/logging
	github.com/go-kratos/kratos/middleware/metadata/v2 => ../middleware/metadata
	github.com/go-kratos/kratos/middleware/metrics/v2 => ../middleware/metrics
	github.com/go-kratos/kratos/middleware/recovery/v2 => ../middleware/recovery
	github.com/go-kratos/kratos/middleware/tracing/v2 => ../middleware/tracing
	github.com/go-kratos/kratos/middleware/validate/v2 => ../middleware/validate
	github.com/go-kratos/kratos/registry/consul/v2 => ../registry/consul
	github.com/go-kratos/kratos/registry/etcd/v2 => ../registry/etcd
	github.com/go-kratos/kratos/registry/nacos/v2 => ../registry/nacos
	github.com/go-kratos/kratos/registry/zookeeper/v2 => ../registry/zookeeper
	github.com/go-kratos/kratos/v2 => ../
)
