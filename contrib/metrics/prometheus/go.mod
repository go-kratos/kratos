module github.com/go-kratos/kratos/contrib/metrics/prometheus/v2

go 1.16

require (
	github.com/go-kratos/kratos/v2 v2.5.0
	github.com/prometheus/client_golang v1.12.2
	github.com/prometheus/common v0.37.0
)

replace github.com/go-kratos/kratos/v2 => ../../../
