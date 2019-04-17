# Kratos

Kratos是[bilibili](https://www.bilibili.com)开源的一套Go微服务框架，包含大量微服务相关框架及工具。主要包括以下组件：

* [http框架blademaster(bm)](doc/wiki-cn/blademaster.md)：基于[gin](https://github.com/gin-gonic/gin)二次开发，具有快速、灵活的特点，可以方便的开发中间件处理通用或特殊逻辑，基础库默认实现了log&trace等。
* [gRPC框架warden](doc/wiki-cn/warden.md)：基于官方gRPC封装，默认使用[discovery](https://github.com/bilibili/discovery)进行服务注册发现，及wrr和p2c(默认)负载均衡。
* [dapper trace](doc/wiki-cn/dapper.md)：基于opentracing，全链路集成了trace，我们还提供dapper实现，请参看：[dapper敬请期待]()。
* [log](doc/wiki-cn/logger.md)：基于[zap](https://github.com/uber-go/zap)的field方式实现的高性能log库，集成了我们提供的[log-agent敬请期待]()日志收集方案。
* [database](doc/wiki-cn/database.md)：集成MySQL&HBase&TiDB的SDK，其中TiDB使用服务发现方案。
* [cache](doc/wiki-cn/cache.md)：集成memcache&redis的SDK，注意无redis-cluster实现，推荐使用代理模式[overlord](https://github.com/bilibili/overlord)。
* [kratos tool](doc/wiki-cn/kratos-tool.md)：kratos相关工具量，包括项目快速生成、pb文件代码生成、swagger文档生成等。

我们致力于提供完整的微服务研发体验，整合相关框架及工具后，微服务治理相关部分可对整体业务开发周期无感，从而更加聚焦于业务交付。对每位开发者而言，整套Kratos框架也是不错的学习仓库，可以了解和参考到[bilibili](https://www.bilibili.com)在微服务方面的技术积累和经验。

# 快速开始

```shell
go get -u github.com/bilibili/kratos/tool/kratos
kratos init
```

`kratos init`会快速生成基于kratos库的脚手架代码，如生成[kratos-demo](https://github.com/bilibili/kratos-demo)

```shell
cd kratos-demo/cmd
go build
./cmd -conf ../configs
```

打开浏览器访问：[http://localhost:8000/kratos-demo/start](http://localhost:8000/kratos-demo/start)，你会看到输出了`Golang 大法好 ！！！`

[快速开始](doc/wiki-cn/quickstart.md)

# Document

[简体中文](doc/wiki-cn/summary.md)

-------------

*Please report bugs, concerns, suggestions by issues, or join QQ-group 716486124 to discuss problems around source code.*
