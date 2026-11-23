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

Kratos 是一套轻量级 Go 微服务框架，围绕传输、Middleware、注册发现、配置、日志、编码和代码生成提供清晰的基础能力，让业务代码保持聚焦。

## 功能特性

- 以 Protobuf 为中心定义 API，并生成 HTTP/gRPC 代码。
- 统一的 [Transport](https://go-kratos.dev/zh-cn/docs/component/transport/overview) 抽象，支持 [HTTP](https://go-kratos.dev/zh-cn/docs/component/transport/http) 和 [gRPC](https://go-kratos.dev/zh-cn/docs/component/transport/grpc)。
- 可组合的 [Middleware](https://go-kratos.dev/zh-cn/docs/component/middleware/overview)，覆盖 Recovery、Logging、Validation、Tracing、Metrics、Auth 等场景。
- 插件化的 [Registry](https://go-kratos.dev/zh-cn/docs/component/registry)、[Config](https://go-kratos.dev/zh-cn/docs/component/config) 和 [Encoding](https://go-kratos.dev/zh-cn/docs/component/encoding) 能力。
- 基于标准库 `log/slog` 的日志能力，OpenTelemetry 扩展由 contrib 包提供。
- 统一的 Metadata、Errors、Validation、OpenAPI 和代码生成工作流。
- contrib 生态提供注册中心、配置中心、Middleware、编码和可观测性等可选集成。

## 安装

### 环境要求

- [Go](https://go.dev/dl/) 1.25 或更高版本
- [protoc](https://github.com/protocolbuffers/protobuf)
- [protoc-gen-go](https://github.com/protocolbuffers/protobuf-go)

### 安装 CLI

```shell
go install github.com/go-kratos/kratos/cmd/kratos/v3@latest
kratos upgrade
```

## 创建服务

```shell
kratos new helloworld
cd helloworld
go mod tidy
kratos run
```

服务启动后访问 `http://localhost:8000/helloworld/kratos`。

如果需要从 proto 开始生成服务代码：

```shell
kratos proto add api/helloworld/helloworld.proto
kratos proto client api/helloworld/helloworld.proto
kratos proto server api/helloworld/helloworld.proto -t internal/service
go generate ./...
kratos run
```

## 使用示例

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

## v3 迁移

Kratos v3 进一步降低核心依赖，并将历史上隐式的行为显式化。升级生产服务前请先阅读 [v2 到 v3 迁移指南](docs/migration/v2-to-v3_zh.md)。

## 扩展阅读

- [文档](https://go-kratos.dev/zh-cn/docs/getting-started/start)
- [示例](https://github.com/go-kratos/examples)
- [项目模板](https://github.com/go-kratos/kratos-layout)
- [v2 到 v3 迁移指南](docs/migration/v2-to-v3_zh.md)
- [贡献指南](https://go-kratos.dev/zh-cn/docs/community/contribution)

## 开发

```shell
make test
make lint
```

## 社区

- [文档站点](https://go-kratos.dev/zh-cn)
- [微信群](https://github.com/go-kratos/kratos/issues/682)
- [Discord](https://discord.gg/BWzJsUJ)
- [GitHub Discussions](https://github.com/go-kratos/kratos/discussions)
- QQ 群：716486124

## 安全

如果你发现 Kratos 存在安全漏洞，请发送邮件到 go-kratos@googlegroups.com。安全问题会在公开披露前以私密方式处理。

## 贡献者

感谢所有为 Kratos 做出贡献的开发者。贡献流程请参考 [Kratos 贡献指南](https://go-kratos.dev/zh-cn/docs/community/contribution)。

<a href="https://github.com/go-kratos/kratos/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=go-kratos/kratos" alt="Kratos contributors">
</a>

## 致谢

以下项目对 Kratos 的设计有重要影响：

- [go-kit/kit](https://github.com/go-kit/kit)
- [go-micro](https://github.com/asim/go-micro)
- [google/go-cloud](https://github.com/google/go-cloud)
- [go-zero](https://github.com/zeromicro/go-zero)
- [beego](https://github.com/beego/beego)

## License

Kratos 基于 [MIT license](./LICENSE) 开源。
