![kratos](docs/images/kratos.png)

[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![Build Status](https://github.com/go-kratos/kratos/workflows/Go/badge.svg)](https://github.com/go-kratos/kratos/actions)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-kratos/kratos/v2)](https://pkg.go.dev/github.com/go-kratos/kratos/v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-kratos/kratos)](https://goreportcard.com/report/github.com/go-kratos/kratos)
[![Discord](https://img.shields.io/discord/766619759214854164?label=chat&logo=discord)](https://discord.gg/BWzJsUJ)

Translations: [English](README.md) | [简体中文](README_zh.md)

# Kratos

Kratos 一套轻量级 Go 微服务框架，包含大量微服务相关框架及工具。  

> 名字来源于:《战神》游戏以希腊神话为背景，讲述由凡人成为战神的奎托斯（Kratos）成为战神并展开弑神屠杀的冒险历程。

## Goals

我们致力于提供完整的微服务研发体验，整合相关框架及工具后，微服务治理相关部分可对整体业务开发周期无感，从而更加聚焦于业务交付。对每位开发者而言，整套 Kratos 框架也是不错的学习仓库，可以了解和参考到微服务方面的技术积累和经验。

### Principles

* 简单：不过度设计，代码平实简单；
* 通用：通用业务开发所需要的基础库的功能；
* 高效：提高业务迭代的效率；
* 稳定：基础库可测试性高，覆盖率高，有线上实践安全可靠；
* 健壮：通过良好的基础库设计，减少错用；
* 高性能：性能高，但不特定为了性能做 hack 优化，引入 unsafe ；
* 扩展性：良好的接口设计，来扩展实现，或者通过新增基础库目录来扩展功能；
* 容错性：为失败设计，大量引入对 SRE 的理解，鲁棒性高；
* 工具链：包含大量工具链，比如 cache 代码生成，lint 工具等等；

## Features
* [APIs]((examples/helloworld/helloworld))：协议通信以 HTTP/gRPC 为基础，通过 Protobuf 进行定义；
* [Errors](examples/errors/api)：通过 Protobuf 的 Enum 作为错误码定义，以及工具生成判定接口；
* [Metadata](examples/metadata)：在协议通信 HTTP/gRPC 中，通过 Middleware 规范化服务元信息传递；
* [Config](examples/config)：支持多数据源方式，进行配置合并铺平，通过 Atomic 方式支持动态配置；
* [Logger](examples/log)：标准日志接口，可方便集成三方 log 库，并可通过 fluentd 收集日志；
* [Metrics](examples/metrics)：统一指标接口，可以实现各种指标系统，默认集成 Prometheus；
* [Tracing](examples/traces)：遵循 OpenTelemetry 规范定义，以实现微服务链路追踪；
* [Encoding](encoding)：支持 Accept 和 Content-Type 进行自动选择内容编码；
* [Transport](transport/transport.go)：通用的 [HTTP](examples/http/middlewares)/[gRPC](examples/middleware/main.go) 传输层，实现统一的 [Middleware](middleware) 插件支持；
* [Registry](examples/registry)：实现统一注册中心接口，可插件化对接各种注册中心；
* [Validation](examples/validate): 通过Protobuf统一定义校验规则，并同时适用于HTTP/gRPC服务.
* [SwaggerAPI](https://github.com/go-kratos/swagger-api/blob/main/examples/helloworld/server/main.go): 通过集成第三方[Swagger插件](https://github.com/go-kratos/swagger-api)能够自动生成Swagger API json并启动一个内置的Swaager UI服务.

## Getting Started
### Required
- [go](https://golang.org/dl/)
- [protoc](https://github.com/protocolbuffers/protobuf)
- [protoc-gen-go](https://github.com/protocolbuffers/protobuf-go)

### Installing
##### go get 安装：
```
go get github.com/go-kratos/kratos/cmd/kratos/v2@latest
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

# 生成proto模板
kratos proto add api/helloworld/helloworld.proto
# 生成proto源码
kratos proto client api/helloworld/helloworld.proto
# 生成server模板
kratos proto server api/helloworld/helloworld.proto -t internal/service

# 生成所有proto源码、wire等等
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
* [Examples](./examples)
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
