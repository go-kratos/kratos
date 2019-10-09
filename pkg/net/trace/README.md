# net/trace

### 项目简介

1. 提供Trace的接口规范
2. 提供 trace 对Tracer接口的实现，供业务接入使用

### 接入示例

1. 启动接入示例

```go
trace.Init(traceConfig) // traceConfig is Config object with value.
```

2. 配置参考

- 上报数据到 jaeger-agent 目前只支持 `compact` Thrift protocol, 添加以下 flags

```
-trace=jaeger+udp://127.0.0.1:6831
```

- 上报数据到 jaeger-collector 通过 HTTP 协议, 添加以下 flags
```
# NOTE: 不可指定 Path 默认上报到 http://jaeger-collector:14268/api/traces
-trace=jaeger+http://jaeger-collector:14268
```

- 上报到 zipkin

```go
import (
    "github.com/bilibili/kratos/pkg/net/trace/zipkin"
)

func main() {
    zipkin.Init(&zipkin.Config{...your config})
}
```

### 测试

1. 执行当前目录下所有测试文件，测试所有功能
