#### 稿件审核后台接口

##### Version 1.29.4
> 1.修改监控数据的接口URL
> 2.修复稿件监控列表数据隐藏PGC的问题

##### Version 1.29.3
> 1.移除app配置

##### Version 1.29.2
> 1.使用up-service grpc

##### Version 1.29.1
> 1.移除稿件无用报表

##### Version 1.29.0
> 1.通用批量tag接口 添加tag或者删除tag接口  x/admin/videoup/archive/batch/tag
> 1.1支持频道回查 								form_list = channel_review
> 1.2支持adminBind/upBind(默认走adminBind)    is_up_bind=true
> 1.3支持同步隐藏tag 							sync_hidden_tag=true
> 1.4因为稿件服务不cache tags 也不需要发force_sync.未来计划砍掉审核库archive.tag

##### Version 1.28.1
> 1.修复稿件、视频搜索的account grpc 报错的问题

##### Version 1.28.0
> 1.联合投稿 admin

##### Version 1.27.8
> 1.升级账号API为grpc服务,简化配置

##### Version 1.27.7
> 1.批量修改分区支持刷数据模式  仅通知稿件服务  force_sync

##### Version 1.27.6
> 1.1000条登入日志不够用，扩充为10000条

##### Version 1.27.5
> 1.去掉稿件列表封面的默认值

##### Version 1.27.4
> 1.修复搜索列表封面URL路径错误的问题

##### Version 1.27.3
> 1.兼容videoup-task-admin 用于判断复审的参数. 待前端完全迁移到v2,删除本项目中任务代码

##### Version 1.27.2
> 1.修复稿件稿件回查列表数据移除导致重复数据的bug

##### Version 1.27.1
> 1.修复稿件搜索接口多个状态传字符串的bug
> 2.修复稿件搜索接口回查列表字段传递不正确的问题
> 3.稿件搜索接口添加关键字匹配度优先的功能

##### Version 1.27.0
> 1.稿件ugc付费流程

##### Version 1.26.3
> 1.修复全部稿件列表批量搜索aid时，不能翻页的问题
> 2.修复监控列表

##### Version 1.26.2
> 1.second_round 消息支持开关邮件发送 默认发邮件

##### Version 1.26.1
> 1.将稿件搜索从manager-v4迁移过来

##### Version 1.26.0
> 1.新增属性位 21位 家长模式

##### Version 1.25.8
> 1.first_round消息附带UP主粉丝数 fans 方便视频云做转码优先级

##### Version 1.25.7
> 1.恢复pgc的第18位属性设置AttrBitBadgepay

##### Version 1.25.6
> 1.删除pgc的第18位属性设置AttrBitBadgepay

##### Version 1.25.5
> 1.修复版权搜索返回结构不正确导致空指针的问题

##### Version 1.25.4
> 1.修复非视频列表也进行鉴权的bug

##### Version 1.25.3
> 1.去除视频监控统计中稿件状态为-100的数据
> 2.解决视频监控结果列表翻页的bug

##### Version 1.25.2
> 1.去除redis keys命令的使用

##### Version 1.25.1
> 1.将权重的redis独立，不影响稿件主流程

##### Version 1.25.0
> 1.添加稿件停留数量统计接口

##### Version 1.24.7
> 1.人工操作评论开关移除，全部改为videoup-report-job的状态联动

##### Version 1.24.6
> 1.修改版权搜索报错信息

##### Version 1.24.5
> 1.二审到三审新增商单通道，(支持商单三审及特色分区三审)

##### Version 1.24.4
> 1.hbase v2

##### Version 1.24.3
> 1.fixbug 修复二审round 10 state-6 非特殊分区 提交审核 变成三审bug

##### Version 1.24.2
> 1.频道rpc返回结构修改

##### Version 1.24.1
> 1.一审视频任务质检添加不支持异步，避免与task_utime上报的竞争

##### Version 1.24.0
> 1.一审视频任务审核支持任务质检添加

##### Version 1.23.22
> 1.隔离预发布redis 队列

##### Version 1.23.21
> 1.定时发布表增加软删除字段deleted_at

##### Version 1.23.20
> 1.增加版权接口配置，解决copyright接口经常504的问题

##### Version 1.23.19
> 1.新增aitrack接口:获取相似稿件aid

##### Version 1.23.18
> 1.修复写视频操作记录，cid不存在时的bug
> 2.修复查询视频时，id参数中有空白字符ParseInt报错的问题

##### Version 1.23.17
> 1.搜索关键字等级统一改成low
> 2.版权接口返回前30条数据

##### Version 1.23.16
> 1.配合谷安将视频、版权搜索接口从v2迁移至v3。提供给前端的接口也从manager-v4迁移至Videoup-admin。

##### Version 1.23.15
> 1.新增频道信息查询接口

##### Version 1.23.14
> 1.将地区策略组加限制的policy_id>1 同步到稿件attr 13bit(取消限制同理)

##### Version 1.23.13
> 1.支持签约up主的报备邮件

##### Version 1.23.12
> 1.新增批量频道审核接口:新增或删除tag，可能触发频道回查

