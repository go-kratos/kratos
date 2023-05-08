module github.com/go-kratos/kratos/contrib/log/aliyun/v2

go 1.16

require (
	github.com/aliyun/aliyun-log-go-sdk v0.1.44
	github.com/go-kratos/kratos/v2 v2.6.1
	google.golang.org/protobuf v1.30.0
)

replace (
	github.com/go-kratos/kratos/v2 => ../../../
	github.com/gogo/protobuf v1.3.1 => github.com/gogo/protobuf v1.3.2
)
