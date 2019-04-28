GOFILES=$(shell find . -name "*.go" -type f)
EXAMPLE_HELLOWORLD_DIR=examples/helloworld
PROTO_FILE=$(EXAMPLE_HELLOWORLD_DIR)/api/helloworld.proto
VENDOR_DIR=$(shell echo $$GOPATH | cut -d':' -f1)/src/go-common/vendor

protoc-gen-bm: $(GOFILES)
	go install .

protoc-gen-go: $(PROTO_FILE)
	protoc -I. -I$(VENDOR_DIR) --go_out=plugins=grpc:. $<

examples-helloworld: $(PROTO_FILE) protoc-gen-bm protoc-gen-go
	protoc -I. -I$(VENDOR_DIR) --bm_out=logtostderr=1,v=10:. $<

examples-helloworld-jsonpb: $(PROTO_FILE) protoc-gen-bm protoc-gen-go
	protoc -I. -I$(VENDOR_DIR) --bm_out=logtostderr=1,v=10,jsonpb=true:. $<

test:
	go test ./...

.PHONY: protoc-gen-bm protoc-gen-go examples-helloworld test
