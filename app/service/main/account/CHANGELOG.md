### account service

### Version 7.17.0
> 1. 修复 mc sdk

### Version 7.15.0
> 1. profile 增加 is_tourist 字段
> 2. 移除 account-notify-job

### Version 7.14.2
> 1. 修正 gRPC 节操修改逻辑

### Version 7.14.1
> 1. goRpc方法添加过期备注

### Version 7.14.0
> 1. gRPC 重构

### Version 7.13.0
> 1. 增加 VIP 相关 gRPC 接口

### Version 7.12.1
> 1. 使用优先队列延迟删除缓存
> 2. 增加队列长度统计

### Version 7.11.0
> 1. gRPC服务迁移到v1

### Version 7.10.3
> 1. add vip decide method 
> 2. 实名认证空数据
> 3. 实名认证空指针检查

### Version 7.10.1
> 1. separate privacy http client

### Version 7.10.0
> 1. 增加 /privacy api

### Version 7.9.0
> 1. rpc.ProfileWithStat3 vip信息增加VipPayType
### Version 7.8.6
> 1. 封禁信息改为调用member

### Version 7.8.5
> 1. 删缓存错误以 cause 为准

### Version 7.8.4
> 1. 整理删缓存日志

### Version 7.8.3
> 1. remove cache save

### Version 7.8.2
> 1. logging delete cache error

### Version 7.8.1
> 1. 删除缓存日志

### Version 7.8.0
> 1. 去除显式 ip 参数传递

### Version 7.7.0
> 1. 修改部分基础库

### Version 7.6.0
> 1. 通过名字搜索用户信息走搜索接口

### Version 7.5.0
> 1. cards和myinfo接口调用go相关服务

### Version 7.4.0
> 1. 经验节操修改调用member

### Version 7.3.10
> 1.使用新版 warden
> 2.去除无用 error wrap

### Version 7.3.0
> 1.使用新版 warden

### Version 7.2.0
> 1.添加 GRPC 接口

#### Version 7.1.0
> 1.集成新的官方认证
> 1.官方认证兼容老接口结构

#### Version 7.0.0
> 1.删除大量无用代码
> 2.规范和统一字段名
> 3.使用最新bm和cache gen

#### Version 6.21.2
> 1.minify wallet cache duration  
> 2.remove statsd  

#### Version 6.21.1
> 1.reduce wallet cache duration  

#### Version 6.21.0
> 1.call account thin APIs for info, infos, myinfo, userinfo, profile.  
> 2.parallel call dependencies.  

#### Version 6.20.0
> 1.add userinfo rpc  

#### Version 6.19.0
> 1.del cache for action cleanAccessKey  

#### Version 6.18.0 
> 1.添加cardbyname接口  

#### Version 6.17.0
> 1.del unused code  

#### Version 6.16.0
> 1.add member proxy config  

#### Version 6.16.0
> 1.add member proxy mid end with 92  

#### Version 6.15.0
> 1.profile 硬币信息直接读取coin-service  

#### Version 6.14.0
> 1.更新b+缓存  

#### Version 6.14.0 
> 1.添加userinfo api接口 

#### Version 6.13.0
> 1.删除relation代码  

#### Version 6.12.0 
> 1.mid获取手机号及用户关系  

#### Version 6.11.0
> 1.add thin info rpc  

#### Version 6.10.4
> 1.fix wallet cache expiration  

#### Version 6.10.3
> 1.add vip info cache2 del  

#### Version 6.10.2
> 1.add vip info cache del  

#### Version 6.10.1
> 1.升级基础库  

#### Version 6.10.0
> 1.update http client conf and usage  
 
#### Version 6.9.0
> 1.支持分组缓存清理  

#### Version 6.8.0
> 1.myinfo接口调用account-java的获取用户信息myinfo内网接口  

#### Version 6.7.0
> 1.迁移info card 到member-service  

#### Version 6.6.0
> 1.add rpc token and group  

#### Version 6.5.0
> 1.add login empty cache  

#### Version 6.4.0
> 1.add http method get vip info  

#### Version 6.3.0
> 1.add rpc service method info3  

#### Version 6.2.0
> 1.adjust internal http router group  

#### Version 6.1.0
> 1.add new mc identify cluster  

#### Version 6.0.0
> 1.merge account-service into kratos  

#### Version 5.16.0
> 1.add unit tests for dao and service  

#### Version 5.15.0
> 1.add new mc identify cluster  

#### Version 5.14.0
> 1.更新mc基础包缓存  

#### Version 5.13.0
> 1.删除relation redis代码  

#### Version 5.12.5
> 1. 修复syslog panic error  

#### Version 5.12.4
> 1.修复go-common memcache stat err 

#### Version 5.12.3
> 1.添加memcache错误监控

#### Version 5.12.2
> 1.修复profile -404错误

#### Version 5.12.1
> 1.修复wallet的缓存错误

#### Version 5.12.0
> 1.http接口改为内网接口

#### Version 5.11.7
> 1.修复relation关注升级

#### Version 5.11.6
> 1.平滑发版
> 2.fix RESTFul vipinfo

#### Version 5.11.5
> 1.add mc statsd

#### Version 5.11.4
> 1.fix rpc server stat init

#### Version 5.11.3
> 1.fix rpc stat init

