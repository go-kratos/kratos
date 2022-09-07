module github.com/go-kratos/kratos/contrib/opensergo/v2

go 1.17

require (
	github.com/go-kratos/kratos/v2 v2.5.0
	github.com/opensergo/opensergo-go v0.0.0-20220331070310-e5b01fee4d1c
	golang.org/x/net v0.0.0-20220520000938-2e3eb7b945c2
	google.golang.org/genproto v0.0.0-20220519153652-3a47de7e79bd
	google.golang.org/grpc v1.46.2
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/go-playground/form/v4 v4.2.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	golang.org/x/sync v0.0.0-20220513210516-0976fa681c29 // indirect
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/yaml.v3 v3.0.0 // indirect
)

replace github.com/go-kratos/kratos/v2 => ../../
