#### archive rpc service

##### Version 6.47.6
> 1.iPhone 5.36版本不吐拜年祭单品稿件的秒开地址

##### Version 6.47.5
> 1.拜年祭单品视频不吐秒开地址

##### Version 6.47.4
> 1.view接口只对archive做强判断

##### Version 6.47.3
> 1.fix view

##### Version 6.47.2
> 1.迁移gorpc方法到gpc

##### Version 6.47.1
> 1.view接口增加staff信息

##### Version 6.47.0
> 1.grpc增加注释

##### Version 6.46.6
> 1.更新稿件缓存增加联合投稿部分

##### Version 6.46.5
> 1.Dislike强制0

##### Version 6.46.4
> 1.增加高能看点、bgm、联合投稿attribute

##### Version 6.46.3
> 1.删除RPC中的addShare方法

##### Version 6.46.2
> 1.分享不发databus消息

##### Version 6.46.1
> 1.fix package 

##### Version 6.46.0
> 1.issues #403 大仓库项目目录结构改进

##### Version 6.45.2
> 1.加参数控制pgc吐playurl

##### Version 6.45.1
> 1.pgc不吐playurl

##### Version 6.45.0
> 1.拦截archive miss不存在的稿件

##### Version 6.44.4
> 1.接入account grpc

##### Version 6.44.3
> 1.秒开拦截辣鸡参数

##### Version 6.44.2
> 1.秒开qn白名单+大会员清晰度降级

##### Version 6.44.1
> 1.计数默认返回值修改

##### Version 6.44.0
> 1.增加UGCPay标识

##### Version 6.43.9
> 1.dash格式加codecid

##### Version 6.43.8
> 1.修复冷门up主投稿列表稿件不全的bug

##### Version 6.43.7
> 1.初始化缓存日志

##### Version 6.43.6
> 1.初始化分区缓存的时候使用context.Background

##### Version 6.43.5
> 1.修改重置up信息的缓存逻辑

##### Version 6.43.4
> 1.增加日志观察job异步databus消息是否发送成功

##### Version 6.43.3
> 1.增加日志观察账号RPC服务的返回是否正常

##### Version 6.43.2
> 1.增加日志观察账号昵称头像是否为空

##### Version 6.43.1
> 1.秒开接口空字段不吐

##### Version 6.43.0
> 1.秒开接口增加dash字段

##### Version 6.42.2
> 1.同步IsNormal,AttrVal方法

##### Version 6.42.1
> 1.接入grpc

##### Version 6.42.0
> 1.添加dao层ut

##### Version 6.42.0
> 1.调整目录

##### Version 6.41.1
> 1.fix cache ctx

##### Version 6.41.0
> 1.issue 249 metadata ip

##### Version 6.40.10
> 1.增加视频云的fnver，fnval字段返回

##### Version 6.40.9
> 1.更新bvc pb文件

##### Version 6.40.8
> 1.优化重新生成账号缓存的逻辑

##### Version 6.40.7
> 1.账号接口请求失败时，走databus慢慢更新

##### Version 6.40.6
> 1.透传视频云fnval,fnver字段

##### Version 6.40.5
> 1.账号老是不知道刷什么东西

##### Version 6.40.4
> 1.增加是否可以投屏

##### Version 6.40.3
> 1.透传投屏信息

##### Version 6.40.2
> 1.支持地区限制

##### Version 6.40.1
> 1.稿件描述走缓存

##### Version 6.40.0
> 1.稿件更新缓存bugfix  
> 2.重置retag

##### Version 6.39.3
> 1.增加autoplay字段

##### Version 6.39.2
> 1.HTTP接口增加稿件分辨率字段

##### Version 6.39.1
> 1.PB接口体增加json字段输出

##### Version 6.39.0
> 1.批量MC接口优化

##### Version 6.38.1
> 1.分辨率0,0,0不做处理

##### Version 6.38.0
> 1.稿件增加分辨率字段

##### Version 6.37.19
> 1.批量稿件接口代码优化

##### Version 6.37.18
> 1.conn close fix

##### Version 6.37.17
> 1.增加缓存容错

##### Version 6.37.16
> 1.bvc灰度

##### Version 6.37.15
> 1.redis set expire -> setex

##### Version 6.37.14
> 1.无投稿的用户只缓存10分钟  

##### Version 6.37.13
> 1.分享行为的databus key从aid改为mid

