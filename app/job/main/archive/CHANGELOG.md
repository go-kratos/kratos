#### 稿件字段同步

##### Version 2.13.7
> 1.去除track表迁移first_pass表

##### Version 2.13.6
> 1.同步联合创作人archive_staff

##### Version 2.13.5
> 1.修复err后重复retry tranResult & -->chan

##### Version 2.13.4
> 1.走account的grpc

##### Version 2.13.3
> 1.扩展archive-service的配置

##### Version 2.13.2
> 1.账号信息空判断

##### Version 2.13.1
> 1.账号databus有重复消息

##### Version 2.13.0
> 1.账号notify走databus，记录日志

##### Version 2.12.13
> 1.账号notify走databus

##### Version 2.12.12
> 1.增加databus消费，重新生成用户昵称和头像缓存

##### Version 2.12.11
> 1.修改第一P的判断逻辑

##### Version 2.12.10
> 1.同步稿件时判断是否为vupload

##### Version 2.12.9
> 1.更新archive表的视频id，分辨率字段

##### Version 2.12.8
> 1.新表字段修改

##### Version 2.12.7
> 1.增加分辨率字段

##### Version 2.12.6
> 1.全量新表

##### Version 2.12.5
> 1.灰度新表

##### Version 2.12.4
> 1.切换BM

##### Version 2.12.3
> 1.计算稿件总时长

##### Version 2.12.2
> 1.聚合所有topic消费  

##### Version 2.12.1
> 1.发送statDm-T消息  

##### Version 2.12.0
> 1.迁移到主站目录下  

##### Version 2.11.1
> 1.弹幕计数做聚合  

##### Version 2.11.0
> 1.弹幕计数  

##### Version 2.10.4
> 1.删除所有插入track的代码  

##### Version 2.10.3
> 1.配置offset，停止track记录  

##### Version 2.10.2
> 1.补充单元测试  

##### Version 2.10.2
> 1.改进报警文案  

##### Version 2.10.1
> 1.pgc异步  

##### Version 2.10.0
> 1.增加监控  
> 2.pgc与ugc分开  

##### Version 2.9.4
> 1.增加无限塞回重试逻辑  

##### Version 2.9.3
> 1.增加force_sync消息  

##### Version 2.9.2
> 1.处理first_round消息  

##### Version 2.9.1
> 1.全量通过databus更新稿件缓存  

##### Version 2.9.0
> 1.通过databus消息更新稿件缓存  
> 2.灰度30%  

##### Version 2.8.1
> 1.调整 passed 然后又改了回来 = =

##### Version 2.8.0
> 1.archive_track 迁移至 hbase  
> 2.调整 passed 逻辑改查 archive_oper 表  

##### Version 2.7.1
> 1.databus通知吐出dynamic字段  

##### Version 2.7.0
> 1.删除hzxs缓存清楚  

##### Version 2.6.1
> 1.archive增加dynamic字段  

##### Version 2.6.0
> 1.去除pgc的databus订阅  

##### Version 2.5.2
> 1.回滚，video relation表数据有误  

##### Version 2.5.1
> 1.result库支持video分表逻辑(灰度10%的稿件)   

##### Version 2.5.0
> 1.删除cdn代码  
> 2.完善更新db的日志  

##### Version 2.4.1
> 1.切新的httpclient  

##### Version 2.4.0
> 1.同步result库判断微调  

##### Version 2.3.0
> 1.删除purge cdn逻辑  

##### Version 2.2.7
> 1.routine数量配置化&优化报警逻辑  

##### Version 2.2.6
> 1.多routine更新数据库&缓存&每分钟报警  

##### Version 2.2.5
> 1.cid为0不同步result  

##### Version 2.2.4
> 1.add error log  

##### Version 2.2.3
> 1.fix cids 0  

##### Version 2.2.1
> 1.修复pgc同步result逻辑  

##### Version 2.2.0
> 1.video缓存更新逻辑优化  

##### Version 2.1.2
> 1.pgc第二次过审不同步result库  

##### Version 2.1.1
> 1.修复videos插入失败的bug   

##### Version 2.1.0
> 1.模块分层  

##### Version 2.0.0
> 1.大仓库版本，依赖新的go-common  
> 2.archive_result库的archive_video表有数据更新会调archive-service接口更新分P详情信息  

##### Version 1.14.0
> 1.databus send error 时无限重试  

##### Version 1.13.0
> 1.暂时去掉HK节点的通知  
> 2.增加notify的databus  

##### Version 1.12.0
> 1.增加https页面的cdn purge  

##### Version 1.11.0
> 1.增加result库的partition消费监控  

##### Version 1.10.0
> 1.update field增加group2  

##### Version 1.9.1
> 1.archive的channal根据aid取余，每个channal只有一个goroutine消费，避免消费乱序  

##### Version 1.9.0
> 1.升级go-common&go-business  
> 2.兼容manager后台修改稿件归属的mid缓存  

##### Version 1.8.0
> 1.所有缓存清理走archive_result库  

##### Version 1.7.15
> 1.迁移archive-service的group2缓存到result库  

