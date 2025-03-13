module github.com/go-kratos/kratos/contrib/log/aliyun/v2

go 1.23.0

toolchain go1.23.1

require (
	github.com/aliyun/aliyun-log-go-sdk v0.1.97
	github.com/go-kratos/kratos/v2 v2.8.4
	google.golang.org/protobuf v1.36.5
)

require (
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/go-kit/kit v0.13.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/pierrec/lz4 v2.6.0+incompatible // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/lint v0.0.0-20190930215403-16217165b5de // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/tools v0.0.0-20210106214847-113979e3529a // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)

replace (
	github.com/go-kratos/kratos/v2 => ../../../
	github.com/gogo/protobuf v1.3.1 => github.com/gogo/protobuf v1.3.2
)
