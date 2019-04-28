### dm2的Gateway服务。


#### Version V3.3.17
> 1. 点赞接入grpc

#### Version V3.3.16
> 1. 字幕和蒙版，迁移出view调用
> 2. 创作中心最近1000条弹幕，30天强限制

#### Version V3.3.15
> 1. 创作中心最新1000条弹幕 limit 30d

#### Version V3.3.14
> 1. 创作中心最新1000条弹幕 接入搜索

#### Version V3.3.13
> 1. fix弹幕发送限速

#### Version V3.3.12
> 1. 弹幕发送黑名单

#### Version V3.3.11
> 1. use fanout

#### Version V3.3.10
> 1. 保护弹幕之前检查弹幕状态

#### Version V3.3.9
> 1. view aid allow empty

#### Version V3.3.8
> 1. account grpc

#### Version V3.3.7
> 1. 弹幕广播使用broadcast

#### Version V3.3.6
> 1. fix mysql count<0

#### Version V3.3.5
> 1. 去掉sync.pool

#### Version V3.3.4
> 1. 弹幕发送接入ai反垃圾限制

#### Version V3.3.3
> 1. localcache add subject
> 2. 分段弹幕大数据降级

#### Version V3.3.2
> 1. view 接口 localcache

#### Version V3.3.1
> 1. fix  ipv4 collect

#### Version V3.3.0
> 1. 弹幕mobile view
> 2. 抽出appview的逻辑
> 3. bfs上传使用sdk

#### Version V3.2.3
> 1. 字幕提交 优化提示
> 2. 付费稿件，up主可以发送弹幕
> 3. 弹幕协管操作次数限制，返回文案

#### Version V3.2.2
> 1. rebuild

#### Version V3.2.1
> 1. 字幕sql bug

#### Version V3.2.0
> 1. 弹幕广播添加限流

#### Version V3.1.25
> 1. 字幕添加开关

#### Version V3.1.24
> 1. remove addSlashes

#### Version V3.1.23
> 1. 添加开关配置

#### Version V3.1.22
> 1. 垃圾弹幕收集开关

#### Version V3.1.21
> 1. 弹幕发送接入行为日志

#### Version V3.1.20
> 1. spy service 降级

#### Version V3.1.19
> 1. figure ignore notfound error

#### Version V3.1.18
> 1.  rpc error catch

#### Version V3.1.17
> 1.  season id 换成grpc

#### Version V3.1.16
> 1.  filter service 换成grpc
> 2.  弹幕发送添加垃圾弹幕屏蔽

#### Version V3.1.15
> 1.  移除web蒙版旧接口

#### Version V3.1.14
> 1.  web蒙版扩展

#### Version V3.1.13
> 1.  弹幕点赞提示文案修改

#### Version V3.1.12
> 1.  增加弹幕发送时间限制

#### Version V3.1.11
> 1.  location的Zone方法改为Info方法

#### Version V3.1.10
> 1.  字幕web端添加原文语言返回

#### Version V3.1.9
> 1.  去掉字幕灰度的接口

#### Version V3.1.8
> 1.  修改字幕举报description

#### Version V3.1.7
> 1.  字幕举报兼容没有信用分

#### Version V3.1.6
> 1.  字幕举报添加metadata

#### Version V3.1.5
> 1.  添加付费视频弹幕发送认证

#### Version V3.1.4
> 1. 新增弹幕时记录用户ip、port

#### Version V3.1.3
> 1.  添加字幕举报

#### Version V3.1.2
> 1.  分段弹幕不刷新redis

#### Version V3.1.1
> 1.  添加字幕字数限制

#### Version V3.1.0
> 1.  字幕语言新增删除状态

#### Version V3.0.8
> 1.  5.30下掉弹幕蒙版

#### Version V3.0.7
> 1.  ajax接口xss

#### Version V3.0.6
> 1.  字幕验证时间戳+1s,兼容视频毫秒时长

#### Version V3.0.5
> 1.  字幕作者添加videoname返回

#### Version V3.0.4
> 1.  发弹幕去掉aid，cid一致检测

#### Version V3.0.3
> 1.  添加敏感词检测

#### Version V3.0.2
> 1.  下掉移动端飘窗弹幕

#### Version V3.0.1
> 1.  弹幕发送校验aid，cid关联

#### Version V3.0.0
> 1. 新增分段弹幕json接口
> 2. 新增历史弹幕json接口
> 3. 新增全段弹幕json接口

#### Version V2.10.16
> 1.  http添加trace

