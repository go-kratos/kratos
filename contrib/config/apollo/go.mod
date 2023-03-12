module github.com/go-kratos/kratos/contrib/config/apollo/v2

go 1.16

require (
	github.com/apolloconfig/agollo/v4 v4.3.0
	github.com/go-kratos/kratos/v2 v2.4.0
)

require (
	github.com/spf13/viper v1.11.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/go-kratos/kratos/v2 => ../../../