##### Version 1.23.11
> 1.添加up主信用分稿件端数据上报 upcredit pub
> 2.移除音频库业务

##### Version 1.23.10
> 1.tag同步绑定一级、二级分区
> 2.稿件审核单个提交，移除tag落库和tag变更日志


##### Version 1.23.9
> 1.添加热门回查功能

##### Version 1.23.8
> 1.去掉redis的大key: task_weight

##### Version 1.23.7
> 1.移除net/http/parse

##### Version 1.23.6
> 1.频道回查、稿件审核单个提交、稿件批量修改属性支持频道禁止
> 2.tag接口：无脑回查从属性位，改为archive_recheck表

##### Version 1.23.5
> 1.优化tag接口：从频道回查列表进入且提交的稿件会重置属性位，不管实时查询是否为频道回查

##### Version 1.23.4
> 1.移除第一次过审过渡代码

##### Version 1.23.3
> 1.changelog版本号错误

##### Version 1.23.2
> 1.修复taskweight的redis过期不生效
> 2.同步videoup-job内TimeFormat修改

##### Version 1.23.1
> 1.新增单独保存稿件tag的接口，指定条件下触发频道回查

##### Version 1.23.0
> 1.添加策略组相关接口

##### Version 1.22.11
> 1.修复权重变更日志发布时间错误

##### Version 1.22.10
> 1.去掉任务释放时的大事务
> 2.时间类型使用timeformat简化
> 3.从manager-admin读取uid和uname，不直接读数据库

##### Version 1.22.9
> 1.稿件私单日志增加私单old和new数据

##### Version 1.22.8
> 1.修复日志缺少用户名,cookie里面不存uid和uname，只存sessionid

##### Version 1.22.7
> 1.修复task_dispatch死锁

##### Version 1.22.6
> 1.审核登入日志接搜索平台

##### Version 1.22.5
> 1.稿件信息追踪的用户编辑区，只显示变更的标题、封面、简介、分P

##### Version 1.22.4
> 1.登出不释放第一条任务，该任务延迟5分钟释放
> 2.去掉for update语句使用

##### Version 1.22.3
> 1.取消私单编辑日志记录 (私单计数方案换成up-service)

##### Version 1.22.2
> 1.权重日志添加定时发布时间

##### Version 1.22.1
> 1.在线用户列表不查询上次登出时间
> 2.权重配置添加生效时间范围

##### Version 1.22.0
> 1.稿件bgm管理

##### Version 1.21.5
> 1.修改稿件私单日志格式

##### Version 1.21.4
> 1.登出时不校验已登陆

##### Version 1.21.3
> 1.忽略多次登入多次登出的错误

##### Version 1.21.2
> 1.使用auth.permit设置uid,不直接读取cookie里面的uid

##### Version 1.21.1
> 1.稿件商单接入日志平台

##### Version 1.21.0
>1.权重配置新增按照分区和投稿来源
>2.移植任务释放，任务延迟，以及用户登入登出接口
>3.新增权重分值查看和配置接口

##### Version 1.20.5
> 1.添加PolicyID逻辑
> 2.添加ApplyID逻辑

##### Version 1.20.4
> 1.日志上报优化

##### Version 1.20.2
> 1.稿件、视频接入日志平台

##### Version 1.20.1
> 1.稿件修改封面逻辑兼容 相对路径截取

##### Version 1.20.0
> 1.添加任务权重

##### Version 1.19.5
> 1.修改path

##### Version 1.19.4
> 1.稿件修改消息助手 modify flag

##### Version 1.19.3
> 1.archive/batch 支持修改copyright 开关 flag_copyright

##### Version 1.19.2
> 1.使用account-service v7

##### Version 1.19.1
> 1.add contributor

##### Version 1.19.0
> 1.私单全量写操作日志

##### Version 1.18.6
> 1.去除statsd

##### Version 1.18.5
> 1.重构第一次过审，查询双写

##### Version 1.18.4
> 1.私单+活动稿件进入私单四审

##### Version 1.18.3
> 1.新版稿件、视频信息追踪

##### Version 1.18.2
> 1.对接B博动态

##### Version 1.18.1
> 1.非特殊分区私单定时发布时，稿件进入四审

##### Version 1.18.0
> 1.二审系统通知写操作记录
> 2.first_round 发送消息带稿件分区ID
> 3.私单流量TAG流量写操作记录

##### Version 1.17.2
> 1.error返回补漏

##### Version 1.17.1
> 1.引入blademaster

##### Version 1.17.0
> 1.pgc接口限流

##### Version 1.16.4
> 1.一审提交 将archive_video_relation state设置为0  允许修复已经被删除的视频 （与archive_video行为一致）

##### Version 1.16.3
> 1.将attribute第13位改成"是否限制地区"
> 2.添加了goconvey test

##### Version 1.16.2
> 1.视频审核若变更属性，触发邮件报备

##### Version 1.16.1
> 1.fixbug 私单二期 flow_design 私单流量 支持稿件详情页修改

##### Version 1.16.0
> 1.新增私单二期业务

##### Version 1.15.2
> 1.添加了打点表的数据接口