#### Version V2.10.15
> 1.  聚合字幕mc

#### Version V2.10.14
> 1.  弹幕历史更换business

#### Version V2.10.13
> 1. 更新弹幕广告接口URI

#### Version V2.10.12
> 1.  字幕web和mobile分开

#### Version V2.10.11
> 1.  发送弹幕 ecode返回值校正

#### Version V2.10.10
> 1.  rpc去掉语言验证

#### Version V2.10.9
> 1.  字幕缓存穿透bug

#### Version V2.10.8
> 1.  搜索空数据兼容

#### Version V2.10.7
> 1. view 接口添加singlegroup

#### Version V2.10.6
> 1. 增加字幕接口

#### Version V2.10.5
> 1. broadcast转义符导致错误 fix

#### Version V2.10.4
> 1. 弹幕点赞增加up主mid

#### Version V2.10.3
> 1. 蒙版缓存使用mc
> 2. use rpc.ServerConfig

#### Version V2.10.2
> 1. rebuild master

#### Version V2.10.1
> 1. 搜索索引更新接口对接sdk v3

#### Version V2.10.0
> 1. 新增弹幕广告
> 2. 分段弹幕缓存出错不再回源db

#### Version V2.9.13
> 1. 弹幕蒙版三期

#### Version V2.9.12
> 1. 优化历史弹幕日期索引查询

#### Version V2.9.11
> 1. remoteIP 更改

#### Version V2.9.10
> 1. 移除创作中心无用接口

#### Version V2.9.9
> 1. 内部接口新增recent和search

#### Version V2.9.8
> 1. identify迁移verify和auth

#### Version V2.9.7
> 1. 创作中心弹幕增加aid

#### Version V2.9.6
> 1. 修复高级弹幕购买状态bug

#### Version V2.9.5
> 1. 修复高级弹幕购买状态bug

#### Version V2.9.4
> 1. 修复创作中心弹幕搜索权限bug

#### Version V2.9.3
> 1.  nil dm content

#### Version V2.9.2
> 1. 移动端弹幕蒙版

#### Version V2.9.1
> 1. fix verify

#### Version V2.9.0
> 1. 移除视频维度的屏蔽词
> 2. 增加up主针对高级弹幕申请的配置
> 3. 增加up主搜索弹幕列表使用es sdk v3

#### Version V2.8.10
> 1. 提高卡片弹幕灰度比例

#### Version V2.8.9
> 1. fix recent dm panic

#### Version V2.8.8
> 1. fix rows.Err() in rows.Scan

#### Version V2.8.7
> 1. 卡片弹幕实验策略变更

#### Version V2.8.6
> 1. 卡片弹幕改为从header里获取Buvid

#### Version V2.8.5
> 1. 移动端卡片弹幕配合天马分组实验

#### Version V2.8.4
> 1. 重新构建master

#### Version V2.8.3
> 1. 历史弹幕搜索使用v3接口

#### Version V2.8.2
> 1. 更新稿件分批获取方法

#### Version V2.8.1
> 1. update discovery appid

#### Version V2.8.0
> 1. 使用discovery rpc client

#### Version V2.7.11
> 1. 移动端ajax弹幕直接返回空

#### Version V2.7.10
> 1. archive  duration缓存穿透

#### Version V2.7.9
> 1. redis fix hmget

#### Version V2.7.8
> 1. 缓存区分字幕弹幕和普通弹幕

#### Version V2.7.7
> 1. count的变更改为操作新库

#### Version V2.7.6
> 1. childpool的变更改为操作新库

#### Version V2.7.5
> 1. 重新构建master

#### Version V2.7.4
> 1. 移除ajax弹幕中的Mode7弹幕
> 2. 优化普罗米修斯的缓存上报

#### Version V2.7.3
> 1. 重新构建master

#### Version V2.7.2
> 1. 弹幕计数迁移到job

#### Version V2.7.1
> 1. 历史弹幕pagesize=5000

#### Version V2.7.0
> 1. update infoc sdk

#### V2.6.3
> 1. 弹幕subject缓存穿透

#### V2.6.2
> 1. 海外用户校验

#### V2.6.1
> 1. 弹幕蒙版相关接口添加

#### V2.6.0
> 1. 对接发号器

#### V2.5.12
> 1. 新增level1用户弹幕长度超过20提示

#### V2.5.11
> 1. 高级弹幕购买需要参数优化
> 2. 高级弹幕相关接口结构体优化
> 3. user filter列表RPC调用结构体优化

#### V2.5.10
> 1. 弹幕发送增加反垃圾过滤备注

