module github.com/go-kratos/kratos/contrib/config/etcd/v2

go 1.16

require (
	github.com/go-kratos/kratos/v2 v2.2.1
	go.etcd.io/etcd/client/v3 v3.5.0
	google.golang.org/grpc v1.44.0
)

replace github.com/go-kratos/kratos/v2 => ../../../