#### Version 5.11.2
> 1.升级RESTFul Get

#### Version 5.11.1
> 1.prom监控bugfix

#### Version 5.11.0
> 1.接入prom监控

#### Version 5.10.0  
> 1.添加attention返回所有关注列表，包括悄悄关注  

#### Version 5.9.0
> 1.增加黑名单rpc  
> 2.更新vendor  

#### Version 5.8.0
> 1.封装relation-service接口 

#### Version 5.7.5
> 1.fixed monitor ping
> 2.修改配置文件加载错误提示

#### Version 5.7.4
> 1.接入新的配置中心
> 2.接入rpc DiscoverOff配置
> 3.优化清除access cache逻辑

#### Version 5.7.3
> 1.account 添加用户auditInfo rpc接口

#### Version 5.7.2

> 1.接入新的配置中心
> 2.rpc 注册配置修改

#### Version 5.7.1

> 1.no change

#### Version 5.7.0

> 1.cache callback优化
> 2.升级dao层使用Get2
> 3.memcahed err 不再回写缓存
> 4.拆分account.go
> 5.去掉service请求dao的异常日志
> 6.删除不再使用的代码

#### Version 5.6.2

> 1.profile兼容逻辑

#### Version 5.6.1

> 1.no commit(docker file missing!!!)

#### Version 5.6.0

> 1.tw-proxy local deploy

#### Version 5.5.0

> 1.cal/del 维护新的mc集群  
> 2.删除thrift rpc代码  

#### Version 5.4.0

> 1.go-common stat小包合并大包  
> 2.go-common修复base64编码bug  

#### Version 5.3.0

> 1.更新最新的net/rpc  

#### Version 5.2.1

> 1.go-common修复trace race bug  

#### Version 5.2.0

> 1.支持优先从本地加载配置  
> 2.更新vendor，支持最新的rpc库  

#### Version 5.1.4

> 1.AccessInfo接口，cache缓存失败return改为回源帐号  
> 2.Secret接口走内存cache  

#### Version 5.1.3

> 1.profile接口支持VerifyOrUserGet  

#### Version 5.1.2

> 1.wallet过期时间改为CouponDueTime  
> 2.修复access删除bug  

#### Version 5.1.1

> 1.修复wallet清缓存bug  
> 2.更新vendor  

#### Version 5.1.0

> 1.支持平滑切流量和重启  

#### Version 5.0.0

> 1.net/rpc改为为golang/rpcx  

#### Version 4.5.0

> 1.接入配置中心  

#### Version 4.3.3

> 1.wallet字段json格式修正  

#### Version 4.3.2

> 1.idfSvc改成service实现,并支持两个新方法Access和Verify  
> 2.升级vendor  

#### Version 4.3.1

> 1.修复ak获取到用户存缓存0的bug  

#### Version 4.3.0

> 1.新增支持accesskey白名单  
> 2.支持b币新接口  

#### Version 4.2.1

> 1.去掉帐号错误码转换  

#### Version 4.2.0

> 1.添加accessInfo接口和获取secret的http接口  

#### Version 4.1.1

> 1.升级zk的包版本到最新  

#### Version 4.1.0

> 1.增加thrift的auth调用  

#### Version 4.0.1

> 1.更新所有匿名rpc client为默认user  

#### Version 4.0.0

> 1.引入vendor  
> 2.rpc修改syslog和上报  

#### Version 3.7.0

> 1.新增verify和secret的rpc接口用于identify  

#### Version 3.6.1

> 1.新增获取移动端头图的接口  

#### Version 3.6.0

> 1.VIP信息聚合  
> 2.新增名片信息和空间信息的聚合接口  
> 3.新增查询用户信息和是否关注up主的聚合接口  

#### Version 3.5.1

> 1.新增profile http接口  
> 2.新增批量获取relation  

#### Version 3.5.0

> 1.新增获取up主和粉丝的关系rpc接口  

#### Version 3.4.0

> 1.新增批量获取up主和粉丝的关系  
> 2.cache更新优化  

#### Version 3.3.1

> 1.fix relation bug

#### Version 3.3.0

> 1.增加用户是否关注up主的http接口  
> 2.增加批量获取用户信息接口的http接口  

##### Version 3.2.2

> 1.trace v2

##### Version 3.2.1

> 1.http client读写超时设置分开  

##### Version 3.2.0  

> 1.修复elk日志  
> 2.支持trace v2  

##### Version 3.1.0  

> 1.使用go-common/xlog  
> 2.修改删除cache接口为get  

##### Version 3.0.1

> 1.修复relation接口返回nil  

##### Version 3.0.0

> 1.context使用官方接口  
> 2.添加coin rpc  

##### Version 2.4.0

> 1.添加昵称查询card rpc接口  

##### Version 2.3.0

> 1.添加服务发现  

##### Version 2.2.2

> 1.修改配置信息，更合理  

##### Version 2.2.1

> 1.修复rpc返回指针赋值  

##### Version 2.2.0

> 1.优化  
> 2.add elk  
> 3.add trace id  
> 4.modify appkey  
> 5.add async update cache  
> 6.add friend relation  

##### Version 2.1.0

> 1.add tracer  

##### Version 1.1.0

> 1.基于go-common重构  

##### Version 1.0.0

> 1.初始化完成用户信息基础查询功能
