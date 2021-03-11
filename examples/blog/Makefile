GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
PROTO_FILES=$(shell find . -name *.proto)
KRATOS_VERSION=$(shell go mod graph |grep go-kratos/kratos/v2 |head -n 1 |awk -F '@' '{print $$2}')
KRATOS=$(GOPATH)/pkg/mod/github.com/go-kratos/kratos/v2@$(KRATOS_VERSION)


.PHONY: init
init:
	go get -u github.com/go-kratos/kratos/cmd/kratos/v2
	go get -u github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2
	go get -u github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

.PHONY: proto
proto:
	protoc --proto_path=. \
           --proto_path=$(KRATOS)/api \
           --proto_path=$(KRATOS)/third_party \
           --proto_path=$(GOPATH)/src \
           --go_out=paths=source_relative:. \
           --go-grpc_out=paths=source_relative:. \
           --go-http_out=paths=source_relative:. \
           --go-errors_out=paths=source_relative:. $(PROTO_FILES)

.PHONY: run
run:
	cd cmd/blog/ && go run .

.PHONY: ent
ent:
	cd internal/data/ && ent generate ./ent/schema

.PHONY: generate
generate:
	go generate ./...

.PHONY: build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: test
test:
	go test -v ./... -cover

