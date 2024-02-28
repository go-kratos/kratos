module github.com/go-kratos/kratos/contrib/metrics/prometheus/v2

go 1.19

require (
	github.com/go-kratos/kratos/v2 v2.7.2
	github.com/prometheus/client_golang v1.17.0
	github.com/prometheus/common v0.45.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/matttproud/golang_protobuf_extensions/v2 v2.0.0 // indirect
	github.com/prometheus/client_model v0.4.1-0.20230718164431-9a2bf3000d16 // indirect
	github.com/prometheus/procfs v0.11.1 // indirect
	golang.org/x/sys v0.13.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)

replace github.com/go-kratos/kratos/v2 => ../../../
