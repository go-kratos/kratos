module github.com/go-kratos/kratos/contrib/config/etcd/v2

go 1.16

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-kratos/kratos/v2 v2.1.2
	github.com/stretchr/testify v1.7.0
	go.etcd.io/etcd/client/v3 v3.5.0
	google.golang.org/grpc v1.42.0
)

replace github.com/go-kratos/kratos/v2 => ../../../
