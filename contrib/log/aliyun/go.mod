module github.com/go-kratos/kratos/contrib/log/aliyun/v2

go 1.22.0

require (
	github.com/aliyun/aliyun-log-go-sdk v0.1.75
	github.com/go-kratos/kratos/v2 v2.8.3
	google.golang.org/protobuf v1.35.2
)

require (
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/go-kit/kit v0.10.0 // indirect
	github.com/go-logfmt/logfmt v0.5.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/klauspost/compress v1.17.8 // indirect
	github.com/pierrec/lz4 v2.6.0+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.uber.org/atomic v1.5.0 // indirect
	golang.org/x/lint v0.0.0-20190930215403-16217165b5de // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

replace (
	github.com/go-kratos/kratos/v2 => ../../../
	github.com/gogo/protobuf v1.3.1 => github.com/gogo/protobuf v1.3.2
)
