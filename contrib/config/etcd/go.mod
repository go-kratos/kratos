module github.com/go-kratos/kratos/contrib/config/etcd/v2

go 1.16

require (
	github.com/go-kratos/kratos/v2 v2.6.3
	go.etcd.io/etcd/client/v3 v3.5.8
	google.golang.org/grpc v1.56.1
)

replace github.com/go-kratos/kratos/v2 => ../../../
