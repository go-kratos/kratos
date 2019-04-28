# videoup-job

##### Version 1.31.1
>1.修改新人稿件过审的链接

##### Version 1.31.0
>1.联合投稿对接动态

##### Version 1.30.2
>1.升级账号API为grpc服务,简化配置

##### Version 1.30.1
>1.规范化waitGroup使用

##### Version 1.30.0
>1.投稿发动态支持投票业务

##### Version 1.29.2
>1.活动被管理员取消时，发送系统通知
>2.新用户首个稿件过审，发送系统通知

##### Version 1.29.1
>1.支持稿件发动态支持LBS

##### Version 1.29.0
>1.支持ugc付费流转

##### Version 1.28.2
>1.fixbug稿件封面存储支持第三方源转存到bfs/archive 替换https:

##### Version 1.28.1
>1.稿件封面存储支持第三方源转存到bfs/archive 替换https:

##### Version 1.28.0
>1.稿件封面存储支持第三方源转存到bfs/archive

##### Version 1.27.9
>1.家长模式支持一审视频聚合到稿件

##### Version 1.27.8
>1.AI补全封面时不再指定第一P

##### Version 1.27.7
>1.切片落库，没有耗时操作 走同步写库逻辑

##### Version 1.27.6
>1.评论开关迁移至videoup-report-job，改为稿件状态联动

##### Version 1.27.5
>1.稿件动态pub 新增show=2标识粉丝动态禁止属性

##### Version 1.27.4
>1.邮件迁移至videoup-report-job

##### Version 1.27.3
>1.archiveState橙色通过剔除活动稿件

##### Version 1.27.2
>1.archiveState 区分审核修改导致的稿件聚合，取消审核导致的修复待审case

##### Version 1.27.1
>1.新增bvc消费延迟短信告警
>2.去除bfs_videoshot双写

##### Version 1.27.0
>1.新增queueRedis 处理耗时任务 目前已经迁入切片逻辑

##### Version 1.26.13
>1.同时消费mail的新旧list, 约定时间点从旧list切到新list

##### Version 1.26.12
>1.split mail to new list

##### Version 1.26.11
>1.定时发布表增加软删除字段deleted_at

##### Version 1.26.10
>1.签约up主报备邮件跟其他报备邮件隔离

##### Version 1.26.9
>1.tag同步全部迁移到videoup-report-job

##### Version 1.26.8
>1.新增转码错误XcodeFailCodes
>2.fix videoshotAdd imgURL 验证

##### Version 1.26.7
>1.fix videoshotAdd binURL 验证

##### Version 1.26.6
>1.videoshotAdd binURL 验证

##### Version 1.26.5
>1.移除视频和稿件删除时的任务删除

##### Version 1.26.4
>1.移除创作姬推送消息

##### Version 1.26.3
>1.完全迁移一审任务

##### Version 1.26.2
>1.移除 manager.logger无用表

##### Version 1.26.1
> 1.迁移一审任务

##### Version 1.26.0
>1.彻底移除dede库

##### Version 1.25.28
> 1.xcode_hd_finish 支持视频meta width,height,rotate
> 2.profile error (rpc timeout) 时设置 profile=nil 防止默认值 0 bug

##### Version 1.25.27
> 1.移除第一次过审过渡代码

##### Version 1.25.26
> 1.修复redis hash未正常过期
> 2.规范timeformat，兼容(0001-01-01 00:00:00)的零时格式

##### Version 1.25.25
> 1.sql 修复类型不识别导致任务分发失败
> 2.定时删除过期task_dispatch_extend

##### Version 1.25.24
> 1.修复formattime时间解析错误

##### Version 1.25.23
> 1.databus 消息加recover防止无限失败
> 2.移除bvcAllSub

##### Version 1.25.22
> 1.移除shotjob 下沙废弃的消费者

##### Version 1.25.21
> 1.fixbug  修复二转失败时，消息推送panic

##### Version 1.25.20
> 1.fixbug  retry开评论剔除掉评论无法开启case
> 2.fixbug  优化syncBVC  retry逻辑
> 3.fixbug  优化retry 重试间隔为200毫秒

##### Version 1.25.19
> 1.拦截三查阈值为0 情况

##### Version 1.25.18
> 1.迁移bm

##### Version 1.25.17
> 1.邮件发送结果存行为日志

##### Version 1.25.16
> 1.增加视频没有找到的校验

