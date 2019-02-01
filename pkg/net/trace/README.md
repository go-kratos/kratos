# net/trace

## 项目简介
1. 提供Trace的接口规范
2. 提供 trace 对Tracer接口的实现，供业务接入使用

## 接入示例
1. 启动接入示例
    ```go
    trace.Init(traceConfig) // traceConfig is Config object with value.
    ```
2. 配置参考
    ```toml
    [tracer]
    network = "unixgram"
    addr = "/var/run/dapper-collect/dapper-collect.sock"
    ```

## 测试
1. 执行当前目录下所有测试文件，测试所有功能
