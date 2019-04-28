# article

##### Version 2.0.4
> 1. fix read split

##### Version 2.0.3
> 1. fix reCount article list read 

##### Version 2.0.2
> 1. fix reason
> 2. add keywords of view

##### Version 2.0.1
> 1. 喜欢 => 点赞
> 2. 关键词

##### Version 2.0.0
> 1. 专栏支持重复编辑

##### Version 1.52.3
> 1. 适配蓝版app

##### Version 1.51.2
> 1. 长评接入

##### Version 1.50.1
> 1. 迁移capture
> 2. 点赞加入反作弊
> 3. xss data-vote-id

##### Version 1.50.0
> 1. 修复点赞bug

##### Version 1.49.9
> 1. 修复并发bug

##### Version 1.49.8
> 1. 修改推荐数据源

##### Version 1.49.7
> 1. 增加推荐池分区统计
> 2. 增加时间删选

##### Version 1.49.6
> 1. 详情页推荐迭代
> 2. 点赞增加up mid

##### Version 1.49.5
> 1. 添加用户阅读心跳接口，记录写入redis交给job处理

##### Version 1.49.4
> 1. 添加用户阅读心跳接口，记录写入redis交给job处理

##### Version 1.49.3
> 1. fix nil pointer bug

##### Version 1.49.2
> 1. 开放创建文集的限制

##### Version 1.49.1
> 1. 左右滑动推荐内容策略调整

##### Version 1.48.15
> 1. 使用metadata中的remoteIP

##### Version 1.48.14
> 1. 去除tag服务无用代码

##### Version 1.48.13
> 1. 通知增加content字段
> 2. 新加草稿箱数量接口
> 3. 增加推荐作者列表接口
> 4. viewInfo增加转发数

##### Version 1.48.12
> 1. 一周年接口
> 2. 图片上传支持uri

##### Version 1.48.11
> 1. stat表过滤预加载

##### Version 1.48.10
> 1. 修改空数据error为info

##### Version 1.48.9
> 1. 去除重复error

##### Version 1.48.8
> 1. 空buvid不调用天马
> 2. 去除重复error
> 3. -404将为WARN

##### Version 1.48.7
> 1. 修改天马数据空日志级别为warning


##### Version 1.48.6
> 1. 历史记录过滤预加载

##### Version 1.48.5
> 1. 左右滑动修改天马获取数量

##### Version 1.48.4
> 1. 添加天马入口策略
> 2. 忽略articleSlide数据上报
> 3. 天马数据获取异常降级

##### Version 1.48.1
> 1. 增加创作中心稿件管理接口
> 2. viewInfo增加上下篇文章信息

##### Version 1.48.1
> 1. 添加databus阅读数回查

##### Version 1.47.3
> 1. 添加hbase share1，share0
> 2. 添加scheme bilibili

##### Version 1.47.2
> 1. 去掉debug和trace配置

##### Version 1.47.1
> 1. use new infoc

##### Version 1.47.0
> 1. 使用bm

##### Version 1.46.1

> 1. 升级tools/cache

##### Version 1.46.0

> 1. 生成分区数据重构: 生成分区数据(最新文章 排序列表)移动到article-job中 
> 2. 作者权限缓存重构: 由redis移动到mc中
##### Version 1.45.4

> 1. 修复xss漏洞

##### Version 1.45.3

> 1. 修复作者列表不全的问题

##### Version 1.45.2

> 1. 优化计算分区排序列表

##### Version 1.45.0

> 1. 支持过期热点显示 
> 2. 热点过滤禁止分区文章
> 3. 增加详情页banner
> 4. 支持同时开展多个活动

##### Version 1.44.0

> 1. 文章和文集接入历史记录

##### Version 1.43.2

> 1. 天马加日志

##### Version 1.43.1

> 1. 天马降级

##### Version 1.43.0

> 1. 用户文章列表接口增加文集字段
> 2. 增加上传图片API
> 3. 天马加入推荐运营位

##### Version 1.42.0

> 1. 增加web端文章文集列表接口
 
##### Version 1.41.3

> 1. 重构作者缓存 移除redis pie

##### Version 1.41.1

> 1. 修复创作中心列表页无数据的问题
> 2. 修复裁剪图片验证不全的问题