##### Version 1.25.15
> 1.修正权重取值错误

##### Version 1.25.14
> 1.新增上海云立方VideoshotpvSub消费者

##### Version 1.25.13
> 1.添加最小权重值，防止任务一直降权
> 2.修复普通任务权重分值取错误

##### Version 1.25.12
> 1.修复特殊任务的权重计算错误

##### Version 1.25.11
> 1.定时任务和普通任务权重互斥
> 2.添加权重配置的有效时间
> 3.无效的任务进行删除

##### Version 1.25.10
> 1.指派任务均匀分配
> 2.任务权重新增按分区和投稿来源配置
> 3.权重分值可配置化

##### Version 1.25.9
> 1. videoshotDown error -404 不用重试直接失败

##### Version 1.25.8
> 1. 清理一波job 各种错误的error

##### Version 1.25.7
> 1. 稿件审核后台计数接入statView-T

##### Version 1.25.6
> 1. 每日10点删除一个月前的任务以及三个月前日志

##### Version 1.25.5
> 1. bfs_videoshotpv双写

##### Version 1.25.4
> 1. 取消指派任务的插队

##### Version 1.25.3
> 1. 修改权重分值

##### Version 1.25.2
> 1. 迁移目录到main

##### Version 1.25.1
> 1. up消息小助手明确是回查阶段 稿件开放状态下有修改才发通知

##### Version 1.25.0
> 1. 将老videoshot迁移进来  
> 2. 移除无用的dede同步  

##### Version 1.24.3
> 1. 兼容任务的指派插队逻辑

##### Version 1.24.2
> 1. 稿件修改消息小助手

##### Version 1.24.1
> 1. 修复rpc请求过于频繁

##### Version 1.24.0
> 1. 添加一审任务权重

##### Version 1.23.18
> 1. account v7

##### Version 1.23.17
> 1. add contributor reviewer + author

##### Version 1.23.16
> 1. archive_first_pass log添加aid

##### Version 1.23.15
> 1. fix  archive_forbid 聚合bug

##### Version 1.23.14
> 1. fix  redigo nil err
> 2. fix  视频删除后对xocde_sd_finish消息不聚和稿件状态

##### Version 1.23.13
> 1.第一次过审重构，查询双写

##### Version 1.23.12
> 1.邮件发送失败后，所有收件人单独发送

##### Version 1.23.11
> 1.新增syncBVC 详细日志 支持aid追溯
> 2.二审修复待审提交不调用syncBVC
> 3.ugc cid 转 pgc 时通知不可播

##### Version 1.23.10
> 1.去掉使用Covers的地方，全部使用AICovers

##### Version 1.23.9
> 1.粉丝动态推送支持传递dynamic

##### Version 1.23.8
> 1.fixbug 解决ArchiveNotify-T 时间兼容

##### Version 1.23.7
> 1.对接B博动态

##### Version 1.23.6
> 1.被删除分p不进待审和任务分发

##### Version 1.23.5
> 1.去掉xcode_sd_finish 里 补发first_round消息逻辑。理由是不清真，补丁功能

##### Version 1.23.4
> 1.删除分P 不再更新dede.dm_index dm_indexdata

##### Version 1.23.3
> 1.非特殊分区私单定时发布时，稿件进入四审

##### Version 1.23.2
> 1.增加bvc video日志

##### Version 1.23.1
> 1.增加bvc日志

##### Version 1.23.0
> 1.重构job对接新表

##### Version 1.22.3
> 1.增加一审开放、编辑稿件无改动切为开放的时候，同步res库
> 2.一审非开放时，不再关闭评论

##### Version 1.22.2
> 1.将稿件属性第13位hideclick改成limit_area
> 2.添加goconvey test

#### Version 1.22.1
> 1.私单四审 修改稿件聚合支持私单业务

#### Version 1.22.0
> 1.私单四审

#### Version 1.21.13
> 1.syncBVC 通知视频可播 判断条件换成 IsNormal == false

#### Version 1.21.12
> 1.video Resolutions HD  新增 112

#### Version 1.21.12
> 1.prom统计分发，二审计数

#### Version 1.21.11
> 1.自动获取稿件封面的源优先为ai，视频云次之

#### Version 1.21.10
> 1.up_from=6的合作嵌套稿件流程变更
> 2.pgc稿件不进ugc的回查流程

