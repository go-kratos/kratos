#### tag 标签系统   

#### Version 4.12.15
> 1.去除/ranking/hots、/ranking/bangumi、/ranking/region 接口

#### Version 4.12.14
> 1.规避tag go rpc count map 并发读写

#### Version 4.12.13
> 1.tag下线无流量API
> 2.tag去除memcahce依赖

#### Version 4.12.12
> 1.修复tag去除特殊字符和空格顺序问题.

#### Version 4.12.11
> 1.upbind增加只修改默认隐藏tag，不修改显示tag.
> 2.迁移cache.Cache 到fanout

#### Version 4.12.10
> 1.迁移account-service的gorpc至grpc api.

#### Version 4.12.9
> 1.支持grpc api.
> 2.增加grpc client.go.

#### Version 4.12.8
> 1.删除废弃接口/article/set_keywords的调用.

#### Version 4.12.7
> 1.增加频道版本展示露出HTTP接口.

#### Version 4.12.6
> 1.修复国际版频道ChannelDetail nil情况.

#### Version 4.12.5
> 1.增加国际版标识.
> 2.频道国际版数据展示过滤.

#### Version 4.12.4
> 1.AI接口增加移动端详情页和频道H5页面来源区分.

#### Version 4.12.3
> 1.增加相似频道功能
> 2.增加单tag 数据查询
> 3.迁移频道回查接口到grpc api.

#### Version 4.12.2
> 1.优化x/internal/tag/archive/multi/tags接口
> 2.切/tag/archive/tags接口到tags grpc api.

#### Version 4.12.1
> 1.修复tag id&name infos 函数 请求count读放大.

#### Version 4.12.0
> 1.优化闭包数据.
> 2. like& hate 

#### Version 4.11.13
> 1.迁移频道至interface层，使用tag-service grpc api.

#### Version 4.11.12
> 1.迁移视频详情页举报到grpc.

#### Version 4.11.11
> 1.workflow-admin平台化，减少业务依赖.

#### Version 4.11.10
> 1.迁移视频详情页举报到gorpc.

#### Version 4.11.9
> 1.update count cache.

#### Version 4.11.8
> 1.迁移点赞、点踩、举报、日志举报功能的spam功能到interface层.

#### Version 4.11.7
> 1.举报、点赞、点踩合并card & block 接口，使用card接口.

#### Version 4.11.6
> 1.使用grpc，迁移举报新增接口到grpc.

#### Version 4.11.5
> 1.去除conf.Conf.Tag.ChanneLimit

#### Version 4.11.4
> 1.频道广场页需求，指定返回频道数，和每个频道下的稿件数.

#### Version 4.11.3
> 1.频道管理后台自定义排序，移动端频道数据重新设定排序规则.

#### Version 4.11.2
> 1.频道接入workflow举报.

#### Version 4.11.1
> 1.增加RPC.TagTop接口，聚合tag info & similar tag.

#### Version 4.11.0
> 1.删除普通用户绑定专栏tag接口

#### Version 4.10.27
> 1.fix admin bind default bugs.

#### Version 4.10.26
> 1.删除c.RemoteIP()
> 2.use DefaultServer & Engine.Start.
> 3.use RPC.NewServer.

#### Version 4.10.25
> 1.频道回查命中多规则，以及显示命中tag列表.

#### Version 4.10.24
> 1.AI新增频道报错提示.

#### Version 4.10.23
> 1.迁移封禁和小黑屋接口到新接口

#### Version 4.10.22
> 1.使用新的鉴权auth和verify.

#### Version 4.10.21
> 1.增加up主默认绑定一二级分区和管理员审核绑定一二级分区.

#### Version 4.10.20
> 1.增加举报内容有效分以及举报次数计数.

#### Version 4.10.19
> 1.修复活动tag增加报错，tag已存在、tag不存在.

#### Version 4.10.18
> 1.频道回查增加频道名称等信息.

#### Version 4.10.17
> 1.优化serice异步逻辑.

#### Version 4.10.16
> 1.使用公共配置.

#### Version 4.10.15
> 1.删除dao注释代码.   

