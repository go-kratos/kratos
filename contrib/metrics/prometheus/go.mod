module github.com/go-kratos/kratos/contrib/metrics/prometheus/v2

go 1.19

require (
	github.com/go-kratos/kratos/v2 v2.7.3
	github.com/prometheus/client_golang v1.18.0
	github.com/prometheus/common v0.46.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)

replace github.com/go-kratos/kratos/v2 => ../../../
