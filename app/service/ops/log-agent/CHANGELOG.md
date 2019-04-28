### log-agent
##### 2.1.0(20181228)
> 1. 支持文件采集

##### 2.0.6(20181224)
> 1. priority==high的时候不采样

##### 2.0.5(20181224)
> 1. 支持grpc与lancer交互
> 2. 支持priority=high标识高优先级日志

##### 2.0.4(20181212)
> 1. fix fmt import

##### 2.0.3(20181205)
> 1. 调整日志默认最长为32K

##### 2.0.2(20181128)
> 1. 支持lancer route table

##### 2.0.1(20181120)
> 1. fix: revocer from bufio.Write error

##### 2.0.0(20181110)
> 1. 插件化架构
> 2. 支持非日志类数据上报（非000161）
> 3. flowmonitor 增加kind字段，标识是否为错误

##### 1.1.8.2(20180830)
> 1. httpstream查看日志

##### 1.1.8.1(20180821)
> 1. flush buf first when recycle conn caused by expired

##### 1.1.8.0(20180820)
> 1. lancer新上报协议

##### 1.1.7.2(20180810)
> 1. 修复sendChan 中bytes.buffer 的data race问题

##### 1.1.7.1(20180726)
> 1. 支持日志聚合发送
> 2. 重构连接池：并发发送 + buf和conn隔离

##### 1.1.7(20180726)
> 1. 修复buf size
> 2. 修复getConn时错误的处理，nil不会被putConn

##### 1.1.6(20180621)
> 1. 校验日志长度是否大于_logLancerHeaderLen

##### 1.1.4(20180619)
> 1. 增加telnet功能，本地流式查看日志
> 2. 重构日志流监控

##### 1.1.3(20180604)
> 1. 修改日志级别策略，先检查路由表再检查app_id合规性

##### 1.1.2(20180529)
> 1. 增加路由表功能：根据日志分级策略路由到不同的logid
> 2. 批量接收日志（collector_tcp.sock）

##### 1.1.1
> 1. 通过conn pool连接到lancer

