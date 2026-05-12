module github.com/go-kratos/kratos/contrib/log/zap/v3

go 1.22

require (
	github.com/go-kratos/kratos/v3 v3.0.0
	go.uber.org/zap v1.26.0
)

require go.uber.org/multierr v1.11.0 // indirect

replace github.com/go-kratos/kratos/v3 => ../../../
