module github.com/go-kratos/kratos/contrib/metrics/prometheus/v2

go 1.16

require (
	github.com/go-kratos/kratos/v2 v2.6.1
	github.com/prometheus/client_golang v1.15.0
	github.com/prometheus/common v0.42.0
)

replace github.com/go-kratos/kratos/v2 => ../../../
