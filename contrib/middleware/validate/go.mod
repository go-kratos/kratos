module github.com/go-kratos/kratos/contrib/middleware/validate/v2

go 1.23.0

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.5-20250307204501-0409229c3780.1
	github.com/bufbuild/protovalidate-go v0.9.2
	github.com/envoyproxy/protoc-gen-validate v1.2.1
	github.com/go-kratos/kratos/v2 v2.8.4
	google.golang.org/protobuf v1.36.5
)

require (
	cel.dev/expr v0.22.0 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/google/cel-go v0.24.1 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	golang.org/x/exp v0.0.0-20250305212735-054e65f0b394 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250313205543-e70fdf4c4cb4 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250313205543-e70fdf4c4cb4 // indirect
	google.golang.org/grpc v1.67.1 // indirect
)

replace github.com/go-kratos/kratos/v2 => ../../../
