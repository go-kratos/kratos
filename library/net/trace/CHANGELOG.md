### trace sdk

##### Version 4.0.0
> 1. 修改日志协议，使用 protobuf 序列化 span
> 2. 修改 interface  添加 Tag, Log Method

##### Version 3.2.2
> 1. 修复sync.Pool 未采样的trace ,未放回pool的bug

##### Version 3.2.1
> 1. 去掉链路日志最后加的换行，目前解析不需要换行

##### Version 3.2.0
> 1. use sync.Pool for new trace.
> 2. spanID use local calc.

##### Version 3.1.1
> 1. 去掉comment interface 支持
> 2. 替换comment 中分隔符为空""

##### Version 3.1.0
> 1. update user to caller
> 2. rpc client init trace info

##### Version 3.0.1
> 1. 去掉title里面的host
> 2. 过滤不采样的url

##### Version 3.0.0
> 1. 修改日志协议，将一个span由之前发送两次改为只发送一次
> 2. 输出流由syslog 改为自定义实现
> 3. 将初始化配置整合到一个config里面

##### Version 2.0.0
> 1. 优化实现逻辑，将之前实现由暴露一个公开的结构体改为提供接口和实现，业务可自行实现接口

##### Version 1.0.0
> 1. 提供dapper接入sdk,对基础组件,rpc,http 调用进行封装，上报链路信息
