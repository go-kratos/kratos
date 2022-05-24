module github.com/go-kratos/kratos/contrib/opensergo/v2

go 1.17

require (
	github.com/go-kratos/kratos/v2 v2.3.0
	github.com/opensergo/opensergo-go v0.0.0-20220331070310-e5b01fee4d1c
	golang.org/x/net v0.0.0-20220520000938-2e3eb7b945c2
	google.golang.org/genproto v0.0.0-20220519153652-3a47de7e79bd
	google.golang.org/grpc v1.46.2
	google.golang.org/protobuf v1.28.0
)

replace github.com/go-kratos/kratos/v2 => ../../
