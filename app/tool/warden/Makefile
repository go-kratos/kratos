generator/templates/assets.go: $(shell find generator/templates -name "*.tmpl")
	cd generator/templates/; go generate

build: generator/templates/assets.go
	go build -o warden

install: generator/templates/assets.go
	go install go-common/app/tool/warden

clean:
	rm warden
