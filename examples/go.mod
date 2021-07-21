module github.com/go-kratos/kratos/examples

go 1.16

require (
	entgo.io/ent v0.8.0
	github.com/BurntSushi/toml v0.3.1
	github.com/envoyproxy/protoc-gen-validate v0.6.1
	github.com/gin-gonic/gin v1.7.2
	github.com/go-kratos/consul v0.1.2
	github.com/go-kratos/etcd v0.1.0
	github.com/go-kratos/gin v0.1.0
	github.com/go-kratos/kratos/v2 v2.0.1
	github.com/go-kratos/nacos v0.1.0
	github.com/go-kratos/prometheus v0.0.0-20210522055322-137e29e7cf47
	github.com/go-kratos/swagger-api v0.1.8
	github.com/go-kratos/zookeeper v0.0.0-20210721021028-9e7498619b89
	github.com/go-redis/redis/extra/redisotel v0.3.0
	github.com/go-redis/redis/v8 v8.9.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/google/wire v0.5.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/hashicorp/consul/api v1.9.1
	github.com/labstack/echo/v4 v4.3.0
	github.com/longXboy/grpc-gateway/v2 v2.0.0-20210707031540-bd2d73d86cee
	github.com/nacos-group/nacos-sdk-go v1.0.7
	github.com/nicksnyder/go-i18n/v2 v2.1.2
	github.com/prometheus/client_golang v1.9.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	go.etcd.io/etcd/client/v3 v3.5.0-beta.4
	go.opentelemetry.io/otel v1.0.0-RC1
	go.opentelemetry.io/otel/exporters/jaeger v1.0.0-RC1
	go.opentelemetry.io/otel/sdk v1.0.0-RC1
	go.opentelemetry.io/otel/trace v1.0.0-RC1
	go.uber.org/zap v1.18.1
	golang.org/x/text v0.3.6
	google.golang.org/genproto v0.0.0-20210701191553-46259e63a0a9
	google.golang.org/grpc v1.39.0
	google.golang.org/protobuf v1.27.1
)

replace github.com/go-kratos/kratos/v2 => ../