#### V2.5.9
> 1. add csrf

#### V2.5.8
> 1. update bm engine

#### V2.5.7
> 1. 弹幕池变更为0时减小move_cnt

#### V2.5.6
> 1. 修复关闭分段弹幕时gzip bug

#### V2.5.5
> 1. up主拉黑用户接口增加限流
> 2. 历史弹幕接口及历史弹幕日期索引接口增加限流

#### V2.5.4
> 1. up主拉黑用户接口更新

#### V2.5.3
> 1. 重构历史弹幕

#### V2.5.2
> 1. up主屏蔽词添加接口兼容协管添加

#### V2.5.1
> 1. 修复关闭分段弹幕不生效的bug

#### V2.5.0
> 1. 迁移blademaster

#### V2.4.32
> 1. 移除屏蔽词sql中的排序
> 2. 过滤弹幕编辑中dmid=0的参数

#### V2.4.31
> 1. 修复dmid为空的操作日志

#### V2.4.30
> 1. 修复稿件mid变化导致的bug

#### V2.4.29
> 1. 变更弹幕池时更新move_count in dm_subject

#### V2.4.28
> 1. 使用弹幕主题中的mid做up主权限校验
#### V2.4.27
> 1. 添加高级弹幕返回（存储bfs）

#### V2.4.26
> 1. fix archiveinfos 批量100请求

#### V2.4.25
> 1. 新增弹幕属性变更rpc方法
> 2. 弹幕状态变更支持举报脚本删除

#### V2.4.24
> 1. dm_subject新增字段

#### V2.4.23
> 1. rank<=15000的用户设置字幕弹幕时增加条数限制

#### V2.4.22
> 1. 修复用户屏蔽词添加的bug
> 2. up主屏蔽词接口改为json参数

#### V2.4.21
> 1. 视频时长缓存使用memcached

#### V2.4.20
> 1. 用户添加屏蔽词接口改为批量新增接口

#### V2.4.19
> 1. 协管拉黑用户时日志记录弹幕内容

#### V2.4.18
> 1. 高级弹幕购买迁移dm -> dm2

#### V2.4.17
> 1. 迁移main目录

#### V2.4.16
> 1. up主屏蔽词对接新库
> 2. 重构屏蔽词服务

#### V2.4.15
> 1. 去掉oplogproc弹幕发送日志时间戳

#### V2.4.14
> 1. 修复OpLog弹幕发送日志时间戳

#### V2.4.13
> 1. 修复OpLog弹幕日志消费队列溢出

#### V2.4.12
> 1. 修复OpLog弹幕日志时间戳生成

#### V2.4.11
> 1. 修复弹幕风纪委bug
> 2. 弹幕实名制增加开关以及白名单

#### V2.4.10
> 1. 使用account-service v7

#### V2.4.9
> 1. 修复弹幕监控计数问题

#### V2.4.8
> 1. 移除剧透弹幕相关逻辑

#### V2.4.7
> 1. up主在自己的视频下发送顶端、底端弹幕不再受等级的限制

#### V2.4.6
> 1. 修复用户撤回判断弹幕逻辑

#### V2.4.5
> 1. 重构弹幕监控

#### V2.4.4
> 1. 添加发送屏蔽状态日志

#### V2.4.3
> 1. up主屏蔽词屏蔽类型 painc 修复

#### V2.4.2
> 1. up主屏蔽词屏蔽类型修改

#### V2.4.1
> 1. 分段弹幕开启实名制

#### V2.4.0
> 1. 添加业务操作infoc2日志
> 2. test 去除simplejson 依赖

#### V2.3.4
> 1. 修复特殊弹幕管理的bug

#### V2.3.3
> 1. 弹幕xml替换非法字符

#### V2.3.2
> 1. 新增弹幕状态变更rpc method
> 2. 新增弹幕池变更rpc method

#### V2.3.1
> 1. 新增用户端弹幕管理
> 2. 移除无用代码

#### V2.3.0
> 1. 新增弹幕计数rpc service
> 2. 提供弹幕发送、实名制开关状态给详情页

#### v2.2.23
> 1. localcache oid 请求大数据

#### v2.2.22
> 1. 修改弹幕池和弹幕状态平行权限的问题

#### v2.2.21
> 1. 修复cache.Save context

#### V2.2.20
> 1. 使用UserInfo替换Card

#### V2.2.19
> 1. 使用UserInfo替换MyInfo

#### V2.2.18
> 1. localcache oid不再请求大数据

