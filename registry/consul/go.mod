module github.com/go-kratos/kratos/registry/consul/v2

go 1.16

require (
	github.com/go-kratos/kratos/v2 v2.0.5
	github.com/hashicorp/consul/api v1.9.1
	github.com/stretchr/testify v1.7.0
)

replace (
	github.com/go-kratos/kratos/v2 v2.0.5 => ../../../kratos
)
