# 日志基础库

## 概览
基于[zap](https://github.com/uber-go/zap)的field方式实现的高性能log库，提供Info、Warn、Error日志级别；  
并提供了context支持，方便打印环境信息以及日志的链路追踪，在框架中都通过field方式实现，避免format日志带来的性能消耗。

## 配置选项

| flag   | env   |      type      |  remark |
|:----------|:----------|:-------------:|:------|
| log.v | LOG_V |  int | 日志级别：DEBUG:0 INFO:1 WARN:2 ERROR:3 FATAL:4 |
| log.stdout | LOG_STDOUT | bool | 是否标准输出：true、false|
| log.dir | LOG_DIR | string | 日志文件目录，如果配置会输出日志到文件，否则不输出日志文件 |
| log.agent | LOG_AGENT | string | 日志采集agent：unixpacket:///var/run/lancer/collector_tcp.sock?timeout=100ms&chan=1024 |
| log.module | LOG_MODULE | string | 指定field信息 format: file=1,file2=2. |
| log.filter | LOG_FILTER | string | 过虑敏感信息 format: field1,field2. |

## 使用方式
```go
func main() {
  // 解析flag
  flag.Parse()
  // 初始化日志模块
  log.Init(nil)
  // 打印日志
  log.Info("hi:%s", "kratos")
  log.Infoc(Context.TODO(), "hi:%s", "kratos")
  log.Infov(Context.TODO(), log.KVInt("key1", 100), log.KVString("key2", "test value")
}
```

## 扩展阅读
* [log-agent](log-agent.md)

