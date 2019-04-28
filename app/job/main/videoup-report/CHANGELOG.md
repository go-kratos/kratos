#### videoup-report-job
##### Versio 1.5.38
>1.使用up-service grpc

##### Version 1.5.37
>1.移除无用报表释放redis资源

##### Version 1.5.36
>1.激励回查添加错误重试

##### Version 1.5.35
>1.account使用grpc

##### Version 1.5.34
>1.规范化waitGroup使用

##### Version 1.5.33
> 1.修复up group返回json格式解析失败的问题

##### Version 1.5.32
> 1.修复UP Group API链接的bug

##### Version 1.5.31
> 1.增加激励回查白名单逻辑

##### Version 1.5.30
> 1.消费二审消息支持邮件开关

##### Version 1.5.29
> 1.修复复审任务删除失败

##### Version 1.5.28
> 1.增加激励回查逻辑

##### Version 1.5.27
> 1.记录任务的参数typeid,upfrom,upgroup到缓存,方便任务复审判断

##### Version 1.5.26
> 1.取消稿件活动的同时把活动tag也去掉

##### Version 1.5.25
> 1.添加视频monitor日志，有些视频不出-30的问题

##### Version 1.5.24
> 1.修复因稿件addit为空，导致不统计的bug

##### Version 1.5.23
> 1.稿件监控忽略PGC稿件

##### Version 1.5.22
> 1.增加邮件超限队列的消耗机会

##### Version 1.5.21
> 1.修改task使用secondary redis
> 1.权重参数从hash改为string

##### Version 1.5.20
> 1.增加视频审核监控

##### Version 1.5.19
> 1.增加权重前后时间日志

##### Version 1.5.18
> 1.订阅二审消息去重复开评论

##### Version 1.5.17
> 1.去除多次与list

##### Version 1.5.16
> 1.修复评论冻结问题

##### Version 1.5.15
> 1.状态改变会联动评论开关

##### Version 1.5.14
> 1.邮件从videoup-job迁移过来
> 2.邮件快慢分离，且统一控制发送api调用频率为5s/次

##### Version 1.5.13
> 1.hbase v2

##### Version 1.5.12
> 1.修改频道回查rpc的返回结构

##### Version 1.5.11
> 1.增加稿件状态停留统计

##### Version 1.5.10
>1.feature: 升级bm，初始化使用engine.Start

##### Version 1.5.9
> 1.定时发布表增加软删除字段deleted_at

##### Version 1.5.8
> 1.将已回查、待回查的日志记录到archive_oper的remark字段，与频道回查一致，方便统计报表

##### Version 1.5.7
> 1.tag同步绑定一级、二级分区名
> 2.videoup-job的tag同步全部迁移到本项目
> 3.编辑稿件时,分区修改触发tag同步
> 4.活动稿件不进入频道回查

##### Version 1.5.6
> 1.热门回查稿件忽略已存在的aid

##### Version 1.5.5
> 1.添加稿件热门回查功能

##### Version 1.5.4
> 1.去掉redis的大key: task_weight
> 2.生成任务记录日志不使用事务,允许日志记录失败

##### Version 1.5.3
> 1.基础库升级

##### Version 1.5.2
> 1.add_archive/modify_archive消息触发tag同步、频道回查落库、开启频道禁止

##### Version 1.5.1
> 1.完全迁移一审任务

##### Version 1.5.0
> 1.迁移一审任务

##### Version 1.4.10
> 1.修复report-job redis zadd hot empty key

##### Version 1.4.9
> 1.调整关闭顺序，避免日志遗漏

##### Version 1.4.8
> 1.使用blademaster

##### Version 1.4.7
> 1.archive_track.remark新增动态描述记录

##### Version 1.4.6

> 1.修复掉SQL中有or导致任务等待时候报表数据错误

##### Version 1.4.5

> 1.从一审任务等待时长报表中去掉定时发布的数据

##### Version 1.4.4

> 1.fix map 并发写panic
> 2.调整视频吞吐报表落库频率为5分钟一次

##### Version 1.4.3

> 1.fix close channel bug

##### Version 1.4.2

> 1.迁移path main

##### Version 1.4.1

> 1.去除statsd

##### Version 1.4.0

> 1.迁入稿件追踪和分发打点

##### Version 1.3.10

> 1.添加视频审核耗时的redis缓存

##### Version 1.3.9

> 1.添加统计10分钟内视频审核总耗时，总耗时=一转耗时+一审耗时+二转耗时+分发耗时

##### Version 1.3.8

> 1.解决视频进审数据统计bug
> 2.解决databus old message json.Unmarshal error bug

##### Version 1.3.7

> 1.使用监听binlog的方式统计一转、二转、分发耗时

##### Version 1.3.6

> 1.添加了一转、二转、分发耗时统计

##### Version 1.3.5

> 1.给一二三查打点数据添加uid字段

##### Version 1.3.4

> 1.修改archive的state字段类型

##### Version 1.3.3

> 1.从archive_oper表获取进入一二三查的时间

##### Version 1.3.2

> 1.添加一二三查耗时打点

##### Version 1.3.1

> 1.修复panic  

##### Version 1.3.0

> 1.into kratos  

##### Version 1.2.0

> 1.给map增加并发互斥锁  

##### Version 1.1.1

> 1.增加二审稿件移区统计  

##### Version 1.1.0

> 1.增加一审待审进审量统计  

##### Version 1.0.0

> 1.增加一审 task 等待耗时数据打点  