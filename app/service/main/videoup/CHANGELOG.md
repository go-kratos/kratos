#### 稿件内部接口

##### Version 2.2.0
>1.添加获取稿件附加属性数据的HTTP接口

##### Version 2.1.2
>1.联合投稿更新消息

##### Version 2.1.1
>1.获取未通过审核稿件列表

##### Version 2.1.0
>1.联合投稿批量修改staff接口
>2.保障staff state 1/2

##### Version 2.0.0
>1.联合投稿

##### Version 1.33.3
>1.IPv6upgrade

##### Version 1.33.2
>1.投稿支持新功能灰度   使用方式是 需要灰度的功能块 走一下 s.checkGrayMid方法  同时在用户组21里维护，清空该组就等于全量
>2.投稿支持bgm 属性位 23bit

##### Version 1.33.1
>1.投稿发动态支持投票业务

##### Version 1.33.0
>1.投稿支持bgm属性位  bit 23

##### Version 1.32.6
>1.升级并精简化自定义错误码  

##### Version 1.32.5
>1.修复和基础库错误码不兼容的问题,需重新发版   

##### Version 1.32.4
>1.具体到某一分P的错误统一沿用interface的VideoError 

##### Version 1.32.3
>1.添加视频分辨率稿件信息接口

##### Version 1.32.2
>1.app 移动投稿支持LBS信息上报

##### Version 1.32.1
>1.videoup/ugc/edit/mission支持绑定活动tag 或者取消稿件活动
>2.videoup/flow/list/juge 支持最多200个aid查询

##### Version 1.32.0
>1.投稿支持UGC付费

##### Version 1.31.1
>1.投稿消息 add_archive 新增up_from 投稿来源信息

##### Version 1.31.0
>1.新增 流量管理接入pgc禁止项  /videoup/flow/entry/oid
>2.新增 流量管理支持批量查询禁止项明细列表  /videoup/flow/list/judge

##### Version 1.30.0
>1.新增 稿件分享动态配置查询接口  /videoup/setting/dynamic

##### Version 1.29.10
>1.dao unittest redo


##### Version 1.29.9
>1.fix 修复多p编辑超时

##### Version 1.29.8
>1.fix 修复超时

##### Version 1.29.7
>1.feature 多分p稿件编辑进行异步并行处理

##### Version 1.29.6
>1.obtain/cid 和 assiginCid 加分布式锁保证cid 与 filename 一一对应

##### Version 1.29.5
>1.升级ip获取的方式，使用metadata.RemoteIP

##### Version 1.29.4
> 1.定时发布表增加软删除字段deleted_at

##### Version 1.29.3
> 1.去除代码模式
> 2.增加prom监控
> 3.增加filename过期的分p报错信息

##### Version 1.29.2
> 1.增加视频和音频绑定的异步接口

##### Version 1.29.1
> 1.videoup/views pagesize上限为50
> 2.modify_archive支持分区修改

##### Version 1.29.0
> 1.新增接口修改稿件活动信息  /ugc/edit/mission

##### Version 1.28.12
> 1.修复基础库bug

##### Version 1.28.11
> 1.移除net/http/parse

##### Version 1.28.10
> 1.移除频道回查的inner_attr/attribute重置逻辑
> 2.modify_archive新增tag变动和新增分P的参数

##### Version 1.28.9
> 1.去除dm_index同步与dede库依赖

##### Version 1.28.8
> 1.增加稿件查询简易版

##### Version 1.28.7
> 1./videoup/views 新增批量aids 查询最大值 20

##### Version 1.28.6
> 1.busSendMsg 拦截异常filename 防止进入无限重试redis 队列

##### Version 1.28.5
> 1.投稿filename 48小时错误码
> 2.去掉无意义的error日志

##### Version 1.28.4
> 1.新增投稿filename 48小时超时验证

##### Version 1.28.3
> 1.fix 私单属性判断导致的查询double

##### Version 1.28.2
> 1.换源或改变tag，重置频道回查属性位

##### Version 1.28.1
> 1.fix自定义错误,兼容通用的错误码接口

##### Version 1.28.0
> 1.迁移bm

##### Version 1.27.14
> 1.去掉UGC稿件描述(desc)必传的校验

##### Version 1.27.13
> 1.移除提交稿件author字段

##### Version 1.27.12
> 1.Video结构体返回duration到json

##### Version 1.27.11
> 1.专栏推荐禁止取消操作加入完整性逻辑 禁止业务不存在情况的取消逻辑

##### Version 1.27.10
> 1.专栏推荐禁止支持取消操作

##### Version 1.27.9
> 1.fix pad filename 404 loop pad bug

##### Version 1.27.8
> 1.迁移main path

