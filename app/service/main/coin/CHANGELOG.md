#### 硬币基础服务

#### V2.16.6
> 1. 修复硬币修改-104问题

#### V2.16.5
> 1. 异步add cache
#### V2.16.4
> 1. 修复投币强依赖redis的问题

#### V2.16.3
> 1. add nolog
#### V2.16.2
> 1. add type to log

#### V2.16.1
> 1. 使用member grpc 替换老的java 接口

#### V2.16.0
> 1. 异步发消息

#### V2.15.1
> 1. 移动grpc目录

#### V2.15.0
> 1. remove net/ip

#### V2.14.1
> 1. bazel编译 
#### V2.14.0
> 1. 拜年祭需求
#### V2.13.1
> 1. 修复缓存miss时可以重复投币的bug
#### V2.13.0
> 1. 移除掉对稿件的依赖
#### V2.12.0
> 1. 增加item/coins接口

#### V2.11.0
> 1. log增加oid

#### V2.10.4
> 1. 初始化dao层ut
> 2. 移除hbase

#### V2.10.3
> 1. 去掉remote ip

#### V2.10.2
> 1. gorpc不注册zk

#### V2.10.1
> 1. 修改proto文件

#### V2.10.0
> 1. 使用新verify
#### V2.9.2
> 1. 移动grpc目录

#### V2.9.1
> 1. 修复addCoins bug

#### V2.9.0
> 1. 使用新的目录结构
> 2. service统一名称
> 3. 重构grpc 增加req的验证
> 4. grpc增加client代码

#### V2.8.0
> 1. 使用es sdk查询日志

#### V2.7.0
> 1. 增加内网投币接口

#### V2.6.0
> 1. 新增grpc

#### V2.5.0
> 1. 去掉account service调用
#### V2.4.1
> 1. 去掉debug和trace配置

#### V2.4.0
> 1. 统一业务标志

#### V2.3.1
> 1. 用户投币列表加空缓存

#### V2.3.0
> 1. remove stat-T articleStat-T topic

#### V2.2.0
> 1. add todayExp rpc

#### V2.1.1
> 1. 修复硬币记录查询时间问题

#### V2.1.0
> 1. 硬币记录查询切换到es
> 2. 修改硬币数接口支持保存操作人

#### V2.0.7
> 1. 升级tools/cache

#### V2.0.6
> 1. 使用bm

#### V2.0.5
> 1. 增加register

#### V2.0.4
> 1. 修复小数问题

#### V2.0.3
> 1. 改硬币增加info日志

#### V2.0.2
> 1. 修复小数问题

#### V2.0.1
> 1. 修复热加载问题

### Version 2.0.0

> 1. 代码整体重构 增加prom监控 移除无用的代码
> 2. 硬币数*100的转换下沉到dao层去做 对service透明
> 3. business放入配置中 不再写死
> 4. 投币日志接入行为日志平台
> 5. 投币接口去掉登录态校验
> 6. 缓存工具重构代码
> 7. 修复经验先加缓存再加经验的问题 改为先加经验再加缓存
> 8. 投币改为事务操作 其他表的修改 放入job去做
> 9. 投币改为走消息队列 不再用异步channel防止消息丢失
> 10. 投币时去掉多余的投币总数查询sql
> 11. 投币记录增加日志转换功能 提供转换后的日志
> 12. 使用common/cache替换之前项目自己实现的异步channel
> 13. 修复up主硬币为负数 用户无法投币的问题

### Version 1.26.0

> 1. 增加UserLog rpc

### Version 1.25.8

> 1. 去掉反作弊日志上报

### Version 1.25.7

> 1. 修复帐号未激活的问题 

### Version 1.25.6
> 1. 使用account-service v7  

### Version 1.25.4
> 1. 修复第一次投币双倍的问题

### Version 1.25.3
> 1. 修复音乐验证问题

### Version 1.25.2
> 1. 硬币记录只返回一周数据  
> 2. 修复投币记录负数溢出的问题

### Version 1.25.1
> 1.fix coin zero val   

### Version 1.25.0
> 1.change user coin to memcache  

### Version 1.24.1
> 1.删除经验补偿逻辑  

### Version 1.23.0
> 1.添加checkzero检验  
> 2.增加硬币数返回值  

### Version 1.22.0
> 1.添加list接口  

### Version 1.21.0
> 1.添加修改硬币rpc方法  

### Version 1.20.0
> 1.硬币去除java依赖  

### Version 1.19.0
> 1.open auth  

### Version 1.18.3
> 1.修复时间戳转换  

### Version 1.18.2
> 1.修复硬币记录  

