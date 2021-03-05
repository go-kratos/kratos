# Kratos Layout

## Install Kratos
```
go get github.com/go-kratos/kratos/cmd/kratos
go get github.com/go-kratos/kratos/cmd/protoc-gen-go-http
go get github.com/go-kratos/kratos/cmd/protoc-gen-go-errors

# from source
cd cmd/kratos && go install
cd cmd/protoc-gen-go-http && go install
cd cmd/protoc-gen-go-errors && go install
```
## Create a service
```
# create a template project
kratos new helloworld

cd helloworld
# Add a proto template
kratos proto add api/helloworld/helloworld.proto
# Generate the source code of service by proto file
kratos proto service api/helloworld/helloworld.proto -t internal/service

make proto
make build
make test
```
