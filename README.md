![kratos](docs/images/kratos.png) 

[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![Build Status](https://github.com/go-kratos/kratos/workflows/Go/badge.svg)](https://github.com/go-kratos/kratos/actions)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-kratos/kratos/v2)](https://pkg.go.dev/github.com/go-kratos/kratos/v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-kratos/kratos/v2)](https://goreportcard.com/report/github.com/go-kratos/kratos/v2)
[![Discord](https://img.shields.io/discord/766619759214854164?label=chat&logo=discord)](https://discord.gg/BWzJsUJ)

Translations: [English](README.md) | [简体中文](README_zh.md)

# Kratos
Kratos is a microservice-oriented governance framework implements by golang, which offers convenient capabilities to help you quickly build a bulletproof application from scratch.


>The name is inspired by the game God of War which is based on Greek myths, tells the Kratos from mortals to become a God of War and launches the adventure of killing god.


## Goals

Kratos boosts your productivity. With the integration of excellent resources and further support, you can get rid of most issues you might encounter in the field of distributed systems and software engineering, so that you can focus on business delivery. For each developer, Kratos is an ideal one for learning as well. You can navigate the diverse knowledge of microservice and accumulate engineering experience.

### Principles

* **Simple**: Appropriate design, plain and simple code.
* **General**: Cover the various utilities for business development.
* **Reliable**: Higher testability and greater test coverage of base libs validated in the production environment.
* **Robust**: Base libs are designed reasonably to reduce misuse.
* **High-performance**: We give you optimal performance without using the hack-way approach like adding the *unsafe* package and guarantee compatibility and stability at the same time.
* **Expandable**: Appropriate API design, you can expand utilities such as base libs to meet your further requirements.
* **Fault-tolerance**: Designed against failure, enhance the understanding and exercising of SRE within Kratos to achieve more robustness.
* **Toolchain**: Includes an extensive toolchain, such as cache code generation, lint tools, etc.

## Features
* APIs: HTTP/gRPC based transport and Protobuf defined communication protocol.
* Errors: We use ProtoBuf Enum to define error code and generate error handle code.
* Metadata: Normalize the service metadata transmission by excellent middleware.
* Config:  Multi-data source Support, well arranged, dynamic configuration via *Atomic* package.
* Logger: Standard log API, easily integrate with third-party log lib, *Fluentd* logs collection.
* Metrics: *Prometheus* integration by default. Furthermore, with the unified Metrics interface, you can implement your metrics system more flexible
* Tracing: Complete micro-service link tracing followed by *OpenTracing* specification.
* Encoding: Support *Accept* and *Content-Type* for auto content encoding.
* Transport: Common HTTP/ GRPC transport layer offers you powerful Middleware support.
 * Registry: One pluggable API for the different registry.

## Getting Started
### Required
- [go](https://golang.org/dl/)
- [protoc](https://github.com/protocolbuffers/protobuf)
- [protoc-gen-go](https://github.com/protocolbuffers/protobuf-go)

### Installing
```
go get github.com/go-kratos/kratos/cmd/kratos/v2@latest
```
### Create a service
```
# create project template
kratos new helloworld

cd helloworld
# generate Proto template
kratos proto add api/helloworld/helloworld.proto
# generate Proto source code
kratos proto client api/helloworld/helloworld.proto
# generate server template
kratos proto server api/helloworld/helloworld.proto -t internal/service

# Generate all proto source code, wire, etc.
go generate ./...
# compile
go build -o ./bin/ ./...
# run
./bin/helloworld -conf ./configs
```

### Kratos Boot
```
import "github.com/go-kratos/kratos/v2"
import "github.com/go-kratos/kratos/v2/transport/grpc"
import "github.com/go-kratos/kratos/v2/transport/http"

httpSrv := http.NewServer(http.Address(":8000"))
grpcSrv := grpc.NewServer(grpc.Address(":9000"))

app := kratos.New(
    kratos.Name("kratos"),
    kratos.Version("latest"),
    kratos.Server(httpSrv, grpcSrv),
)
app.Run()
```

## Related

* [Docs](https://go-kratos.dev/)
* [Examples](https://github.com/go-kratos/examples)
* [Service Layout](https://github.com/go-kratos/kratos-layout)

## Community
* [Wechat Group](https://github.com/go-kratos/kratos/issues/682)
* [Discord Group](https://discord.gg/BWzJsUJ)
* QQ Group: 716486124

## License
Kratos is MIT licensed. See the [LICENSE](./LICENSE) file for details.
