.PHONY: all

doc:
	godoc `pwd`

webdoc:
	godoc -http=:44444

format: 
	go fmt

test:
	go test -race 
	
check: format test

benchmark:
	go test -bench=. -benchmem

coverage:
	go test -coverprofile=coverage.out
	go tool cover -html="coverage.out"

lint: format
	go get github.com/alecthomas/gometalinter
	gometalinter --install
	gometalinter 

default: lint test