#### Version 4.10.15
> 1.修复bind接口split数组.

#### Version 4.10.14
> 1.AI接口增加Build参数上报.

#### Version 4.10.13
> 1.迁移BM框架.

#### Version 4.10.12
> 1.频道分类自定义排序.

#### Version 4.10.11
> 1.限制频道举报、up主删除操作.

#### Version 4.10.10
> 1.避开ai推荐接口error
> 2.增加频道、灾备标记

#### Version 4.10.9
> 1.去除频道缓存依赖，使用内存依赖

#### Version 4.10.8
> 1.channel mc .

#### Version 4.10.7
> 1.fix mc eof.

#### Version 4.10.6
> 1.fix ai推荐 缓存失效问题.

#### Version 4.10.5
> 1.name 过滤不可见字符.

#### Version 4.10.4
> 1.修改部分频道缓存设计和回源逻辑.
> 2.增加单飞模式.

#### Version 4.10.3
> 1.去掉bilibili_tag MySQL 库.     

#### Version 4.10.2
> 1.视频详情页下增加频道数据.
> 2.增加请求AI接口超时设置.

#### Version 4.10.1
> 1. 限制非管理员、up主用户删除和添加tag.
> 2. 增加频道tag不允许点赞点踩.
> 3. 增加频道数据加入memcache数据.

#### Version 4.10.0
> 1. 增加频道业务 

#### Version 4.9.46
> 1.rank_result          

#### Version 4.9.45
> 1.ranking  

#### Version 4.9.45
> 1.tag fix bazel build

#### Version 4.9.44
> 1.tag group 

#### Version 4.9.44
> 1.WhiteUser & limitResource  

#### Version 4.9.44
> 1.停止tag下视频轮训补数据          

#### Version 4.9.43
> 1.修复admin bind tag     

#### Version 4.9.42
> 1.删除无效的迁移archive tag         

#### Version 4.9.41
> 1.add report partID        

#### Version 4.9.38
> 1.archive UpArcBind tagListUser    

#### Version 4.9.37
> 1.archive rpc bind tag        

#### Version 4.9.36
> 1.platform admin 迁移           

#### Version 4.9.35
> 1.专栏&音乐 绑定tag 迁移      

#### Version 4.9.34
> 1.专栏 绑定tag 迁移         

#### Version 4.9.33
> 1.user bind & up create tag         
    
#### Version 4.9.32
> 1. 删掉tag_action 逻辑代码         

#### Version 4.9.31
> 1. 删除无效写archive tag 逻辑代码,合并批量创建tag    

#### Version 4.9.30
> 1. 删除无效写tag 逻辑代码,合并批量创建tag        

#### Version 4.9.29
> 1. 迁移tag 写新库            

#### Version 4.9.26
> 1. topic tag fix       

#### Version 4.9.25
> 1. 迁移tag 写状态         

#### Version 4.9.24
> 1. 删除无效的mtag info逻辑代码  

#### Version 4.9.23
> 1. 删除无效的mtag info逻辑代码  

#### Version 4.9.22
> 1. 删除无效的tag info逻辑代码  

#### Version 4.9.21
> 1. 删除无效的action tag 逻辑代码  

#### Version 4.9.20
> 1. 删除无效的sub tag 逻辑代码  

#### Version 4.9.19
> 1. 修复aciton map逻辑       

#### Version 4.9.18
> 1. 增加aciton map      

#### Version 4.9.17
> 1. archive like&hate tag state      

#### Version 4.9.16
> 1. like & hate & sub & cancel    

#### Version 4.9.15
> 1.迁移点踩点中的resource         

#### Version 4.9.14
> 1. 点踩点赞逻辑保留新版本逻辑,保留老版本接口url&param           

#### Version 4.9.13
> 1. add admin add&del& lock 接口 

#### Version 4.9.12
> 1.change account http to rpc
 
#### Version 4.9.11
> 1.fix tag batch info return data.

#### Version 4.9.10
> 1.fix tag infos  
> 2.fix res infos not tags         

#### Version 4.9.9
> 1.read subtag all in rpc          

