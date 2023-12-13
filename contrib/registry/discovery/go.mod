module github.com/go-kratos/kratos/contrib/registry/discovery/v2

go 1.19

require (
	github.com/go-kratos/kratos/v2 v2.7.2
	github.com/go-resty/resty/v2 v2.10.0
	github.com/pkg/errors v0.9.1
)

require golang.org/x/net v0.17.0 // indirect

replace github.com/go-kratos/kratos/v2 => ../../../