##### Version 6.37.12
> 1.fix row close

##### Version 6.37.11
> 1.清理share代码

##### Version 6.37.10
> 1.删除videoshot add接口

##### Version 6.37.9
> 1.RPC不需要token了  By 郝冠伟确认

##### Version 6.37.8
> 1.增加register  

##### Version 6.37.7
> 1.删除冗余代码  

##### Version 6.37.6
> 1.使用bm  

##### Version 6.37.5
> 1.删除多余配置  

##### Version 6.37.4
> 1.增加批量获取up投稿数量的http接口  

##### Version 6.37.3
> 1.删除limit模块  

##### Version 6.37.2
> 1.第一次分享改发databus  

##### Version 6.37.1
> 1.使用account-service v7  

##### Version 6.37.0
> 1.迁移到主站目录下  

##### Version 6.36.16
> 1.取消强制开关  

##### Version 6.36.15
> 1.提供给B+的秒开接口强行不返回playurl  

##### Version 6.36.14
> 1.参数长度调整为200    

##### Version 6.36.13
> 1.补充UnitTest  

##### Version 6.36.12
> 1.配置文件增加开关选项，控制是否请求视频云获取播放信息  

##### Version 6.36.11
> 1.firstCid只用vupload，外部源不缓存   

##### Version 6.36.10
> 1.优化秒开代码  

##### Version 6.36.9
> 1.Archive3结构体增加第一P的cid，供后续业务扩展使用  

##### Version 6.36.8
> 1.attr增加地区限制  

##### Version 6.36.7
> 1.share数双写新databus  

##### Version 6.36.6
> 1.接bvc的pb接口  

##### Version 6.36.5
> 1.Convey test  

##### Version 6.36.4
> 1.BFS改回来    

##### Version 6.36.3
> 1.BFS的封面图强制返回https  

##### Version 6.36.2
> 1.删除attrbithideclick相关代码

##### Version 6.36.1
> 1.删除like的相关代码与配置  

##### Version 6.36.0
> 1.删除like的相关代码与配置  

##### Version 6.35.1
> 1.修改透传的player字段名  

##### Version 6.35.0
> 1.attr的第十二位改成 IsPorder 私单标记

##### Version 6.34.0
> 1.增加player接口  

##### Version 6.33.1
> 1.批量接口增加参数日志  

##### Version 6.33.0
> 1.增加maxAID的接口  

##### Version 6.32.9
> 1.删除废弃的代码  

##### Version 6.32.8
> 1.删除废弃的RPC server端代码  

##### Version 6.32.7
> 1.删除冗余代码  

##### Version 6.32.6
> 1.增加prom db  

##### Version 6.32.5
> 1.内置prom  

##### Version 6.32.4
> 1.使用内置prom  

##### Version 6.32.3
> 1.修复缓存miss时少吐数据的问题  

##### Version 6.32.1
> 1.Video3 RPC  

##### Version 6.32.0
> 1.兼容客户端传多次点赞  

##### Version 6.31.3
> 1.统一修改errgroup包路径  

##### Version 6.31.2
> 1.attr的第九位改成 isPGC  

##### Version 6.31.1
> 1.修改Views3返回值  

##### Version 6.31.0
> 1.删除非internal的对外http接口  

##### Version 6.30.0
> 1.Archive3结构体改为非指针  

##### Version 6.29.1
> 1.补全RPC PB接口，video3  

##### Version 6.29.0
> 1.补全RPC PB接口  

##### Version 6.28.1
> 1.archive增加dynamic字段  

##### Version 6.28.0
> 1.增加up主推荐视频的RPC接口  

##### Version 6.27.0
> 1.删除pgc相关逻辑  

##### Version 6.26.0
> 1.delete Movie2 AidByCid  

##### Version 6.25.0
> 1.add Page3 pb rpc  

##### Version 6.24.2
> 1.upArcs & upsArcs pb  

##### Version 6.24.2
> 1.rpc删除likes2接口  

##### Version 6.24.1
> 1.rpc增加likes3的pb接口  

##### Version 6.24.0
> 1.rpc增加stat，stats的pb接口  

##### Version 6.23.0
> 1.rpc 增加archive3,archives3的pb接口  
> 2.rpc 删除废弃的videos2,videosByCids2,CidByEpIDs2等方法  
> 3.pgc接口只吐电影信息  

