![kratos](doc/img/kratos3.png)

[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![Build Status](https://travis-ci.org/bilibili/kratos.svg?branch=master)](https://travis-ci.org/bilibili/kratos)
[![GoDoc](https://godoc.org/github.com/bilibili/kratos?status.svg)](https://godoc.org/github.com/bilibili/kratos)
[![Go Report Card](https://goreportcard.com/badge/github.com/bilibili/kratos)](https://goreportcard.com/report/github.com/bilibili/kratos)

# Kratos

Kratos是[bilibili](https://www.bilibili.com)开源的一套Go微服务框架，包含大量微服务相关框架及工具。  

> 名字来源于:《战神》游戏以希腊神话为背景，讲述由凡人成为战神的奎托斯（Kratos）成为战神并展开弑神屠杀的冒险历程。

## Goals

我们致力于提供完整的微服务研发体验，整合相关框架及工具后，微服务治理相关部分可对整体业务开发周期无感，从而更加聚焦于业务交付。对每位开发者而言，整套Kratos框架也是不错的学习仓库，可以了解和参考到[bilibili](https://www.bilibili.com)在微服务方面的技术积累和经验。

## Features
* HTTP Blademaster：核心基于[gin](https://github.com/gin-gonic/gin)进行模块化设计，简单易用、核心足够轻量；
* GRPC Warden：基于官方gRPC开发，集成[discovery](https://github.com/bilibili/discovery)服务发现，并融合P2C负载均衡；
* Cache：优雅的接口化设计，非常方便的缓存序列化，推荐结合代理模式[overlord](https://github.com/bilibili/overlord)；
* Database：集成MySQL/HBase/TiDB，添加熔断保护和统计支持，可快速发现数据层压力；
* Config：方便易用的[paladin sdk](doc/wiki-cn/config.md)，可配合远程配置中心，实现配置版本管理和更新；
* Log：类似[zap](https://github.com/uber-go/zap)的field实现高性能日志库，并结合log-agent实现远程日志管理；
* Trace：基于opentracing，集成了全链路trace支持（gRPC/HTTP/MySQL/Redis/Memcached）；
* Kratos Tool：工具链，可快速生成标准项目，或者通过Protobuf生成代码，非常便捷使用gRPC、HTTP、swagger文档；

## Quick start

### Requirments

Go version>=1.13

### Installation
```shell
GO111MODULE=on && go get -u github.com/bilibili/kratos/tool/kratos
cd $GOPATH/src
kratos new kratos-demo
```

通过 `kratos new` 会快速生成基于kratos库的脚手架代码，如生成 [kratos-demo](https://github.com/bilibili/kratos-demo) 

### Build & Run

```shell
cd kratos-demo/cmd
go build
./cmd -conf ../configs
```

打开浏览器访问：[http://localhost:8000/kratos-demo/start](http://localhost:8000/kratos-demo/start)，你会看到输出了`Golang 大法好 ！！！`

[快速开始](doc/wiki-cn/quickstart.md)  [kratos工具](doc/wiki-cn/kratos-tool.md)

## Documentation

> [简体中文](doc/wiki-cn/summary.md)  
> [FAQ](doc/wiki-cn/FAQ.md)  

## License
Kratos is under the MIT license. See the [LICENSE](./LICENSE) file for details.

-------------

*Please report bugs, concerns, suggestions by issues, or join QQ-group 716486124 to discuss problems around source code.*
