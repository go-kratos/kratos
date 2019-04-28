#### Basic service of historical function

#### Version v5.9.11
> 1.调整 Histories businesses 

#### Version v5.9.10
> 1.屏蔽aid =-1|| ugc&pgc cid==0 

#### Version v5.9.9
> 1.屏蔽aid =-1   

#### Version v5.9.8
> 1.无脑式删除两遍     

#### Version v5.9.7
> 1.web aid pgc 删除   

#### Version v5.9.6
> 1. aid pgc 删除   

#### Version v5.9.5
> 1. position 返回NothingFound

#### Version v5.9.4
> 1. Business id 转换

#### Version v5.9.3
> 1. 删除番剧  fix 

#### Version v5.9.2
> 1. 同步删除   
> 2. 删除番剧  

#### Version v5.9.1
> 1. read tidb view at

#### Version v5.9.0
> 1. read tidb 

#### Version v5.8.5
> 1. 去除v1 目录  
> 2. 删除pgc http   

#### Version v5.8.4
> 1. NewEpIndex  
> 2. 清空操作加入行为日志   

#### Version v5.8.3
> 1. 重新打包docker  

#### Version v5.8.2
> 1. ToViewOverMax  

#### Version v5.8.1
> 1. 删除hbase 90天之前的数据   
> 2. 时间相同 按照aid排序    

#### Version v5.8.0
> 1. toview ugcPay     

#### Version v5.7.8
> 1. 修复hide

#### Version v5.7.7
> 1. NewestEpIndex - NewEpShow 
> 2. 配置文件控制迁移番剧grpc   

#### Version v5.7.6
> 1. gorpc  service    

#### Version v5.7.5
> 1. bangumi  grpc   

#### Version v5.7.4
> 1. 增加clear 内网接口   

#### Version v5.7.3
> 1. 修复 delete 无business bug

#### Version v5.7.1
> 1. delete 改为同步删除

#### Version v5.7.0
> 1. 接入history-service 双写

#### Version v5.6.11
> 1.grpc.  

#### Version v5.6.10
> 1.rpc discovery.  

#### Version v5.6.9
> 1.删除remote ip.  

#### Version v5.6.8
> 1.删除PlayProgress-T 生产方代码 .  

#### Version v5.6.7
> 1.屏蔽上报数据为离线类型 .   

#### Version v5.6.6
> 1.mid .   

#### Version v5.6.5
> 1.hbase 落库    

#### Version v5.6.4
> 1. 重新打镜像        

#### Version v5.6.3
> 1. auth&verify     

#### Version v5.6.2
> 1. docker 

#### Version v5.6.1
> 1. 缓存数据填充business
   
#### Version v5.6.0
> 1.合并播放进度    

#### Version v5.5.2
> 1. 修改delete rpc 参数类型
#### Version v5.5.1
> 1. 修复查询bug

#### Version v5.5.0
> 1. 查询历史记录接口增加分类字段
> 2. 删除历史记录接口使用business标志
> 3. 增加清空列表rpc

#### Version v5.4.1
> 1. 优化游标未找到的情况  

#### Version v5.4.0
> 1. 统一业务标志
    
#### Version v5.3.2
> 1.接入漫画     
> 2.增加内网批量删除接口          

#### Version v5.3.1
> 1.update hbase sdk

#### Version v5.3.0
> 1.update infoc sdk

#### Version v5.2.11
> 1.DELETE kafka  

#### Version v5.2.10
> 1.单独获取直播|文章 历史记录列表           

#### Version v5.2.9
> 1.兼容app 播放器上报参数缺少情况       

#### Version v5.2.8
> 1.分批调用收藏RPC接口
> 2.修改内网查用户历史记录接口

#### Version v5.2.7
> 1.增加安卓端TV设备类型

#### Version v5.2.6
> 1.迁移收藏接口至rpc接口isFavs

#### Version v5.2.5
> 1.fix toviews adds, 已存在列表后，添加失败.
> 2.增加管理员查看用户的历史记录和稍后待看.
> 3.修改用户加经验调整顺序，redis操作前.

#### Version v5.2.4
> 1.增加默认时间

#### Version v5.2.3
> 1.fix web 单删

#### Version v5.2.2
> 1.增加report 内网接口
> 2.用户唯一rowkey  
> 3.增加live数据    
> 4.删除source参数         

#### Version v5.2.1
> 1.迁出kafka
> 2.toview去掉分区限制逻辑        

#### Version v5.2.0
> 1.迁移用户增加经验通道至databus        

#### Version v5.1.5
> 1.迁移至business/interface/main目录    

#### Version v5.1.4
> 1.迁移至BM框架  
> 2.删除bangumi 老接口  
> 3.bangumi 批量请求       

#### Version v5.1.3
> 1.databus client 重新打包      

#### Version v5.1.2
> 1.删除statsd   

#### Version v5.1.1
> 1.更新bangumi信息      

#### Version v5.1.0
> 1.去掉老hbase        

#### Version v5.0.1
> 1.fix mid      

#### Version v5.0.0
> 1.历史记录平台化     

#### Version v4.10.2
> 1.修复history热点key问题

#### Version v4.10.1
> 1.兼容android app端subtype      

#### Version v4.10.0
> 1.独立出toview cache conf      

#### Version v4.9.14
> 1.修复缓存context     

#### Version v4.9.13
> 1.返回数据增加sub_type字段  

#### Version v4.9.12
> 1.report增加sub_type字段  

#### Version v4.9.10 
> 1.升级archive3接口  

#### Version v4.9.9
> 1.trace测试

#### Version v4.9.8
> 1.trace测试

#### Version v4.9.7
> 1.trace测试

