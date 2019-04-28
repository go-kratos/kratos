#### favorite

### Version v7.5.1
> 1. close senstive

### Version v7.5.0
> 1. 播单2.0

### Version v7.4.3
> 1. 收藏夹权限问题修复

### Version v7.4.1
> 1. dao 目录迁移

### Version v7.4.0
> 1. grpc client v1迁移到api

### Version v7.3.0
> 1. 播单支持管理员修改

### Version v7.2.0
> 1. medialist fix bug

### Version v7.0.1
> 1. medialist

### Version v7.0.0
> 1. medialist

### Version v6.0.3
> 1. 最近三条收藏 force index

### Version v6.0.1
> 1. 并发读db获取最近三条收藏

### Version v6.0.0
> 1. 收藏夹停止双写

### Version v4.11.7
> 1. fix rank service err

### Version v4.11.5
> 1. fix grpc server err

### Version v4.11.4
> 1. recentfavs bit

### Version v4.11.3
> 1. fix del bug

### Version 4.11.1
> 1. fix rank pagesize

### Version 4.11.0
> 1. riot and rank search

### Version 4.10.7
> 1. add db rows.error

### Version 4.10.6
> 1. change dir to api

### Version 4.9.15
> 1. upgrade to grpc

### Version 4.9.14
> 1. ut dao test

### Version 4.9.13
> 1. meta data ip

### Version 4.9.11
> 1. change indentify to grpc

### Version 4.9.10
> 1. rm video service

### Version 4.9.8
> 1. common config

### Version 4.9.7
> 1. fix two same folder

### Version 4.9.6
> 1. batchOids ylf db

### Version 4.9.5
> 1. fix batchOids unFavedBit

### Version 4.9.4
> 1. fix oidsCount mc

### Version 4.9.3
> 1. migrate relation redis

### Version 4.9.2
> 1. add api oids count
> 2. add api get 1000 favorites 
> 3. fix marshal panic 

### Version 4.9.1
> 1. fix favoreds null

### Version 4.9.0
> 1. update infoc sdk

##### Version 4.8.13
> 1.add oid count api

##### Version 4.8.12
> 1.fix http isFavoreds

##### Version 4.8.11
> 1.fix index ix_fid_state_mtime

##### Version 4.8.10
> 1.force index uk_fid_oid_type

##### Version 4.8.9
> 1.add read mysql for push service

##### Version 4.8.9
> 1.add read mysql for push service

##### Version 4.8.7
> 1.fix relation recent oids cache

##### Version 4.8.6
> 1.优化relation cache

##### Version 4.8.4
> 1.增加 userFolder rpc

##### Version 4.8.3
> 1.增加register

##### Version 4.8.2
> 1.add users api for push service
> 2.add fav and del fav move to job

##### Version 4.8.1
> 1.add isFavs rpc

##### Version 4.8.0
> 1.migrate main path

##### Version 4.7.4
> 1.使用account-service v7  

##### Version 4.7.3
> 1.fix cleanvideo http response

##### Version 4.7.2
> 1.add bm verfiy setMid

##### Version 4.7.1
> 1.fix batch archives3 

##### Version 4.7.0
> 1.migrate blademaster

##### Version 4.6.11
> 1.fix fav dao test

##### Version 4.6.10
> 1.delete statsd

##### Version 4.6.9
> 1.fix save cache context

##### Version 4.6.8
> 1.fix is faved

##### Version 4.6.7
> 1.fav bit redis key mid reverse

##### Version 4.6.6
> 1.迁移 accRPC profile2 myinfo 到 userinfo

##### Version 4.6.5
> 1.fix folder cover by removed unnormal archive's cover

##### Version 4.6.4
> 1.get archives by pages from arc rpc

##### Version 4.6.2
> 1.topic fav read platform table

##### Version 4.6.1
> 1.fix topic fav double write

##### Version 4.6.0
> 1.platform api support 

##### Version 4.5.7
> 1.fix mysql clsoe  

##### Version 4.5.6
> 1.fix err1 for double write topic to platform

##### Version 4.5.5
> 1.remove folder description \n filtering rule
> 2.fix add fav topic verify
> 3.double write topic to platform

##### Version 4.5.2
> 1.change ecode NoLogin to RequestErr
> 2.httpClient topic upgrade
> 3.remove innerRouter and localRouter

##### Version 4.5.1
> 1.fix get default folder

##### Version 4.5.0
> 1.收藏夹信息存储换用mc

##### Version 4.4.3
> 1.add redis expire

##### Version 4.4.2
> 1.fix invalid cache  

##### Version 4.4.1
> 1.folder count 接口

##### Version 4.4.0
> 1.support playlist

##### Version 4.3.11
> 1.archiveRPC archive3

##### Version 4.3.8
> 1.http rounter config

##### Version 4.3.7
> 1.支持supervisor
> 2.添加收藏 取消收藏放到service

##### Version 4.3.5
> 1.httpclient 去掉appkey

##### Version 4.3.2
> 1.全站实名制支持

##### Version 4.3.1
> 1.账号封禁下的收藏操作提示信息修正

##### Version 4.3.0
> 1.一键清理失效视频

##### Version 4.2.11
> 1.优化添加收藏的反作弊上报

