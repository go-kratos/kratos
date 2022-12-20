module github.com/go-kratos/kratos/contrib/config/etcd/v2

go 1.16

require (
	github.com/go-kratos/kratos/v2 v2.4.0
	go.etcd.io/etcd/client/v3 v3.5.4
	google.golang.org/grpc v1.51.0
)

replace github.com/go-kratos/kratos/v2 => ../../../
