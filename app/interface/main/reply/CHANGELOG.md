### 评论的Gateway服务

#### Version 7.0.11
> 1.修复filter过滤cd导致的bug

#### Version 7.0.9
> 1.增加调整热评

#### Version 7.0.8
> 1.修复report oid和rpid 不一致

#### Version 7.0.7
> 1.fix -400

#### Version 7.0.6
> 1.fix reist bug

#### Version 7.0.5
> 1.fix set state bug

#### Version 7.0.4
> 1.优化依赖服务调用

#### Version 7.0.3
> 1.增加bnj视频热评数

#### Version 7.0.2
> 1.fix duplicate

#### Version 7.0.1
> 1.修复判断逻辑

#### Version 7.0.0
> 1.折叠评论及评论新接口

#### Version 6.1.12
> 1.account grpc

#### Version 6.1.11
> 1. 相簿去掉topics

#### Version 6.1.9
> 1. remove user action table

#### Version 6.1.8
> 1.no topics for type=17

#### Version 6.1.7
> 1.fil out deleted reply

#### Version 6.1.6
> 1.max_id

#### Version 6.1.5
> 1.下掉limit

#### Version 6.1.4
> 1.修改hots batch接口 mid非必传

#### Version 6.1.3
> 1.修改hots batch接口增加mid

#### Version 6.1.2
> 1.abtest

#### Version 6.0.41
> 1.国际版接口返回空取消

#### Version 6.0.40
> 1.bnj topics and reply_admin_log

#### Version 6.0.39
> 1.国际版表情限制取消

#### Version 6.0.38
> 1.bbq批量热评接口

#### Version 6.0.37
> 1.火鸟举报

#### Version 6.0.34
> 1.workflow接入火鸟

#### Version 6.0.33
> 1.notice判断

#### Version 6.0.32
> 1.增加对热门评论点赞数的判断

#### Version 6.0.31
> 1.up置顶被管理员删除后仅up可见

#### Version 6.0.30
> 1.up置顶被删除后不可见

#### Version 6.0.29
> 1.小黄条支持模板

#### Version 6.0.28
> 1.location的Zone方法改为Info方法

#### Version 6.0.27
> 1.新增是否是热评的接口

#### Version 6.0.26
> 1.fix bug

#### Version 6.0.25
> 1.针对国际版做表情兼容

#### Version 6.0.24
> 1.针对BBQ和火鸟放开账号限制

#### Version 6.0.23
> 1.修复根评论小于20条时不显示热门评论

#### Version 6.0.22
> 1.评论cd调整，删除cd code如果人证成功

#### Version 6.0.15
> 1.已删除的评论不被举报

#### Version 6.0.12
> 1.add new type

#### Version 6.0.11
> 1.行为日志增加端口信息

#### Version 6.0.7
> 1.remove like user action

#### Version 6.0.6
> 1.reply workflow

#### Version 6.0.3
> 1.新的type

#### Version 6.0.2
> 1.remove emoji for ios

#### Version 6.0.1
> 1.reply dialog, 修复parent

#### Version 6.0.0
> 1.reply dialog, 5.33

#### Version 5.28.17
> 1.去掉spam检查为了测试

#### Version 5.28.16
> 1.reply_record es

#### Version 5.28.16
> 1.reply_record es

#### Version 5.28.15
> 1.新算法5条

#### Version 5.28.13
> 1.热评5条

#### Version 5.28.12
> 1.MC回源改成从库

#### Version 5.28.10
> 1.去掉5条热评ABtest

#### Version 5.28.7
> 1.点赞增加upmid字段

#### Version 5.28.5
> 1.log.Warn

#### Version 5.28.4
> 1.去掉remote ip，改为context metadata

#### Version 5.28.3
> 1.为reply record增加html escape功能

#### Version 5.28.2
> 1.添加过虑空白字符

#### Version 5.28.1
> 1.增加火鸟项目type

#### Version 5.28.0
> 1.增加feature flag

#### Version 5.27.18
> 1.去掉remote ip，改为context metadata

#### Version 5.27.17
> 1.new trace

#### Version 5.27.16
> 1.子评论游标接口，定位时兼容rootID
> 2.升级bm server

#### Version 5.27.16
> 1.灰度trace

#### Version 5.27.14
> 1.change identify to grpc

#### Version 5.27.13
> 1.配置文件
> 2.xint move to model

#### Version 5.27.11
> 1.fix buvid get method

#### Version 5.27.10
> 1.reply hots ab test

#### Version 5.27.8
> 1.fix golint

#### Version 5.27.7
> 1.注册用户开放评论点赞点踩

#### Version 5.27.6
> 1.过滤评论新增mid

#### Version 5.27.5
> 1.点赞去重

#### Version 5.27.4
> 1.过滤评论入库

#### Version 5.27.3
> 1.新增行为日志

