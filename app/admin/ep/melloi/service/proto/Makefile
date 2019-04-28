SRCS := $(shell find . -name '*.go')

BIN_DEPS := \
	github.com/golang/lint/golint \
	github.com/kisielk/errcheck \
	honnef.co/go/tools/cmd/staticcheck \
	honnef.co/go/tools/cmd/unused

.PHONY: all
all: test

.PHONY: deps
deps:
	go get -d -v ./...

.PHONY: updatedeps
updatedeps:
	go get -d -v -u -f ./...

.PHONY: bindeps
bindeps:
	go get -v $(BIN_DEPS)

.PHONY: updatebindeps
updatebindeps:
	go get -u -v $(BIN_DEPS)

.PHONY: testdeps
testdeps: bindeps
	go get -d -v -t ./...

.PHONY: updatetestdeps
updatetestdeps: updatebindeps
	go get -d -v -t -u -f ./...

.PHONY: install
install: deps
	go install ./...

.PHONY: golint
golint: testdeps
	@# TODO: readd cmd/proto2gql when fixed
	@#for file in $(SRCS); do
	for file in $(shell echo $(SRCS) | grep -v cmd/proto2gql); do \
		golint $${file}; \
		if [ -n "$$(golint $${file})" ]; then \
			exit 1; \
		fi; \
	done

.PHONY: vet
vet: testdeps
	go vet ./...

.PHONY: testdeps
errcheck: testdeps
	errcheck ./...

.PHONY: staticcheck
staticcheck: testdeps
	staticcheck ./...

.PHONY: unused
unused: testdeps
	unused ./...

.PHONY: lint
# TODO: readd errcheck and unused when fixed
#lint: golint vet errcheck staticcheck unused
lint: golint vet staticcheck

.PHONY: test
test: testdeps lint
	go test -race ./...

.PHONY: clean
clean:
	go clean -i ./...

integration:
	PB=y go test -cover