user	:=	$(shell whoami)
rev 	:= 	$(shell git rev-parse --short HEAD)

# GOBIN > GOPATH > INSTALLDIR
GOBIN	:=	$(shell echo ${GOBIN} | cut -d':' -f1)
GOPATH	:=	$(shell echo $(GOPATH) | cut -d':' -f1)
BIN		:= 	""

TOOLS_SHELL="./hack/tools.sh"
# golangci-lint
LINTER := bin/golangci-lint

# check GOBIN
ifneq ($(GOBIN),)
	BIN=$(GOBIN)
else
	# check GOPATH
	ifneq ($(GOPATH),)
		BIN=$(GOPATH)/bin
	endif
endif

$(LINTER): 
	curl -SL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s latest

all:
	@cd cmd/kratos && go build && cd - &> /dev/null
	@cd cmd/protoc-gen-go-errors && go build && cd - &> /dev/null
	@cd cmd/protoc-gen-go-http && go build && cd - &> /dev/null

.PHONY: install
install: all
ifeq ($(user),root)
#root, install for all user
	@cp ./cmd/kratos/kratos /usr/bin
	@cp ./cmd/protoc-gen-go-errors/protoc-gen-go-errors /usr/bin
	@cp ./cmd/protoc-gen-go-http/protoc-gen-go-http /usr/bin
else
#!root, install for current user
	$(shell if [ -z $(BIN) ]; then read -p "Please select installdir: " REPLY; mkdir -p $${REPLY};\
	cp ./cmd/kratos/kratos $${REPLY}/;cp ./cmd/protoc-gen-go-errors/protoc-gen-go-errors $${REPLY}/;cp ./cmd/protoc-gen-go-http/protoc-gen-go-http $${REPLY}/;else mkdir -p $(BIN);\
	cp ./cmd/kratos/kratos $(BIN);cp ./cmd/protoc-gen-go-errors/protoc-gen-go-errors $(BIN);cp ./cmd/protoc-gen-go-http/protoc-gen-go-http $(BIN); fi)
endif
	@which protoc-gen-go &> /dev/null || go get google.golang.org/protobuf/cmd/protoc-gen-go
	@which protoc-gen-go-grpc &> /dev/null || go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
	@which protoc-gen-validate  &> /dev/null || go get github.com/envoyproxy/protoc-gen-validate
	@echo "install finished"

.PHONY: uninstall
uninstall:
	$(shell for i in `which -a kratos | grep -v '/usr/bin/kratos' 2>/dev/null | sort | uniq`; do read -p "Press to remove $${i} (y/n): " REPLY; if [ $${REPLY} = "y" ]; then rm -f $${i}; fi; done)
	$(shell for i in `which -a protoc-gen-go-grpc | grep -v '/usr/bin/protoc-gen-go-errors' 2>/dev/null | sort | uniq`; do read -p "Press to remove $${i} (y/n): " REPLY; if [ $${REPLY} = "y" ]; then rm -f $${i}; fi; done)
	$(shell for i in `which -a protoc-gen-validate | grep -v '/usr/bin/protoc-gen-go-errors' 2>/dev/null | sort | uniq`; do read -p "Press to remove $${i} (y/n): " REPLY; if [ $${REPLY} = "y" ]; then rm -f $${i}; fi; done)
	@echo "uninstall finished"

.PHONY: clean
clean:
	@${TOOLS_SHELL} tidy
	@echo "clean finished"

.PHONY: fix
fix: $(LINTER)
	@${TOOLS_SHELL} fix
	@echo "lint fix finished"

.PHONY: test
test:
	@${TOOLS_SHELL} test
	@echo "go test finished"

.PHONY: test-coverage
test-coverage:
	@${TOOLS_SHELL} test_coverage
	@echo "go test with coverage finished"	

.PHONY: lint
lint: $(LINTER)
	@${TOOLS_SHELL} lint
	@echo "lint check finished"