#### Version 5.27.2
> 1.修复子游标定位不足一页时展示够

#### Version 5.27.1
> 1.修复子游标定位传根评论

#### Version 5.27.0
> 1.修复子游标接口返回非根评论
> 2.添加db从库读取

#### Version 5.26.3
> 1.修复表情包remark

#### Version 5.26.2
> 1.修复表情包为null

#### Version 5.26.1
> 1.vip评论表情由请求http改为查询本地DB

#### Version 5.25.7
> 1.rebuild master

#### Version 5.25.6
> 1.es迁移v3

#### Version 5.25.5
> 1.rebuild master

#### Version 5.25.4
> 1.fix

#### Version 5.25.3
> 1.删除非块级配置

#### Version 5.25.2
> 1.添加大数据过虑开关

#### Version 5.25.1
> 1.fix点赞新平台

#### Version 5.25.0
> 1.点赞新平台

#### Version 5.24.22
> 1.添加vegas限流

#### Version 5.24.21
> 1.平台日志修改配置，删除rpc common conf

#### Version 5.24.20
> 1.平台日志换用Databus

#### Version 5.24.19
> 1.点赞双写

#### Version 5.24.18
> 1.修复bug

#### Version 5.24.17
> 1.新增business支持

#### Version 5.24.16
> 1.撤销点赞双写

#### Version 5.24.15
> 1.修复context deadline

#### Version 5.24.14
> 1.增加漫画类型

#### Version 5.24.13
> 1.点赞报错忽略

#### Version 5.24.12
> 1.点赞双写

#### Version 5.24.11
> 1.修复已经删除评论可以添加

#### Version 5.24.10
> 1.修复冻结ecode

#### Version 5.24.9
> 1.修复bm internal context mid取不到

#### Version 5.24.8
> 1.修复重构中出现的问题

#### Version 5.24.7
> 1.修复迁移到bm

#### Version 5.24.6
> 1.迁移到bm

#### Version 5.24.5
> 1.迁移到databus

#### Version 5.24.4
> 1.修复举报日志index类型

#### Version 5.24.3
> 1.修复举报日志content

#### Version 5.24.2
> 1.增加用户举报日志以及用户置顶评论和取消置顶评论日志

#### Version 5.24.0
> 1.修复基础库no rpc client

#### Version 5.24.0
> 1.游标model重构
> 1.添加子评论列表游标接口

##### Version 5.23.5
> 1.增加日志输出用于统计

#### Version 5.23.4
> 1.修复state接口-404

#### Version 5.23.3
> 1.使用discovery

#### Version 5.23.2
> 1.HTTP新鉴权修复

#### Version 5.23.1
> 1.HTTP使用新鉴权

#### Version 5.23.0
> 1.评论分段回源

#### Version 5.22.1
> 1.修复子评论已经删除跳转还出现

#### Version 5.22.0
> 1.增加对-1的兼容

#### Version 5.21.7
> 1.修复up稿件对等级、rank限制

#### Version 5.21.6
> 1.使用account-service v7

#### Version 5.21.5
> 1.添加评论举报理由新增“青少年不良信息”
> 2.用户等级需大于等于Lv2才可发表评论，up主不受限制发评论

#### Version 5.21.4
> 1.修改已经删除评论还在列表中

#### Version 5.21.3
> 1.修复空缓存时间为30s

#### Version 5.21.2
> 1.修复reply/log显示问题

#### Version 5.21.1
> 1.修复游标热门跟分页不一致

#### Version 5.21.0
> 1.添加topic标签支持
> 2.删除无用配置

#### Version 5.20.7
> 1.Clone reply fix
> 2.Call new captch api

#### Version 5.20.6
> 1.CacheSave context fix

#### Version 5.20.5
> 1.添加依赖服务白名单

#### Version 5.20.4
> 1.添加依赖服务白名单

#### Version 5.20.3
> 1.更换myinfo为userinfo rpc

#### Version 5.20.2
> 1.修复add mc为nil

#### Version 5.20.1
> 1.添加hmget参数校验
> 2.fix reply roots为null的问题

#### Version 5.20.0
> 1.优化置顶评论

#### Version 5.19.12
> 1.添加size小于等于0时校验

#### Version 5.19.11
> 1.修复置顶评论和根评论覆盖bug

#### Version 5.19.10
> 1.游标优化第三方请求次数

#### Version 5.19.9
> 1.游标避免子评论为空时回源

#### Version 5.19.8
> 1.游标插入楼层排序修复

#### Version 5.19.7
> 1.修复prom nouser

#### Version 5.19.6
> 1.游标角标越界修复

#### Version 5.19.5
> 1.游标角标越界修复

#### Version 5.19.4
> 1.记录脏数据日志

#### Version 5.19.3
> 1.游标修改空指针bug

