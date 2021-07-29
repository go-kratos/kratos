package echo

//go:generate protoc --proto_path=. --proto_path=../../third_party --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. --go-http_out=paths=source_relative:. echo_service.proto stream.proto