##### Version 1.27.7
> 1.fix 私单字段数据强一致性 porder official brandName bug

##### Version 1.27.6
> 1.add contributor

##### Version 1.27.5
> 1.流量禁止应用支持粉丝动态 AttrBitNoPushBplus 20

##### Version 1.27.4
> 1.去除statsd

##### Version 1.27.3
> 1.增加稿件审核"阻塞"状态

##### Version 1.27.2
> 1.支持async edit

##### Version 1.27.1
> 1.完善flow go lint

##### Version 1.27.0
> 1.提供非稿件业务流量mid入口 /videoup/flow/entry/mid
> 2.私单稿件默认设置前端展示 show_front=1

##### Version 1.26.0
> 1.删除稿件校验mid
> 2.编辑稿件针对新增分P验证必须是新投视频

##### Version 1.25.5
> 1.porder/config/list 新增rank 排序
> 2.投稿支持日文标识

##### Version 1.25.4
> 1./porder/arc/list 返回所有字段，按照id倒叙 

##### Version 1.25.3
> 1.pgc 增加视频云pubagent

##### Version 1.25.2
> 1.将稿件属性第13位hideclick改成limit_area

##### Version 1.25.1
> 1.修复新视频查询err
> 2.统一给成长计划的返回数据为{}而不是nil，应成长计划要求

##### Version 1.25.0
> 1.私单二期业务
> 2.私单流量聚合

##### Version 1.24.2
> 1.fix bug: 如果分P没有做任何修改并且只修改稿件信息，只添加稿件历史，不添加分P历史  

##### Version 1.24.1
> 1. update:添加和编辑的时候，批量插入视频的历史记录，提高响应速度 

##### Version 1.24.0
>1.私单配置返回给videoup-interface和creative-interface  
>2.私单按照时间区间查询，提供给商业产品部门的成长计划项目组
>3.添加接口查稿件打回理由关联的申诉tag_id


##### Version 1.23.1
> 1. 添加视频审核时长等级接口

##### Version 1.23.0
> 1. video相关全部切到新表

##### Version 1.22.2
> 1. ugc投稿只发送ugc_submit消息

##### Version 1.22.1
> 1. ugc投稿，只发送sync_cid消息，sync_cid中submit=1

##### Version 1.22.0
> 1. 支持pgc drm投稿

##### Version 1.21.2
> 1. 实验室投稿发送ugc_submit消息

##### Version 1.21.1
> 1.投稿不再自动添加海外禁止

##### Version 1.21.0
> 1.新增 pgc topic PGCVideoup2Bvc

##### Version 1.20.4
> 1.将稿件attribute中的第9位作为is_pgc字段

##### Version 1.20.3
> 1.fix 代码模式下不执行syncCid

##### Version 1.20.2
> 1.开放videoup/view 接口 src_type属性

##### Version 1.20.1
> 1.fix重试监控报警

##### Version 1.20.0
> 1.新增重试监控报警

##### Version 1.19.9
> 1.去掉稿件推荐双写的逻辑， 去除dede老库的依赖

##### Version 1.19.8
> 1.synccid异步发消息逻辑修改

##### Version 1.19.7
> 1.删除稿件动作不更新视频关联关系（-100）

##### Version 1.19.6
> 1.fix pgc提交发submit消息 未注册分类

##### Version 1.19.5
> 1.pgc提交发submit消息  

##### Version 1.19.4
> 1.稿件回溯编辑archive_video_relation state=0强制

##### Version 1.19.3
> 1.稿件支持 dynamic 特性  

##### Version 1.19.2
>1.hot fix: 强制 archive_video 与 archive_video_relation id一致  

##### Version 1.19.1
> 1.推荐表， 由于权限不足，老表操作使用创作中心的账号  

##### Version 1.19.0
> 1.推荐表， 双写，双更新，查询先查老数据，再查新数据   

##### Version 1.18.0
>1.新增cid查询接口  /videoup/query/cid  
>2.新增流量融合的禁止项用户组聚合禁止项应用到投稿  
>3.新增流量融合的禁止项用户组聚合数据查询接口 /videoup/up/forbid  

##### Version 1.17.3
>1.新增查询特殊用户组接口  /videoup/up/special  


##### Version 1.17.2
>1.hot fix: 同一个filename sync cid两次  

##### Version 1.17.1
> 1.去除appkey参数  

##### Version 1.17.0
> 1. 新增生成cid的接口  
> 2. 去除老的dm_index的依赖  

##### Version 1.16.3
> 1. 增加网安的md5数据导入数据库的内部接口, 接口地址:/videoup/ns/Md5  

##### Version 1.16.2
> 1. 用户修改转载来源聚合稿件状态   

