### 协管服务 (assist-service)

#### 2.0.23
>1.update account api from gorpc to grpc  

#### 2.0.22 
>1.完善dao层的单元测试

#### 2.0.21
>1.升级中间件值verify和bm 

#### 2.0.20
>1.fix bug：校验协管个数的时候以Up主的mid为主业务查询 

#### 2.0.19
>1.fix bug: syncChannel异步任务超时重新设置,添加缓存清理的日志 

#### 2.0.18
>1.merge master 分支，since 20180530

#### 2.0.17
>1.更新mc的生成算法，前面添加前缀:"assist_relation_mid_" 

#### 2.0.16
>1.fix bug: baldemaster默认返回json格式字符串，不返回空字符 

#### 2.0.15
>1.升级http组件到baldemaster

#### 2.0.14
> 1. 增加register

#### 2.0.13
> 1. 来自space空间产品线的要求: 返回当前被委任为协管的up主列表的时候，添加当前up主的vip信息 

#### 2.0.12
> 1. 按照服务树的规则将代码目录迁移到service/main下 

#### 2.0.11
> 1. protect assist ups cards3 rpc call 

#### 2.0.10
> 1. 修改account-service v7  

#### 2.0.9
> 1. v2.0.8 => v2.0.9      

#### 2.0.8
> 1.去掉无用代码，其实Profile2的RPC接口并没有调用，也不需要添加     

#### 2.0.7
> 1.添加对HttpClient增加key,secret的支持    

#### 2.0.6
> 1.business enhancement: 协管服务更新大仓库版本，并添加对dapper tracer的支持  

#### 2.0.5
> 1.Fix Bug: 取消Sleep 2秒, 会造成连续的延迟消费

#### 2.0.4
> 1.增加关系链的判断逻辑，action为update,并且是两个case   
> 2.调整fid和mid的顺序，fid永远是之前被关注的那个人，mid是去关注他人的人  
> 3.不关注表取mod的分配表名，只关注单向的关系链状态流转  

#### 2.0.3
> 1.添加Relation-T Topic消费数据的日志记录  
> 2.fix databus sub的bug, 不能直接return

#### 2.0.2
> 1.给空间的 rpc Ups提供官方认证的信息    

#### 2.0.1
> 1.rpc server开启Handshake, 支持rpc token   

#### 2.0.0
> 1.支持直播的房管功能    
> 2.为[空间]提供允许用户主动退出骑士团的HTTP接口，/x/internal/assist/exist，包括给space空间的rpc接口    
> 3.监听DataBus Topic:Relation-T, 在Up主动移除粉丝的同时监听消息，异步地接触其对应的骑士团资格  
> 4.已经被封禁的协管账号不允许被添加为骑士团  
> 5.同理4，如果之前是骑士，之后账号被封禁，那么也不允许进行后续对应的业务操作  
> 6.添加协管Model里面支持直播操作的常量，type=3; action=8/9
> 7.完善添加协管和增加日志时候的限制规则，如下：

         a: 单个协管关系操作同一类型的日志业务，每天不得超过100次
         b: 协管最多不超过10个
         c: 每天单个MID任命上限为100次，另外需要给Up主发系统通知， 不能超过100次/天
         d: 每日任命同一用户不超过2次
> 8.添加接口:按照mid(UP主mid)和assmids(批量协管)返回日志计数信息    
> 9.添加接口:获取当前用户是哪个up主的协管的mid集合，按照创建时间倒叙排, 添加key为assUps_的Rds缓存，加速查询并减轻数据库的查询压力,包括给space空间的rpc接口    

#### 1.1.0
> 1.启用新的统一的手机验证实名制限制策略   
> 2.更新govendor包， go-business和go-common
> 3.添加任命和卸载协管的时候，发【站内消息-系统通知】的内容

#### 1.0.8
> 1.去掉协管添加时候的实名认证限制   

#### 1.0.7
> 1.http /assist/logs 接口添加计数信息, 用于动态计数日志个数  

#### 1.0.6
> 1.更新缓存空数组，防止空值穿透到db
> 2.fix rpc IDs返回值错误的bug， 必须返回指针，必须是指针  

#### 1.0.5
> 1.添加Prom的基础组件依赖监控代码    

#### 1.0.4
> 1.memcache 判断有误，判断等于并制空  
> 2.info接口返回count值, 方便第三方业务自己根据设定的阈值来进行判断操作  

#### 1.0.3
> 1.更新vendor包go-business  
> 2.协管服务添加业务常量ActCancelDisUser,指代协管可以帮助Up主取消屏蔽用户  
> 3.添加协管的RPC接口:RPC.AssistIDs  

#### 1.0.2
> 1.go-common升级到6.17.0
> 2.为log-agent做准备
> 3.重构prom的接入方式

#### 1.0.1
> 1.合并到develop分支  

#### 1.0.0
> 1.初始化功能