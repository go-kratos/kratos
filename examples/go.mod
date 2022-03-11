module github.com/SeeMusic/kratos/examples

go 1.16

require (
	entgo.io/ent v0.9.0
	github.com/BurntSushi/toml v0.3.1
	github.com/envoyproxy/protoc-gen-validate v0.6.2
	github.com/gin-gonic/gin v1.7.7
	github.com/go-kratos/gin v0.1.0
	github.com/SeeMusic/kratos/contrib/config/apollo/v2 v2.0.0-20210901080230-515b71ec9061
	github.com/SeeMusic/kratos/contrib/metrics/prometheus/v2 v2.0.0-00010101000000-000000000000
	github.com/SeeMusic/kratos/contrib/registry/consul/v2 v2.0.0-00010101000000-000000000000
	github.com/SeeMusic/kratos/contrib/registry/discovery/v2 v2.0.0-00010101000000-000000000000
	github.com/SeeMusic/kratos/contrib/registry/etcd/v2 v2.0.0-00010101000000-000000000000
	github.com/SeeMusic/kratos/contrib/registry/nacos/v2 v2.0.0-00010101000000-000000000000
	github.com/SeeMusic/kratos/contrib/registry/polaris/v2 v2.0.0
	github.com/SeeMusic/kratos/contrib/registry/zookeeper/v2 v2.0.0-00010101000000-000000000000
	github.com/SeeMusic/kratos/v2 v2.2.0
	github.com/go-kratos/swagger-api v1.0.0
	github.com/go-redis/redis/extra/redisotel v0.3.0
	github.com/go-redis/redis/v8 v8.11.2
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang-jwt/jwt/v4 v4.2.0
	github.com/google/wire v0.5.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.5.0
	github.com/hashicorp/consul/api v1.12.0
	github.com/labstack/echo/v4 v4.6.3
	github.com/nacos-group/nacos-sdk-go v1.0.9
	github.com/nicksnyder/go-i18n/v2 v2.1.2
	github.com/polarismesh/polaris-go v1.0.1
	github.com/prometheus/client_golang v1.12.0
	github.com/segmentio/kafka-go v0.4.27
	github.com/sirupsen/logrus v1.8.1
	github.com/soheilhy/cmux v0.1.4
	go.etcd.io/etcd/client/v3 v3.5.1
	go.opentelemetry.io/otel v1.3.0
	go.opentelemetry.io/otel/exporters/jaeger v1.3.0
	go.opentelemetry.io/otel/sdk v1.3.0
	go.opentelemetry.io/otel/trace v1.3.0
	go.uber.org/zap v1.19.1
	golang.org/x/text v0.3.7
	google.golang.org/genproto v0.0.0-20220126215142-9970aeb2e350
	google.golang.org/grpc v1.44.0
	google.golang.org/grpc/examples v0.0.0-20220105183818-2fb1ac854b20 // indirect
	google.golang.org/protobuf v1.27.1
	gorm.io/driver/mysql v1.3.2
	gorm.io/gorm v1.23.1
)

replace (
	github.com/buger/jsonparser => github.com/buger/jsonparser v1.1.1
	github.com/SeeMusic/kratos/contrib/config/apollo/v2 => ../contrib/config/apollo
	github.com/SeeMusic/kratos/contrib/metrics/prometheus/v2 => ../contrib/metrics/prometheus
	github.com/SeeMusic/kratos/contrib/registry/consul/v2 => ../contrib/registry/consul
	github.com/SeeMusic/kratos/contrib/registry/discovery/v2 => ../contrib/registry/discovery
	github.com/SeeMusic/kratos/contrib/registry/etcd/v2 => ../contrib/registry/etcd
	github.com/SeeMusic/kratos/contrib/registry/nacos/v2 => ../contrib/registry/nacos
	github.com/SeeMusic/kratos/contrib/registry/polaris/v2 => ../contrib/registry/polaris
	github.com/SeeMusic/kratos/contrib/registry/zookeeper/v2 => ../contrib/registry/zookeeper

	github.com/SeeMusic/kratos/v2 => ../
)
