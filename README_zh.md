<p align="center"><a href="https://go-kratos.dev/" target="_blank"><img src="https://github.com/go-kratos/kratos/blob/main/docs/images/kratos-large.png?raw=true"></a></p>

<p align="center">
<a href="https://github.com/go-kratos/kratos/actions"><img src="https://github.com/go-kratos/kratos/workflows/Go/badge.svg" alt="Build Status"></a>
<a href="https://pkg.go.dev/github.com/go-kratos/kratos/v2"><img src="https://pkg.go.dev/badge/github.com/go-kratos/kratos/v2" alt="GoDoc"></a>
<a href="https://codecov.io/gh/go-kratos/kratos"><img src="https://codecov.io/gh/go-kratos/kratos/master/graph/badge.svg" alt="codeCov"></a>
<a href="https://goreportcard.com/report/github.com/go-kratos/kratos"><img src="https://goreportcard.com/badge/github.com/go-kratos/kratos" alt="Go Report Card"></a>
<a href="https://github.com/go-kratos/kratos/blob/main/LICENSE"><img src="https://img.shields.io/github/license/go-kratos/kratos" alt="License"></a>
<a href="https://github.com/avelino/awesome-go"><img src="https://awesome.re/mentioned-badge.svg" alt="Awesome Go"></a>
<a href="https://discord.gg/BWzJsUJ"><img src="https://img.shields.io/discord/766619759214854164?label=chat&logo=discord" alt="Discord"></a>
</p>
<p align="center">
<a href="https://trendshift.io/repositories/3233" target="_blank"><img src="https://trendshift.io/api/badge/repositories/3233" alt="go-kratos%2Fkratos | Trendshift" style="width: 250px; height: 55px;" width="250" height="55"/></a>
<a href="https://www.producthunt.com/posts/go-kratos?utm_source=badge-featured&utm_medium=badge&utm_souce=badge-go-kratos" target="_blank"><img src="https://api.producthunt.com/widgets/embed-image/v1/featured.svg?post_id=306565&theme=light" alt="Go Kratos - A Go framework for microservices. | Product Hunt" style="width: 250px; height: 54px;" width="250" height="54" /></a>
</p>

Translations: [English](README.md) | [简体中文](README_zh.md)

# Kratos

Kratos 一套轻量级 Go 微服务框架，包含大量微服务相关功能及工具。  

> 名字来源于:《战神》游戏以希腊神话为背景，讲述奎托斯（Kratos）由凡人成为战神并展开弑神屠杀的冒险经历。

## Goals

我们致力于提供完整、全面的微服务研发体验，通过整合相关框架和工具，微服务治理部分能够无缝融入整个业务开发周期，使开发者更加专注于业务交付。
对每位开发者而言，Kratos 是非常不错的学习仓库，可以了解和参考微服务领域的技术积累和经验。

### Principles

* 简单：不过度设计，代码平实简单；
* 通用：通用业务开发所需要的基础库的功能；
* 高效：提高业务迭代的效率；
* 稳定：基础库可测试性高，覆盖率高，有线上实践安全可靠；
* 健壮：通过良好的基础库设计，减少错用；
* 高性能：性能高，但不特定为了性能做 hack 优化，引入 unsafe ；
* 扩展性：良好的接口设计，来扩展实现，或者通过新增基础库目录来扩展功能；
* 容错性：为失败设计，大量引入对 SRE 的理解，鲁棒性高；
* 工具链：包含大量工具链，比如 cache 代码生成，lint 工具等等。

## Features

* [APIs](https://go-kratos.dev/docs/component/api) ：协议通信以 HTTP/gRPC 为基础，通过 Protobuf 进行定义；
* [Errors](https://go-kratos.dev/docs/component/errors/) ：通过 Protobuf 的 Enum 作为错误码定义，以及工具生成判定接口；
* [Metadata](https://go-kratos.dev/docs/component/metadata) ：在协议通信 HTTP/gRPC 中，通过 Middleware 规范化服务元信息传递；
* [Config](https://go-kratos.dev/docs/component/config) ：支持多数据源方式，进行配置合并铺平，通过 Atomic 方式支持动态配置；
* [Logger](https://go-kratos.dev/docs/component/log) ：标准日志接口，可方便集成三方 log 库，并可通过 fluentd 收集日志；
* [Metrics](https://go-kratos.dev/docs/component/middleware/metrics) ：统一指标接口，可以实现各种指标系统，默认集成 Prometheus；
* [Tracing](https://go-kratos.dev/docs/component/middleware/tracing) ：遵循 OpenTelemetry 规范定义，以实现微服务链路追踪；
* [Encoding](https://go-kratos.dev/docs/component/encoding) ：支持 Accept 和 Content-Type 进行自动选择内容编码；
* [Transport](https://go-kratos.dev/docs/component/transport/overview) ：通用的 [HTTP](https://go-kratos.dev/docs/component/transport/http) /[gRPC](https://go-kratos.dev/docs/component/transport/grpc) 传输层，实现统一的 [Middleware](https://go-kratos.dev/docs/component/middleware/overview) 插件支持；
* [Registry](https://go-kratos.dev/docs/component/registry) ：实现统一注册中心接口，可插件化对接各种注册中心；
* [Validation](https://go-kratos.dev/docs/component/middleware/validate): 通过 Protobuf 统一定义校验规则，并同时适用于 HTTP/gRPC 服务；
* [SwaggerAPI](https://go-kratos.dev/docs/guide/openapi): 通过集成第三方 [Swagger 插件](https://github.com/go-kratos/swagger-api) 能够自动生成 Swagger API 文档并启动内置的 Swagger UI服 务。

## Getting Started

### Required

- [go](https://golang.org/dl/)
- [protoc](https://github.com/protocolbuffers/protobuf)
- [protoc-gen-go](https://github.com/protocolbuffers/protobuf-go)

### Installing

##### go install 安装：

```
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
kratos upgrade
```

##### 源码编译安装：

```
git clone https://github.com/go-kratos/kratos
cd kratos
make install
```

### Create a service

```
# 创建项目模板
kratos new helloworld

cd helloworld
# 拉取项目依赖
go mod download

# 生成 proto 模板
kratos proto add api/helloworld/helloworld.proto
# 生成 proto 源码
kratos proto client api/helloworld/helloworld.proto
# 生成 server 模板
kratos proto server api/helloworld/helloworld.proto -t internal/service

# 生成所有 proto 源码、wire 等等
go generate ./...

# 运行程序
kratos run
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
* Website:  [go-kratos.dev](https://go-kratos.dev)
* QQ Group: 716486124

## WeChat Official Account

![kratos](docs/images/wechat.png)

## Conventional commits

提交信息的结构应该如下所示:

```text
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

提交信息应按照下面的格式:
- fix: simply describe the problem that has been fixed
- feat(log): simple describe of new features
- deps(examples): simple describe the change of the dependency
- break(http): simple describe the reasons for breaking change

## Sponsors and Backers

![kratos](docs/images/alipay.png)

## License

Kratos is MIT licensed. See the [LICENSE](./LICENSE) file for details.
