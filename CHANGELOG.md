### Go-Common
### Version 7.19
1. 更新vendor下github.com/xanzy/go-gitlab包

### Version 7.18
1. 删除无用的git文件

### Version 7.17
1. 添加 WithContext 方法, Conn 待实现
2. 使用 traceConn 实现 trace 埋点
3. GetMulti 检查 key 列表不为空, 空列表直接返回空 map
4. legalKey 检查 key 不为空

### Version 7.16.8
1. 修改saga check无需在master分支跑

### Version 7.16.7
1. pipeline编译暂时去掉merge master

### Version 7.16.5
1. pipeline增加bazel局部编译job

### Version 7.16.4
1. pipeline增加unitest的执行stage

### Version 7.16.3
1. gometalinter暂时不检查pb

### Version 7.16.3
1. 修正gometalinter的bug
2.合并make update到COMPILE里

### Version 7.16.2
1. 增加saga-admin的功能

### Version 7.16.1
1. 将编译次数改为一次
2. 增加gometalinter

### Version 7.15.2
1. 树状结构展示个业务方权限责任人

### Version 7.15.1
1. 去掉发邮件的stage
2. 调整makeupdate为独立的stage

### Version 7.14.3
1. 由于历史代码太多无法通过lint，所以暂时允许lint失败
2. 增加将lint不通过的文件列表打印出来

### Version 7.14.2
1. Runner触发以后，先Merge Master，然后触发Compile

### Version 7.14.1
1. 为bazel info添加-k参数

### Version 7.14.0 
1. memcache large value storage

### Version 7.13.2
1. 修复ecode.OK也打Error日志  

### Version 7.13.1 
1. 调整cache chan prom数据记录时机   

### Version 7.13.0 
1. 全链路传递超时时间  
 
  ### Version 7.13.1
1. HTTP Client 升级，修复此升级导致的identity基础库bug

 ### Version 7.13.0
1. HTTP Client 升级，通用参数统一配置

### Version 7.12.0
1. log保留\n，替换\r为空  

### Version 7.11.0
1. 迁移account interface到Kratos

### Version 7.10.1
1. 修复hbase依赖btree

### Version 7.10.0
1. 迁移account service到Kratos

### Version 7.9.0
1. 增加stack 信息记录，
2. 修改statHandler为一个方法调用

### Version 7.8.1
1. 更新ecode的readme和changelog

### Version 7.8.0
1. 增加pkg/errors 包，用于记录错误信息堆栈信息
2. 在ecode的example中增加error使用example

### Version 7.7.2
1. 修复了memcache存入数据, Object & Value 为nil的情况

### Version 7.7.1
1. 修复redigo
2. 修复了databus sdk Commit的饥饿问题

### Version 7.7.0
1. 迁移golang库里的redigo到go-common里的cache/redis
2. 修复了cache/redis 普罗米修斯耗时、异常上报

### Version 7.6.0
1. 新增spy service rpc client
2. 增加history

### Version 7.5.0
> 1. cache/memcache 升级支持protobuf
> 2. cache/memcache 破坏性增加了conn.Scan去掉了item.Scan方法
> 3. business/client/identify 缓存gob改成了protobuf
> 4. vendor新增了github.com/gogo/protobuf的依赖


### Version 7.4.2
> 1. 更新vendor里golang库到最新版

### Version 7.4.1
> 1. 更新ecode文档

### Version 7.4.0
> 1. ecode 获取code message由从数据库全量更新改为通过接口增量更新
> 2. 升级配置，不兼容老的版本，参考 http://info.bilibili.co/pages/viewpage.action?pageId=3684076

### Version 7.3.0
> 1.支持vendor
> 2.继承了location-service

### Version 7.2.0
> 1.big-repo ，修改business目录  


### Version 7.1.0
> 1.添加secure model  
> 2.修改location model  

### Version 7.0.0
> 1.合并go-business
> 2.合并rouer
> 3.拆分interceptor

### Version 6.24.5
> 1.http client的breaker状态变更支持上报prometheus  

### Version 6.24.4
> 1.rpc server支持recover

### Version 6.24.3
> 1.修复databus中offset为0不能commit的问题