##### Version 1.15.1
> 1.切换到video新表

##### Version 1.15.0
> 1.增加video新表相关操作

##### Version 1.14.10
> 1.合作方嵌套（up_from=6）稿件流程变更

##### Version 1.14.9
> 1.稿件tag去重防止触发邮件报备。

##### Version 1.14.8
> 1.一审打回/锁定时，adminChange=true，即可以发送报备邮件。

##### Version 1.14.7
> 1.解决视频自动锁定的一些已知bug。查找399指派任务时去除state=2的任务。

##### Version 1.14.6
> 1.添加视频自动锁定【欧美电影】，【日本电影】，【其他国家】，【港台剧】，【海外剧】

##### Version 1.14.5
> 1.upos实验室上传的稿件，发送ugc_first_round消息

##### Version 1.14.4
> 1.审核流程变更

##### Version 1.14.3
> 1.将稿件attribute中的第9位作为is_pgc字段

##### Version 1.14.2
> 1.恢复 http track video 接口  /va/track/video

##### Version 1.14.1
> 1.恢复passed 逻辑查 archive_track 表

##### Version 1.14.0
> 1.archive_track 迁移至 hbase  
> 2.调整 passed 逻辑改查 archive_oper 表  

##### Version 1.13.3
> 1.稿件支持 dynamic 特性

##### Version 1.13.2
> 1.移除APPkey参数

##### Version 1.13.1
> 1.移除 archive_video_track 相关逻辑  

##### Version 1.13.0
> 1.天马流量融合私单  
> 2.开放预览和橙色通过才可进私单回查  

##### Version 1.12.5
> 1. readme去掉空格  

##### Version 1.12.4
> 1. 同步new_video的时候，增加ctime字段  

##### Version 1.12.3
> 1.修复二审提交批量attr异常改的allowtag属性  
> 2.修复redirectURL不能修改的bug  

##### Version 1.12.2
> 1.修复二审提交批量attr异常改的allowtag属性  
> 2.修复redirectURL不能修改的bug  

##### Version 1.12.1
> 1.二审增加默认更新mtime  

##### Version 1.12.0
> 1.archive_video拆表双写逻辑  
> 2.删除老的批量接口逻辑  

##### Version 1.11.0
> 1.批量提交一审／二审／修改稿件attr／移动稿件分区，及一、二审的稿件和视频的审核日志  

##### Version 1.10.2
> 1.新增商业平台更新定时发布时间  

##### Version 1.10.1
> 1.fix 批量审核，延迟发布被删除的bug  

##### Version 1.10.0
> 1.简介描述相关功能  

##### Version 1.9.5
> 1.删除延迟发布的逻辑  
> 2.fix私单打回流程bug  
> 3.去掉一审的二禁、三禁用户attr修改  

##### Version 1.9.4
> 1.活动稿件私单流程问题  

##### Version 1.9.3
> 1.兼容活动私单稿件bug  

##### Version 1.9.2
> 1.兼容活动databus的bug  

##### Version 1.9.1
> 1.私单活动稿件case fix  

##### Version 1.9.0
> 1.私单相关功能  
> 2.一审二禁用户增加动态／推荐限制  
> 3.修复redirectURL取消不了的bug  

##### Version 1.8.4
> 1.二审视频操作(增删改)  

##### Version 1.8.3
> 1.二审 修改判断job是否发邮件的逻辑  

##### Version 1.8.2
> 1.二审 判断是否需要发邮件并databus通知job  

##### Version 1.8.1
> 1.增加一审 task 等待耗时打点数据拉取接口

##### Version 1.8.0
> 1.修改attr的跳转链接时日志描述  

##### Version 1.7.9
> 1.fix log—agent的bug  

##### Version 1.7.8
> 1.水印迁移（新老兼容)  

##### Version 1.7.7
> 1.修复发布时间修改  

##### Version 1.7.6
> 1.二审遗漏note字段  

##### Version 1.7.5
> 1.入参统一使用ap  

##### Version 1.7.4
> 1.fix三审三查bug  

##### Version 1.7.3
> 1.稿件二审代码优化  

##### Version 1.7.2
> 1.视频会员可见驱动稿件会员可见  

##### Version 1.7.1
> 1.一审消息添加分区ID  

##### Version 1.7.0
> 1.添加二审接口  

##### Version 1.6.2
> 1.重发消息队列key修改  

##### Version 1.6.1
> 1.new_identify  

##### Version 1.6.0
> 1.新增一审提交修改三位用户attr的值  

##### Version 1.5.0
> 1.删除dede库的老逻辑  

##### Version 1.4.0
> 1.新增PGC属性修改接口  
> 2.新增PGC封禁接口  

##### Version 1.3.1
> 1.fix bug 增加重试协程  

##### Version 1.3.0
> 1.接入新的配置中心  

##### Version 1.2.0
> 1.一审视频  
> 2.追踪信息  

##### Version 1.1.0
> 1.增加修改attr逻辑  
> 2.增加商业产品修改attr接口  

##### Version 1.0.1
> 1.PGC自动日志表修改  

##### Version 1.0.0
> 1.PGC自动过审