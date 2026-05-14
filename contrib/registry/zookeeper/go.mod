module github.com/go-kratos/kratos/contrib/registry/zookeeper/v3

go 1.25.0

require (
	github.com/go-kratos/kratos/v3 v3.0.0
	github.com/go-zookeeper/zk v1.0.4
	golang.org/x/sync v0.20.0
)

replace github.com/go-kratos/kratos/v3 => ../../../
