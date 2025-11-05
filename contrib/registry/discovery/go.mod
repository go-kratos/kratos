module github.com/go-kratos/kratos/contrib/registry/discovery/v2

go 1.22

require (
	github.com/go-kratos/kratos/v2 v2.9.0
	github.com/go-resty/resty/v2 v2.11.0
	github.com/pkg/errors v0.9.1
)

require golang.org/x/net v0.33.0 // indirect

replace github.com/go-kratos/kratos/v2 => ../../../
