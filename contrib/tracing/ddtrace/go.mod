module github.com/go-kratos/kratos/contrib/tracing/ddtrace/v2

go 1.16

require (
	github.com/go-kratos/kratos/v2 v2.3.1
	gopkg.in/DataDog/dd-trace-go.v1 v1.38.1
)

replace github.com/go-kratos/kratos/v2 => ../../../
