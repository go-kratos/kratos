### 好友关系

### Version 2.29.0
> 1. 添加grpc

### Version 2.28.0
> 1. 相同关注接口

### Version 2.27.0
> 1. 加关注按appkey限流

### Version 2.26.1
> 1. tag 不存在直接跳过

### Version 2.26.0
> 1. 去除 RemoteIP 方法
> 2. refine test

### Version 2.25.0
> 1. 增加新粉丝通知频率限制

### Version 2.24.4
> 1. block业务迁移依赖member-service

### Version 2.24.3
> 1. 成就奖励获取限制非正式会员

### Version 2.24.2
> 1. 修正取消悄悄关注日志

### Version 2.24.1
> 1. 修正粉丝成就粉丝数

### Version 2.24.0
> 1. 粉丝成就接口

### Version 2.23.9
> 1. 旧Audit接口下线，引发的小重构

### Version 2.23.8
> 1. 去除 SQL Prepare

### Version 2.23.7
> 1. 修正 context 使用

### Version 2.23.6
> 1. 去除 xints

### Version 2.23.5
> 1. 去除删除属性的错误

### Version 2.23.4
> 1. 使用 bm Bind

### Version 2.23.3
> 1. 对比缓存不一致删缓存

### Version 2.23.2
> 1. 增加黑名单日志

### Version 2.23.1
> 1. 修复关注日志中 source 不显示的问题

### Version 2.23.0
> 1. 上报关注日志

### Version 2.22.1
> 1. 忽略 redis nil error

### Version 2.22.0
> 1. service 直接用来统计最近粉丝数

### Version 2.21.0
> 1. remove global hot

### Version 2.20.5
> 1. restrict local cache

### Version 2.20.4
> 1. configureable local cache least follower

### Version 2.20.3
> 1. ignore nil stat in local cache

### Version 2.20.2
> 1. local cache stat

### Version 2.20.1
> 1. ignore delete stat cache key not found error

### Version 2.20.0
> 1. move stat cache to memcahched

### Version 2.19.0

> 1. support bm.  

### Version 2.18.0

> 1. update infoc sdk

### Version 2.17.5
> 1. 修改 proto package，合并 proto 文件

### Version 2.17.4
> 1.add register

### Version 2.17.3
> 1.增加 followers 和 followings 缓存计数
> 2.修正 gitignore

### Version 2.17.2
> 1.去掉account-service的使用  

### Version 2.17.1
> 1.关注提示功能
> 2.remove statsd

### Version 2.17.0
> 1.关注推荐功能

### Version 2.16.0
> 1.修复stat 闭包错误  

### Version 2.15.1
> 1.修复空缓存无效  

### Version 2.15.0 
> 1.添加b+特殊关注接口  

### Version 2.14.2
> 1.修复stat缓存  

### Version 2.14.1
> 1.token interceport  

### Version 2.14.0
> 1.打开handshake token验证  

### Version 2.13.1
> 1.关注限制：修改code，未绑定手机且用户未转正时候禁止关注、悄悄关注操作

### Version 2.13.0
> 1.关注限制：未绑定手机且用户未转正时候禁止关注、悄悄关注操作

### Version 2.12.0
> 1.限制mid<=0

### Version 2.11.0
> 1.特殊关注

### Version 2.10.0
> 1.业务场景下关注引导提示  
### Version 2.9.1
> 1.分组排序  

### Version 2.9.0
> 1.关注分组  

### Version 2.8.0
> 1.监控特定用户被关注，返回关注成功

### Version 2.7.3
> 1.空缓存优化  

### Version 2.7.2
> 1.停掉写入member log

### Version 2.7.0
> 1.使用protobuf  

### Version 2.6.7
> 1.fix 添加悄悄关注 不增加粉丝数  

### Version 2.6.6 
> 1.接入新trace  

### Version 2.6.5
> 1.fix addmlog err report

### Version 2.6.3
> 1.判断mid，fid > 0

### Version 2.6.0
> 1.迁移到大仓库

### Version 2.5.1
> 1.恢复重复操作错误提示

### Version 2.5.0
> 1.开启gzip和gob

### Version 2.4.0
> 1.更新mc基础包压缩缓存

### Version 2.3.0
> 1.日志上报数据平台（infoc方式）

### Version 2.2.2
> 1.接入antispam

### Version 2.2.1
> 1.接入rpc平滑发版

### Version 2.2.0
> 1.更新依赖包
> 2.接入普罗米修斯

### Version 2.1.0
> 1.拆分黑名单数量限制

### Version 2.0.15
> 1.memberlog增加src字段

### Version 2.0.14
> 1.写入memberlog表

### Version 2.0.13
> 1.修复uint64转int64溢出

### Version 2.0.12
> 1.加回audit接口
> 1.兼容旧数据允许用户取消拉黑自己

### Version 2.0.11
> 1.删除ctime

### Version 2.0.10
> 1.兼容好友attr

### Version 2.0.9
> 1.改回mtime排序

### Version 2.0.8
> 1.增加ctime

### Version 2.0.7
> 1.修复关注排序

### Version 2.0.6

> 1.黑名单列表增加stat判断

### Version 2.0.5
> 1.修改expire

### Version 2.0.4
> 1.删除黑名单双向关系

### Version 2.0.3

> 1.修改url为internal

### Version 2.0.1

> 1.修复setstat设置stat计数
> 2.添加互相关注关系到粉丝列表

### Version 2.0.1

> 1.暂时注释掉audit以及速率限制接口

### Version 1.0.2

> 1.接入新的配置中心

### Version 1.0.1

> 1.更改/monitor/ping

#### Version 1.0.0
> 1.基础api
