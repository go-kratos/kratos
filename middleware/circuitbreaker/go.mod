module github.com/go-kratos/kratos/middleware/circuitbreaker/v2

go 1.17

require (
	github.com/go-kratos/kratos/v2 v2.0.5
	github.com/go-kratos/sra v0.0.0-20210905065551-b690f7ef1d3e
	github.com/golang/protobuf v1.5.2 // indirect
	google.golang.org/genproto v0.0.0-20210805201207-89edb61ffb67 // indirect
	google.golang.org/grpc v1.39.1 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)

replace github.com/go-kratos/kratos/v2 => ../../../kratos