#### Version 4.9.8
> 1.fix admin bind tag 16001         

#### Version 4.9.7
> 1.使用account-service v7
 
#### Version 4.9.6 
> 1.读tag info by id    
> 2.fix info     

#### Version 4.9.5 
> 1.迁移到main 目录   

#### Version 4.9.4  
> 1.将dao archive rpc 迁移到service           

#### Version 4.9.3  
> 1.删除statsd     

#### Version 4.9.2   
> 1.删除无用的http 接口。             

#### Version 4.9.1   
> 1.迁移部分archive rpc 接口。     

#### Version 4.9.0   
> 1.第三方http接口脱离tag interface 迁移至service             

#### Version 4.8.19
> 1.修复自定义更新排序tag参数tids为空时判断     

#### Version 4.8.18   
> 1.fix arc_log user role     

#### Version 4.8.17   
> 1.archive tag log      
 
#### Version 4.8.16   
> 1.增加tag自定义排序功能

#### Version 4.8.15   
> 1.增加话题批量创建tag接口  
> 2.修改tag名字长度限制
> 3.增加tag批量查询接口(json传参)

#### Version 4.8.14   
> 1.fix closure rate     

#### Version 4.8.13   
> 1.fix account card to info        

#### Version 4.8.12   
> 1.fix sub tag sort     

#### Version 4.8.11   
> 1.迁移流量查询多条tag         

#### Version 4.8.10   
> 1.迁移流量策略调整         

#### Version 4.8.9   
> 1.迁移流量taginfo service降级容错    

#### Version 4.8.7   
> 1.迁移流量tag计数     

#### Version 4.8.6   
> 1.迁移流量专栏页下tag     

#### Version 4.8.6   
> 1.迁移流量视频详情页下tag     

#### Version 4.8.5   
> 1.迁移流量tagName 

#### Version 4.8.4   
> 1.重新打包     

#### Version 4.8.3   
> 1.迁移流量用户订阅tag     

#### Version 4.8.2   
> 1.迁移流量可配置,去掉report表查询事务       

#### Version 4.8.1    
> 1.修复sql带空格   

#### Version 4.8.0    
> 1.tag信息迁移千分之一流量  

#### Version 4.7.12
> 1.修复分区页和tag详情页下的分页数据丢失(ZREM cache aid & update state=3)

#### Version 4.7.11
> 1.fix nologin code  

#### Version 4.7.10
> 1.增加mysql日志   
> 2.修改context   
> 3.修改filter接口  

#### Version 4.7.9
> 1.增加up查询检测tag接口    

#### Version 4.7.8
> 1.fix slice bounds out of range    

#### Version 4.7.7
> 1.最新视频补数据增加redisErrNil 

#### Version 4.7.6
> 1.最新视频补数据 

#### Version 4.7.5
> 1.redis Receive 计数  

#### Version 4.7.4
> 1.修复sub nil空值&Slice   

#### Version 4.7.3
> 1.修复热门分区下tag最新视频
> 2.迁移到archive3 rpc接口    

#### Version 4.7.2
> 1.修复stmt初始化   

#### Version 4.7.1
> 1.rename SQL名称   

#### Version 4.7.0
> 1.合并dao层优化数据库链接    

#### Version 4.6.2
> 1.去掉绑定错误-404    

#### Version 4.6.1
> 1.修复bug-十九大     

#### Version 4.6.0
> 1.迁移至大仓库    
> 2.增加十九大屏蔽    

#### Version 4.5.2
> 1.增加视频tag实名认证

#### Version 4.4.2
> 1.增加dynamic-servicey依赖接口

#### Version 4.3.2
> 1.增加log agent

#### Version 4.2.2
> 1.修改创作中心和后台调用返回-404

#### Version 4.2.1
> 1.修改同义词tag数据格式

#### Version 4.2.0
> 1.提供同义词tag接口给大数据

#### Version 4.1.2
> 1.优化HGETALL

#### Version 4.1.1
> 1.修复已经关注bug

#### Version 4.1.0
> 1.关注换为redis hash缓存

#### Version 4.0.0
> 1.修改tag返回错误码

