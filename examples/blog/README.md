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
# create project template
kratos new blog

cd helloworld
# download modules
go mod download

# generate Proto template
kratos proto add api/blog/blog.proto
# generate Proto source code
kratos proto client api/blog/blog.proto
# generate server template
kratos proto server api/blog/blog.proto -t internal/service、

# generate all proto source code, wire, etc.
go generate ./...
```
