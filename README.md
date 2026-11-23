<p align="center"><a href="https://go-kratos.dev/" target="_blank"><img src="https://github.com/go-kratos/kratos/blob/main/docs/images/kratos-large.png?raw=true"></a></p>

<p align="center">
<a href="https://github.com/go-kratos/kratos/actions"><img src="https://github.com/go-kratos/kratos/workflows/Go/badge.svg" alt="Build Status"></a>
<a href="https://pkg.go.dev/github.com/go-kratos/kratos/v3"><img src="https://pkg.go.dev/badge/github.com/go-kratos/kratos/v3" alt="GoDoc"></a>
<a href="https://deepwiki.com/go-kratos/kratos"><img src="https://img.shields.io/badge/DeepWiki-go--kratos%2Fkratos-blue.svg" alt="DeepWiki"></a>
<a href="https://codecov.io/gh/go-kratos/kratos"><img src="https://codecov.io/gh/go-kratos/kratos/master/graph/badge.svg" alt="codeCov"></a>
<a href="https://goreportcard.com/report/github.com/go-kratos/kratos"><img src="https://goreportcard.com/badge/github.com/go-kratos/kratos" alt="Go Report Card"></a>
<a href="https://github.com/go-kratos/kratos/blob/main/LICENSE"><img src="https://img.shields.io/github/license/go-kratos/kratos" alt="License"></a>
<a href="https://github.com/avelino/awesome-go"><img src="https://awesome.re/mentioned-badge.svg" alt="Awesome Go"></a>
<a href="https://discord.gg/BWzJsUJ"><img src="https://img.shields.io/discord/766619759214854164?label=chat&logo=discord" alt="Discord"></a>
</p>

<p align="center">
<a href="https://trendshift.io/repositories/3233" target="_blank"><img src="https://trendshift.io/api/badge/repositories/3233" alt="go-kratos%2Fkratos | Trendshift" style="width: 250px; height: 55px;" width="250" height="55"></a>
<a href="https://www.producthunt.com/posts/go-kratos?utm_source=badge-featured&utm_medium=badge&utm_souce=badge-go-kratos" target="_blank"><img src="https://api.producthunt.com/widgets/embed-image/v1/featured.svg?post_id=306565&theme=light" alt="Go Kratos - A Go framework for microservices. | Product Hunt" style="width: 250px; height: 54px;" width="250" height="54"></a>
</p>

Translations: [English](README.md) | [简体中文](README_zh.md)

# Kratos

Kratos is a lightweight Go framework for building cloud-native microservices. It provides small, explicit APIs for transport, middleware, registry, configuration, logging, encoding, and code generation so applications can focus on business logic.

## Features

- API-first development with Protobuf and generated HTTP/gRPC code.
- Unified [transport](https://go-kratos.dev/docs/component/transport/overview) layer for [HTTP](https://go-kratos.dev/docs/component/transport/http) and [gRPC](https://go-kratos.dev/docs/component/transport/grpc).
- Composable [middleware](https://go-kratos.dev/docs/component/middleware/overview) for recovery, logging, validation, tracing, metrics, auth, and more.
- Pluggable [registry](https://go-kratos.dev/docs/component/registry), [configuration](https://go-kratos.dev/docs/component/config), and [encoding](https://go-kratos.dev/docs/component/encoding) components.
- Standard-library `log/slog` based logging with OpenTelemetry extensions in contrib packages.
- Consistent metadata, errors, validation, OpenAPI, and code-generation workflows.
- A contrib ecosystem for optional integrations such as registries, config stores, middleware, encodings, and observability.

## Installation

### Requirements

- [Go](https://go.dev/dl/) 1.25 or later
- [protoc](https://github.com/protocolbuffers/protobuf)
- [protoc-gen-go](https://github.com/protocolbuffers/protobuf-go)

### Install the CLI

```shell
go install github.com/go-kratos/kratos/cmd/kratos/v3@latest
kratos upgrade
```

## Create a Service

```shell
kratos new helloworld
cd helloworld
go mod tidy
kratos run
```

Visit `http://localhost:8000/helloworld/kratos` after the service starts.

For a fuller generated service flow:

```shell
kratos proto add api/helloworld/helloworld.proto
kratos proto client api/helloworld/helloworld.proto
kratos proto server api/helloworld/helloworld.proto -t internal/service
go generate ./...
kratos run
```

## Usage Example

```go
package main

import (
	"github.com/go-kratos/kratos/v3"
	"github.com/go-kratos/kratos/v3/transport/grpc"
	"github.com/go-kratos/kratos/v3/transport/http"
)

func main() {
	httpSrv := http.NewServer(http.Address(":8000"))
	grpcSrv := grpc.NewServer(grpc.Address(":9000"))

	app := kratos.New(
		kratos.Name("helloworld"),
		kratos.Version("v1.0.0"),
		kratos.Server(httpSrv, grpcSrv),
	)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
```

## v3 Migration

Kratos v3 reduces core dependencies and makes previously implicit behavior explicit. Review the [v2 to v3 migration guide](docs/migration/v2-to-v3.md) before upgrading production services.

## Further Reading

- [Documentation](https://go-kratos.dev/docs/getting-started/start)
- [Examples](https://github.com/go-kratos/examples)
- [Project Layout](https://github.com/go-kratos/kratos-layout)
- [v2 to v3 Migration Guide](docs/migration/v2-to-v3.md)
- [Community Contribution Guide](https://go-kratos.dev/docs/community/contribution)

## Development

```shell
make test
make lint
```

## Community

- [Documentation](https://go-kratos.dev/en)
- [WeChat Group](https://github.com/go-kratos/kratos/issues/682)
- [Discord Group](https://discord.gg/BWzJsUJ)
- [Discussions](https://github.com/go-kratos/kratos/discussions)

## Security

If you discover a security vulnerability in Kratos, please contact go-kratos@googlegroups.com. Security reports are handled privately before disclosure.

## Contributors

Thank you for contributing to Kratos. The contribution guide is available in the [Kratos documentation](https://go-kratos.dev/docs/community/contribution).

<a href="https://github.com/go-kratos/kratos/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=go-kratos/kratos" alt="Kratos contributors">
</a>

## Acknowledgments

The following projects influenced Kratos design:

- [go-kit/kit](https://github.com/go-kit/kit)
- [go-micro](https://github.com/asim/go-micro)
- [google/go-cloud](https://github.com/google/go-cloud)
- [go-zero](https://github.com/zeromicro/go-zero)
- [beego](https://github.com/beego/beego)

## License

Kratos is open-sourced software licensed under the [MIT license](./LICENSE).
