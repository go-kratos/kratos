module github.com/go-kratos/kratos/contrib/log/tencent/v2

go 1.19

require (
	github.com/go-kratos/kratos/v2 v2.7.2
	github.com/tencentcloud/tencentcloud-cls-sdk-go v1.0.2
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/klauspost/compress v1.15.1 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	go.uber.org/atomic v1.9.0 // indirect
)

replace github.com/go-kratos/kratos/v2 => ../../../