#### Version 1.21.9
> 1.视频云消息触发统一新表状态

#### Version 1.21.8
> 1.prom统计一转二转及一审计数

#### Version 1.21.7
> 1.修复panic

#### Version 1.21.6
> 1.prom统计一转二转耗时

#### Version 1.21.5
> 1.发送创作姬的推送消息中"点击查看"从URL改成文字显示

#### Version 1.21.4
> 1.发送创作姬的推送消息

#### Version 1.21.3
> 1.支持一审打回/锁定发送报备邮件,统一邮件发送逻辑

#### Version 1.21.2
> 1.支持私单报备邮件

#### Version 1.21.1
> 1.新版prom databus,retry prom 计数

#### Version 1.21.0
> 1.新增 databus,retry prom 计数

#### Version 1.20.4
> 1.upos实验室上传的稿件，消耗ugc_first_round消息

#### Version 1.20.3
> 1.将稿件attribute中的第9位作为is_pgc字段

#### Version 1.20.2
> 1.审核流程变更

#### Version 1.20.1
> 1.邮件报备新增 181,177 两个大区支持

#### Version 1.20.0
> 1.去除dede同步逻辑

#### Version 1.19.15
> 1.切换封面消息topic到Bvc2VuSub

#### Version 1.19.14
> 1.添加一审导致稿件禁止时审核消息(first_round_forbid)

#### Version 1.19.13
> 1.添加自动过审、定时发布的消息（auto_open、delay_open）

#### Version 1.19.12
> 1.修复 异步发送email

#### Version 1.19.11
> 1.异步发送email

#### Version 1.19.10
> 1.恢复passed 逻辑查 archive_oper表判断稿件是否过审

#### Version 1.19.9
> 1.新增水印下载失败导致转码失败的code

#### Version 1.19.8
> 1.调整 passed 逻辑改查 archive_oper 表

#### Version 1.19.7
> 1.archiveState fix forbid bug

#### Version 1.19.6
> 1.dedeSync 直接异步重试,去掉同步执行逻辑

#### Version 1.19.5
> 1.去除APPkey参数

#### Version 1.19.4
> 1.修复待审时点击量打点  
> 2.调整进入三查的计算规则  

#### Version 1.19.3
> 1.封面地址上传之后直接上传到redis,不需要直接插入到DB，后续 rename && drop table  
> 2.xcodeSDFinish渣清封面也重构到redis里面获取数据  
> 3.评论info日志内去除error文案  

#### Version 1.19.2
> 1.添加定时发布三审逻辑  

#### Version 1.19.1
> 1.video表更新使用hash64索引  

#### Version 1.19.0
> 1.archive_video拆表双写逻辑  

#### Version 1.18.0
> 1.task_oper_history迁移  
> 2.archive_edit_history/archive_video_edit_history数据清理  

#### Version 1.17.1
> 1.报警电话号配置  

#### Version 1.17.0
> 1.删除同步result  

#### Version 1.16.2
> 1.调整reply check status的错误级别   

#### Version 1.16.1
> 1.pgc的result同步  

#### Version 1.16.0
> 1.去除所有prepare  
> 2.接入prom  

#### Version 1.15.9
> 1.为去掉多余封面数据做准备，同步封面数据到Redis，缓存15天,为创作中心做准备  

#### Version 1.15.8
> 1.修复重发first_round消息xcode值错误  

#### Version 1.15.7
> 1.bvc通知修复多P、删除等情况  

#### Version 1.15.6
> 1.一二审消息bvc通知增加稿件状态判断  

#### Version 1.15.5
> 1.是否PGC的判断改为是否番剧  

#### Version 1.15.4
> 1.调整bvc通知按cid状态绝对通知  

#### Version 1.15.3
> 1.添加如果是付费稿件则不调用bvc逻辑  

#### Version 1.15.2
> 1.定时发布调用bvc  
> 2.secondRound补加是否pgc的判断  

#### Version 1.15.1
> 1.稿件状态变更和评论状态变更操作日志  
> 2.稿件评论开关无限重试  
> 3.删除延迟发布  
> 4.增加视频attribute的搜索禁止聚合到archive表  
> 5.去掉二审attribute位的聚合  

#### Version 1.15.0
> 1.videoup-job合入大仓库  

#### Version 1.14.0
> 1.result库更新时增加forward  

#### Version 1.13.0
> 1.私单round列表聚合  