##### Version 1.41.0

> 1. 文集增加默认图片
> 2. 文集管理界面增加文集字数和阅读数
> 3. 商品页接口增加类型字段 支持价格待定类型

##### Version 1.40.1

> 1. 文集缓存用mc工具重构

##### Version 1.40.0

> 1. 专栏草稿写库方式重构
> 2. 专栏投稿改为强依赖tag
> 3. 分批获取搜索表数据

##### Version 1.39.0

> 1. 专栏阅读数防刷需求

##### Version 1.38.0

> 1. web端接入天马

##### Version 1.37.2

> 1. 修复粉丝数问题

##### Version 1.37.1

> 1. 使用account-service v7

##### Version 1.37.0

> 1. lv2以上的用户开放投稿权限

##### Version 1.36.0

> 1. 增加热点标签功能

##### Version 1.35.0

> 1. 文集列表支持排序

##### Version 1.34.1

> 1. 修复空指针问题

##### Version 1.34.0

> 1. 空间页增加文集入口+文集详情落地页
> 2. 文集新增发布时间 阅读数 字数 简介 文章数信息

##### Version 1.33.2

> 1. 修复更多文章过滤问题

##### Version 1.33.1

> 1. 修复更多文章过滤问题

##### Version 1.33.0

> 1. 文集用缓存框架重构
> 2. 上报收藏infoc日志

##### Version 1.32.0

> 1. 增加up主投稿地址跳转接口

##### Version 1.31.1

> 1. 修复文集过滤问题 

##### Version 1.31.0

> 1. 专栏首页/推荐页面接入天马

##### Version 1.30.4

> 1. 提醒状态增加新手引导类型
> 2. 修复文集排序问题

##### Version 1.30.3

> 1. 修复文集bug 过滤删除文章
> 2. 文集update_time改为文集本身修改时间 
> 3. 修复tag绑定问题

##### Version 1.30.2

> 1. 写死移动端第五个分区

##### Version 1.30.1

> 1. fixed panic bug

##### Version 1.30.0

> 1. 增加文集功能
> 2. 待审视频up没有写专栏权限 

##### Version 1.29.7

> 1. 增加推荐接口pn限制 

##### Version 1.29.6

> 1. 增加推荐接口pn限制 
> 2. 增加收藏和点赞时验证id存在

##### Version 1.29.5

> 1. 支持公告分版本平台显示

##### Version 1.29.4

> 1. ups文章列表 接口不调用infos rpc

##### Version 1.29.3

> 1. ups文章列表 接口不调用infos rpc

##### Version 1.29.2

> 1. fix isAuthor bug

##### Version 1.29.1 

> 1. 异步载入作者列表

##### Version 1.29.0

> 1. 活动开始后自动增加活动tag
> 2. 增加活动时成为创作者API
> 3. 更改查询作者权限方案 不存内存 存redis 解决延时问题
> 4. 支持草稿无标题
> 5. 增加勋章字段

##### Version 1.28.1

> 1. 增加推荐池列表接口

##### Version 1.28.0

> 1. 创作中心逻辑迁移到专栏

##### Version 1.27.1

> 1. 修复更多文章接口返回null的问题

##### Version 1.27.0

> 1. isAuthor接口增加是否封禁字段
> 2. 增加更多文章接口
> 3. home接口增加热门排行
> 4. viewinfo增加稍后再看和小窗播放选项
> 5. 增加用户引导状态功能

##### Version 1.26.7

> 1.增加isAuthor接口

##### Version 1.26.6

> 1.避免分区排序批量回源

##### Version 1.26.5
> 1.修复author bug

##### Version 1.26.4
> 1.修复banner bug

##### Version 1.26.2
> 1.修复票务接口bug

##### Version 1.26.1
> 1.增加notice接口

##### Version 1.26.0
> 1.增加recommends/plus接口

##### Version 1.25.2
> 1.增加origin_image_urls字段

##### Version 1.25.1
> 1.增加动态简介字段

##### Version 1.25.0
> 1.增加番剧/音乐/商品/票务卡片接口

##### Version 1.24.9
> 1.读流量接入点赞服务

##### Version 1.24.8
> 1.修复客户端参数bug

