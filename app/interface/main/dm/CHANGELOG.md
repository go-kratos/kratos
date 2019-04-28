### dm的Gateway服务

#### Version 3.18.2
> 1. fix dmid =0

#### Version 3.18.1
> 1. remove cache.Cache

#### Version 3.18.0
> 1. account service使用grpc

#### Version 3.17.12
> 1. 修改创作中心弹幕举报排序字段

#### Version 3.17.11
> 1. 创作中心弹幕举报增加排序

#### Version 3.17.10
> 1. 修复up主弹幕保护时的bug

#### Version 3.17.9
> 1. 未处理的弹幕审核不发送通知的

#### Version 3.17.8
> 1. 弹幕举报稿件搜索添加二审

#### Version 3.17.7
> 1. rebuild master

#### Version 3.17.6
> 1. 弹幕举报deleted字段更新
> 2. 获取稿件信息去除多余代码

#### Version 3.17.5
> 1. up主处理举报弹幕验证owner

#### Version 3.17.4
> 1. 弹幕举报结果为空不返回null

#### Version 3.17.3
> 1. 弹幕举报稿件以及更新接口对接搜索v3

#### Version 3.17.2
> 1. 移除视频维度的屏蔽词

#### Version 3.17.1
> 1. rebuild master

#### Version 3.17.0
> 1. 切换identify 为auth和verify

#### Version 3.16.4
> 1. fix rows.Err()

#### Version 3.16.3
> 1. 弹幕举报搜索使用sdk v3

#### Version 3.16.2
> 1. 修复保护弹幕不能举报的bug
> 2. 移除DMOld结构体

#### Version 3.16.1
> 1. 使用公共配置

#### Version 3.16.0
> 1. 弹幕举报对接搜索新索引
> 2. 优化部分notify相关代码
> 3. 移除多余http client

#### Version 3.15.3
> 1. 重新构建master

#### Version 3.15.2
> 1. 修复弹幕转移的bug

#### Version 3.15.1
> 1. 变更弹幕转移时目标视频不存在时的提示信息

#### Version 3.15.0
> 1. 移除无用的接口

#### Version 3.14.4
> 1. 修复弹幕举报隐藏恢复

#### Version 3.14.3
> 1. 修复弹幕保护列表接口

#### Version 3.14.2
> 1. 高级弹幕购买需要参数优化
> 2. 高级弹幕相关接口RPC调用结构体优化
> 3. user filter列表RPC调用结构体优化

#### Version 3.14.1
> 1. remove cors handler

#### Version 3.14.0
> 1. http接口迁移bm
> 2. 移除协管无用接口

#### Version 3.13.18
> 1. 用户屏蔽词增加接口增加防刷

#### Version 3.13.17
> 1. 修复弹幕转移的bug

#### Version 3.13.16
> 1. 修改弹幕转移状态变量名

#### Version 3.13.15
> 1. 修改隐藏时间为20h

#### Version 3.13.14
> 1. 使用弹幕主题表中的mid

#### Version 3.13.13
> 1. 修复弹幕状态修改参数传错bug

#### Version 3.13.12
> 1. fix archiveinfos 批量100请求

#### Version 3.13.11
> 1. global filter 因为没有数据，http层直接return空

#### Version 3.13.10
> 1. 修复up主拉黑用户cid为空的bug

#### Version 3.13.9
> 1. 根据举报有效分自动删除或隐藏弹幕

#### Version 3.13.8
> 1. 协管保护弹幕直接调RPC.EditDMAttr

#### Version 3.13.7
> 1. dm_subject新增字段

#### Version 3.13.6
> 1. 添加举报有效分

#### Version 3.13.5
> 1. 高级弹幕购买迁移dm -> dm2

#### Version 3.13.4
> 1. up主屏蔽编辑时增加oid

#### Version 3.13.3
> 1. 彻底移除表dm_index_filter

#### Version 3.13.2
> 1. 迁移main目录

#### Version 3.13.1
> 1. 屏蔽词迁移至dm2服务

#### Version 3.13.0
> 1. up 主屏蔽词重构

#### Version 3.12.6
> 1. 使用account-service v7 

#### Version 3.12.5
> 1. 修复用户撤回判断弹幕逻辑

#### Version 3.12.4
> 1. 修改高级弹幕列表限制

#### Version 3.12.3
> 1. 弹幕举报新增举报理由

#### Version 3.12.2
> 1. 修复弹幕转移导致的panic

#### Version 3.12.1
> 1. test 去除 simplejson 依赖
> 1. 修改弹幕状态和弹幕池接口调用
> 2. 获取若干dao层数据获取sql

#### Version 3.12.0
> 1. 将变更弹幕状态、变更弹幕池相关方法都转向dm2
> 2. 弹幕撤回临时先传up主mid

#### Version 3.11.9
> 1. 修复 cache.Save context

#### Version 3.11.8
> 1. 优化弹幕日志接口

#### Version 3.11.7
> 1. 弹幕转移日志和重试接口

#### Version 3.11.6
> 1. 高级弹幕请求信息中 url 修改

#### Version 3.11.5
> 1. 高级弹幕申请列表优化

#### Version 3.11.4
> 1. up主添加屏蔽词验证正则

#### Version 3.11.3
> 1. 高级弹幕申请列表添加aid

#### Version 3.11.2
> 1. 兼容屏蔽词双份缓存

#### Version 3.11.1
> 1. 兼容filter表脏数据

#### Version 3.11.0
> 1. 调整项目代码结构

#### Version 3.10.6
> 1. 调整弹幕转移部分代码

#### Version 3.10.5
> 1. 修复高级弹幕购买并发问题
> 2. 优化部分代码结构

