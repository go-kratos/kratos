regenerate:
	go install go-common/vendor/github.com/golang/protobuf/protoc-gen-go
	protoc --go_out=paths=source_relative:. gogo.proto

restore:
	cp gogo.pb.golden gogo.pb.go

preserve:
	cp gogo.pb.go gogo.pb.golden
