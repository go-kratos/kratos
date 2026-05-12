module github.com/go-kratos/kratos/contrib/log/fluent/v2

go 1.24

require github.com/fluent/fluent-logger-golang v1.10.1

require (
	github.com/go-kratos/kratos/v2 v2.9.2
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/tinylib/msgp v1.6.4 // indirect
)

replace github.com/go-kratos/kratos/v2 => ../../..
