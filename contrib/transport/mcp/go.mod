module github.com/go-kratos/kratos/contrib/transport/mcp/v2

go 1.23

toolchain go1.24.6

require (
	github.com/go-kratos/kratos/v2 v2.9.0
	github.com/mark3labs/mcp-go v0.23.0
)

require (
	github.com/go-playground/form/v4 v4.2.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/go-kratos/kratos/v2 => ../../../
