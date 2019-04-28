#### 收藏夹job

### Version v7.4.1
> 1. fix mv bug

### Version v7.4.1
> 1. fix copy init bug

### Version v7.4.0
> 1. 播单2.0

### Version v7.3.4
> 1. copy move multiDel clean操作创建收藏列表缓存

### Version v7.3.3
> 1. 修改mysql index force

### Version v7.3.0
> 1. 移除无用代码 binlogDao, videoDao

### Version v7.2.1
> 1. 播单播放数防刷

### Version v7.1.0
> 1. 删除videoDao
> 2. 下线binlog

### Version v7.0.2
> 1. change max seq

### Version v7.0.0
> 1. medialist

##### Version 6.0.1
> 1. 播单计数

##### Version 6.0.0
> 1. 停止双写老表

##### Version 4.8.9
> 1. 拜年祭平台收藏聚合支持

##### Version 4.8.8
> 1. 拜年祭aid收藏聚合支持

##### Version 4.8.7
> 1. fix rows err

##### Version 4.8.6
> 1. del recentfavs mc

##### Version 4.8.5
> 1. fix relations slow sql

### Version 4.8.1
> 1. add db rows.error

##### Version 4.7.21
> 1. common config

##### Version 4.7.19
> 1. rm video folder redis
 
##### Version 4.7.18
> 1. add batch oids mc

##### Version 4.7.17
> 1. rm databus stat-t

##### Version 4.7.16
> 1. migrate bm

##### Version 4.7.15
> 1. set oid count mc

##### Version 4.7.14
> 1. fix index ix_fid_state_mtime

##### Version 4.7.13
> 1. add read mysql

##### Version 4.7.12
> 1. enhance binlog cache

##### Version 4.7.11
> 1.fix binlog cache

##### Version 4.7.9
> 1.del cover support migration

##### Version 4.7.8
> 1.add relation binlog to update stat and add user

##### Version 4.7.6
> 1.rm cache err return

##### Version 4.7.5
> 1.account v7

##### Version 4.7.4
> 1.move to main folder
> 1.log error level fix

##### Version 4.7.3
> 1.fix delFolderSQL state

##### Version 4.7.2
> 1.fix replace sql

##### Version 4.7.1
> 1.fix binlog UpdateRelationSQL

##### Version 4.7.0
> 1.migrate archive data by binlog

##### Version 4.6.8
> 1.add dao fav unit test

##### Version 4.6.7
> 1.fix cache save context

##### Version 4.6.6
> 1.fix is faved

##### Version 4.6.5
> 1.fav bit redis key mid reverse

##### Version 4.6.4
> 1.fix count -1 sql err

##### Version 4.6.3
> 1.add batch job sleeptime

##### Version 4.6.2
> 1.new stats databus

##### Version 4.6.0
> 1.platform api support  

##### Version 4.5.3
> 1.fix mc folderkey

##### Version 4.5.2
> 1.fix folder nil

##### Version 4.5.1
> 1.fix 播单count更新

##### Version 4.5.0
> 1.收藏夹信息存储换用mc

##### Version 4.4.1
> 1.add redis expire

##### Version 4.4.0
> 1.support playlist

##### Version 4.3.6
> 1.archiveRPC archive3

##### Version 4.3.5
> 1.databus stats

##### Version 4.3.3
> 1.aidsByFid sql force index

##### Version 4.3.2
> 1.fix arc stats 

##### Version 4.3.1
> 1.in条件语句去掉prepare

##### Version 4.3.0
> 1.一键清理失效视频

##### Version 4.2.4
> 1.fix article map nil

##### Version 4.2.3
> 1.fix mid negative number

##### Version 4.2.0
> 1.迁移大仓库

##### Version 4.1.1
> 1.fix chan close

##### Version 4.1.0
> 1.job kafka改用databus

##### Version 4.0.3
> 1.200收藏加1硬币

#### Version 4.0.2
> 1.fix 收藏计数

##### Version 4.0.1
> 1.用户收藏标志缓存写入

##### Version 4.0.0
> 1.收藏平台化

##### Version 3.3.5
> 1.修复aid计数

##### Version 3.3.4
> 1.修复多chan处理mid  

##### Version 3.3.3
> 1.修复收藏计数bug  

##### Version 3.3.2
> 1.去掉计数相关hbase,dedeDB,kafka  

##### Version 3.3.1
> 1.monitor ping

##### Version 3.3.0
> 1.升级vendor

##### Version 3.2.6
> 1.修复bug:收藏视频批量移动后缓存不更新

##### Version 3.2.5 
> 1.redis bit添加用户是否有收藏缓存  

##### Version 3.2.4
> 1.添加默认收藏夹缓存  

##### Version 3.2.3
> 1.升级vendor  

##### Version 3.2.2
> 1.收藏计数双写到databus  

##### Version 3.2.1  
> 1.修复删除收藏夹收藏视屏缓存错误  

##### Version 3.2.0 
> 1.修复hbase计数  
> 2.收藏夹维护aid全部收藏数  

##### Version 3.1.0 
> 1.获取视屏对应的所有收藏夹  

##### Version 3.0.0  
> 1.添加govendor支持    
> 2.修复缓存expire  

##### Version 2.0.2 
> 1.修复删除收藏硬币扣除  

##### Version 2.0.1
> 1.添加job监控  
> 2.增加更新老稿件表收藏数  

##### Version 2.0.0
> 1.收藏夹job重构  