##### Version 6.22.8
> 1.rpc 增加view3的pb接口  

##### Version 6.22.7
> 1.http video接口走PB  

##### Version 6.22.6
> 1.http archive、archives、page全量开放，异步更新page缓存  

##### Version 6.22.5
> 1.http/view接口全量pb缓存预热  

##### Version 6.22.4
> 1.流量扩大得到aid%10<5走PB 
> 2.分P的http接口也走pb  
 
##### Version 6.22.3
> 1.aid%10<3走pb  

##### Version 6.22.2
> 1.http stat/stats接口全量走pb  

##### Version 6.22.1
> 1.pb的func/model/struct/service等全面改名为数字3结尾  

##### Version 6.22.0
> 1.pb bugfix  

##### Version 6.21.1
> 1.archive http 接口 aid%10=1的走pb  

##### Version 6.21.0
> 1.增加limiter限流  

##### Version 6.20.0
> 1.增加批量views接口，aids限制为20个  

##### Version 6.19.0
> 1.直播限制50个  

##### Version 6.18.0
> 1.cids接口不直接return  

##### Version 6.17.0
> 1.videoshot接口试水pkg/errors  

##### Version 6.16.0
> 1.增加全区7天内最新稿件  

##### Version 6.15.0
> 1.redis errnil return  

##### Version 6.14.0
> 1.likes相关数据落库  
> 2.增加likes列表的RPC接口  

##### Version 6.13.0
> 1.bilibili_archive库全都读写分离  

##### Version 6.12.0
> 1.upspass score bugfix  

##### Version 6.11.0
> 1.upsPass接口增加copyright  

##### Version 6.10.0
> 1.添加获取单P信息的http接口(包含description字段)  
> 2.修改原获取单P信息的service层逻辑  
> 2.添加获取长简介的http和rpc接口  

##### Version 6.9.0
> 1.升级go-common

##### Version 6.8.0
> 1.升级go-common  
> 2.迁移model到项目中  

##### Version 6.7.0
> 1.增加主站排行榜专用接口  

##### Version 6.6.1
> 1.memcache json  

##### Version 6.6.0
> 1.memcache gob  

##### Version 6.5.0
> 1.增加点赞相关RPC接口  

##### Version 6.4.0
> 1.升级go-common&go-business  
> 2.videoshot rpc 增加aid参数  

##### Version 6.3.0
> 1.videosho接口增加aid参数  

##### Version 6.2.4
> 1.http context fix  

##### Version 6.2.3
> 1.升级go-business  
> 2.manager后台变更稿件归属mid时,变更相应缓存  

##### Version 6.2.2
> 1.增加upspassed rpc方法  

##### Version 6.2.1
> 1.rpc video2 nil fix  

##### Version 6.2.0
> 1.所有稿件&视频走新archive_result数据库  
> 2.升级go-common&go-business  

##### Version 6.1.21
> 1.删除SetStatCache2接口  

##### Version 6.1.20
> 1.修复可能导致panic的问题  

##### Version 6.1.19
> 1.增加http的typelist接口  
>>>>>>> develop

##### Version 6.1.18
> 1.修改ci配置  

##### Version 6.1.17
> 1.增加account清楚缓存时的参数  

##### Version 6.1.16
> 1.增加昵称&头像更新后的缓存清理逻辑  

##### Version 6.1.15
> 1.增加无脑生成view&click缓存

##### Version 6.1.14
> 1.ci配置分支

##### Version 6.1.13
> 1.删掉zlimit相关残留代码  
> 2.archive和archives接口返回archive_report_result中is_show等于1的result  

##### Version 6.1.12
> 1.去掉老的dede  

##### Version 6.1.11
> 1.分区表走新分区  
> 2.增加RPC获取所有type的方法   
> 3.升级go-common和go-business  

##### Version 6.1.10
> 1.升级go-common和go-business  
> 2.修改prom写法  

##### Version 6.1.9
> 1.修复闭包缓存  

##### Version 6.1.8
> 1.重发ci  

##### Version 6.1.7
> 1.修复分类缓存  

##### Version 6.1.6
> 1.修复videoshot nil 导致panic  

##### Version 6.1.5
> 1.修复videoshot nil 导致panic  

##### Version 6.1.4
> 1.page字段走自增形式  

##### Version 6.1.3
> 1.增加auth  

##### Version 6.1.2
> 1.增加prom  