##### Version 1.24.7
> 1.修复banner index 从0开始的问题
> 2.修复作者更多文章没有过滤禁止分发的问题

##### Version 1.24.6
> 1.修复up主缓存不过期的问题
> 2.修复up主列表为空时空缓存不生效的问题
> 3.applyinfo去掉未登录错误
##### Version 1.24.4
> 1.接入点赞服务
##### Version 1.24.1
> 1.封禁用户不能投稿

##### Version 1.24.0
> 1.支持自定义投稿上限    
> 2.viewinfo增加能否分享字段  
> 3.支持关闭通过投稿视频获得的专栏投稿权限 

##### Version 1.23.9
> 1.修复排序时间bug

##### Version 1.23.8
> 1.修复duration计算bug

##### Version 1.23.7
> 1.修复过期排序数据问题

##### Version 1.23.6
> 1.修复同一用户多次投诉问题

##### Version 1.23.5
> 1.修复专栏锁定再次提交的问题  

##### Version 1.23.4
> 1.排序缓存过滤禁止分发的内容 

##### Version 1.23.3
> 1.修复addView接口 && home banner

##### Version 1.23.2
> 1.修复空间接口bug   

##### Version 1.23.1
> 1.修复banner上报字段  

##### Version 1.23.0
> 1.banner增加客户端上报所需的字段  
> 2.增加客户端浏览数据上报http接口    
> 3.增加获取更多文章的rpc接口  
> 4.up主文章列表增加排序字段  
> 5.内网接口增加验证  

##### Version 1.22.5
> 1.fix rand panic问题 

##### Version 1.22.4
> 1.修复article josn序列化问题  

##### Version 1.22.3
> 1. article model 换pb
> 2. log换log context  

##### Version 1.22.2
> 1. 修复收藏失效文档导致数量不对的问题  
> 2. 排行榜换新接口  

##### Version 1.22.1
> 1.修复点赞消息链接 

##### Version 1.22.0
> 1.增加UP主点赞消息通知      
> 2.新增批量获取metas的接口供前端专栏卡片使用  
> 3.recommends和home接口聚合是否点赞的状态  
> 4.viewinfo增加image_urls字段  
> 5.文章详情页Ta的更多文章模块，展示前一篇+后三篇文章      

##### Version 1.21.2
> 1.修复排行榜缓存更新问题 
> 2.增加获取up主专栏列表的内网接口       
##### Version 1.21.1
> 1.跟随基础库升级http client       

##### Version 1.21.0
> 1.排行榜增加score和note字段      
> 2.增加获取最新投稿数量的rpc接口   
> 3.增加批量获取是否点赞的rpc接口  

##### Version 1.20.5
> 1.返回先发后审的文章        

##### Version 1.20.4
> 1.修复推荐角标问题      

##### Version 1.20.3
> 1.增加人工智能部门需要的点击日志上报      

##### Version 1.20.2
> 1.修复排序缓存bug      

##### Version 1.20.1
> 1.修复收藏bug     

##### Version 1.20.0
> 1.增加web需要的收藏接口    
> 2.支持首页降级  
> 3.增加阅读数排序 
> 4.增加排行榜功能  
> 5.banner从读库改成rpc接口  
> 6.排序缓存改为最近3周数据 最新投稿全量缓存

##### Version 1.19.4
> 1.fix 推荐页面 panic bug   

##### Version 1.19.3
> 1.视频up主自动获得发专栏权限   

##### Version 1.19.2
> 1.回源时禁止分发不出现在作者的缓存中里    

##### Version 1.19.1
> 1.更新时禁止分发不出现在作者的缓存中里    

##### Version 1.19.0
> 1.新增是否是作者的RPC接口    
> 2.Recommends RPC接口增加cid字段    
> 3.专栏管理增加按照阅读数、收藏数、投币数排序的接口逻辑  
> 4.viewinfo增加is_author字段
> 5.recommends接口增加rec_text字段     
> 6.改变推荐分区为推荐池 每次随机取 

##### Version 1.18.4
> 1.修复创作中心草稿分页错误

##### Version 1.18.3
> 1.修复作者表不更新的问题

##### Version 1.18.2
> 1.fix UserGet脏数据导致的panic问题