#### Version 5.19.2
> 1.游标删除多余返回参数(previous/next link)

#### Version 5.19.1
> 1.游标添加跳转接口
> 2.添加批量返回评论数新接口
> 3.internal router添加热门评论接口

#### Version 5.19.0
> 1.添加一二审转审
> 2.增加举报日志

#### Version 5.18.4
> 1.接入filter-service新接口

#### Version 5.18.3
> 1.支持seq rpc sharing.

#### Version 5.18.1
> 1.使用取号器id int32

#### Version 5.18.0
> 1.使用新seq取号器

#### Version 5.17.0
> 1.添加热评SEO接口
> 2.添加音乐播单类型

#### Version 5.16.4
> 1.游标接口修复返回错误码

#### Version 5.16.3
> 1.游标接口添加评论数
> 2.游标sql去掉attr

#### Version 5.16.2
> 1. 将多次断言合并为一次

#### Version 5.16.1
> 1. 删除多余的 ecode.NoLogin

#### Version 5.16.0
> 1.添加评论记录列表接口
> 2.修复relation的mid为0传透

#### Version 5.15.1
> 1.数据库回源优化

#### Version 5.15.0
> 1.添加举报信用评分

#### Version 5.14.0
> 1.公告添加海外判断

#### Version 5.13.1
> 1.添加游标接口

#### Version 5.12.2
> 1.修复验证码mc close

#### Version 5.12.1
> 1.修复定位评论分页数据

#### Version 5.12.0
> 1.定位评论添加热门和置顶

#### Version 5.11.0
> 1.二级评论楼主和热评楼主加关注关系

#### Version 5.10.4
> 1.修复发评论效验无效

#### Version 5.10.0
> 1.添加举报相关接口

#### Version 5.9.2
> 1.修复过虑ecode check

#### Version 5.9.1
> 1.修复账号active判断

#### Version 5.9.0
> 1.添加实名认证

#### Version 5.8.1
> 1.修复是评论注册缓存

#### Version 5.8.0
> 1.迁移大仓库，更新memcache
> 2.修复重复举报ttl

#### Version 5.7.1
> 1.修复内容为空
> 2.修复Version ip冻结状态下无法评论

#### Version 5.7.0
> 1.去除topic老库依赖

#### Version 5.6.2
> 1.修复过虑限制

#### Version 5.6.1
> 1.修复评论列表

#### Version 5.6.0
> 1.重构评论列表
> 2.添加先审后发

#### Version 5.4.0
> 1.添加粉丝勋章名称及其他参数

#### Version 5.3.3
> 1.修改过滤逻辑

#### Version 5.3.1
> 1.升级filter-serVersion ice ecode

#### Version 5.3.0
> 1.添加对反垃圾支持

#### Version 5.2.2
> 1.修复关闭评论区，up主无法删除评论

#### Version 5.2.1
> 1.添加获取评论删除日志配置入口
> 2.新增内部接口(info,count,minfo,mcount)

#### Version 5.2.0
> 1.评论日志删除管理入口
> 2.去掉老黑名单http接口
> 3.协管批量获取接口替换

#### Version 5.1.0
> 1.新增web端jsonp调用emojs表情列表接口(CDN缓存优化)

#### Version 5.0.5
> 1.修复后台批量通过接口

#### Version 5.0.4
> 1.修复评论内容表情转码问题

#### Version 5.0.3
> 1.粉丝列表intimacy类型值溢出

#### Version 5.0.2
> 1.评论项目int值溢出，粉丝接口普罗米修斯监控

#### Version 5.0.1
> 1.添加获取粉丝勋章列表日志和业务类型

#### Version 5.0.0
> 1.评论添加获取粉丝勋章列表

#### Version 4.9.1
> 1.修复评论添加协管员日志

#### Version 4.9.0
> 1.评论添加协管员删除评论逻辑

#### Version 4.8.7
> 1.修改Version endor中的Shopify依赖冲突

#### Version 4.8.6
> 1.升级go-common,go-business依赖（普罗米修斯，logagent）

#### Version 4.8.5
> 1.添加文章类型

#### Version 4.8.4
> 1.修复数据年份

#### Version 4.8.3
> 1.修复mcount接口返回禁止评论

#### Version 4.8.2
> 1.修复敏感词打码

#### Version 4.8.1
> 1.敏感视频关闭评论＆优化

#### Version 4.8.0
> 1.添加周年庆所需用户信息字段

#### Version 4.7.9
> 1.添加稿件发评论黑名单列表

#### Version 4.7.8
> 1.网监需求，我和你一样也不能理解

#### Version 4.7.7
> 1.黑名单rpc降级

#### Version 4.7.6
> 1.评论获取黑名单调用rpc接口

#### Version 4.7.5
> 1.降级黑名单

#### Version 4.7.4
> 1.修改分页获取黑名单数量

