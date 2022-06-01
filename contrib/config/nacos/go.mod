module github.com/go-kratos/kratos/contrib/config/nacos/v2

go 1.16

require (
	github.com/go-kratos/kratos/v2 v2.3.0
	github.com/nacos-group/nacos-sdk-go v1.0.9
	gopkg.in/yaml.v3 v3.0.1
)

replace github.com/go-kratos/kratos/v2 => ../../../

replace github.com/buger/jsonparser => github.com/buger/jsonparser v1.1.1
