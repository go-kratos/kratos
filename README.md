![kratos](docs/images/kratos.png)

[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![Build Status](https://github.com/go-kratos/kratos/workflows/Go/badge.svg)](https://github.com/go-kratos/kratos/actions)
[![GoDoc](https://godoc.org/github.com/go-kratos/kratos?status.svg)](https://godoc.org/github.com/go-kratos/kratos)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-kratos/kratos)](https://goreportcard.com/report/github.com/go-kratos/kratos)
[![Discord](https://img.shields.io/discord/766619759214854164?label=chat&logo=discord)](https://discord.gg/BWzJsUJ)

# Kratos

Kratos 一套轻量级 Go 微服务框架，包含大量微服务相关框架及工具。  

> 名字来源于:《战神》游戏以希腊神话为背景，讲述由凡人成为战神的奎托斯（Kratos）成为战神并展开弑神屠杀的冒险历程。

## Goals

我们致力于提供完整的微服务研发体验，整合相关框架及工具后，微服务治理相关部分可对整体业务开发周期无感，从而更加聚焦于业务交付。对每位开发者而言，整套Kratos框架也是不错的学习仓库，可以了解和参考到微服务方面的技术积累和经验。

## Features
* APIs - 协议通信以 HTTP/gRPC 为基础，通过 Protobuf 进行定义；
* Errors - 通过 Protobuf 定义对应的 Enum 作为错误码，以及国际化信息支持；
* Metadata - 在协议通信 HTTP/gRPC 中，规范化服务元信息传递；
* Config - 通过KeyValue方式实例，对多种配置源进行铺平，Atomic方式动态配置支持；
* Logger - 标准日志接口，可方便集成三方 log库，或者通过 fluentd 收集日志；
* Metrics - 统一指标接口定义，可以默认集成 Prometheus；
* Tracing - 通过 OpenTracing 规范，进行实现微服务链路追踪；
* Encoding - 支持Accept和Content-Type进行自动选择内容编码；
* Transport - 通用的 HTTP/gRPC 传输层，实现统一的 Middleware 插件支持；
* Server - 进行基础的 Server 层封装，统一以 Options 方式配置使用；

## Getting Started
### Required
- [go](https://golang.org/dl/)
- [protoc](https://github.com/protocolbuffers/protobuf)
- [protoc-gen-go](https://github.com/protocolbuffers/protobuf-go)

### Install Kratos
```
# 安装生成工具
go get github.com/go-kratos/kratos/cmd/kratos
go get github.com/go-kratos/kratos/cmd/protoc-gen-go-http
go get github.com/go-kratos/kratos/cmd/protoc-gen-go-errors

# 或者通过 Source 安装
cd cmd/kratos && go install
cd cmd/protoc-gen-go-http && go install
cd cmd/protoc-gen-go-errors && go install
```
### Create a service
```
# 创建项目模板
kratos new helloworld

cd helloworld
# 生成proto模板
kratos proto add api/helloworld/helloworld.proto
# 生成service模板
kratos proto service api/helloworld/helloworld.proto -t internal/service

# 生成api下所有proto文件
make proto
# 编码cmd下所有main文件
make build
# 进行单元测试
make test
```

## Service Layout
* [Service Layout](https://github.com/go-kratos/service-layout)

## Community
* [Wechat Group](https://github.com/go-kratos/kratos/issues/682)
* [Discord Group](https://discord.gg/BWzJsUJ)
* QQ Group: 716486124

## License
Kratos is MIT licensed. See the LICENSE file for details.