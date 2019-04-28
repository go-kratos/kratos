### dapper-service
##### Version 3.1.1
>1. 如果没有配置就不启动 kafka collector

##### Version 3.1.0
>1. 添加 kafka collector
>2. 忽略旧 sdk 的 span 采样点

##### Version 3.0.1
>1. 修复 HTTP client opreation_name 过多的问题

##### Version 3.0.0
>1. dapper 重构, 重新设计存储格式，添加了 influxdb

##### Version 2.0.4
>1. 允许通过 family 搜索 span

##### Version 2.0.3
>1. depends 接口添加 cache

##### Version 2.0.2
>1. 优化 es 查询, 缩短默认时间范围到1小时，避免超时

##### Version 2.0.1
>1. 修复 es mapping 不一致无法查询的问题

##### Version 2.0
>1.dapper 重构

##### Version 1.3.6
>1.优化 collect 日志避免写满磁盘

##### Version 1.3.5
>1. 升级 hbase client

##### Version 1.3.4
>1. 删除无用配置
>2. 接bm

##### Version 1.3.3
>1. 迁移目录

##### Version 1.3.2
>1. 移除statsd 模块

##### Version 1.3.1
>1. 修复循环依赖的接口

##### Version 1.3.0
>1. collect 支持处理聚合日志，协议unxigram改为unix

##### Version 1.2.0

>1. 增加没有服务依赖的接口，单独返回
>2. 增加循环依赖服务接口
>3. 增加服务组件依赖接口，查询服务最近12小时依赖的组件和服务title

##### Version 1.1.1

> 1.修改服务依赖图数据，没有调用方的服务自己依赖自己

##### Version 1.1.0

> 1.span时间修复

##### Version 1.0.0

> 1.初始化完成dapper服务的基本功能,agent和collect 分别是日志客户端收集器。收集trace发送到dapper