### Version 6.24.2
> 1.增加syscall/signal 对 macos(darwin) 的支持

### Version 6.24.1
> 1.强制要求http和rpc client设置breaker，否则会运行panic

### Version 6.24.0
> 1.去disconf，使用config-service SDK作为唯一Client

### Version 6.23.0
> 1.memcache新增序列化和压缩

> 2.新版memcache接口

> 3.net/http新增错误普罗米修斯上报

<b>memcache不再兼容，带有破坏性修改！！！！！！！</b>

### Version 6.22.2
> 1.去掉vendor

### Version 6.22.1
1. 修复 syslog 在linux环境下，空指针错误

### Version 6.22.0
1.新增vendor支持第三方依赖包

### Version 6.21.1
1.fix mc Stat, 以及增加单元测试

### Version 6.21.0
兼容了windows，编译：
1. 增加了Windows上Signal信号处理的Fake方法；
2. 增加了Syslog兼容的Fake方法；

喜欢windows开发的同学，可以

syslog -> go-common/syslog（syslog日志收集）；

os/Signal ->go-common/os/signal，syscall -> go-common/syscall（信号处理）；


### Version 6.20.0
> 1.迁移golang库中的gomemcache，交由go-common/cache/memcache维护；
> 2.优化了net/trace包内私有方法；

### Version 6.19.2
> 1.修复database/sql Stmt函数漏初始化db变量导致的panic

### Version 6.19.1
> 1.解决先前版本readme的冲突  

### Version 6.19.0
> 1.add RESTful httpclient

### Version 6.18.0 
> 1.修复mysql lifetime,迁移mysql配置  

### Version 6.17.1

> 1.修复log-agent sdk收集日志中有换行符未转义的bug

### Version 6.17.0

> 1.新增log-agent日志收集sdk，以unix socket方式发送日志

### Version 6.16.0

> 1.修改httpconf  
> 1.改为读写锁读取配置  

### Version 6.15.0

> 1.修改net/netutil熔断器支持全局开关

### Version 6.14.0

> 1.config sdk增加读取appoint参数，用作回退时读取指定配置文件
### Version 6.14.0

> 1.config sdk增加读取appoint参数，用作回退时读取指定配置文件

### Version 6.13.0

> 1.调整router handler参数，将函数内部join pattern改为外部传入完整pattern

### Version 6.12.1

> 1.修复db 事务初始化的bug  

### Version 6.12.0

> 1.stat支持prometheus功能，实现统计和监控  

### Version 6.11.0
> 1. 对reids进行了修改，以后不依赖conf包了，配置直接写在redis本包
 
### Version 6.10.0
> 1. 增加rpc sharding

### Version 6.9.0

> 1.配置中心client 启动参数增加token字段，区分应用和环境

### Version 6.8.0

> 1.依赖zookeeper的rpc client由连接池改为单连接  
> 2.breaker新增了callback，通知状态变更

### Version 6.7.2
> 1.修改zookeeper注册参数  

### Version 6.7.1

> 1.配置中心增加获得配置文件路径方法

### Version 6.7.0

1.fix rpc权重为0时，client不创建长连接  
2.rpc增加配置是否注册zookeeper  

### Version 6.6.4

> 1.fix mc expire max ttl

### Version 6.6.3

> 1. 将配置中心启动参数设置成和disconf的一样

### Version 6.6.2

> 1. 优化了net/http Client的buffer过小导致的syscall过多

### Version 6.6.1

> 1.fix http client超时设置不准确的问题，去掉了读包体和反序列化的时间  

### Version 6.6.0
> 1.rpc Broadcast 添加reply参数,支持对任意方法进行广播  

### Version 6.5.2

> 1.fix 新版配置中心和老版本init冲突问题

### Version 6.5.1

> 1.fix rpc Boardcast的bug  

### Version 6.5.0
> 1. 新版本配置中心conf/Client  

### Version 6.4.1
> 1. 修复remoteip获取  

### Version 6.4.0
> 1. 去除rpcx  

### Version 6.3.1

> 1.fix配置文件名覆盖的问题  

### Version 6.3.0
> 1. net/rpc支持了Boardcast广播调用

