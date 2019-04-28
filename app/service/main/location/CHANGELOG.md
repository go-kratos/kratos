### location-service

#### Version 6.1.1
> 1.account grpc   

#### Version 6.1.0
> 1.迁移gRPC目录  

#### Version 6.0.0
> 1.使用IPIP的IPv4库  
> 2.使用IPIP的IPv6库  
> 3.去掉rpc的Zone、Zones、Info2接口  

#### Version 5.8.6
> 1.删除dao层无用代码  
> 2.VPN匿名库去掉本地文件逻辑  
> 3.重构authByPids(旧接口已无调用)  

#### Version 5.8.5
> 1.单IP查询接口判断nil  

#### Version 5.8.4
> 1.IP查询接口增加country_code  

#### Version 5.8.3
> 1.视频云IP库和旧库的加载逻辑对齐  

#### Version 5.8.2
> 1.live zk空配置不注册

#### Version 5.8.1 - 2018.09.29
###### Features
> 1.判断IP是否为VPN地址  
> 2.增加批量查询gid权限接口  

#### Version 5.8.0 - 2018.08.14
###### Features
> 1.优先使用二进制IP库  

#### Version 4.8.0 - 2018.08.08
###### Features
> 1.identify为grpc，切换verify  

#### Version 4.7.4 - 2018.07.10
###### BugFix
> 1.修复grpc方法的返回值是nil导致的panic  

#### Version 4.7.3 - 2018.07.09
###### Features
> 1.grpc增加注册liveZK逻辑  

#### Version 4.7.2 - 2018.07.03
###### Features
> 1.增加grpc接口  

#### Version 4.7.1 - 2018.06.05
###### Features
> 1.policy_item、policy_group表增加软删除过滤  

#### Version 4.7.0 - 2018.05.03
###### Features
> 1.http切bm  

#### Version 4.6.0 - 2018.04.24
###### Features
> 1.rpc lient 增加discovery new方法  

#### Version 4.5.2 - 2018.04.12
###### Features
> 1.接discovery，添加register接口  

#### Version 4.5.1 - 2018.03.25
###### Features
> 1.使用account-service v7  

#### Version 4.5.0 - 2018.03.12
###### Features
> 1.添加archive2接口(http和rpc)  

#### Version 4.4.0 - 2018.03.07
###### Features
> 1.添加AuthPIDs接口，判断ipaddr是否符合规则pids  

#### Version 4.3.1 - 2017.10.22
###### Features
> 1.archive和group接口，如果传的是局域网IP，使用CDNIP  

#### Version 4.3.0 - 2017.09.27
###### Features
> 1.重构zone和check的http、rpc接口（第一版写的什么辣鸡玩意儿！！)  
> 2.删除dat格式IP库的解析逻辑  
> 3.删除下载IP库文件接口  
> 4.删除获取IP库全部数据的rpc接口  
> 5.去掉ecode初始化  

#### Version 4.2.0 - 2017.09.08
###### Features
> 1.添加根据pids和ip获取稿件权限的RPC和HTTP接口  

#### Version 4.1.0 - 2017.09.08
###### Features
> 1.添加根据pids和ip获取稿件权限的RPC和HTTP接口  

#### Version 4.0.0 - 2017.07.21
###### Features
> 1.business/model/location移到项目目录中(business/service/main/location/model)  
> 2.添加CONTRIBUTORS.md和CHANGELOG.md    

#### Version 3.0.0
###### Features
> 1.location-service meger into Kratos/business/service  

#### Version 2.0.0
###### Features
> 1.更新go-common v7.0.0  
> 2.去掉go-business 依赖  
> 3.添加rpc方法Info2和Infos2（查询IP信息并返回完整zoned_id组）  

#### Version 1.4.4
###### Features
> 1.完善prom监控  

#### Version 1.4.3
###### Features
> 1.添加prom监控  

#### Version 1.4.1
###### Features
> 1.修改http层Verify初始化方法  

#### Version 1.4.0
###### Features
> 1.重新接新版配置中心  
> 2.添加获取指定aid的地区限制规则API  
> 3.添加获取指定规则组的地区限制规则API  

#### Version 1.3.9
###### Features
> 1.暂时改回旧版配置中心  

#### Version 1.3.8
###### Features
> 1.接新配置中心  
> 2.更新了rpc的调用  

#### Version 1.3.6
###### Features
> 1.逻辑重构  
> 2.单IP查询接口  
> 3.多IP查询接口  
> 4.规则查询接口(根据ip、aid、policy_id、group_id)  
> 5.是否可观看查询接口  
> 6.根据group_id批量查询可观看的zone_ids  
> 7.添加强制刷新规则缓存接口  
> 8.添加管理接口，自动下载(运维提供地址)IP库文件并重新加载(只支持.dat文件)  
> 9.添加.dat和.txt双IP库文件识别，根据不同IP库文件调用不同的加载方法和查询方法  

#### Version 1.3.3
###### Features
> 1.多IP查询接口

#### Version 1.3.1
###### Features
> 1.升级go-common v6.2.4,go-business v2.4.3, golang v2.9.1 

#### Version 1.3.0
###### Features
> 1.增加zone 接口代替check接口.  
> 2.升级go-common v5.2.2,go-business location, golang v2.6.0  

#### Version 1.2.0
###### Features
> 1.接入配置中心和CI系统.

#### Version 1.1.0
###### Features
> 1.升级go-common v5.0.0,golang v2.6.0,go-business v2.2.1  

#### Version 1.0.0
###### Features
> 1.升级go-common v4.4.1,golang v2.5.3,go-business v2.1.0  
> 2.增加下载IP库

#### Version 0.1.0
###### Features
> 1.接入配置中心  
> 2.go-common develop分支  

###### Bug Fixes
> 1.修复查找算法  

#### Version 0.0.1
###### Features
> 1.查询IP接口  
> 2.支持vendor