##### Version 1.18.1
> 1.fix 空指针问题

##### Version 1.18.0
> 1.新增一二级分区排序功能
> 2.新增首页banner位
> 3.viewinfo接口增加作者昵称   

##### Version 1.17.4
> 1.修复viewinfo panic问题   

##### Version 1.17.3
> 1.fix创作中心计数问题   
> 2.完善保存文章和草稿时增加日志  

##### Version 1.17.2
> 1.保存文章和草稿时增加日志，方便查问题   

##### Version 1.17.1
> 1.修复更新推荐bug   

##### Version 1.17
> 1.新增更新作者缓存RPC接口  

##### Version 1.16.4
> 1.文章新增缓存增加日志  
> 2.修复重启程序关闭顺序问题  

##### Version 1.16.3
> 1."更早文章"接口不显示禁止分发文章  
> 2.推荐接口增加from字段  

##### Version 1.16.2
> 1.fix AddCache panic

##### Version 1.16.1
> 1.推荐接口支持按aid分页获取    

##### Version 1.16.0
> 1.增加分区禁止属性   
> 2.推荐位改为逆序排列      

##### Version 1.15.0
> 1.添加"更早文章"接口  
> 2.增加coin计数
> 3.meta中增加分类列表
> 4.喜欢功能做成双向表  

##### Version 1.14.0
> 1.投诉计数  
> 2.增加是否展示推荐新投稿配置项     

##### Version 1.13.1
> 1.修复推荐不显示最新投稿问题    

##### Version 1.13.0
> 1.新增收藏RPC接口  
> 2.修复数据上报问题  
> 3.修复推荐角标没有取消的问题   

##### Version 1.12.0
> 1.修复草稿列表时间问题  

##### Version 1.11.0
> 1.接入大仓库  

##### Version 1.10.0
> 1.添加申请成为专栏作者的入口  

##### Version 1.9.2
> 1.添加单个获取文章的HTTP接口

##### Version 1.9.1
> 1.修改viewinfo增加title字段

##### Version 1.9.0
> 1.添加批量获取文章的HTTP接口

##### Version 1.8.8
> 1.优化查询tag

##### Version 1.8.7
> 1.移除addView接口合并到viewinfo接口中

##### Version 1.8.6
> 1.修复创作中心获取文章内容调用表不对的问题  

##### Version 1.8.5
> 1.修改文章列表不获取tag  

##### Version 1.8.4
> 1.调整添加、更新、删除文章缓存  
> 2.修复作者列表空缓存问题   
> 3.增加推荐标志

##### Version 1.8.3
> 1.修复创作中心文章撤回后文章内容丢失

##### Version 1.8.2
> 1.http接口返回[]

##### Version 1.8.1
> 1.优化草稿列表，加读缓存逻辑
> 2.ptimes多返回值重构 
> 3.收藏http接口返回[]

##### Version 1.8.0
> 1.增加投诉API  
> 2.过滤落库存表
> 3.其他细微优化

##### Version 1.7.2
> 1.使用rpc token

##### Version 1.7.1
> 1.草稿分表

##### Version 1.7.0
> 1.长文本过滤从job转移过来   
> 2.正文增加是否过滤的标志位    
> 3.增加每日发布文章次数限制      

##### Version 1.6.0
> 1.修复stat解析bug   
> 2.作者信息增加认证和挂件字段   

##### Version 1.5.0
> 1.stat 单个查询优化  
> 2.去掉推荐读写锁改为chan  
> 3.点赞缓存优化mc改为redis hash  
> 4.切换成新版过滤rpc  

##### Version 1.4.0

> 1.添加草稿异步写db
> 2.文章过滤放在job里做

##### Version 1.3.1

> 1.修复stats信息被置0的bug  

##### Version 1.3.0

> 1.长文本过滤优化  
> 2.修复稿件分区改变最新文章缓存未清理的bug  
> 3.修复pn为0导致panic的bug  

##### Version 1.2.1

> 1.稿件列表接口添加根据状态过滤

##### Version 1.2.0

> 1.增加批量获取稿件信息接口 使用过滤rpc

##### Version 1.1.1

> 1.修复过滤bug

##### Version 1.0.0

> 1.专栏项目初始化  
