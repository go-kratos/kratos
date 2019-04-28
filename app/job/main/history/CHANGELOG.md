#### history-job 

##### Version 1.11.4
> 1. tidb client 

##### Version 1.11.3
> 1. 重新打包 

##### Version 1.11.2
> 1. 去掉删除逻辑

##### Version 1.11.1
> 1. 调整删除逻辑
##### Version 1.11.0
> 1. 调整删除逻辑: 增加互斥锁 按时间分批删除数据
##### Version 1.10.2
> 1. 配置重试时间
##### Version 1.10.1
> 1. 无限重试insert
##### Version 1.10.0
> 1. 删除用户过多的历史记录
##### Version 1.9.1
> 1. 调整burst
##### Version 1.9.0
> 1. 限制写入速度
##### Version 1.8.3
> 1. 每次都commit
##### Version 1.8.2
> 1. 调整上报日志
##### Version 1.8.1
> 1. 增加prom上报
##### Version 1.8.0
> 1. 相同的mid聚合在一起
##### Version 1.7.3
> 1. 增加ignore参数
##### Version 1.7.2
> 1. rebase master
##### Version 1.7.1
> 1. 调整删除语句 改为根据mtime删除
> 2. 限制TIDB 写入qps
> 3. 修复落后时重复写入的问题
> 4. 支持忽略databus消息

##### Version 1.6.1
> 1. 去掉错误重试
##### Version 1.6.0
> 1. job写入数据库

##### Version 1.5.3
> 1. 消费databus同步commit

##### Version 1.5.1
> 1. 调整调用rpc超时时间
> 2. 分批调用flush rpc 接口

##### Version 1.5.0
> 1. 接入history-service

##### Version 1.4.8
> 1.异步写hbase  fix   

##### Version 1.4.7
> 1.异步写hbase           

##### Version 1.4.6
> 1.重新构建镜像    

##### Version 1.4.5
> 1.迁移 bm fix  

##### Version 1.4.4
> 1.迁移 bm   

##### Version 1.4.3
> 1.迁移main目录

##### Version 1.4.2
> 1.去除statsd

##### Version 1.4.1
> 1.去除identify

##### Version 1.4.0
> 1.接入新配置中心

##### Version 1.3.0
> 1.优化聚合

##### Version 1.2.5
> 1.修复消费跟不上

##### Version 1.2.4
> 1.修复聚合丢数据

##### Version 1.2.3
> 1.bug 

##### Version 1.2.2
> 1.增加 mid白名单

##### Version 1.2.1

##### Version 1.2.0
> 1.异步聚合数据

##### Version 1.1.0
> 1.播放进度聚合databus

##### Version 1.0.0
> 1.rpc聚合数据
