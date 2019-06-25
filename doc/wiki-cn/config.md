# config

## 介绍
初看起来，配置管理可能很简单，但是这其实是不稳定的一个重要来源。  
即变更管理导致的故障，我们目前基于配置中心（config-service）的部署方式，二进制文件的发布与配置文件的修改是异步进行的，每次变更配置，需要重新构建发布版。  
由此，我们整体对配置文件进行梳理，对配置进行模块化，以及方便易用的paladin config sdk。

## 环境配置

| flag   | env   | remark |
|:----------|:----------|:------|
| region | REGION | 部署地区，sh-上海、gz-广州、bj-北京 |
| zone | ZONE | 分布区域，sh001-上海核心、sh004-上海嘉定 |
| deploy.env | DEPLOY_ENV | dev-开发、fat1-功能、uat-集成、pre-预发、prod-生产 |
| deploy.color | DEPLOY_COLOR | 服务颜色，blue（测试feature染色请求） |
| - | HOSTNAME | 主机名，xxx-hostname |

全局公用环境变量，通常为部署环境配置，由系统、发布系统或supervisor进行环境变量注入，并不用进行例外配置，如果是开发过程中则可以通过flag注入进行运行测试。

## 应用配置

| flag   | env   |      default      |  remark |
|:----------|:----------|:-------------|:------|
| appid | APP_ID | - | 应用ID |
| http | HTTP | tcp://0.0.0.0:8000/?timeout=1s | http 监听端口 |
| http.perf | HTTP_PERF | tcp://0.0.0.0:2233/?timeout=1s | http perf 监听端口 |
| grpc | GRPC | tcp://0.0.0.0:9000/?timeout=1s&idle_timeout=60s | grpc 监听端口 |
| grpc.target | - | - | 指定服务运行：<br>-grpc.target=demo.service=127.0.0.1:9000 <br>-grpc.target=demo.service=127.0.0.2:9000 |
| discovery.nodes | DISCOVERY_NODES | - | 服务发现节点：127.0.0.1:7171,127.0.0.2:7171 |
| log.v | LOG_V |  0 | 日志级别：<br>DEBUG:0 INFO:1 WARN:2 ERROR:3 FATAL:4 |
| log.stdout | LOG_STDOUT | false | 是否标准输出：true、false|
| log.dir | LOG_DIR | - | 日志文件目录，如果配置会输出日志到文件，否则不输出日志文件 |
| log.agent | LOG_AGENT | - | 日志采集agent：<br>unixpacket:///var/run/lancer/collector_tcp.sock?timeout=100ms&chan=1024 |
| log.module | LOG_MODULE | - | 指定field信息 format: file=1,file2=2. |
| log.filter | LOG_FILTER | - | 过虑敏感信息 format: field1,field2. |

基本为一些应用相关的配置信息，通常发布系统和supervisor都有对应的部署环境进行配置注入，并不用进行例外配置，如果开发过程中可以通过flag进行注入运行测试。

## 业务配置
Redis、MySQL等业务组件，可以使用静态的配置文件来初始化，根据应用业务集群进行配置。

## 在线配置
需要在线读取、变更的配置信息，比如某个业务开关，可以实现配置reload实时更新。

## 扩展阅读

[paladin配置sdk](config-paladin.md)  

-------------

[文档目录树](summary.md)