### Version 1.18.1 
> 1.补硬币记录  
### Version 1.18.0 
> 1.登陆奖励接口去重  

### Version 1.17.2
> 1.修复日志读取scan bug

### Version 1.17.1
> 1.增加用户ip获取  

### Version 1.17.0
> 1.双写硬币计数流  

### Version 1.16.2
> 1.音乐投币加经验  

### Version 1.16.1
> 1.ingore duplicate key.  

### Version 1.16.0
> 1.增加音频投币reason  

### Version 1.15.0 
> 1.添加稿件硬币数接口  

### Version 1.14.1
> 1. add log proc  

### Version 1.14.0
> 1. 用户硬币重构  

### Version 1.13.0 
> 1.添加硬币防刷  

### Version 1.12.1
> 1.修复rpc参数  

### Version 1.12.0
> 1.接入archive pb  
> 2.去除ding  

### Version 1.11.2
> 1.修复redis，没有设置过期时间  
> 2.优化稿件投币缓存，只展示最近20个稿件投币  

### Version 1.11.1
> 1.修复分表sharding  

### Version 1.11.0
> 1.迁移大仓库  

### Version 1.10.2 
> 1.修改硬币模板  

#### Version 1.10.1
> 1. 修复投币奖励模板  

#### Version 1.10.0 
> 1.文章投币 

#### Version 1.9.6

> 1.规范spy header

#### Version 1.9.5

> 1.兼容帐号硬币接口

#### Version 1.9.4

> 1.上报反作弊数据

#### Version 1.9.3

> 1.切换发布平台

#### Version 1.9.2

> 1.升级基础包

#### Version 1.9.1

> 1.修复硬币可以扣为负数

#### Version 1.9.0

> 1.稿件统计计数－大数据更新接口

#### Version 1.8.9

> 1.粉丝勋章

#### Version 1.8.8

> 1.update go-common and go-business

#### Version 1.8.7

> 1.接入prom监控  

#### Version 1.8.6

> 1.一天内同一个稿件的投币自动合并计数  

#### Version 1.8.3

> 1.剔除archive hbase计数和DB更新  

#### Version 1.8.2

> 1.更改missch为1024  

#### Version 1.8.0

> 1.接入新配置中心  

#### Version 1.7.8

> 1.修复hbase context  

#### Version 1.7.7

> 1.升级go-business到v2.7.0  
> 2.更新monitorPing  

#### Version 1.7.6

> 1.添加查询用户投币rpc接口

#### Version 1.7.5

> 1.剔除archivePGC上报大数据

#### Version 1.7.4

> 1.剔除account-service依赖

#### Version 1.7.2

> 1.支持配置中心版本

#### Version 1.7.1

> 1.更新go-common、golang、go-business包


#### Version 1.7.0

> 1.硬币库从老库迁移到新库  
> 2.新增稿件投币总数表  

#### Version 1.6.3

> 1.去掉databus v1  
> 2.去掉写daily库  
> 3.ding异步  

#### Version 1.6.2

> 1.升级vendor  

#### Version 1.6.1

> 1.升级vendor  
> 2.硬币数双写  

#### Version 1.6.0

> 1.databus双写V1和V2  
> 2.2016/11/25 00:00后写推送投币记录给coin-job  

#### Version 1.5.2

> 1.trace v2  

#### Version 1.5.1

> 1.fix添加硬币bug  

#### Version 1.5.0

> 1.临时版本2016/11/01 00:00之后投币将不再增加upper经验  
> 2.取消特权判断  

#### Version 1.4.0

> 1.硬币数写hbase由增量写改为写绝对值  

#### Version 1.3.1

> 1.databus消息新增type和pgc_id两个字段

#### Version 1.3.0

> 1.投币记录推送到databus给大数据动态使用  

#### Version 1.2.8

> 1.优化日志显示  

#### Version 1.2.7

> 1.修改投币上报描述  

#### Version 1.2.6

> 1. 修复稿件计数  

#### Version 1.2.5

> 1.增加稿件状态判断  

#### Version 1.2.4

> 1.删除多余参数  

#### Version 1.2.3

> 1.修复bug，初始化max_coin  

#### Version 1.2.2

> 1.允许PC端没绑定手机进行投币

#### Version 1.2.1

> 1.更新稿件投币动态  

#### Version 1.2.0

> 1.细化返回错误码  
> 2.接入trace2  
> 3.增加upper主视频收益

#### Version 1.1.1

>  1.兼容list请求参数错误  

####  Version 1.1.0

>  1.获取当日投币经验  

####  Version 1.0.0

>  1.投币业务重构  
>  2.投币业务接口  
>  3.投币历史记录查询  