#### Version 1.12.0
> 1.聚合attr信息时,同时判断archive的attr  

#### Version 1.11.1
> 1.fix sendMessage bug  
> 2.fix bvc api log type  

#### Version 1.11.0
> 1.一审通过无脑同步result  

#### Version 1.10.0
> 1.分发完成时，定时稿件也做同步逻辑  

#### Version 1.9.2
> 1.bvc接口三次重试失败后添加到redis队列无限重试  
> 2.某cid只要在archive_video表中存在至少一条state为0或10000并且xcode为6的数据，即认为此cid可播放  
> 3.分发完成，非PGC调bvc接口同步cid是否可公开播放情况  

#### Version 1.9.1
> 1.bugfix 分发完成，如果是pgc的投稿，一定同步result  
> 2.修改isPGC方法，根据aid通过addit表的upfrom字段判断  

#### Version 1.9.0
> 1.一审、二审的时候，查询cid的审核结果并调用bvc接口通知是否可以播放  
> 2.二审提交时，针对粉丝大于10万的UP和优质UP主，如果有数据修改，发邮件通知相应负责人  
> 3.分区点击量阈值和待审分区使用新表archive_config  
> 4.一级分区使用新表archive_type  
> 5.接入log agent  
> 6.go-common从v6.11.0更新到v6.17.1  
> 7.评论开关判断：如果是修复待审(-6)、用户提交(-30)、用户删除(-100),不调评论接口(保持原状态)  
> 8.二审去掉聚合forbidden attribute的逻辑  

#### Version 1.8.11
> 1.archive_result在update时,强制更新mtime  

#### Version 1.8.10
> 1.second无脑聚合attr  
> 2.活动结构json结构修改  

#### Version 1.8.9
> 1.调整archive_result库的插入顺序,先video后archive  

#### Version 1.8.8
> 1.databus的second_round消息增加reply、send_notify、mission_id字段，修改评论和发消息的对应逻辑  
> 2.修改二审调用reply接口、发消息、调活动平台接口的逻辑  
> 3.修改定时发布调活动平台接口的逻辑  
> 4.修复result表videos数量的bug  
> 5.增加monitor/ping   

#### Version 1.8.7
> 1.addArchive添加判断: 数据库查询archive是否为nil  
> 2.databus消息的route判断的前后添加log  

#### Version 1.8.6
> 1.任务被up删除，如果任务延迟也需要删除  

#### Version 1.8.5
> 1.dede_archives_history替换为新表archive_edit_history  
> 2.二审过审和延迟发布时，调活动平台接口添加稿件信息  
> 3.修改tag接口  

#### Version 1.8.4
> 1.定期迁移任务表数据  

#### Version 1.8.3
> 1.一审处理reply和tag  

#### Version 1.8.2
> 1.增加archive_result结果表同步  

#### Version 1.8.1
> 1.second_round增加无脑聚合视频时长的逻辑  

#### Version 1.8.0
> 1.task任务分发  

#### Version 1.7.3
> 1.增加用户删除稿件时,清理task表的逻辑

#### Version 1.7.2
> 1.增加用户删除视频时,清理task表的逻辑

#### Version 1.7.1
> 1.shot时判断v是否为nil

#### Version 1.7.0
> 1.增加v5.0稿件生态任务分派逻辑
> 2.增加新版配置中心  
> 3.稿件封面图走新消息直接落库  
> 4.多次收到sd fiinish时 增加first round重发逻辑  
> 5.增加插入videoshot的逻辑    
> 6.同步老表的dede_arctiny的逻辑

#### Version 1.6.9
> 1.修复删除分P时,视频总时长未更新的bug  

#### Version 1.6.8
> 1.修复一审后archive round不变更问题  

#### Version 1.6.7
> 1.在synccid之后删除48h判断的filename  

#### Version 1.6.6
> 1.删除track无用代码  

#### Version 1.6.5
> 1.聚合video的forbidattr到archive中  

#### Version 1.6.4
> 1.videoup databus 增加幂等  
> 2.round定时99改为1小时的维度   
> 3.封面截图3->9张   

#### Version 1.6.3
> 1.大于半小时的消息不消费  

#### Version 1.6.2
> 1.修复商单state&round驱动(定时发布导致开放时)  

#### Version 1.6.1
> 1.修复商单state&round驱动  