#### Version v4.9.6  
> 1.修复webToview data 为nil  
> 2.同步番剧数据增加类型和子类型  
> 3.增加hbase deadline-exceeded日志  

#### Version v4.9.5  
> 1.去掉调用稿件AidByCid  

#### Version v4.9.4  
> 1.优化用户无稍后再看数据问题  

#### Version v4.9.3  
> 1.用户稍后再看数据迁移      

#### Version v4.9.2  
> 1.用户稍后再看数据双写  

#### Version v4.9.1  
> 1.迁移用户开关状态  
> 2.info hbase库family fix                    

#### Version v4.9.0  
> 1.迁移hbase数据,不支持回滚      

#### Version v4.8.2  
> 1.增加双读hbase  

#### Version v4.8.1  
> 1.增加新hbase日志  

#### Version v4.8.0  
> 1.增加双写hbase   

#### Version v4.7.3
> 1.web端带看列表增加bangumi信息   

#### Version v4.7.2
> 1.获取待看列表cache map

#### Version v4.7.1
> 1.移除稿件过滤逻辑 

#### Version v4.7.0
> 1.增加批量增加待看列表功能  

#### Version v4.6.0
> 1.merge history into kratos

##### Version v4.5.10

> 1.fix map key  
> 2.增加反作弊上报日志    
> 3.冷方式删除番剧多集   

##### Version v4.5.9

> 1.fix 移动端批量同步缓存

##### Version v4.5.8

> 1.待看列表删除大于30s

##### Version v4.5.7

> 1.接入log-agent

##### Version v4.5.6

> 1.更新go-common v6.16.0，go-business v2.21.4  
> 2.接入新prom
> 3.删除delChan增加监控    

##### Version v4.5.5

> 1.优化databus key  
> 2.databus增加字段  
> 3.fix 版权问题

##### Version v4.5.4

> 1.修复hbase TODO.  

##### Version v4.5.2

> 1.读hbase降级
> 2.去除hbase双写
> 3.升级go-common go-business

##### Version v4.5.1

> 1.fix hbase表名

##### Version v4.5.0

> 1.hbase表双写

##### Version v4.4.3

> 1.prom

##### Version v4.4.2

> 1.Identify&csrf

##### Version v4.4.1

> 1.修复aids 空间
> 2.接入prom  

##### Version v4.4.0

> 1.去掉批量查询分P信息  
> 2.修复P数，待看列表的总数  
> 3.接入proms 

##### Version v4.3.0

> 1.增加web端待看列表接口  
> 2.待看列表增加page分批信息  
> 3.增加总P数   
> 4.待看列表aid>0   

##### Version v4.2.5

> 1.接入新配置中心 

##### Version v4.2.4

> 1.修复AddHistories 
> 2.增加rpc 

##### Version v4.2.2

> 1.平滑发布 

##### Version v4.2.0

> 1.删除report进度表 

##### Version v4.1.2

> 1.去除report表逻辑 
> 2.databus key mid%100 

##### Version v4.1.1

> 1.修复toview列表数量超限  

##### Version v4.0.21

> 1.增加白名单日志 
> 2.修复了ctx声明 

##### Version v4.0.20

> 1.修复mid为0不保存历史记录   

##### Version v4.0.4

> 1.缓存数量限制 

##### Version v4.0.3

> 1.修复缓存数量限制 

##### Version v3.5.0

> 1.job异步持久化 
> 2.重构历史记录 

##### Version v3.4.5

> 1.修改monitor/ping调用方式  

##### Version v3.4.4

> 1.更新http接口支持内外网流量 

##### Version v3.4.3 

> 1.增加收藏标识 
> 2.历史记录暂停功能  

##### Version v3.4.3 

> 1.增加收藏标识  
> 2.历史记录暂停功能  

##### Version 3.4.0 

> 1.接入ci配置中心   
> 2.只记录登陆用户上报   
> 3.添加inner接口   

##### Version 3.3.1 

> 1.兼容ios参数float  

##### Version 3.3.0  

> 1.增加cid查询稿件信息  
> 2.更新identity  
> 3.天马稍后再看  

##### Version 3.2.1

> 1.更新vendor依赖  

##### Version 3.2.0

> 1.更新vendor依赖  

##### Version 3.1.0

> 1.番剧播放离线同步  

##### Version 3.0.0

> 1.govendor支持  
> 2.go-business依赖  
> 3.增加habse平监控  

##### Version 2.7.2

> 1.清除pay，全部返回0  

##### Version 2.7.1

> 1.修复历史记录条数截断  

##### Version 2.7.0

> 1.movie字段特殊化处理  
> 2.暂时去除稿件状态判断，兼容安卓  
> 3.删除redis和hbase1 全量走hbase2 

##### Version 2.6.2

> 1.修复移动端循环加载历史记录  

##### Version 2.6.1

> 1.添加开关控制播放历史流量走新Hbase表  

##### Version 2.6.0

> 1.添加trace_v2  
> 2.修改hbase存储方式双写  

##### Version 2.5.0

> 1.添加获取播放历史稿件id列表的接口

##### Version 2.4.1

> 1.[bug]判断页码页大小  
> 2.[bug]hbase scan需要加上分隔符"|"  

##### Version 2.4.0

> 1.持久化由redis迁移至hbase

##### Version 2.3.0

> 1.添加服务发现

##### Version 2.2.1

> 1.优化配置

##### Version 2.2.0

> 1.优化
> 2.add tracer id
> 3.add elk

##### Version 2.1.0

> 1.add tracer

##### Version 2.0.0

> 1.api 2.0

##### Version 1.2.0

> 1.reconstruction

##### Version 1.1.0

> 1.按照go-common重构

##### Version 1.0.0

> 1.重构历史
