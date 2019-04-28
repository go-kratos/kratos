# v4.10.0
1. 敏感词及时生效

# v4.9.0
1. main danmu ai filter

# v4.8.5
1. 解除敏感词任意风险等级硬编码, retag

# v4.8.4
1. 解除敏感词任意风险等级硬编码

# v4.8.3
1. white mid ignore 0

# v4.8.2
1. global white mid list

# v4.8.1
1. 补充 filter grpc 接口返回字段

# v4.8.0
1. 添加 hit v3 接口，返回 risk level 信息

# v4.7.0
1. add /Hit grpc api
2. add /MHit grpc api
3. 优化context使用方式

# v4.6.4
1. edit AI host URL

# v4.6.3
1. update hbase sdk

# v4.6.2
1. add /dao ut

# v4.6.1
1. discovery register grpc

# v4.6.0
1. 直播弹幕接入ai filter

# v4.5.0
1. 增加 grpc api
2. 干掉RPCClient配置
3. 修复deprecated :
    bm/identify --> bm/verify
    http.ServerConfig --> bm.ServerConfig

# v4.4.1
1. fix hit cache

# v4.4.0
1. 命中日志上报数据平台

# v4.3.2
1. delay chan call reply api

# v4.3.1
1. fix aifilter ai score bug

# v4.3.0
1. http router 切换为bm
2. 下线接口 :
    /x/filter*
    /x/internal/filter/rubbish*

# v4.2.2
1. fix mid=0

# v4.2.1
1. ai垃圾信息过滤模型优化

# v4.2.0
1. ai垃圾信息过滤模型接入

# v4.1.2
1. area 选项可配置化

# v4.1.1
1. add register

# v4.1.0
1. area维度走pcre-jit, key维度走官方regexp
2. make update

# v4.0.5
1. 忽略反垃圾报错

# v4.0.4
1. 优化反垃圾生效策略

# v4.0.3
1. pcre-jit update

# v4.0.2
1. update replace

# v4.0.1
1. 对于长文本加缓存

# v4.0.0
1. use pcre-jit

# v3.12.9
1. rpc model 解耦(解决依赖方构建依赖pcre)

# v3.12.8
1. 修复 pcre panic
2. 对于 pcre match 返回值异常，走默认不匹配逻辑

# v3.12.7
1. 更改正则引擎为PCRE-JIT
2. 去掉baselibParam、bizlibParam参数
3. 去掉msg为空报错
4. 批量过滤改为并行执行
5. 业务小文本加上mc cache

# v3.12.6
1. 支持b+反垃圾
2. 支持b+过滤返回全level

# v3.12.5
1. 取消文本分段及分布式过滤
2. 优化分段逻辑（备用）

# v3.12.4
1. fix log
2. 对直播简介和hit的大文本使用分布式过滤

# v3.12.3
1. 对直播简介做长文本分段过滤

# v3.12.2
1. config gomaxproce
# v3.12.1
> 1.remove stastd
> 2.优化regexp编译错误日志

# v3.12.0
> 1.直播接入重构
> 2.修改/x/internal/filter/test(GET -> POST)
> 3.下线/x/internal/filter/key
> 4.下线/x/internal/filter/key/dm

# v3.11.1
> 1.修复弹幕后台命中测试接口

# v3.11.0
> 1.弹幕统一过滤接口
> 2.弹幕反垃圾接入

# v3.10.1
> 1.增加若干测试用例

# v3.10.0
> 1.删除HTTP接口/x/internal/filter/v2
> 2.删除HTTP接口/x/internal/filter/v2/post
> 3.删除HTTP接口/x/filter/v2
> 4.删除HTTP接口/x/filter/v2/post
> 5.修改HTTP接口/x/internal/filter
> 6.修改HTTP接口/x/internal/filter/post
> 7.修改HTTP接口/x/internal/filter/multi
> 8.修改HTTP接口/x/internal/filter/mpost
> 9.修改RPC接口RPC.Filter

# v3.9.0
> 1.去掉反垃圾功能，转为接入antispam

# v3.8.7
> 1.update hit log

# v3.8.6
> 1.优化专栏hit接口

# v3.8.5
> 1.优化identify

# v3.8.4
> 1.update Xlog to Log

# v3.8.3
> 1.加入清理敏感词开关

# v3.8.2
> 1.清理已删除敏感词异步化

# v3.8.1
> 1.增加启动清理已删除敏感词

# v3.8.0
> 1.新增敏感词通知search

# v3.7.5
> 1.优化admin逻辑和代码,并解决一些已知bug