#### Version 3.10.4
> 1. 移除拒绝高级弹幕购买相关接口
> 2. 修复高级弹幕缓存刷新

#### Version 3.10.3
> 1. 合并dao 层mc
> 2. 接入普罗米休斯

#### Version 3.10.2
> 1. 从dm项目中去除读comment库相关逻辑
> 2. 从稿件获取cid的mid

#### Version 3.10.1
> 1. 屏蔽词列表接口优化

#### Version 3.10.0
> 1. 优化重构部分函数
> 2. 整理冗余model

#### Version 3.9.4
> 1. 更新高级弹幕购买状态接口

#### Version 3.9.3
> 1. remove ecode.NoLogin

#### Version 3.9.2
> 1. 移除单元测试代码(暂时)

#### Version 3.9.2
> 1. 移除无用的最新弹幕消息代码
> 2. 移除rpc server
> 3. 移除refresh index 逻辑
> 4. 移除无用的ugc databus

#### Version 3.9.1
> 1. 修改filter使用的mc

#### Version 3.9.0
> 1. 创作中心屏蔽词列表和修改

#### Version 3.8.12
> 1. change archives2->archives3

#### Version 3.8.11
> 1. 变更弹幕购买提示

#### Version v3.8.10
> 1. Archive2 to Archive3

#### Version 3.8.4
> 1. 修复ecode message="0"的bug

#### Version 3.8.3
> 1. 弹幕自动删除条件限制

#### Version 3.8.2
> 1. 修复mysql使用

#### Version 3.8.1
> 1. 弹幕转移添加重试逻辑

#### Version 3.8.0
> 1.迁移大仓库

#### Version 3.7.9
> 1.高级弹幕购买重构

#### Version 3.7.7
> 1. 创作中心up主弹幕举报显示一二审理由

#### Version 3.7.6
> 1. remove mysql stmt

#### Version 3.7.5
> 1. 修复global filter rule 空缓存穿透

#### Version 3.7.4
> 1. 升级基础库
> 2. fix memcached链接泄露

#### Version 3.7.3
> 1. 修复协管日志

#### Version 3.7.2
> 1. 修改协日志
> 2. 修改添加up主屏蔽接口

#### Version 3.7.1
> 1. 自动删除时发信息的bug修复

#### Version 3.7.0
> 1. 弹幕举报超过10次自动删除

#### Version 3.6.2
> 1. 使用go-common/cache包进行异步缓存回刷屏蔽词缓存
> 2. 增加一些单元测试

#### Version 3.6.1
> 1. 修复一些错误日志

#### Version 3.6.0
> 1. 协管

#### Version 3.5.1
> 1. 修复弹幕保护由于sql关键字导致的错误

#### Version 3.5.0
> 1. 加入普罗米修斯监控
> 2. 接入logagent
> 3. 修复go-common中的一个bug
> 4. 弹幕已举报后重复举报则状态不再改变

#### Version 3.4.3
> 1. 修复弹幕撤回一个BUG
> 2. 普罗米休息监控
> 3. logagent

#### Version 3.3
> 1. 弹幕举报根据理由分成一审和二审

#### Version 3.2.12
> 1. 删除推荐弹幕接口
> 2. 删除一些没用的代码

#### Version 3.2.11
> 1. 我的弹幕功能接入搜索
> 2. 删除大量无用代码

#### Version 3.2.10
> 1. 移除实时弹幕大数据过滤

#### Version 3.2.9
> 1. 修复保护弹幕通知的地址错误

#### Version 3.2.8
> 1. 弹幕举报一二审对接搜索修改

#### Version 3.2.7
> 1. 更新正则解析器
> 2. 下线弹幕云推荐的功能

#### Version 3.2.6
> 1. 修复 up 主最新弹幕中存在的 bug
> 2. 限定正则屏蔽词解析的次数,防止出现无限解析的 bug

#### Version 3.2.5
> 1. 弹幕举报一二审

#### Version 3.2.4
> 1. 移除屏蔽词中使用的select for update的SQL语句
> 2. 变更配置文件，弹幕保护写市北主库

#### Version 3.2.3

> 1. 申请保护弹幕

#### Version 3.1.4

> 1. 新增的屏蔽规则已存在时返回屏蔽规则内容

#### Version 3.1.3

> 1. 弹幕实时id增加过期时间
> 2. 弹幕云屏蔽的参数配置通过配置文件设定
> 3. 当从cache层获取屏蔽规则出错时，不再直接返回-500，而是尝试从db load数据

#### Version 3.1.2

> 1. 禁止重复添加用户id到云屏蔽

#### Version 3.1.1

> 1. 升级go-business

#### Version 3.1.0

> 1. 修复屏蔽词可以重复添加的bug
> 2. 修复屏蔽词内容为空的bug
> 3. 更新go-common go-business

#### Version 3.0.6

> 1. 增加弹幕云屏蔽相关功能
> 2. 增加稿件弹幕计数功能,将弹幕计数写入hbase
> 3. 增加弹幕信息查询接口

#### Version 3.0.3

> 1. 实时过虑添加白名单，配置cid不走推荐弹幕和过虑弹幕。

#### Version 3.0.2

> 1. hbase里面没有推荐弹幕，直接从redis读取弹幕ID，如果不够从DB中补充（<实时弹幕ID）

#### Version 3.0.1

> 1. 支持up主最近弹幕
> 2. 支持没有推荐弹幕时，实时过虑弹幕从DB中加载

#### Version 3.0.0

> 1.redis保存最近cid弹幕ID
> 2.hbase中获取大数据推荐弹幕ID
> 3.通过redis弹幕ID和hbase弹幕ID从数据库中获取弹幕内容