##### Version 1.7.14
> 1.增加monitor/ping   

##### Version 1.7.13
> 1.group1使用archive_result库更新缓存  

##### Version 1.7.12
> 1.bugfix,任务分发:稿件修改分区，写redis逻辑错误    

##### Version 1.7.11
> 1.移除所有group1的调用  

##### Version 1.7.10
> 1.bugfix,修复track的redis  

##### Version 1.7.9
> 1.bugfix,修复track的redis  

##### Version 1.7.8
> 1.bugfix,修复sql语句  

##### Version 1.7.7
> 1.增加香港节点的配置  

##### Version 1.7.6
> 1.已删除的稿件不记录track  

##### Version 1.7.5
> 1.dede_arctype替换为archive_type并修改相关逻辑  
> 2.添加杭州group  

##### Version 1.7.5
> 1.视频track变更增加标题简介  

##### Version 1.7.4
> 1.去掉feed push  

##### Version 1.7.3
> 1.升级go-common和go-business  
> 2.接入新版配置中心  

##### Version 1.7.2
> 1.增加proc处理  
> 2.修改track video待审重复计数bug  

##### Version 1.7.1
> 1.增加track video  

##### Version 1.7.0
> 1.增加track信息记录  
> 2.计数更新硬币数  

##### Version 1.6.24
> 1.增加后台报表统计  
> 2.track表备注增加活动id  
> 3.增加upcache错误重试  

##### Version 1.6.23
> 1.增加接口失败重试  

##### Version 1.6.22
> 1.稿件变更无脑add/del  
> 2.基于无脑变更，去掉发布时间和属性变化通知  

##### Version 1.6.21
> 1.升级vendor  

##### Version 1.6.20
> 1.增加archive-service group2 缓存增量更新  

##### Version 1.6.19
> 1.修复日志错误  

##### Version 1.6.18
> 1.移除stat的databus消费

##### Version 1.6.17
> 1.移除评论注册  

##### Version 1.6.16
> 1.archive缓存走RPC  

##### Version 1.6.15
> 1.根据state状态记录access变更  

##### Version 1.6.14
> 1.track去除没必要的记录  

##### Version 1.6.13
> 1.内部接口地址规范化  

##### Version 1.6.12
> 1.track兼容round改动2  

##### Version 1.6.11
> 1.track兼容round的改动  

##### Version 1.6.10
> 1.track增加attr和access  

##### Version 1.6.9
> 1.增加attr不在列表输出的判断  

##### Version 1.6.8
> 1.增加纪录片的表同步  

##### Version 1.6.7
> 1.up更改后通知评论替换subject的mid  

##### Version 1.6.6
> 1.修复track时间错误  

##### Version 1.6.5
> 1.修复xcode判断错误  

##### Version 1.6.4
> 1.增加PGC表同步逻辑  

##### Version 1.6.3
> 1.临时去除attr发邮件判断  

##### Version 1.6.2
> 1.增加track追踪  

##### Version 1.6.1
> 1.修复insert时的json解析错误  

##### Version 1.6.0

> 1.稿件依赖databus  

##### Version 1.5.7

> 1.kafka to databus

##### Version 1.5.6

> 1.增加insert事件 评论subject注册  

##### Version 1.5.5

> 1.增加monitor监控  

##### Version 1.5.4

> 1.fix sub databus bug  

##### Version 1.5.3

> 1.刷新cdn加条件  

##### Version 1.5.2

> 1.稿件计数消费databus  

##### Version 1.5.1

> 1.更新stat cache的Rpc为2  

##### Version 1.5.0

> 1.番剧和电影状态变更需要发送邮件  

##### Version 1.4.4

> 1.[consumer]去除archive binlog调用tag change接口

##### Version 1.4.3

> 1.[consumer]稿件推送动态改为过审就推  

##### Version 1.4.2

> 1.[consumer]修复tag过审同步问题  

##### Version 1.4.1

> 1.[consumer]增加sleep控制消费速度  

##### Version 1.4.0

> 1.[consumer]增加稿件变动purge  

##### Version 1.3.0

> 1.[consumer]增加archive的事件字段变更，新老字段  

##### Version 1.2.0

> 1.[consumer]修改archive的事件注册，分通过、不通过、无差别删除cache  

##### Version 1.1.5

> 1.[consumer]修改archive的cache更新接口  

##### Version 1.1.4

> 1.[consumer]修复没过审视频被推tag动态问题  

##### Version 1.1.3

> 1.[consumer]修正通过tag服务批量获取tag信息签名错误  
> 2.[consumer]修改调用tag change 接口mid错误  

##### Version 1.1.2

> 1.[consumer]feed调用历史数量判断修改，防止feed不成功  

##### Version 1.1.1

> 1.[consumer]feed调用加稿件审核历史判断  

##### Version 1.1.0

> 1.[consumer]增加分区视频列表cache更新和删除  
> 2.[consumer]修复改tag的更新  

##### Version 1.0.0

> 1.[producer]binlog同步  
> 2.[consumer]信息消费更新缓存等