### Version 6.2.5
> 1. net/rpc支持了group路由策略

### Version 6.2.4
> 1. 优化了statsd批量发包

### Version 6.2.3
> 1. 修复了trace comment 在annocation的bug

### Version 6.2.2
> 1. 优化了net/rpc反射带来的性能问题
> 2. net/rpc内置了ping

### Version 6.2.1
> 1. 临时加回net/rpcx, TODO remove
> 2. net/trace.Trace2 奔溃和race修复

### Version 6.2.0
> 1. 去除了net/rpcx

### Version 6.1.3
> 1. 新增了memcache Get2/Gets

### Version 6.1.2
> 1. net/rpc使用CPU个数建立连接

### Version 6.1.1
> 1. 兼容net/rpc server的Client trace传递

### Version 6.1.0
> 1. 升级databus sdk，注意配置文件有变更

#### Version 6.0.0

> 1. xtime->time, xlog->log perf->net/http/perf
> 2. rpc支持设置方法级别超时
> 3. rpc支持breaker熔断
> 4. database 修复Row和标准库不兼容，使用database Rows替换标准库的Rows使用
> 5. 新的rpc框架net/rpc
> 6. net/trace支持Family初始化

#### Version 5.2.2

> 1.Zone结构体加json tag  

#### Version 5.2.0

> 1.更改http包名和路径  
> 2.增加http单元测试  
> 3.statd去掉hostname  
> 4.ip结构体增加isp字段  

#### Version 5.1.2

> 1.xip改为支持对象访问，去掉全局对象和函数  

#### Version 5.1.1

> 1.修复上报trace的位置  

#### Version 5.1.0

> 1.支持熔断  
> 2.rpc server判断zk是否注册  
> 3.修复Infoc连接重连  
> 4.xhttp xrpc xweb改为httpx rpcx webx  
> 5.修复trace level的bug  

#### Version 5.0.0

> 0.注意一定要使用Go1.7及以上版本  
> 1.用golang/rpcx替换官方库  
> 2.使用go1.7的context包  
> 3.增加traceon业务监控上报  
> 4.xhttp中ip方法挪到xip包  
> 5.rpc服务暴露close接口  
> 6.修复ugc配置中心等待30s的bug  
> 7.修复rpc client因权重变更导致panic的bug  
> 8.使用context.WithTimeout替代timer  

#### Version 4.4.1

> 1.日志新增按文件大小rotate  

#### Version 4.4.0

> 1.infoc支持udp和tcp方式  
> 2.去掉stdout、stderr输出到syslog的逻辑  

#### Version 4.3.2

> 1.fix rpc timeout连接泄露的bug  
> 2.rpc单连接改为多连接  

#### Version 4.3.1

> 1.支持从环境变量获取配置  
> 2.syslog支持打印标准输出和错误  

#### Version 4.3.0

> 1.支持配置中心  

#### Version 4.2.0

> 1.修复xredis keys的bug  
> 2.修复xmemcache批量删除bug  
> 3.新增 databus v2 客户端   

#### Version 4.1.3

> 1.trace 优化  
> 2.去掉sp 运营商字段  

#### Version 4.1.2

> 1.trace id改为int64  
> 2.trace http client增加host  
> 3.ip新增运营商字段  

#### Version 4.1.1

> 1.fix kafka monitor  

#### Version 4.1.0

> 1.去掉ecode和router  

### Version 4.0.0

> 1.business移到go-business  
> 2.新增InternalIp()获取本机ip  
> 3.rpc ping加超时  
> 4.增加ecode配置  
> 5.新增支持syslog  

#### Version 3.6.6

> 1.修复xip边界值时死循环问题  

#### Version 3.6.5

> 1.space接口只保留s_img、l_img  
> 2.archive-service新增viewPage的rpc方法  

#### Version 3.6.4

> 1.VIP相关接口及错误码  

### Version 3.6.3

> 1.修复ip递归查找导致的栈溢出  

#### Version 3.6.2

> 1.account-service profile的http接口、批量获取relation接口  
> 2.账号新增official_verify字段  

#### Version 3.6.1

> 1.修复degrade中变量名错误  
> 2.简化redis的auth逻辑，使用option
