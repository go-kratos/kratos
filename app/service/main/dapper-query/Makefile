PROTO_FILE=api/v1/api.proto
VENDOR_DIR=$(shell echo $$GOPATH | cut -d':' -f1)/src/go-common/vendor
GO_FILES=$(shell find . -name "*.go" -type f)

rundev: build
	./cmd/cmd -conf cmd/dapper-query-example.toml -http tcp://127.0.0.1:8000\?timeout=10s -http.perf tcp://127.0.0.1:6006 -log.v 10 -log.stdout

genapi: $(PROTO_FILE)
	protoc -I. -I$(VENDOR_DIR) --gogo_out=plugins=grpc:. $<
	protoc -I. -I$(VENDOR_DIR) --bm_out=logtostderr=1,jsonpb=true:. $<
	protoc -I. -I$(VENDOR_DIR) --swagger_out=logtostderr=1:. $<
	swagger-markdown -i api/v1/api.swagger.json -o api/v1/README.md


build: $(GO_FILES) genapi
	go build -o cmd/cmd cmd/main.go

clean:
	rm -rf bin/*

.PHONY: genapi build clean rundev