#### Version 3.0.3
> 1.修改tag返回错误码

#### Version 3.0.2
> 1.增加批量获取aid下tag列表

#### Version 3.0.1
> 1.修改redis cache log  
> 2.修改router identify

#### Version 3.0.0
> 1.通用tag资源关系
> 2.接入普罗米修斯

#### Version 2.9.3
> 1.回写cache把占位0cache删掉

#### Version 2.9.2
> 1.回写cache把占位-1cache删掉

#### Version 2.9.1
> 1.修复up主tag不能删bug

#### Version 2.9.0

> 1.优化过审同步tag，拆分成2个
> 2.音乐平台绑定tag

#### Version 2.8.8
> 1.增加6.4政策限制

#### Version 2.8.7
> 1.修改一审忽略操作

#### Version 2.8.6
> 1.修改获取我关注tag，is_atten=1

#### Version 2.8.5
> 1.修改批量获取全量aid

#### Version 2.8.4
> 1.优化过审同步tag，拆分成2个
> 1.修复异步cache panic

#### Version 2.8.3
> 1.优化tag分区排序缓存

#### Version 2.8.2
> 1.tag详情页分区筛选，投稿推荐tag

#### Version 2.8.1
> 1.使用videoup更新稿件tag
> 2.修复一个小bug

#### Version 2.8.0
> 1.升级vendor支持csrf

#### Version 2.7.8
> 1.修复过审状态不正常  
> 2.接入最新的配置中心  
> 3.升级vendor

#### Version 2.7.7
> 1.优化tag一二审

#### Version 2.7.6
> 1.修复detail，similar tag状态不正常返回error

#### Version 2.7.5
> 1.修改order by ctime

#### Version 2.7.4
> 1.修改rank redis 稿件缓存排序为pubtime

#### Version 2.7.4
> 1.举报1，2审优化  
> 2.降级相关tag

#### Version 2.7.3
> 1.添加model.TagsArcAdd日志
> 2.后台过审不验证稿件状态

#### Version 2.7.2
> 1.修复判读订阅cache错误

#### Version 2.7.1
> 1.优化订阅zscore数量  
> 2.降级大数据接口


#### Version 2.7.0
> 1.tag部分接口支持rpc  
> 2.优化查询视频tags缓存  
> 3.优化查询订阅tags缓存  
> 4.增加internal内部路由,删除degrade  
> 5.升级vendor  
> 6.detail接口支持分页逻辑  
> 7.增加相似tag换一换接口  
> 8.使用filter-service过滤

#### Version 2.6.7
> 1.升级vendor

#### Version 2.6.5
> 1.小黑屋  
> 2.订阅,视频下tag,cache问题回源db  
> 3.升级go-business  
> 4.优化log  
> 5.删除后台删除单条tag功能  
> 6.活动tag

#### Version 2.6.4
> 1.fix trace owner

#### Version 2.6.3
> 1.修复管理员tag操作不显示

#### Version 2.6.2
> 1.增加旧顶踩接口,兼容app

#### Version 2.6.1
> 1.cache timeout 回源到db  
> 2.修改cache增加expire操作

#### Version 2.6.0

> 1.修复map并发读写  
> 2.升级vendor  
> 3.顶踩入库,增加可取消顶踩  
> 4.增加活动tag,详情页不可增删  
> 5.视频tag操作页面的log显示可选  
> 6.调用大数据接口增加熔断  


#### Version 2.5.8

> 1.优化tag绑定视频sql

#### Version 2.5.7

> 1.接入配置中心

#### Version 2.5.6

> 1.tag修改记录不展示管理员和mid=0的操作

#### Version 2.5.5

> 1.rpc调用参数增加RealIP  
> 2.过审后的tag,up主也可删

#### Version 2.5.4

> 1.修复举报插入逻辑  
> 2.修改视频下tag排序规则,like > tagType > arcTagRole > hate > ctime  
> 3.增加举报表父子节点关系  
> 4.升级vendor

#### Version 2.5.3

> 1.增加tag添加删除错误traceon  

#### Version 2.5.2