##### Version 1.16.1
> 1. 投稿添加视频的操作日志的操作人错误修复  

##### Version 1.16.0
> 1. 合入大仓库  
> 2. 新增简版archive和videos接口  

##### Version 1.15.3
> 1. 同步new_video的时候，增加ctime字段  

##### Version 1.15.2
> 1. desc_format的components改成string类型  

##### Version 1.15.1
> 1. desc_format的typeid修改  

##### Version 1.15.0
> 1. 天马和私单融合方案查询接口  
> 2. 提供archive_desc_format的数据接口  

##### Version 1.14.0
> 1. archive_video 拆表双写逻辑  

##### Version 1.13.1
> 1.新增稿件投诉接口  

##### Version 1.13.0
> 1.去除prepare  
> 2.恢复一审理由聚合  

##### Version 1.12.1
> 1.去除海外禁止特殊分区若干  

##### Version 1.12.0
> 1.UGC新增简介描述功能(不含PGC)  

##### Version 1.11.4
> 1.新增批量通过cid获取稿件aid列表  

##### Version 1.11.3
> 1.设置保留稿件的state为-100  

##### Version 1.11.2
> 1.支持内部保留aid特性  

##### Version 1.11.1
> 1.根据up主的mid，分页获取aids数据  

##### Version 1.11.0
> 1.UGC新增私单功能模块(不含PGC)  
> 2.PGC 编辑稿件删除分P发delete_video消息  

##### Version 1.10.0
> 1.更新go-common,使用log-agent  

##### Version 1.9.6
> 1.去掉一些海外的分区id  

##### Version 1.9.5
> 1.PGC新增外链upfrom  

##### Version 1.9.4
> 1.vendor add burntsushi  

##### Version 1.9.3
> 1.修改重试databus的key   

##### Version 1.9.2
> 1.处理分P为0的prom统计  

##### Version 1.9.1
> 1.细分prom统计  

##### Version 1.9.0
> 1.加入普罗米修斯监控统计日志  

##### Version 1.8.3
> 1.增加机密PGC投稿接口  

##### Version 1.8.2
> 1.限制30天以内的历史纪录回溯  

##### Version 1.8.1
> 1.通过cid获取video信息  

##### Version 1.8.0
> 1.代码模式实现方式改为通过cid  

##### Version 1.7.5
> 1.打回状态非渣转完成不用发modify消息  

##### Version 1.7.4
> 1.新增modify_video的databus消息  

##### Version 1.7.3
> 1.启用新分区表  

##### Version 1.7.2
> 1.修复删除aid老表依赖的bug  

##### Version 1.7.1
> 1.删除aid老表依赖  
> 2.修改活动ID、商单ID、定时发布时间聚合稿件状态  

##### Version 1.7.0
> 1.接入新配置中心  

##### Version 1.6.0
> 1.日志清理和修改  
> 2.添加稿件回溯记录  
> 3.提供稿件回溯记录查询  

##### Version 1.5.1
> 1.删除track功能，移到admin  
> 2.更新vendor  

##### Version 1.5.0
> 1.增加批量稿件信息获取  
> 2.增加稿件TAG更新接口  

##### Version 1.4.0
> 1.去掉没用代码，挪admin中  

##### Version 1.3.3
> 1.更改稿件状态赋值和日志  

##### Version 1.3.2
> 1.修改暂缓状态可编辑  

##### Version 1.3.1
> 1.修复attr日志bug  

##### Version 1.3.0
> 1.删除稿件，同步老表状态  
> 2.移除删除稿件的扣硬币和解绑活动的代码  
> 3.删除稿件，同步视频云删除xcode小于2的视频  
> 4.修改同步删除视频云的route  

##### Version 1.2.4
> 1.修复view接口，rejectreason字段逻辑错误  

##### Version 1.2.3
> 1.pgc外链sohu和hunan的dm_index处理  
> 2.添加稿件，round默认值从-1变为0  

##### Version 1.2.2
> 1.修复pgc外链未插Dm_index的bug  

##### Version 1.2.1
> 1.pgc外链默认二转完成  

##### Version 1.2.0
> 1.增加商单ID和NO_PRINT逻辑  
> 2.PGC自动过审和attr更新  
> 3.商业产品attr位更新  
> 4.修复默认允许添加tag  

##### Version 1.1.1
> 1.支持PGC外链上传  
> 2.分P的audit理由根据状态返回  

##### Version 1.1.0
> 1.增加更新attr  
> 2.PGC自动过审(差相关业务同步逻辑)  
> 3.地区限制记录  

##### Version 1.0.1
> 1.增加recover后事务回滚  

##### Version 1.0.0
> 1.稿件新增编辑接口  