module github.com/go-kratos/kratos/contrib/metrics/datadog/v2

go 1.19

require (
	github.com/DataDog/datadog-go v4.8.3+incompatible
	github.com/go-kratos/kratos/v2 v2.7.2
)

require (
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/stretchr/testify v1.8.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
)

replace github.com/go-kratos/kratos/v2 => ../../../
