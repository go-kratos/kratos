module github.com/go-kratos/kratos/contrib/log/tencent/v2

go 1.24

require (
	github.com/tencentcloud/tencentcloud-cls-sdk-go v1.0.14
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/go-kratos/kratos/v2 v2.9.2
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/klauspost/compress v1.18.6 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	go.uber.org/atomic v1.11.0 // indirect
)

replace github.com/go-kratos/kratos/v2 => ../../..
