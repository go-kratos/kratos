#### report-click 点击上报

##### Version 2.17.6
> 1. 增加heartbeat/mobile上报失败的接口
> 2. heartbeat/mobile中增加build号的prom上报

##### Version 2.17.5
> 1. 修复h5播放上报的cookie中的did无法解析出正确的ftime的问题

##### Version 2.17.4
> 1. h5播放上报的ip改用metadata中的ip

##### Version 2.17.3
> 1. http Shutdown

##### Version 2.17.2
> 1. 将inline_play_heartbeat改为inline_play_to_view 

##### Version 2.17.1
> 1. IPv4 fix 

##### Version 2.17.0
> 1. 新增inline播放超10秒计入播放数
> 2. 原有的auto_play改为inline_play_begin（inline开播）和inline_play_heartbeat（inline播放结束）
> 3. play中新增session参数
> 4. cache.New改为Fanout

##### Version 2.16.15
> 1. IPv6 fix 

##### Version 2.16.14
> 1. IPv6

##### Version 2.16.13
> 1. 动态自动播放 from 711,7111,712,7151,7161,7163,7171,7181 增加播放计数


##### Version 2.16.12
> 1. 删除 Topic:ArchiveClick-T databus
> 2. 删除model 

##### Version 2.16.11
> 1. /x/report/heartbeat/mobile auto_play(2) ===> /x/report/click/android2 & ios 

##### Version 2.16.10
> 1. copy buf 

##### Version 2.16.9
> 1. detail_play_time&list_play_time       

##### Version 2.16.8
> 1. autoPlay   

##### Version 2.16.7
> 1. merge databus   

##### Version 2.16.6
> 1. delete csrf  

##### Version 2.16.5
> 1. RemoteIP()--> metadata.RemoteIP      

##### Version 2.16.4
> 1. h5 click 优化播放数上报           

##### Version 2.16.3 
> 1. 自动播放   

##### Version 2.16.2
> 1. 播放时长返回ts 时间 

##### Version 2.16.1
> 1. mid 0 

##### Version 2.16.0
> 1. bm 

##### Version 2.15.0
> 1. update infoc sdk

##### Version 2.14.6
> 1. 播放时长增加 epid_status play_status user_status 字段     

##### Version 2.14.5
> 1. 播放时长增加 play_mode device from 字段     
> 2. 点击上报增加 play_mode platform device mobi_app 字段    
> 3. 去掉kafka          

##### Version 2.14.4
> 1. 增大服务端超时范围

##### Version 2.14.3
> 1. 修复mobileapp 接口上报参数mobile_app至mobi_app


##### Version 2.14.2 
> 1.account依赖接入discovery


##### Version 2.14.1
> 1.使用account-service v7

##### Version 2.14.0
> 1.增加h5 外链点击上报
> 2.迁移interface main大目录    

##### Version 2.13.2
> 1.删除 老infoc通道上报数据    

##### Version 2.13.1
> 1.删除statsd

##### Version 2.13.0
> 1.点展数据双写kafka   
> 2.增加Android TV点展接口        

##### Version 2.12.8
> 1.增加app端统计上报播放时长需求
> 2.修改infoc2的配置方式   

##### Version 2.12.7
> 1.fix context.TODO()      

##### Version 2.12.6
> 1.rpc.account.card2 to UserInfo   

##### Version 2.12.5
> 1.兼容参数avid ios app   
       
##### Version 2.12.4
> 1.针对移动端上报，传参mid>0，并且没有accesskey或者accesskey校验不通过时使用不同UA标记

##### Version 2.12.3
> 1.针对移动端传accesskey与参数mid不一致UA标记

##### Version 2.12.2
> 1.没有校验通过的和检验通过的用户采用不同的UA进行标记
> 2.增加/x/report/click/outer接口14004警告日志

##### Version 2.12.1
> 1.heartbeat增加新字段playtype上报，兼容playtype和play_type

##### Version 2.12.0
> 1.增加app端点展access_key参数  

##### Version 2.11.4
> 1.修复移动端mid防刷  

##### Version 2.11.3
> 1.心跳日志web端与app端字段个数保持一致   

##### Version 2.11.2
> 1.点展增加接口上报字段（type、sub_type、sid、epid）

##### Version 2.11.1
> 1.增加数据写入Kafka时的异步处理

##### Version 2.11.0
> 1.增加AI&dataPlatfrom 数据上报字段     

##### Version 2.10.4
> 1.删除日志  

##### Version 2.10.3
> 1.增加sub_type参数  

##### Version 2.10.2
> 1.修复web端mid防刷  

##### Version 2.10.1
> 1.flash端账号防刷     

##### Version 2.10.0
> 1.接入大仓库  

##### Version 2.9.0
> 1.优化用户等级   

##### Version 2.5.13
> 1.heartbeat 增加pause字段,click增加buvid字段     

##### Version 2.5.12
> 1.修改为不区分大小写did，兼容客户端      

##### Version 2.5.11

> 1.升级go-common，go-business(去掉panic)  
> 2.接入prom  
> 3.增加对m域名的区分上报  

##### Version 2.5.10

> 1.修复IP字段为空  
> 2.使用GuestConfigPost和GuestPost  

##### Version 2.5.9

> 1.升级infoc2  

##### Version 2.5.8

> 1.升级go-business修复csrf  
> 2.兼容seasonID  

##### Version 2.5.7

> 1.infoc改成一条条发  

##### Version 2.5.6

> 1.增加infoc超时  

##### Version 2.5.5

> 1.修复csrf   

##### Version 2.5.4

> 1.增加infoc错误日志  

##### Version 2.5.3

> 1.修复mid认证   

##### Version 2.5.2

> 1.更新go-business的infoc2  

##### Version 2.5.1

> 1.本地tw  

##### Version 2.5.0

> 1.接入新配置中心
> 2.web/heartbeat 改为GuestPost
> 3.report-click接入历史记录rpc

##### Version 2.4.3

> 1.心跳日志双写一份到下沙  

##### Version 2.4.2

> 1.monitor ping  

##### Version 2.4.1

> 1.平滑发布  

##### Version 2.4.0

> 1.增加epid参数  

##### Version 2.3.4

> 1.增加配置开关  

##### Version 2.3.3

> 1.升级库切identify  

##### Version 2.3.2

> 1.先注释掉播放上报  

##### Version 2.3.2

> 1.先注释掉播放上报  

##### Version 2.3.1

> 1.修复历史进度上报为POST  

##### Version 2.3.0

> 1.添加历史进度上报   
 
##### Version 2.2.3

> 1.修复web心跳buvid  

##### Version 2.2.2

> 1.拜年祭上报兼容  

##### Version 2.2.1

> 1.过滤aid=0的情况  

##### Version 2.2.0

> 1.增加web播放器心跳上报  

##### Version 2.1.1

> 1.升级vendor  

##### Version 2.1.0

> 1.配置中心  
> 2.升级vendor  

##### Version 2.0.3

> 1.新增播放上报字段  

##### Version 2.0.2

> 1.播放上报新增字段  

##### Version 2.0.1

> 1.更新infoc vendor支持tcp发  

##### Version 2.0.0

> 1.更新底层库依赖  
> 2.增加govendor  

##### Version 1.2.1

> 1.修复ChecDid bug

##### Version 1.2.0

> 1.h5上报算法优化  
> 2.web h5上报接口  

##### Version 1.1.0

> 1.支持视频播放上报，每隔30s上报一次只打日志  

##### Version 1.0.0

> 1.report-click 基础功能
