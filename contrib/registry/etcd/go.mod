module github.com/go-kratos/kratos/contrib/registry/etcd/v3

go 1.25.0

require (
	github.com/go-kratos/kratos/v3 v3.0.0
	go.etcd.io/etcd/client/v3 v3.6.11
	google.golang.org/grpc v1.81.0
)

require (
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/coreos/go-systemd/v22 v22.7.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.29.0 // indirect
	go.etcd.io/etcd/api/v3 v3.6.11 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.6.11 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.28.0 // indirect
	golang.org/x/net v0.54.0 // indirect
	golang.org/x/sys v0.44.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260511170946-3700d4141b60 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260511170946-3700d4141b60 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/go-kratos/kratos/v3 => ../../../
