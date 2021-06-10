GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
API_PROTO_FILES=$(shell find api -name *.proto)
KRATOS_VERSION=$(shell go mod graph |grep go-kratos/kratos/v2 |head -n 1 |awk -F '@' '{print $$2}')
KRATOS=$(GOPATH)/pkg/mod/github.com/go-kratos/kratos/v2@$(KRATOS_VERSION)

.PHONY: init
# init env
init:
	go get -u github.com/go-kratos/kratos/cmd/kratos/v2
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2
	go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	go get -u github.com/google/wire/cmd/wire

.PHONY: grpc
# generate grpc code
grpc:
	protoc --proto_path=. \
           --proto_path=$(KRATOS)/third_party \
           --go_out=paths=source_relative:. \
           --go-grpc_out=paths=source_relative:. \
           $(API_PROTO_FILES)

.PHONY: http
# generate http code
http:
	protoc --proto_path=. \
           --proto_path=$(KRATOS)/third_party \
           --go_out=paths=source_relative:. \
           --go-http_out=paths=source_relative:. \
           $(API_PROTO_FILES)

.PHONY: errors
# generate errors code
errors:
	protoc --proto_path=. \
           --proto_path=./third_party \
           --go_out=paths=source_relative:. \
           --go-errors_out=paths=source_relative:. \
           $(API_PROTO_FILES)

.PHONY: proto
# generate internal proto
proto:
	protoc --proto_path=. \
           --proto_path=$(KRATOS)/third_party \
           --go_out=paths=source_relative:. \
           $(INTERNAL_PROTO_FILES)

.PHONY: swagger
# generate swagger file
swagger:
	protoc --proto_path=. \
		--proto_path=$(KRATOS)/third_party \
		--openapiv2_out . \
		--openapiv2_opt logtostderr=true \
		$(API_PROTO_FILES)
		
.PHONY: generate
# generate client code
generate:
	go generate ./...

.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: test
# test
test:
	go test -v ./... -cover

.PHONY: all
# generate all
all:
	make generate;
	make grpc;
	make http;
	make proto;
	make errors;
	make swagger;
	make build;
	make test;

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
