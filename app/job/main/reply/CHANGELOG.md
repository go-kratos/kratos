#### reply-job

#### Version 6.0.1
> 1.修复err

#### Version 6.0.0
> 1.折叠评论，置顶及删除逻辑

#### Version 5.0.16
> 1.account grpc

#### Version 5.0.15
> 1. remove user action table

#### Version 5.0.14
> 1.CD数减半

#### Version 5.0.13
> 1.fix nil sub bug

#### Version 5.0.11
> 1.增加置顶和取消置顶的消息

#### Version 5.0.10
> 1.国际版消息推送

#### Version 5.0.9
>1.bbq

#### Version 5.0.8
> 1.ats 不要每次都从message里拿

#### Version 5.0.7
> 1.点赞限制放宽，10->15

#### Version 5.0.5
> 1.评论cd调整，删除cd code如果人证成功

#### Version 5.0.2
>1.修复mtime

#### Version 5.0.1
>1.reply dialog, 修复sql

#### Version 5.0.0
>1.reply dialog, 5.33

#### Version 4.20.8
>1.fix reply report update

#### Version 4.20.6
>like大于3

#### Version 4.20.3
>删除非必要日志

#### Version 4.20.2
>log.Warn

#### Version 4.20.1
>epid grpc

#### Version 4.19.11
> 1.增加feature flag

#### Version 4.19.9
> 修复回复通知到ats 中

#### Version 4.19.8
> 1.配置文件
> 2.xint move to model

#### Version 4.19.7
> 1.notify新增nativeJump

#### Version 4.19.6
> 1.修复job的golint问题

#### Version 4.19.4
> 1.点赞MaxLike去重

#### Version 4.19.3
> 1.去除点赞双写数据库

#### Version 4.19.2
> 1.去掉子评论老缓存index

#### Version 4.19.0
> 1.双写子评论缓存index

#### Version 4.18.8
> 1.修复通知评论内容长度

#### Version 4.18.7
> 1.回复通知标题换为根评论内容

#### Version 4.18.6
> 1.下线老Stat-T

#### Version 4.18.5
> 1.音频评论消息通知接入

#### Version 4.18.4
> 1.rebuild master

#### Version 4.18.3
> 1.删除非块级配置

#### Version 4.18.2
> 1.fix旧点赞冲突，无法更新

#### Version 4.18.1
> 1.fix点赞新平台

#### Version 4.18.0
> 1.点赞新平台

#### Version 4.17.13
> 1.修复bug

#### Version 4.17.12
> 1.新增business支持

#### Version 4.17.11
> 1.修复添加评论时根评论已删除

#### Version 4.17.10
> 1.迁移到bm

#### Version 4.17.9
> 1.迁移到databus

#### Version 4.17.8
> 1.修复可能出现的deadlock

#### Version 4.17.7
> 1.增加subject mcount

#### Version 4.17.6
> 1.调整顺序，解决死锁问题

#### Version 4.17.5
> 1.回源分段

#### Version 4.17.4
> 1.删除子评论oid index

#### Version 4.17.3
> 1.新子评论index缓存key

#### Version 4.17.2
> 1.修复通知回复mc

#### Version 4.17.1
> 1.评论计数通知

#### Version 4.17.0
> 1.评论分段回源

#### Version 4.16.7
> 1.添加通知native跳转链接

#### Version 4.16.6
> 1.话题正则修改，去掉emoji和link

#### Version 4.16.5
> 1.去掉稿件、动态、小视频、相簿评论数通知
> 2.切换账号v7 rpc

#### Version 4.16.4
> 1.添加稿件、动态、小视频、相簿评论数通知

#### Version 4.16.3
> 1.修复错误封禁通知链接

#### Version 4.16.2
> 1.修复reply_admin_log mid插入问题

#### Version 4.16.1
> 1.异步发送通知

#### Version 4.16.0
> 1.增加了对评论话题的支持

#### Version 4.15.8
> 1.在热门评论索引处添加根评论判断

#### Version 4.15.7
> 1.修改DynamicUrl

#### Version 4.15.6
> 1.修改封禁文案

#### Version 4.15.5
> 1.更换小黑屋通知标题接口

#### Version 4.15.4
> 1.优化置顶评论

##### Version 4.15.3
> 1.动态消息推送

##### Version 4.15.2
> 1.修改评论At的正则表达式

##### Version 4.15.1
> 1.添加StatReply-T计数databus

##### Version 4.15.0
> 1.添加一二审状态处理

##### Version 4.14.5
> 1.修复热评删除后计算时放到缓存中

##### Version 4.14.4
> 1.添加置顶缓存过期时间配置

##### Version 4.14.3
> 1.top评论回源时sql语句增加limit限制

##### Version 4.14.2
> 1.修复置顶时楼层index被删除评论
> 2.优化单次回源时，加载好index缓存
> 3.置顶评论预加载
> 4.修复可能存在两个置顶的问题

##### Version 4.14.1
> 1.添加举报审核databus事件

##### Version 4.14.0
> 1.添加评论记录更新

##### Version 4.13.0
> 1.数据库回源优化

##### Version 4.13.0
> 1.业务计数配置化

##### Version 4.12.3
> 1.修改小黑屋通知链接

##### Version 4.12.2
> 1.切换archive3 rpc接口

