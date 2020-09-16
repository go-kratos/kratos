![kratos](img/kratos3.png)
# Kratos

Kratos是bilibili开源的一套Go微服务框架，包含大量微服务相关框架及工具。  

### Goals

我们致力于提供完整的微服务研发体验，整合相关框架及工具后，微服务治理相关部分可对整体业务开发周期无感，从而更加聚焦于业务交付。对每位开发者而言，整套Kratos框架也是不错的学习仓库，可以了解和参考到bilibili在微服务方面的技术积累和经验。

### Principles

* 简单：不过度设计，代码平实简单
* 通用：通用业务开发所需要的基础库的功能
* 高效：提高业务迭代的效率
* 稳定：基础库可测试性高，覆盖率高，有线上实践安全可靠
* 健壮：通过良好的基础库设计，减少错用
* 高性能：性能高，但不特定为了性能做hack优化，引入unsafe
* 扩展性：良好的接口设计，来扩展实现，或者通过新增基础库目录来扩展功能
* 容错性：为失败设计，大量引入对SRE的理解，鲁棒性高
* 工具链：包含大量工具链，比如cache代码生成，lint工具等等

### Features
* HTTP Blademaster：核心基于[gin](https://github.com/gin-gonic/gin)进行模块化设计，简单易用、核心足够轻量；
* GRPC Warden：基于官方gRPC开发，集成[discovery](https://github.com/bilibili/discovery)服务发现，并融合P2C负载均衡；
* Cache：优雅的接口化设计，非常方便的缓存序列化，推荐结合代理模式[overlord](https://github.com/bilibili/overlord)；
* Database：集成MySQL/HBase/TiDB，添加熔断保护和统计支持，可快速发现数据层压力；
* Config：方便易用的[paladin sdk](config-paladin.md)，可配合远程配置中心，实现配置版本管理和更新；
* Log：类似[zap](https://github.com/uber-go/zap)的field实现高性能日志库，并结合log-agent实现远程日志管理；
* Trace：基于opentracing，集成了全链路trace支持（gRPC/HTTP/MySQL/Redis/Memcached）；
* Kratos Tool：工具链，可快速生成标准项目，或者通过Protobuf生成代码，非常便捷使用gRPC、HTTP、swagger文档；


-------------

> 名字来源于:《战神》游戏以希腊神话为背景，讲述由凡人成为战神的奎托斯（Kratos）成为战神并展开弑神屠杀的冒险历程。

