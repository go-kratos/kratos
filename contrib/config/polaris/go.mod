module github.com/go-kratos/kratos/contrib/config/polaris/v2

go 1.16

require (
	github.com/go-kratos/kratos/v2 v2.6.2
	github.com/polarismesh/polaris-go v1.1.0
)

require (
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/prometheus/client_golang v1.12.2 // indirect
	github.com/prometheus/common v0.35.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
)

replace github.com/go-kratos/kratos/v2 => ../../../