##### Version 4.12.1
> 1.通知规避部分违规违禁内容

##### Version 4.12.0
> 1.添加评论计数通知

##### Version 4.11.0
> 1.databus实时数据流

##### Version 4.10.0
> 1.迁移大仓库，更新memcache

##### Version 4.9.2
> 1.修复点赞数量为0通知

##### Version 4.9.1
> 1.修复点赞慢查询

##### Version 4.9.0
> 1.添加先审后发

##### Version 4.8.5
> 1.修复通知标题为空

##### Version 4.8.4
> 1.修复消息推送http头

##### Version 4.8.3
> 1.修改搜索通知URL

##### Version 4.8.2
> 1.修复mobile通知跳转

##### Version 4.8.1
> 1.修复稿件redirect通知跳转

##### Version 4.8.0
> 1.添加评论通知

##### Version 4.7.2
> 1.调整协管删除评论日志格式

##### Version 4.7.1
> 1.添加协管员删除日志详情(50个字符长度)

##### Version 4.7.0
> 1.新增协管员删除评论逻辑

##### Version 4.6.5
> 1.修复syslog奔溃

##### Version 4.6.4
> 1.修复syslog panic
> 2.修复admin log用户删除operation

##### Version 4.6.3
> 1.更新标准库
> 2.增加稿件类型"article"的消息通知

##### Version 4.6.2
> 1.优化搜索批量更新

##### Version 4.6.1
> 1.被拉黑名单屏蔽通知

##### Version 4.6.0
> 1.新搜索评论列表通知
> 2.添加活动稿件通知链接
> 3.文章计数发送databus

##### Version 4.5.6
> 1.去除老搜索通知

##### Version 4.5.5
> 1.修改热评Version 2.0逻辑

##### Version 4.5.4
> 1.添加ping资源

##### Version 4.5.3
> 1.修复审核移除评论默认状态值

##### Version 4.5.2
> 1.修复运营后台评论状态同步搜索索引

##### Version 4.5.1
> 1.更新健康检查端口

##### Version 4.5.0
> 1.更新lib
> 2.新配置中心

##### Version 4.4.1
> 1.修复举报通知

##### Version 4.4.0
> 1.添加举报用户反馈

##### Version 4.3.3
> 1.修复举报的热门评论

##### Version 4.3.2
> 1.修复点踩使用缓存不计数

##### Version 4.3.1
> 1.修复评论更新时间
> 2.修复监控评论没有推送通知

##### Version 4.3.0
> 1.热评算法调整

##### Version 4.2.2
> 1.举报去除通知索引更新

##### Version 4.2.1
> 1.修复评论删除状态

##### Version 4.2.0
> 1.添加举报一二审

##### Version 4.1.1
> 1.增加封禁时间

##### Version 4.1.0
> 1.小黑屋评论跳转

##### Version 4.0.11
> 1.修改先审后发为先发后审
> 2.番剧评论跳转
> 3.举报不加节操


##### Version 4.0.10
> 1.限制口节操发送消息长度

##### Version 4.0.9
> 1.删除先审后发根评论缓存

#####	Version 4.0.8
> 1.删除番剧评论跳转

##### Version 4.0.7
> 1.删除小视频消息通知

##### Version 4.0.6
> 1.top主动添加置顶

##### v4.0.5
> 1.修改点赞过期时间

##### Version 4.0.4
> 1.踩不发通知。
> 2.从缓存获取reply

##### Version 4.0.3
> 1.修复content内容为空

##### Version 4.0.2
> 1.删除sql attr位操作

##### Version 4.0.1
> 1.修改top置顶为异步

##### Version 4.0.0
> 1.添加评论踩
> 2.添加up置顶评论
> 3.添加up删除评论
> 4.添加评论attr属性字段，隔离属性和状态
> 5.添加subject主体attr字段，记录是否存在置顶等属性
> 6.更新基础库依赖

##### Version 3.3.0
> 1.更新go-common依赖

##### Version 3.2.2
> 1.修复二级评论计数

##### Version 3.2.1
> 1.评论计数双写到databus

##### Version 3.2.0
> 1.更新活动搜索索引
> 2.更新楼层计数必须依赖于db

##### Version 3.1.0
> 1.增加电影区评论跳转
> 2.更新go-business至1.6.0
> 3.更新go-common至4.1.3

##### v3.0.1
> 1.更新基础库依赖包

##### v3.0.0
> 1.番剧评论跳转锚点处理
> 2.增加活动单页评论类型
> 3.修复置顶评论cache count不更新
> 4.govendor支持

##### v2.3.0
> 1.评论点赞通知

##### v2.2.4
> 1.允许隐藏filter状态

##### v2.2.3
> 1.疑似垃圾评论隐藏

##### v2.2.2
> 1.修复疑似垃圾评论按点赞数不显示

##### v2.2.1
>  1.修复置顶平路member为空

##### v2.2.0
> 1.修复评论删除及恢复
> 2.修复疑是垃圾评论不显示

##### v2.1.3
> 1.kafka消息同步推送失败重试

##### v2.1.2
> 1.修复无top评论稿件查询db

##### v2.1.0
> 1.添加jump跳转
> 2.支持画站类型
> 3.修复top点赞无效

##### v2.0.0
> reply-job初始化
