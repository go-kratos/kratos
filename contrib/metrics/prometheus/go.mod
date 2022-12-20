module github.com/go-kratos/kratos/contrib/metrics/prometheus/v2

go 1.16

require (
	github.com/go-kratos/kratos/v2 v2.5.3
	github.com/prometheus/client_golang v1.14.0
	github.com/prometheus/common v0.39.0
)

replace github.com/go-kratos/kratos/v2 => ../../../