> 1.修改删除视频tag逻辑,up主tag不能删,up主删自己tag是算入计数  
> 2.异步更新热门tag,投稿时间排序数据,增加time sleep  
> 3.changeApi 增加 setkeywords

#### Version 2.5.1

> 1.修复tagarc cacheProc 稿件为nil panic  

#### Version 2.5.0

> 1.增加视频tag操作记录举报  
> 2.每次视频,每天,每人,只能删2个  
> 3.升级基础库  
> 4.增加稿件分区变更,更新热门tag cache接口

#### Version 2.4.1

> 1.更新indetify，支持直接访问mc和帐号  

#### Version 2.4.0

> 1.增加删除tag监控上报  

#### Version 2.3.0

> 1.升级vendor  
> 2.healthcheck方法修改  

#### Version 2.2.3
> 1.update vendor,升级  
> 2.去除tag服务changetag反向更新稿件表逻辑  

#### Version 2.2.2
> 1.稿件后台审核调用tag修改接口，修复并发，和日志错乱问题  

#### Version 2.2.1
> 1.顶踩视频tag接口，保留原有接口，并分拆成顶和踩   

#### Version 2.2.0
> 1.修复change接口,新增tag tid=0问题  
> 2.升级vendor:go-business,go-common,golang

#### Version 2.1.2
> 1.up主可自由增删tag不存在等级和稿件up主不能增加tag限制  
> 2.add localrouter ping  
> 3.tag like接口做兼容，新增hate接口  

#### Version 2.1.1

> 1.添加ecode   
> 2.稿件与tag相关操作稿件需要正常状态

#### Version 2.1.0

> 1.修复热门tag不同分区相同tag订阅不同步  
> 2.修复订阅,视频和tag相关cache操作先expire   
> 3.增加detail tag api  
> 4.去除新增tag html转义  

#### Version 2.0.0

> 1.代码重构及表结构分表  
> 2.点赞踩等功能  
> 3.vendor  

#### Version 1.7.5

> 1.fix ridall未初始化的bug  

##### Version 1.7.4

> 1.增加escape转码  

##### Version 1.7.3

> 1.新增支持tag remove事件  
> 2.订阅tags返回订阅总数  
> 3.过滤0宽度的空格字符<200b>  

##### Version 1.7.2

> 1.获取热门tag是否订阅  

##### Version 1.7.1

> 1.修复count不一致bug  
> 2.优化load rank cache逻辑  

##### Version 1.7.0

> 1.新增tag下最新视频  
> 2.新增获取用户订阅的tag下的最新视频  
> 3.新增tag订阅推kafka  

##### Version 1.6.1

> 1.订阅，取消订阅后需清理tag cache    

##### Version 1.6.0

> 1.管理后台相关接口  
> 2.视频tag列表返回是否关注  
> 3.稿件rpc接口settag2  
> 4.修复mc的key中有空格问题，替换为特殊字符  

##### Version 1.5.1

> 1.相关tag接口修复  

##### Version 1.5.0

> 1.添加取消订阅tag刷新动态  
> 2.过滤新增乱码tag  
> 3.更改稿件tag管理change接口去掉验证mid  

##### Version 1.4.1

> 1.排行榜ranking_bangumi接口增加title字段  

##### Version 1.4.0

> 1.排行榜热门tag、热门番剧  
> 2.相似tag，用于移动端tag推荐  
> 3.订阅支持批量  
> 4.增加更新archive tag字段  

##### Version 1.3.0

> 1.支持trace v2  
> 2.[bug]change方法导致出现tagid为0的错误数据  

##### Version 1.2.1

> 1.视频关联tag只返回状态正常的tag  
> 2.up主可以给自己锁定的视频增加tag  
> 3.增加过滤新增tag的前后空格  

##### Version 1.2.0

> 1.新增修改稿件关联tag内部接口  
> 2.兼容最新go-common  

##### Version 1.1.1

> 1.fix memcache 空格bug  

##### Version 1.1.0

> 1.新增热门tag下的最新、最热视频接口  

##### Version 1.0.0

> 1.tag基础功能
