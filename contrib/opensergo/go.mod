module github.com/go-kratos/kratos/contrib/opensergo/v2

go 1.22.0

require (
	github.com/go-kratos/kratos/v2 v2.8.3
	github.com/opensergo/opensergo-go v0.0.0-20220331070310-e5b01fee4d1c
	golang.org/x/net v0.33.0
	google.golang.org/genproto/googleapis/api v0.0.0-20240401170217-c3f982113cda
	google.golang.org/grpc v1.63.2
	google.golang.org/protobuf v1.35.2
)

require (
	github.com/go-playground/form/v4 v4.2.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240401170217-c3f982113cda // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/go-kratos/kratos/v2 => ../../