#### Version 1.6.0
> 1.打开评论增加重试三次的逻辑,3s一次试3次  
> 2.增加定时讲过期稿件的round置为end的逻辑  
> 3.增加round变更记录  
> 4.二审消息同步评论和tag  
> 5.分发中的稿件状态置为-1  

#### Version 1.5.6
> 1.修改内存缓存顺序  

#### Version 1.5.5
> 1.增加round变更日志  

#### Version 1.5.4
> 1.分发完成无脑同步  

#### Version 1.5.3
> 1.去除灰度用户  

#### Version 1.5.2
> 1.修改attr,access,state逻辑  

#### Version 1.5.1
> 1.增加回查进入触发回查中间31  

#### Version 1.5.0
> 1.二审拆分普通/分区二审、增加分区三审，回查拆分分区/社区  

#### Version 1.4.17
> 1.增加delete archive的无脑同步  

#### Version 1.4.16
> 1.同步老表去掉转码中和分发中  
> 2.评论更改内部url  

#### Version 1.4.15
> 1.增加tag调用日志  

#### Version 1.4.14
> 1.活动表换成addit  

#### Version 1.4.13
> 1.评论接口增加mid  

#### Version 1.4.12
> 1.嵌套支持qq/sohu/hunan  

#### Version 1.4.11
> 1.去掉一转无脑同步  

#### Version 1.4.10
> 1.再次修复分发完成前用户编辑稿件直接过审的bug  
> 2.round值修改   

#### Version 1.4.9
> 1.修复分发完成前用户编辑稿件直接过审的bug  

#### Version 1.4.8
> 1.一审后状态变更同步老表  

#### Version 1.4.7
> 1.定时发布再编辑,进待审列表  

#### Version 1.4.6
> 1.转载不做同步处理  

#### Version 1.4.5
> 1.修复转载来源同步到老desc字段  

#### Version 1.4.4
> 1.增加user_delete时同步老表的逻辑  

#### Version 1.4.3
> 1.同步到老库进行html转义  

#### Version 1.4.2
> 1.一审过审增加同步老表，为审核能预览稿件  
> 2.修改第一次不过，再次编辑不进待审自动过bug  

#### Version 1.4.1
> 1.分P创建已提交，稿件也是创建提交  

#### Version 1.4.0
> 1.增加addit source状态的同步
> 2.修复history表body信息同步错误
> 3.增加add archive消息

#### Version 1.3.3
> 1.修改回查阈值最大值加3000  

#### Version 1.3.2
> 1.删除archive_history表的插入

#### Version 1.3.1
> 1.增加外部源同步  

#### Version 1.3.0
> 1.增加定时发布逻辑
> 2.增加同步老表失败重试
> 3.增加recommend逻辑  
> 4.增加firstround时的aid判断  

#### Version 1.2.6
> 1.update go-common  
> 2.merge archive history  

#### Version 1.2.5
> 1.增加delay逻辑

#### Version 1.2.4
> 1.删除archive和video的trace  

#### Version 1.2.3
> 1.稿件attr的海外禁止定位第5bit  

#### Version 1.2.2
> 1.聚合稿件信息时，增加mission的判断  
> 2.增加filename的redis缓存,供videoup-web使用  
> 3.mission 同步错误  

#### Version 1.2.1
> 1.修复ptime时间  

#### Version 1.2.0
> 1.增加delay逻辑  

#### Version 1.1.9
> 1.只有1P的情况下，不同步desc字段  

#### Version 1.1.8
> 1.增加同步老表的flags属性  

##### Version 1.1.7
> 1.修复redis的key前缀  

##### Version 1.1.6
> 1.历史记录只有过审添加  

##### Version 1.1.5
> 1.archive_history date to timestamp  
> 2.dede_archives_history date to timestamp

##### Version 1.1.4

> 1.修复tag同步字段  
> 2.增加archive_history记录  

##### Version 1.1.3

> 1.完善日志格式  

##### Version 1.1.2

> 1.增加消费监控  

##### Version 1.1.1

> 1.修复待审状态也同步到老库数据  
> 2.修复redis比consumer提前关闭导致的redigo: get on closed pool  

##### Version 1.1.0

> 1.增加稿件及视频状态变更记录  

##### Version 1.0.2

> 1.PGC稿件不发UGC_X  

##### Version 1.0.1

> 1.修复不同步作者昵称到老表问题  

##### Version 1.0.0

> 1.初始化重构投稿封面相关job
