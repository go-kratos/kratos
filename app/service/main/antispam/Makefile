print-%: ; @echo $*=$($*)

PREFIX?=/data/app/antispam-service
PROG=antispam-service
PROG_BAK=antispam-service_bak
TEST_SRC=$(wildcard ./**/*_test.go)
GOTESTOPT= -v
SOURCE_DIR=./
SUB_DIR=$(shell ls)

CUR_BRANCH=$(shell git branch | grep \* | cut -d ' ' -f2)
CONF?=cmd/antispam-test.toml

all: build vet lint test_quick clean 

.PHONY: build run lint test vet clean test_quick 

build: $(PROG) 

$(PROG):
	@echo "\nstart build..."
	@go build -o $(PROG) cmd/main.go 

run: clean build 
	@echo "start running in container..."
	@./$(PROG) -conf=$(CONF)

clean:
	@echo "\nstart cleaning up ..."
	@go clean $(GO_FLAGS) $(SOURCE_DIR)
	@rm -rf $(PROG)
	@rm -rf build
	@rm -f `find . -name "*.swp" -print`
	@rm -f `find . -name "*.xlsx" -print`

vet:
	@echo "\nstart running go vet..."
	@for f in $(SUB_DIR); do \
		if [ -d $$f ]; then \
			go vet ./"$$f"/...; \
		fi; \
	done

lint:
	@echo "\nstart running golint..."
	@for f in $(SUB_DIR); do \
		if [ -d $$f ]; then \
			golint $$f; \
		fi; \
	done

test_quick:
	@echo "\nstart running quick unit test..."
	@go test $(GOTESTOPT) ./http/... ./service/... ./util/... 

test:
	@echo "\nstart running full unit test..."
	@go test $(GOTESTOPT) ./dao/...

deploy:
	@echo "start deploying on dev machine..."
	@cp $(PREFIX)/$(PROG) $(PREFIX)/$(PROG_BAK) && go build -o $(PREFIX)/$(PROG) && supervisorctl restart $(PROG)
