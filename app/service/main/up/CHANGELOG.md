#### 创作中心up主身份服务

# v1.11.11
> 1.批量增加查询active tid的接口

# v1.11.10
> 1.loadSpGroupsMids load append not flush  

# v1.11.9
> 1.调整group相关接口 

# v1.11.8
> 1.高能联盟 grpc服务  

# v1.11.7
> 1.切gorpc to grpc服务  

# v1.11.6
> 1.调整分组信息列表为map  

# v1.11.5
> 1.添加特殊用户up主的grpc接口  

# v1.11.4
> 1.添加up主活跃度的列表grpc接口  

# v1.11.3
> 1.change account gorpc to grpc

# v1.11.2
> 1.增加查询active tid的接口

# v1.11.1
> 1.修改speical add为异步逻辑

# v1.11.0
> 1.up主列表相关 + 联合投稿staff      

# v1.10.1
> 1.特殊用户组接口修改错误信息      

# v1.10.0
> 1.up主列表相关    

##### Version 1.9.0
> 1.修改up列表API-查询参数增加活跃度字段

##### Version 1.8.9
> 1.修改人物卡片查询API: up图片和视频按id倒序排序

##### Version 1.8.8
> 1.增加api接口

##### Version 1.8.7
> 1.播放器开关，如果没有找到Up主配置，现在默认为开

##### Version 1.8.6
> 1.增加API: 按 mid 列表查询人物卡片信息

##### Version 1.8.5
> 1.增加获取所有up主mid的服务

##### Version 1.8.4
> 1.增加人物卡片API

##### Version 1.8.3
> 1.fix error log    

##### Version 1.8.2
> 1.修复更新成功触发删除缓存  

##### Version 1.8.1
> 1.修复ut

##### Version 1.8.0
> 1.增加up主开关功能  

##### Version 1.7.10
> 1.增加数据库中不存在时的缓存

##### Version 1.7.9
> 1.修改hbase为hbasev2

##### Version 1.7.8
> 1.更换remoteIP调用

##### Version 1.7.7
> 1.修复没权限返回数据错误的问题
> 2.去掉一处打印过多的日志

##### Version 1.7.6
> 1.修复服务器重启时时DataScheduler导致crash的问题
> 2.增加签约用户访问权限

##### Version 1.7.5
> 1.修复服务器重启时会crash的问题

##### Version 1.7.4
> 1.修复validate tag

##### Version 1.7.3
> 1.替换net/http/parse

##### Version 1.7.2
> 1.重新合入master基础库

##### Version 1.7.1
> 1.增加获取up主特殊用户组时的颜色标签

##### Version 1.7.0
> 1.优化databus消费，使用队列分发策略并行分发且保证顺序消费  

##### Version 1.6.4
> 1.add register

##### Version 1.6.3
> 1.修复proto编译问题

##### Version 1.6.2
> 1.增加用户统计数据接口

##### Version 1.6.1
> 1.调整逻辑：重复Add变为Edit

##### Version 1.6.0
> 1.迁移admin上up special相关接口（1.5.4继承过来）
> 2.增加权限点

##### Version 1.5.4
> 1.迁移admin上up special相关接口

##### Version 1.5.3
> 1.修复getToken导致高CPU占用的问题

##### Version 1.5.2
> 1.更新protobuf的目录结构

##### Version 1.5.1
> 1.增加RPC服务，注册目录为/microservice/up-service/
> 2.增加Info与Special的RPC接口

##### Version 1.5.0
> 1.直播接入

##### Version 1.4.8
> 1.迁移到main目录

##### Version 1.4.7
> 1.fix up nil

##### Version 1.4.6
> 1.解耦异步db更新的频率和routine数量
> 2.确保异步处理databus消息的顺序一致性

##### Version 1.4.5
> 1.迁移videoup-service上的up/special接口到up-service

##### Version 1.4.4
> 1.不消费报警可配置    

##### Version 1.4.3
> 1.五分钟不消费报警  

##### Version 1.4.2
> 1.消费databus，异步更新db

##### Version 1.4.1
> 1.去除statsd

##### Version 1.4.0
> 1.对外接口加个缓存

### Version 1.3.1
> 1.修改up-service接口uri

##### Version 1.3.0
> 1.支持移动端投稿身份校验   

##### Version 1.2.1
> 1.fix databus consume msg type  

##### Version 1.2.0
> 1.接入bm  

##### Version 1.1.0
> 1.去除identity和ecode

##### Version 1.0.0
> 1.稿件up主接入    
