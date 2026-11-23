module github.com/go-kratos/kratos/contrib/middleware/validate/v3

go 1.25.0

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.11-20260415201107-50325440f8f2.1
	github.com/envoyproxy/protoc-gen-validate v1.3.3
	github.com/go-kratos/kratos/v3 v3.0.0
	google.golang.org/protobuf v1.36.11
)

require go.yaml.in/yaml/v3 v3.0.4 // indirect

require (
	buf.build/go/protovalidate v1.2.0
	cel.dev/expr v0.25.2 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/google/cel-go v0.28.1 // indirect
	golang.org/x/exp v0.0.0-20260508232706-74f9aab9d74a // indirect
	golang.org/x/sys v0.44.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260511170946-3700d4141b60 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260511170946-3700d4141b60 // indirect
	google.golang.org/grpc v1.81.0 // indirect
)

replace github.com/go-kratos/kratos/v3 => ../../../
