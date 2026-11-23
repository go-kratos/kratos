module github.com/go-kratos/kratos/contrib/registry/discovery/v3

go 1.25.0

require (
	github.com/go-kratos/kratos/v3 v3.0.0
	github.com/go-resty/resty/v2 v2.17.2
	github.com/pkg/errors v0.9.1
)

require (
	golang.org/x/net v0.54.0 // indirect
	golang.org/x/time v0.15.0 // indirect
)

replace github.com/go-kratos/kratos/v3 => ../../../
