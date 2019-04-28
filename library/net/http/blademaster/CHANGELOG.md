### net/http/blademaster

##### Version 1.1.4
1. 临时移除 httptrace 避免 datarace

##### Version 1.1.3
1. bind 错误设置到context error

##### Version 1.1.2
1. 将 ecode 作为 header 写入

##### Version 1.1.1
1. device 信息加入metadata

##### Version 1.1.0

1. 对压测流量打标，写入md

##### Version 1.0.6

1. 业务错误日志记为 WARN

##### Version 1.0.5

1. 增加 device 中间件

##### Version 1.0.4

1. 增加 metadata 接口，可以获取 Path 和 Method 信息

##### Version 1.0.3

1. 当请求被 CORS 或者 CSRF 模块拒绝后，输出一个 level 为 5 的 Error 日志

##### Version 1.0.2

1. 调整 context.go 里的输出方法参数顺序，改为数据在前，error 在后
2. Context 里增加 JSONMap 方法，用于适配早期数据结构
3. Recovery 里打印 panic 信息到 stderr

##### Version 1.0.1

1. logger 里增加上报用于监控的 caller

##### Version 1.0.0

1. 完成基本功能与测试