#### V2.2.17
> 1. 弹幕发送计数迁移到dm2-job

#### V2.2.16
> 1. 调整弹幕发送时发送给kafka的key
> 2. 分段弹幕从缓存获取dmid时增加限制

#### V2.2.15
> 1. local cache dm subject and video duration

#### V2.2.14
> 1. 移除DmDmpost-T的双写

#### V2.2.13
> 1. 弹幕发送添加反垃圾功能

#### V2.2.12
> 1. 增加禁止发送弹幕功能
> 2. xml header中增加state字段

#### V2.2.11
> 1. 弹幕发送更新 dm_monitor 表

#### V2.2.10
> 1. 热加载local cache 配置

#### V2.2.9
> 1. 修复up正则匹配bug

#### V2.2.8
> 1. 添加弹幕和分段弹幕白名单,做本地缓存
> 2. 给大数据接口增加aid 参数

#### v2.2.7
> 1. 打开弹幕广播功能

#### V2.2.6
> 1. 修复数据库返回值导致的panic

#### V2.2.5
> 1. 修复DmDMpost-T双写
> 2. 关闭发送弹幕接口中的广播，待双写关闭时打开

#### V2.2.4
> 1. 弹幕发送双写 DmDmpost-T

#### V2.2.3
> 1. 迁移弹幕广播功能到弹幕发送接口

#### V2.2.2
> 1. 变更屏蔽词缓存key

#### V2.2.1
> 1. 修复高级弹幕购买缓存清理

#### V2.2.0
> 1. 新增弹幕发送接口

#### V2.1.18
> 1. 批量获取subjects状态筛选

#### V2.1.17
> 1. 去除弹幕镜像功能

#### V2.1.16
> 1. 风纪委弹幕列表sql优化2

#### V2.1.15
> 1. 风纪委弹幕列表sql优化

#### V2.1.14
> 1. subject 被关闭时仍返回弹幕列表
> 2. 更新prom 统计方法

#### V2.1.13
> 1. 注释单元测试(暂时)

#### V2.1.12
> 1. 点赞接口cid别名问题

#### V2.1.11
> 1. fix rec_switch flag

#### V2.1.10
> 1. 最新弹幕新老接口拆分

#### V2.1.9
> 1. 弹幕实名制
> 2. 更改实名白名单

#### V2.1.8
> 1. 弹幕点赞接口

#### V2.1.7
> 1. 移除childpool 3相应逻辑

#### V2.1.6
> 1. 修改最新弹幕返回aid

#### V2.1.5
> 1. 修改最新弹幕字段

#### V2.1.4
> 1. 更改弹幕列表v1回源逻辑

#### V2.1.3
> 1. 新增弹幕分布统计接口

#### V2.1.2
> 1. 新增弹幕弹幕镜像
> 2. 创作中心弹幕列表删除分P不存在bug

#### V2.1.1
> 1. 修复up删除弹幕不更新count
> 2. 同步老数据库dm_index 的 childpool 字段

#### V2.0.11
> 1. 修复缓存bug

#### V2.0.10
> 1. Page2 to Page3
> 2. Video2 to Video3

#### V2.0.9
> 1. 使用singleflight限制并发

#### V2.0.8
> 1. 新增up主最新1000条弹幕

#### V2.0.7
> 1. 新增一波单元测试

#### V2.0.6
> 1. 分段弹幕和非分段弹幕redis缓存分开

#### V2.0.5
> 1. 采用sync.pool优化GC

#### V2.0.4
> 1. 使用bytes.Buffer优化xml生成

#### V2.0.3
> 1. 新增弹幕风纪委功能

#### V2.0.2
> 1. 优化弹幕列表刷新逻辑

#### V2.0.1
> 1. 新增首页ajax弹幕列表接口

#### V2.0.0
> 1. 重构主站弹幕列表

#### V1.0.8
> 1. 依据dmid查询弹幕内容时并发查询

#### V1.0.7
> 1. 调用大数据接口时limit=2*maxlimit
> 2. 根据oid调用大数据多域名接口

#### V1.0.6
> 1. 调整弹幕显示数量为2*maxlimit

#### V1.0.5
> 1. 不再解析大数据json
> 2. 调整大数据接口降级逻辑

#### V1.0.4
> 1. 修复 .so后缀使用时res["code"]问题

#### V1.0.3
> 1. 合并到大仓库

#### V1.0.2
> 1. 添加大数据降级日志

#### V1.0.1
> 1. 修复弹幕上限分页逻辑

#### V1.0.0
> 1. 族群弹幕初始版本
