module github.com/go-kratos/kratos/contrib/registry/nacos2/v2

go 1.16

require (
	github.com/go-kratos/kratos/v2 v2.1.5
	github.com/nacos-group/nacos-sdk-go/v2 v2.0.2
)

replace (
	github.com/buger/jsonparser => github.com/buger/jsonparser v1.1.1
	github.com/go-kratos/kratos/v2 => ../../../
)
