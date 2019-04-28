# article-job

##### Version 1.20.6
> 1. 删除脏数据

##### Version 1.20.5
> 1. fix read split

##### Version 1.20.4
> 1. 添加神马MIP推送

##### Version 1.20.3
> 1. recheck without act

##### Version 1.20.2
> 1. update waitGroup

##### Version 1.20.1
> 1. 长评接入

##### Version 1.19.6
> 1. 文章删除后关闭评论

##### Version 1.19.5
> 1. fix redis

##### Version 1.19.4
> 1. 消费redis中的用户阅读心跳数据，通过infoc上报阅读时长到数据平台

##### Version 1.19.3
> 1. add err log
> 2. fix nil pointer

##### Version 1.19.2
> 1. 优化sql语句

##### Version 1.19.1
> 1. 增加分区作者推荐计算

##### Version 1.18.1
> 1. 浏览数回查 day 或 view 任意为0 认为不生效

##### Version 1.18.0
> 1. 修复 set to setex

##### Version1.17.0
> 1. update infoc sdk

##### Version 1.16.1
> 1. 使用env

##### Version 1.16.0
> 1. 使用bm

##### Version 1.15.1
> 1. 修复生成最新投稿的问题

##### Version 1.15.0

> 1. 生成分区数据重构: 生成分区数据(最新文章 排序列表)移动到article-job中 
> 2. 作者权限缓存重构: 由redis移动到mc中
> 3. 搜索表增加attributes字段

##### Version 1.14.0
> 1. 评论/硬币/收藏 使用通用计数流
> 2. 增加动态消息重试

##### Version 1.13.1
> 1. 重构游戏缓存 remove redis pie
##### Version 1.13.0
> 1. 专栏草稿写库方式重构

### v1.12.0
> 1. 专栏阅读数防刷需求

### v1.11.0
> 定时更新热点标签

### v1.10.0
> 定时更新文集的总阅读数

### v1.9.0
> 专栏过审/打回/锁定向b+发送消息

### v1.8.0
> 专栏过审调用流量管理接口

### v1.7.14
> 修复活动接口

### v1.7.13
> 定期刷新作者列表

### v1.7.12
> 点赞时调用bili的活动排序接口

### v1.7.11
> 接入点赞计数databus

### v1.7.10
> fixed game sync bug

### v1.7.9
> 增加image_urls字段

### v1.7.8
> 接入点赞databus

### v1.7.7
> 修复草稿panicbug

### v1.7.6
> 1.游戏白名单迁移到数据库

### v1.7.5
> 1.接入点赞服务

### v1.7.4
> 1.fixed bug of update stats ctime

### v1.7.3
> 1.fixed close channel and game

### v1.7.2
> 1.游戏的文章变动时候进行通知  

### v1.7.1
> 1.更新http client   

### v1.7.0
> 1.增加阅读数排序   

### v1.6.0
> 1.接入搜索

### v1.5.0
> 1.动态更新排序列表缓存

### v1.4.1
> 1.修改刷新最新文章的时间

### v1.4.0
> 1.监听作者通过/拒绝事件  

### v1.3.4
> 1.草稿重试改成同步，异步重试有时序问题  

### v1.3.3
> 1.草稿添加无限重试  

### v1.3.2
> 1.调整添加、更新、删除文章缓存

### v1.3.1
> 1.刷新CDN添加APP端URL

### v1.3.0
> 1.binlog监听改为filtered_articles

### v1.2.3
> 1.去除filter相关逻辑

### v1.2.2
> 1.databus中的草稿消息去掉content字段

### v1.2.1
> 1.过滤操作放到文章里面做

### v1.2.0
> 1.使用新版过滤rpc  

### v1.1.0
> 1.草稿异步落库  
> 2.去掉conf.go中已经不用的刷新推荐数据的配置

### v1.0.2
> 1.当编辑直接在后台修改文章时更新缓存

### v1.0.1
> 1.加 prometheus 监控

### v1.0.0
> 1.项目初始化