##### Version 4.2.10
> 1.get three aids sql force index

##### Version 4.2.8
> 1.fix rename 默认收藏夹
> 2.fix move arcs sql deadlock

##### Version 4.2.7
> 1.fix folder sort index out of range

##### Version 4.2.6
> 1.fix sort err

##### Version 4.2.5
> 1.infoc nil fix
> 2.文章favmdl包名修改

##### Version 4.2.1
> 1.迁移RPC Client
> 2.整合model

##### Version 4.2.0
> 1.接入大仓库

##### Version 4.1.2
> 1.fix 收藏平台type范围

##### Version 4.1.1
> 1.上报日志到AI

##### Version 4.1.0
> 1.http gotest
> 2.接入log agent
> 3.敏感词接入

##### Version 4.0.4
> 1.antispam fix

##### Version 4.0.3
> 1.prome withtimer

##### Version 4.0.2
> 1.fix fids outof range

##### Version 4.0.1
> 1.search fix

##### Version 4.0.0
> 1.收藏平台化

##### Version 3.8.8
> 1.prometheus

##### Version 3.8.7
> 1.fix antispam path

##### Version 3.8.6
> 1.接入防刷 antispam

##### Version 3.8.5
> 1.添加收藏和删除收藏DB操作移至job

##### Version 3.8.4
> 1.add supervision

##### Version 3.8.3
> 1.话题收藏兼容移动端https

##### Version 3.8.2
> 1.视频收藏接入新版搜索
> 2.添加搜索敏感词errcode

##### Version 3.8.1
> 1.视频收藏接入新版搜索

##### Version 3.8.0
> 1.代码reset
> 2.修复访客状态下私密收藏夹展示问题

##### Version 3.7.5
> 1.修复访客状态下私密收藏夹展示问题

##### Version 3.7.3
> 1.fix 私密收藏夹展示问题

##### Version 3.7.2
> 1.fix私密收藏夹展示问题

##### Version 3.7.1
> 1.最近收藏结构优化

##### Version 3.7.0
> 1.视频列表和最近收藏接入新搜索

##### Version 3.6.5
> 1.CSRF日志  

##### Version 3.6.4
> 1.CSRF日志  

##### Version 3.6.3
> 1.添加CSRF日志  

##### Version 3.6.2
> 1.使用CSRF配置  

##### Version 3.6.1
> 1.新doker镜像  

##### Version 3.6.0
> 1.升级go-common 
> 2.接入新的配置中心

##### Version 3.5.4
> 1.新增话题收藏状态
> 2.话题信息调用活动平台接口

##### Version 3.5.3
> 1.fix 默认排序中默认收藏夹不在首位

##### Version 3.5.2
> 1.tw升级新版本，需要重新打包上线

##### Version 3.5.1
> 1.fix redis setex 断言错误

##### Version 3.5.0
> 1.收藏夹自定义排序

##### Version 3.4.6
> 1.升级go-buisiness
> 2.MonitorPing检测redis和mysql

##### Version 3.4.5
> 1.收藏夹视频批量移动逻辑改成覆盖
> 2.优化 更新收藏夹是否公开状态接口

##### Version 3.4.1
> 1.回滚http接口支持内外网流量

##### Version 3.4.0
> 1.1.更新http接口支持内外网流量

##### Version 3.3.11
> 1.升级go-commom,go-business，修复mo状态码统计  

##### Version 3.3.10
> 1.修复family  

##### Version 3.3.9
> 1.查询DB，map初始化

##### Version 3.3.7
> 1.增肌批量查询收藏视频接口
> 2.升级go-common,golang,go-business

##### Version 3.3.6
> 1.增加返回收藏夹数 

##### Version 3.3.5
> 1.重新构建 

##### Version 3.3.4
> 1.修复显示私有收藏夹  

##### Version 3.3.3
> 1.default 增加是否有收藏bit判断  

##### Version 3.3.2
> 1.redis bit添加用户是否有收藏缓存  
> 2.接入配置中心  

##### Version 3.3.1
> 1.初始化默认收藏夹缓存  

##### Version 3.3.0
> 1.更新identity  
> 2.增加稿件状态判断  
> 3.cover pipeline  

##### Version 3.2.2
> 1.更新vendor依赖  

##### Version 3.2.1
> 1.更新vendor依赖  

##### Version 3.2.0
> 1.更新基础库依赖  
> 2.删除冗余error  

##### Version 3.1.1
> 1.修复收藏夹参数为空  

##### Version 3.1.0
> 1.获取视屏对应的所有收藏夹  
> 2.批量加入/删除 视频到多个收藏夹  
> 3.收藏夹封面增加缓存  

##### Version 3.0.0
> 1.接入trace_v2  
> 2.切换go-business  
> 3.接入ecode统一管理  
> 4.添加govendor支持  
> 5.修复html转义  

##### Version 2.2.3
> 1.支持添加收藏夹时指定属性  

##### Version 2.2.2
> 1.修复在没登录状态查看别人收藏夹视频列表mid为nil问题  

##### Version 2.2.1
> 1.修复cache，先进行expire  

##### Version 2.2.0
> 1.添加内部接口

##### Version 2.1.0
> 1.话题收藏迁移  

##### Version 2.0.0

> 1.收藏夹重构    