##### Version 6.1.1
> 1.批量大小改为60  

##### Version 6.1.0
> 1.计数闭包问题修复  

##### Version 6.0.16
> 1.rows close bug fix  

##### Version 6.0.13
> 1.mc改成永不过期  

##### Version 6.0.12
> 1.计数加aid在json  

##### Version 6.0.11
> 1.增加cache出错不回写  
> 2.增加prom回源统计  

##### Version 6.0.10
> 1.修复点击计数panic  

##### Version 6.0.9
> 1.修复prom参数个数  

##### Version 6.0.8
> 1.增加prom包  

##### Version 6.0.7
> 1.增加memcache随机过期时间  

##### Version 6.0.6
> 1.增加upcount缓存逻辑  

##### Version 6.0.5
> 1.修复chan未设置长度的bug  

##### Version 6.0.4
> 1.修复view2和views2  

##### Version 6.0.3
> 1.修复DB prepare配置  

##### Version 6.0.2
> 1.批量没有默认返回空map  

##### Version 6.0.1
> 1.去掉重复的view接口  

##### Version 6.0.0
> 1.重构-删除无用方法(dede等)  
> 2.重构-优化批量查询  
> 3.重构-优化计数信息缓存  
> 4.增加批量aids获取View信息  
> 5.增加单aid获取view的http接口  
> 6.增加SetStat rpc方法(mc)  

##### Version 5.6.14
> 1.增加rpc接口,全量更新stat数值(redis)  

##### Version 5.6.13
> 1.增加internal/view  

##### Version 5.6.12
> 1.cache回写逻辑  

##### Version 5.6.11
> 1.去除404的header  

##### Version 5.6.10
> 1.click走mysql  

##### Version 5.6.9
> 1.批量计数查不到不设置空值  

##### Version 5.6.8
> 1.修改identity为verify  

##### Version 5.6.7
> 1.修复stat panic  

##### Version 5.6.6
> 1.修复rows.next()  

##### Version 5.6.5
> 1.接入新配置中心
> 2.rpc接口参数校验  
> 3.去除hbase  

##### Version 5.6.4
> 1.升级go-common  

##### Version 5.6.3
> 1.rpc接口支持缓存的修改  

##### Version 5.6.2
> 1.日志错误修复  

##### Version 5.6.1
> 1.archive/page走新表,修改sql

##### Version 5.6.0
> 1.archive/page走新表

##### Version 5.5.1
> 1.videos接口增加Ptitle

##### Version 5.5.0
> 1.RPC增加根据aids获取stat接口

##### Version 5.4.0
> 1.paas发布占用

##### Version 5.3.6
> 1.RPC增加一级分区最新视频与数量接口  
> 2.RPC增加Upcount方法 获取用户投稿总数  
> 3.内部http接口改名  
> 4.升级go-common  

##### Version 5.3.5
> 1.PGC只查status=开放的  

##### Version 5.3.4
> 1.ArcsNoCheck2接口校验,aid为空则直接返回参数错误  

##### Version 5.3.3
> 1.monitor挪到内部接口  

##### Version 5.3.2
> 1.修复批量用户动态panic的bug  
> 2.增加field数量  

##### Version 5.3.1
> 1.统一monitor ping接口  
> 2.修复批量用户动态panic的bug  
> 3.增加field数量  
> 4.分页接口增加兼容性处理  

##### Version 5.3.0
> 1.修复redis cache  

##### Version 5.2.8
> 1.增加根据aids获取seasonid接口 rpc  
> 2.更改up过审稿件sql的排序字段  

##### Version 5.2.7
> 1.up主过审稿件改为pubtime排序  

##### Version 5.2.6
> 1.增加RPC方法,根据mids获取最新投稿
> 2.支持attribute参数,在列表中去除展示
> 3.升级go-common

##### Version 5.2.5
> 1.增加RPC方法,根据aids获取archive聚合信息  

##### Version 5.2.4
> 1.增加RPC方法根据EpID获取cid  
> 2.增加RPC方法根据CID获取video信息  

##### Version 5.2.3
> 1.注释PGCproc方法  

##### Version 5.2.2
> 1.增加RPC的分区信息接口  

##### Version 5.2.1
> 1.注释pgc方法  

##### Version 5.2.0

> 1.升级go-common新版本  
> 2.fix view接口，多次查单个请求改为批量请求  
> 3.conf支持优先从本地加载配置  

