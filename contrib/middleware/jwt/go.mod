module github.com/go-kratos/kratos/contrib/middleware/jwt/v3

go 1.25.0

require (
	github.com/go-kratos/kratos/v3 v3.0.0
	github.com/golang-jwt/jwt/v5 v5.3.1
)

require (
	github.com/go-playground/form/v4 v4.3.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	golang.org/x/sys v0.44.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260511170946-3700d4141b60 // indirect
	google.golang.org/grpc v1.81.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/go-kratos/kratos/v3 => ../../..
