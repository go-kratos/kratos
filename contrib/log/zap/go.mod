module github.com/go-kratos/kratos/contrib/log/zap/v2

go 1.19

require (
	github.com/go-kratos/kratos/v2 v2.7.1
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.8.1
	go.uber.org/zap v1.26.0
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/go-kratos/kratos/v2 => ../../../