#### Version 4.7.3
> 1.修改热评Version 2.0逻辑

#### Version 4.7.2
> 1.新增评论业务类型

#### Version 4.7.1
> 1.接入配置中心

#### Version 4.7.0
> 1.评论新增黑名单过滤接口

#### Version 4.6.7
> 1.去除稿件依赖

#### Version 4.6.3/6
> 1. csrf日志

#### Version 4.6.2
> 1.添加大数过虑降级

#### Version 4.6.1
> 1.大数据过虑改为POST

#### Version 4.6.0
> 1.添加举报记录

#### Version 4.5.4
> 1.修复注册状态覆盖
> 2.修复ipad版本公告奔溃

#### Version 4.5.3
> 1.tw升级新版本，需要重新打包上线

#### Version 4.5.2
> 1.过虑过多回车和空格
> 2.修复子楼层设置踩时添加到热门

#### Version 4.5.1
> 1.添加修改赞踩数

#### Version 4.5.0
> 1.更新热评、举报理由

#### Version 4.4.2
> 1.添加用户信息降级

#### Version 4.4.1
> 1.修复更新subject状态时还没注册

#### Version 4.4.0
> 1.添加举报一二审

#### Version 4.3.7
> 1.更新go-business，整理接口命名规范

#### Version 4.3.6
> 1.添加internal内网流量接口

#### Version 4.3.5
> 1.更新http接口支持内外网流量

#### Version 4.3.3/4
> 1.本地TW docker.

#### Version 4.3.2
> 1.更新golang/go-common/go-business

#### Version 4.3.1
> 1.修复监控评论错误码

#### Version 4.3.0
> 1.过虑添加评论ID
> 2.添加up主mid更新

#### Version 4.2.6
> 1.添加日志封禁时间

#### Version 4.2.5
> 1.批量获取info

#### Version 4.2.4
> 1.删除先审后发
> 2.注册subject置顶state

#### Version 4.2.3
> 1.修复缓存miss判断处理bug

#### Version 4.2.2
> 1.合并二级评论请求为pipeline
> 2.修复楼层重复
> 3.修改先审后发为先发后审核

#### Version 4.2.1
> 1.合并sub请求

#### Version 4.2.0
> 1.优化reply获取

#### Version 4.1.0
> 1.增加realip
> 2.更新Version endor

#### Version 4.0.5
> 1.是否点赞使用HMGET获取

#### Version 4.0.4
> 1.修改点赞过期时间

#### Version 4.0.3
> 1.允许待审评论被删除

#### Version 4.0.2
> 1.修复json panic

#### Version 4.0.1
> 1.修改top置顶为异步加缓存

#### Version 4.0.0
> 1.添加评论踩
> 2.添加up置顶评论
> 3.添加up删除评论
> 4.添加评论attr属性字段，隔离属性和状态
> 5.添加subject主体attr字段，记录是否存在置顶等属性
> 6.更新基础库依赖

#### Version 3.3.0
> 1.接入新的过滤服务
> 2.更新go-common依赖

#### Version 3.2.1
> 1.更新Version endor依赖

#### Version 3.2.0
> 1.更新Version endor依赖

#### Version 3.1.1
> 1.add window0 plat

#### Version 3.1.0
> 1.subject 自动注册

#### Version 3.0.1
> 1.修复up主无法显示

#### Version 3.0.0
> 1.评论root是否点赞处理
> 2.置顶评论不允许删除
> 3.接入错误码管理后台
> 4.goVersion endor支持
> 5.修改go-business依赖

#### Version 2.4.1
> 1.评论门槛等级提高

#### Version 2.4.0
> 1.up主举报自身视频自身评论视为删除
> 2.支持移动端定位跳转

#### Version 2.3.1
> 1.修改json字段deleted
> 2.增加表情删除状态判断

#### Version 2.3.0
> 1.新增Version ip表情包接口
> 2.Version ip添加评论处理
> 3.评论门槛限制

#### Version 2.2.2
> 1.过滤异常utf8编码
> 2.添加subject状态查询接口
> 3.允许隐藏疑似垃圾评论
> 4.subject增加禁止评论状态判断

#### Version 2.2.1
> 1.up主显示隐藏评论

#### Version 2.2.0
> 1.批量举报和隐藏和显示
> 2.修复大数据过滤

#### Version 2.1.6
> 1.修复三个热门评论不显示子评论

#### Version 2.1.5
> 1.修复top评论子回复数

#### Version 2.1.4
> 1.action不存在时，增加无action标志缓存

#### Version 2.1.3
> 1.修复查询评论无内容bug

#### Version 2.1.2
> 1.修复无top评论稿件查询db

#### Version 2.1.1
> 1.修复二级评论查询db

#### Version 2.1.0
> 1.支持画站类型
> 2.置顶评论是否点赞

#### Version 2.0.0
> 1.评论重构
