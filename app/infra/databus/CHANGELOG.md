### Databus

##### Version 3.5.2
1. 返回消费者地址  

##### Version 3.5.1
1. 添加pub耗时监控  

##### Version 3.5.0
1. 添加 HTTP 协议发布接口

##### Version 3.4.0
1. 修复read value中包含换行符读取失败
2. 添加hset指令，支持metadata
3. 添加protobuf返回，mget指定pb或json
4. 添加color染色消息过虑

##### Version 3.3.0
1. 迁移infra

##### Version 3.2.6
1. 修复producer关闭顺序

##### Version 3.2.5
1. 修复集群变更导致的panic

##### Version 3.2.4
1. add dao ut

##### Version 3.2.3
1. 去掉大量日志

##### Version 3.2.2
1. Message 支持返回ts

##### Version 3.2.1
1. 修复sql错误

##### Version 3.2.0
1. databus支持置顶批量拉取消息数

##### Version 3.1.0
1. 删除offset选项，默认使用newest
2. 查询判断auth2的appid不为0

##### Version 3.0.1
1. offset默认改为new

##### Version 3.0.0
1. 迁移大仓库

##### Version 2.10.2
1. 增加错误返回

##### Version 2.10.1
1. 增加register接口

##### Version 2.10.0
1. 使用新配置中心v2版

##### Version 2.9.0
1. 去掉集群配对，使集群无状态
2. 支持topic切换集群，对客户端无感换集群
3. XLog改为 Log

##### Version 2.8.0

1. 升级依赖sarama-cluster到v2.1.10
2. 升级依赖sarama到v1.14.0
3. 增加sub时rebalance notify判断

##### Version 2.7.2

1. 更换auth 方式

##### Version 2.7.1

1. fix 重启中produer生产失败

##### Version 2.7.0

1. 添加prom监控
2. 限制consumer创建个数不超过partition数量

##### Version 2.6.2

1. 设定最大重试次数
2. 一旦有未确认发送成功的消息时则后续消息不允许发送

##### Version 2.6.1

1. 修复锁的使用

##### Version 2.6.0

1. 接入配置中心
2. 增加debug日志

##### Version 2.5.0

1. 兼容log agent
2. 当读连接错误或者客户端主动断开时，不再写连接

##### Version 2.4.5

1. tcp连接设置写超时为5s

##### Version 2.4.4

1. 设置sarama tcp keepalive=30s

##### Version 2.4.3

1. 更改consumer max process time 为 50ms
2. 更改consumer max wait time 为250ms

##### Version 2.4.2

1. 修复rebalance没有踢出老sub
2. 添加统计信息

##### Version 2.4.1

1. 修改为异步链接监听

##### Version 2.4.0

1. 采用mo进行统计

##### Version 2.3.9

1. 移除配置中心并格式化代码

##### Version 2.3.8

1. 多producer

##### Version 2.3.7

1. 增加监控接口

##### Version 2.3.6

1. 接入配置中心，无配置启动
2. auth时新增offset参数(new/old)，允许client指定初始消费位置

##### Version 2.3.5

1. 修复ReadSlice() 导致的bug

##### Version 2.3.4

1. 更改 monitor 统计信息map的key为string类型

##### Version 2.3.3

1. 强制要求业务生产的msg内容必须为json格式，否则可以pub成功，但会导致sub失败
2. kafka消息的value采用json.RawMessage格式，不再对value decode
3. 调整tcp的读写buffer大小，read buffer:64k,write buffer:8k，消息体最多允许64k大小，否则报错
4. 调整授权信息的缓存策略，当查询mysql出错时，不再清空cache，查询次数为5分钟一次

##### Version 2.2.1

1. 在建立客户端和databus的链接后，可以多次auth
2. 兼容redis客户端断开连接时发送的QUIT命令
3. 新增监控信息，pub角色统计生成消息总数、字节数；sub角色统计消费消息总数、字节数、每个分区已消费和已提交的offset
4. 增加读连接超时，新建连接后5s 内 不发auth断开连接;生产消息 20分钟 没消息断开连接;消费消息 40s 内没 mget 断开连接
5. 取消 mset 命令，不再支持设置 partititon 的 offset 进行回滚
6. 使用 govendor 进行第三方包管理

##### Version 2.1.1

1. 移除对 go-common/business/identify,xweb/router,xhttp/router 的依赖
2. 新增mset命令，同时设置多个partition 的 offset 进行消息回滚
3. 修改第三方库sarama-cluster，设置partition 的offset后动态回滚到指定offset位置
4. offset提交改为标记方式，databus自动提交（提交间隔为一秒）

##### Version 2.0.1

1. 用redis mget命令替换smembers保证返回消息的顺序性

##### Version 2.0.0

1. 升级kafka至0.10
2. appkey区分集群,新增授权时需增加业务名字段business
3. auth 时使用dsn协议，key:secret@group/topic=?&role=?
4. redis smembers、set命令进行通信，producer 使用set生成消息;consumer使用smembers批量消费消息、set提交partition offset
5. 不再自动保存partition offset，由consumer自己手动提交
6. 一次smembers目前默认flush 100 条消息;或者100 ms超时时flush一次

##### Version 1.4.1

1.增加日志

##### Version 1.4.0

1.同步写数据
2.同个group+topic允许多个consumer
3.修复subscrible时参数错误可能导致的数据越界

##### Version 1.3.1

1.完善错误日志信息

##### Version 1.3.0

1.完善auth授权
2.添加监控接口

##### Version 1.2.0

1.共用zk链接

##### Version 1.1.0

1. 支持客户端维护offset，并通过第一次请求设置offset
2. 支持主动设置offset到zookeeper

##### Version 1.0.0

1.数据总线
2.TODO auth and close
