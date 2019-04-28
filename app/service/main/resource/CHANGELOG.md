##  内容运营服务   

### Version 2.25.5
##### Features
> 1.增加获取审核态接口  

### Version 2.25.4
##### Features
> 1.科技区右侧推广栏输出标签  

### Version 2.25.3
##### Features
> 1.sidebar增加Language  

### Version 2.25.2
##### Features
> 1.修改resource表查询语句的排序

### Version 2.25.1
##### Features
> 1.更换地区限制方法 

### Version 2.25.0
##### Bugfix
> 1.修复specialCache并发读写的问题  
> 2.不需要的err报错改为warn  

### Version 2.24.8
##### Features
> 1.接入grpc  
> 2.pgc特殊卡片、相关推荐  

### Version 2.24.7
##### Features
> 1.增加获取贴片视频cid接口  

### Version 2.24.6
> 1.location的Zone接口修改为Info  

### Version 2.24.5
##### Features
> 1.新增需要展示标签的位置  
> 2.URL监测功能针对直播URL的安全校验做兼容  

### Version 2.24.4
##### Bugfix
> 1.修改banner返回nil导致空指针异常的问题  

### Version 2.24.3
##### Features
> 1.修改abtest逻辑，大于等于改为大于  

### Version 2.24.2
##### Features
> 1.修改abtest的buvid转换方法  

### Version 2.24.1
##### Features
> 1.banner增加过滤逻辑，降低广告请求频率(参数中的version不为空且等于本地hashcache时，直接返回空)  

### Version 2.24.0
##### Features
> 1.identify为grpc，切换verify  

### Version 2.23.1
##### Features
> 1.siderbar增加红点字段  

### Version 2.23.0
##### Features
> 1.增加移动端“我的”数据接口  
> 2.获取resource数据时，去掉status筛选逻辑  
> 3.增加abtest接口逻辑  
> 4.增加rows.Err  

### Version 2.22.12
##### Features
> 1.URL监控增加限速  

### Version 2.22.11
##### Features
> 1.对被监控的稿件数据做格式兼容处理  

### Version 2.22.10
##### Features
> 1.新增需要获取稿件信息的推广位id  

### Version 2.22.9
##### Features
> 1.增加音频分区推荐卡片接口  

### Version 2.22.8
##### Features
> 1.整合告警信息  
> 2.URL监控告警逻辑简化  
> 3.移动端banner数据增加build过滤逻辑  
> 4.移动端banner输出数据增加stime字段  

### Version 2.22.7
##### Features
> 1.优化稿件自动下线和监测、URL监测告警的逻辑  

### Version 2.22.6
> 1.http default client add timeout

### Version 2.22.5 - 2018.05.25
##### Features
> 3.修改需要加label的位置ID  

### Version 2.22.4 - 2018.05.24
##### Features
> 1.完善内容运营后台的投放内容自动下线和告警的逻辑  
> 2.补充自动下线的各种日志  
> 3.增加需要加label的位置ID  

### Version 2.22.3 - 2018.05.23
##### Features
> 1.去掉投放内容自动下线逻辑中的多余log  

### Version 2.22.2 - 2018.05.23
##### Features
> 1.对URL监控功能中的URL做处理  

### Version 2.22.1 - 2018.05.22
##### Features
> 1.增加稿件和URL监控的开关  

### Version 2.22.0 - 2018.05.22
##### Features
> 1.去掉无用的http接口assignment、defbanner  
> 2.增加url类型模拟请求和告警  
> 3.告警方式变为企业微信  
> 4.根据稿件状态，自动下线内容运营数据  

### Version 2.21.0 - 2018.05.03
##### Features
> 1.http切bm  

### Version 2.20.1 - 2018.04.24
##### Features
> 1.接archive的discovery  
> 2.推广内容告警触发条件改为：稿件状态变更且变更后的状态不是开发浏览  
> 3.告警邮件标题  

### Version 2.20.0 - 2018.04.24
##### Features
> 1.rpc lient 增加discovery new方法  

### Version 2.19.1 - 2018.04.19
##### BugFix
> 1.修复creative_type赋值的问题  

### Version 2.19.0 - 2018.04.18
##### Features
> 1.增加创作中心creative_type字段  

### Version 2.18.2 - 2018.04.12
##### BugFix
> 1.修改banner排序逻辑  

### Version 2.18.1 - 2018.04.12
##### Features
> 1.接discovery，添加register接口  

### Version 2.18.0 - 2018.04.09
##### Features
> 1.增加直播弹幕盒子接口  

### Version 2.17.0 - 2018.02.27
##### Features
> 1.label恢复原有逻辑，不再使用note临时替代  
> 2.增加地区限制过滤  
> 3.修改推荐池投放的优先级  

### Version 2.16.2 - 2018.01.30
##### Features
> 1.未登录贴片增加aid和是否跳转  
> 2.番剧贴片增加aid  

### Version 2.16.1 - 2018.01.16
##### Features
> 1.优化video_ads相关逻辑，删除无用的逻辑和接口  
> 2.番剧贴片接口增加了跳转url字段  

