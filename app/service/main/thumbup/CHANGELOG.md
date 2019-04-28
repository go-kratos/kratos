### thumbup-service

##### Version 1.20.4
> 1. 使用zrange
##### Version 1.20.3
> 1. 异步发消息
##### Version 1.20.2
> 1. 优化 has like接口 查询500条
##### Version 1.20.1
> 1. has like接口返回点赞时间

##### Version 1.20.0
> 1. use tidb

##### Version 1.19.1
> 1. fixed panic bug

##### Version 1.19.0
> 1. 增加grpc接口

##### Version 1.18.0
> 1. 支持切换 tidb 写
> 2. 增加正序查询压测接口
> 3. 增加单查压测接口

##### Version 1.17.0
> 1. 增加压测接口

##### Version 1.16.5
> 1. 恢复tidb
##### Version 1.16.3
> 1. 停掉tidb 双写
##### Version 1.16.2
> 1. 写流量迁移到mysql
##### Version 1.16.1
> 1. 增加读流量来源切换
##### Version 1.16.0
> 1. 拜年祭需求

##### Version 1.15.3
> 1. 写接口增加mysql降级

##### Version 1.15.2
> 1. 增加up主mid更新接口验证
##### Version 1.15.1
> 1. update_upmids接口改为只更新不增加
##### Version 1.15.0
> 1. 增加update_upmids接口
##### Version 1.14.0
> 1. tidb双写
##### Version 1.13.0
> 1. 从tidb获取数据
##### Version 1.12.1
> 1. item忽略mtime索引
##### Version 1.12.0
> 1. databus增加up_mid
##### Version 1.11.1
> 1. 修复pn ps 为0的问题
##### Version 1.11.0
> 1. 新增点赞返回计数信息的RPC接口

##### Version 1.10.0
> 1.  点赞时更新up mid
##### Version 1.9.0
> 1.  增加被点赞人字段
##### Version 1.8.3
> 1.  add ut
##### Version 1.8.2
> 1.  使用新的rpc server
##### Version 1.8.1
> 1.  去除remote ip
##### Version 1.8.0
> 1.  db主从分离

##### Version 1.7.0
> 1.  use new verify
##### Version 1.6.0
> 1.  计数流增加mid字段
##### Version 1.5.2
> 1.  点赞人列表增加mid字段
##### Version 1.5.1
> 1. 修复回源率统计

##### Version 1.5.0
> 1. 支持修改点赞/踩的值

##### Version 1.4.4
> 1. 升级tools/cache

##### Version 1.4.3
> 1. add register

##### Version 1.4.2
> 1. use bm

##### Version 1.4.1
> 1. 修复单飞问题

##### Version 1.4.0
> 1.缓存使用缓存工具重构 
> 2.去掉article缓存包的依赖
> 3.业务表增加用户点赞列表 减少不必要的查询

##### Version 1.3.2
> 1.点赞bug修复

##### Version 1.3.1
> 1.点赞计数聚合写入

##### Version 1.3.0
> 1.发送计数databus信息

##### Version 1.2.9
> 1.发送稿件databus信息

##### Version 1.2.8
> 1.重复点赞的时候报错
> 2.增加点赞总数接口

##### Version 1.2.7
> 1.修复类型转换错误

##### Version 1.2.5
> 1.支持跨业务批量查询点赞

##### Version 1.2.4
> 1.重复点赞的时候更新时间

##### Version 1.2.3
> 1.fix like check bug

##### Version 1.2.2
> 1.stats 接口优化 返回空数据

##### Version 1.2.1
> 1.stats 接口优化 && fix bug

##### Version 1.2.0
> 1.stats 接口优化空缓存 && 去掉json解析

##### Version 1.1.1
> 1.identity bugfix

##### Version 1.1.0
> 1.合并interface到service
> 2.是否点赞接口批量查询优化    

##### Version 1.0.0
> 1.点赞项目初始化  
