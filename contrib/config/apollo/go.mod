module github.com/go-kratos/kratos/contrib/config/apollo/v2

go 1.16

require (
	github.com/apolloconfig/agollo/v4 v4.1.0
	github.com/go-kratos/kratos/v2 v2.2.1
)

require (
	github.com/spf13/afero v1.8.0 // indirect
	github.com/spf13/viper v1.10.1 // indirect
	golang.org/x/sys v0.0.0-20220114195835-da31bd327af9 // indirect
	gopkg.in/ini.v1 v1.66.3 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace github.com/go-kratos/kratos/v2 => ../../../