##### Version 5.1.3

> 1.archives/nocheck接口新增返回返回archive_video和archive_video_audit表数据  
> 2.去掉moment逻辑  

##### Version 5.1.2

> 1.router加入rpcCloser  

##### Version 5.1.1

> 1.忽略video计数错误  

##### Version 5.1.0

> 1.升级配置中心  
> 2.使用公用identify  
> 3.使用统一参数开关  

##### Version 5.0.0

> 1.net/rpc升级为golang/rpcx  

##### Version 4.3.0

> 1.新增rpc获取稿件点击数量  
> 2.新增rpc通过cid查aid  
> 3.更新go-business

##### Version 4.2.5

> 1.分享计数增加databus双写  

##### Version 4.2.4

> 1.新增videoshot接口供管理后台访问  

##### Version 4.2.3

> 1.videoshot接口增加稿件状态校验  

##### Version 4.2.2

> 1.依赖包升级  

##### Version 4.2.1

> 1.修复db使用错误  

##### Version 4.2.0

> 1.添加获取视频详情rpc接口  

##### Version 4.1.4

> 1.fix len(attens) == 0 不能被除  

##### Version 4.1.3

> 1.更新所有匿名rpc client为默认user  

##### Version 4.1.2

> 1.修改syslog日志和上报  

##### Version 4.1.1

> 1.更新go-business为1.3.1  

##### Version 4.1.0

> 1.支持查询pgc信息  
> 2.支持查询用户关注的up主的过审稿件  

##### Version 4.0.0

> 1.go vendor支持  
> 2.go-common/business换成go-business包  
> 3.获取本机ip注册到zk  
> 4.memcache批量获取支持多连接并发  
> 5.新增rpc日志  

##### Version 3.6.1

> 1.修复第一次分享的topic  

##### Version 3.6.0

> 1.新增稿件page信息接口  

##### Version 3.5.1

> 1.修复批量获取cache出错还加入cache问题  

##### Version 3.5.0

> 1.获取稿件列表不检测权限  
> 2.修复稿件分区变更后二级分区最新视频转移分区  

##### Version 3.4.3

> 1.修复二级分区最新视频安装pubdate排序  

##### Version 3.4.2

> 1.新增update稿件cache  

##### Version 3.4.1

> 1.数组越界bug  

##### Version 3.4.0

> 1.增加获取up主投稿列表接口  
> 2.修复增加全量分区视频时变量没有重新初始化bug  
> 3.优化缓存key使均匀分布  

##### Version 3.3.1

> 1.修复最新视频bug：新增视频可见过滤条件：access、attrBitNoWeb、attrBitNoMobile  

##### Version 3.3.0

> 1.增加分区的视频按投稿时间排序  
> 2.新增查询过审记录接口  

##### Version 3.2.3

> 1.修复回复的稿件置首bug  

##### Version 3.2.2

> 1.修复Archive接口cache bug  

##### Version 3.2.1

> 1.修复elk日志  

##### Version 3.2.0  

> 1.稿件添加字段reject_reason  
> 2.修改share接口  
> 3.新增set_tag接口  
> 4.支持trace v2  

##### Version 3.1.0  

> 1.新增获取stat接口  
> 2.新增获取多条stat接口  
> 3.新增stat更新redis接口  
> 4.修改稿件获取stat的方法  
> 5.增加或修改ping方法  
> 6.优化部分代码  

##### Version 3.0.0  

> 1.context使用官方接口  
> 2.添加share计数  
> 3.优化部分代码  

##### Version 2.5.0  

> 1.新增视频缩略图版本号  
> 2.支持视频缩略图更新cid  
> 3.添加up主视频动态接口  

##### Version 2.4.0  

> 1.添加获取用户最新评论稿件以及后台job  
> 2.优化配置  
> 3.添加服务发现  


##### Version 2.3.0  

> 1.添加获取videoshot接口  
> 2.rpc调用bug  

##### Version 2.2.0  

> 1.优化  
> 2.add elk  
> 3.add trace id  
> 4.add haiwai api  
> 5.remove noused code  
> 6.add mid recommend  
> 7.fix some bug  

##### Version 2.1.0  

> 1.add tracer  

##### Version 1.1.0  

> 1.基于go-common重构  

##### Version 1.0.0  

> 1.初始化完成稿件基础查询功能