### Version 2.16.0 - 2018.01.15
##### Features
> 1.播放器控件添加Hash  
> 2.修改resource和bannner的SQL  
> 3.各目录增加单元测试  

### Version 2.15.1 - 2018.01.05
##### BugFix & Features
> 1.优化编辑投放稿件状态变化告警邮件内容  
> 2.修复map并发读写问题  
> 3.修复推荐池会读取历史素材的问题  

### Version 2.15.0 - 2018.01.05
##### Features
> 1.增加番剧获取贴片的接口  

### Version 2.14.3 - 2018.01.02
##### BugFix
> 1.修复resource并发引起的slice越界问题  

### Version 2.14.2 - 2017.12.29
##### BugFix
> 1.兼容旧逻辑，返回给web的assignment数据的weight全部置为0  

### Version 2.14.1 - 2017.12.29
##### Features
> 1.兼容旧逻辑，修改resource_assignment表读的值(position->weight)  

### Version 2.14.0 - 2017.12.29
##### Features
> 1.添加获取播放器控件接口  

### Version 2.13.1 - 2017.12.28
##### Features
> 1.兼容部分旧逻辑  

### Version 2.13.0 - 2017.12.28
##### Features
> 1.根据新的内容运营平台修改resource和banner逻辑  

### Version 2.12.3 - 2017.12.20
##### Features
> 1.修改SQL，确保读的是旧数据。防止新版后台预发添加新数据影响线上服务  

### Version 2.12.2 - 2017.11.20
##### Features
> 1.banner限制总数修改  

### Version 2.12.1 - 2017.11.20
##### Bug
> 1.修正对RPC方法PasterAPP的err的判断  

### Version 2.12.0 - 2017.11.20
##### Bug
> 1.RPC方法DefBanner、Resource、PasterAPP增加返回值nil或err的判断  

### Version 2.11.0 - 2017.11.13
##### Features
> 1.banner接口增加aid参数、增加透传字段ad_extra(替代原来lat、lng)  
> 2.商业广告接口返回值增加透传字段extra  
> 3.提供新接口，获取首页引导图  

### Version 2.10.0 - 2017.11.7
##### Features
> 1.未登录贴片的投放目标ID(aid、season_id、type_id)进行数据库字段拆分  

### Version 2.9.0 - 2017.11.3
##### Features
> 1.未登录贴片的投放目标ID改为支持逗号分隔  

### Version 2.8.0 - 2017.11.1
##### Features
> 1.banner接口增加经纬度字段、增加open_event字段  

### Version 2.7.0 - 2017.10.25
##### Features
> 1.增加获取登录引导贴片的接口  

### Version 2.6.2 - 2017.10.17
##### Bug
> 1.修正bilibili_ads库的video_ads表aid为null的情况下，逗号分隔报错导致初始化失败的问题  

### Version 2.6.1 - 2017.10.17
##### Bug
> 1.修正bilibili_ads库的video_ads表部分数据default NULL导致panic的问题  

### Version 2.6.0 - 2017.10.16
##### Features
> 1.banner接口增加传参(version)  
> 2.修改banner的rpc和http接口返回值  
> 3.banner的dao层修改查询条件  

### Version 2.5.1 - 2017.10.13
##### Bug
> 1.修复banner逻辑初始化res的bug  

### Version 2.5.0 - 2017.10.11
##### Features
> 1.banner接口改为批量接口，接收多个resource_id  
> 2.banner接口的plat参数改为接收调用方传参  
> 3.banner接口增加is_ad参数，判断是否调用广告接口  
> 4.调用广告接口的方法去掉版本和mid判断  

### Version 2.4.0 - 2017.10.11
##### Features
> 1.banner接口不再进行签名校验  
> 2.banner接口返回值结构修改  
> 3.banner调广告接口的逻辑添加mobile_app和build的判断逻辑  
> 4.banner添加rpc接口  

### Version 2.3.0 - 2017.09.27
##### Features
> 1.去掉ecode.Init  
> 2.httpClient请求去掉app  

### Version 2.2.0 - 2017.09.07
##### Features
> 1.app-show的banner逻辑整体迁移进来  

### Version 2.1.0 - 2017.08.28
##### Features
> 1.添加获取bilibili_ads库video_ads表的方法(aid维度)  
> 2.添加获取bilibili_ads库video_ads表的方法(seasonid维度)  
> 2.广告接口调整目录(ad变为cpm)  

### Version 2.0.0 - 2017.08.23
##### Features
> 1.添加全量获取resource表数据的RPC方法  
> 2.添加全量获取resource_assignment表数据的RPC方法  
> 3.添加default_one表查询接口(http/rpc)  
> 4.添加单独和批量查询resource接口(http/rpc)  
> 5.添加单独和批量查询assignment接口(http/rpc)  

### Version 2.0.0 - 2017.07.21
##### Features
> 1.更新go-common v7    
> 2.去掉go-business依赖  

### Version 1.0.0
##### Features
> 1.广告获取接口(针对APP)  
> 2.添加prom监控  