# v3.7.4
> 1.新增反垃圾命中日志

# v3.7.3
> 1.新增禁止添加重复敏感词和白名单

# v3.7.2
> 1.敏感词命中统计prom

# v3.7.1
> 1.弹幕命中日志

# v3.7.0
> 1.优化后台接口：测试、敏感词增删改、白名单增改

# v3.6.0
> 1.新增v2命中查询接口

# v3.5.2
> 1.修复dao panic

# v3.5.1
> 1.剔除model对conf依赖

# v3.5.0
> 1.新增反垃圾接口

# v3.4.0
> 1.优化反垃圾代码
> 2.LimitType新增area进行拆分

# v3.3.1
> 1.修复过滤词生效时间

# v3.3.0
> 1.迁移hbase

# v3.2.0
> 1.修改反垃圾databus推送时机

# v3.1.0
> 1.更精准的过滤词星号替换

# v3.0.11
> 1.分段加载反垃圾字典树

# v3.0.10
> 1.反垃圾剔除无用的数据库查询

# v3.0.9
> 1.修复白名单命中level值错误

# v3.0.8
> 1.提升图文接口性能

# v3.0.7
> 1.所有http接口走verify

# v3.0.6
> 1.新增filter v2接口，优化反垃圾返回值
> 2.Article 接口去掉common过滤

# v3.0.5
> 1.整理filter过滤逻辑

# v3.0.4
> 1.修复panic

# v3.0.3
> 1.修复反垃圾逻辑
> 2.修复MFilterArea,FilterArea逻辑
> 3.增加filter命中log

# v3.0.2
> 1.update internal IP

# v3.0.1
> 1.list接口默认area

# v3.0.0
> 1.屏蔽词添加来源，类型
> 2.屏蔽词添加level 15处理逻辑

# v2.9.4
> 1.修复直接加黑逻辑

# v2.9.3
> 1.增加databus log

# v2.9.2
> 1.修改rubbish filter position

# v2.9.1
> 1.修改ecode

# v2.9.0
> 1.白名单

# v2.8.0
> 1.垃圾过滤

# v2.7.1
> 1.增加批量过滤业务post接口

# v2.6.0
> 1.增加批量过滤post接口

# v2.5.0
> 1.升级vendor，提供不走common过滤rpc

# v2.4.2
> 1.修复部分分区过滤错误

# v2.4.1
> 1.升级vendor

# v2.3.3
> 1.修复弹幕过滤命中返回null rule

# v2.3.2
> 1.变量函数名优化
> 2.修复一些小bug
> 3.增加service部分单元测试

# v2.2.2
> 1.白名单越界优化

# v2.2.1
> 1.并行正则优化后台图文过滤接口

# v2.2.0
> 1.并行正则优化

# v2.1.1
> 1.dm过滤优先进行业务方过滤规则

# v2.1.0
> 1.修改日志
> 2.增加文章过滤

# v2.0.7
> 1.修改锁

# v2.0.6
> 1.修复链表循环问题

# v2.0.5
> 1.修复lru链表需要加锁

# v2.0.4
> 1.修复comment字符长度限制

# v2.0.3
> 1.修复编辑不认证正则合法性

# v2.0.2
> 1.修复后台过滤测试文字显示不全bug,编译正则bug

# v2.0.1
> 1.修复cid到64

# v2.0.0
> 1.增加key过滤
> 2.修改字典树到ac字典树
> 3.升级vendor，接入新的配置中心

# v1.6.3
> 1.过滤接口增加post支持

# v1.6.2
> 1.升级vendor

# v1.6.1
> 1.兼容老后台同步

# v1.6.0
> 1.添加过虑原文保存
> 2.添加过虑原文接口

# v1.5.0
> 1.添加操作记录
> 2.添加编辑功能
> 3.字符串模式自动忽略大小写
> 4.敏感词检测

# v1.4.0
> 1.增加关键词查询

# v1.3.0
> 1.增加批量过滤
> 2.增加返回命中分区

# v1.2.2
> 1.进程退出增加rpc close

# v1.2.1
> 1.增加日志

# v1.2.0
> 1.增加ping方法
> 2.升级vendor

# v1.1.1

> 1.修复局部变量导致没赋值到全局问题

# v1.1.0

> 1.增加rpc服务

# v1.0.3
> 1.修改lelvel返回类型

# v1.0.2

> 1.去掉老的slbchecker

# v1.0.1

> 1.更新vendor
> 2.配置中心

# v1.0.0
> 1.过滤服务初始化
