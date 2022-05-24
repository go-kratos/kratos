module github.com/go-kratos/kratos/contrib/registry/etcd/v2

go 1.16

require (
	github.com/go-kratos/kratos/v2 v2.3.0
	go.etcd.io/etcd/client/v3 v3.5.4
	google.golang.org/grpc v1.46.2
)

replace github.com/go-kratos/kratos/v2 => ../../../
