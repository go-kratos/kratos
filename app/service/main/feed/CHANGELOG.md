### Feed-Service

##### Version 1.13.2
> 1. add ut
##### Version 1.13.1
> 1. 使用新的rpc server

##### Version 1.13.0
> 1. 使用bm

##### Version 1.12.8
> 1.增加register

##### Version 1.12.7
> 1.使用account-service v7  

##### Version 1.12.6
> 1.查找线上问题

##### Version 1.12.5
> 1.使用pb重构 改为调用archives3接口   

##### Version 1.12.4
> 1.升级http client   

##### Version 1.12.3
> 1.换up主缓存key避免上个版本稿件服务导致的脏数据   

##### Version 1.12.1
> 1.代码优化, 去掉mc的time缓存改为存在redis score中 
> 2.未读数不再需要查archives2接口   
> 3.生成feed数据只查询需要的稿件数据 不再全量查询  
> 4.针对up主缓存 读取方式从expire改为ttl防止可能出现的缓存不一致现象

##### Version 1.11.9
> 1.合并大仓库修复commit丢失的问题  

##### Version 1.11.1
> 1.fix 稿件map不存在导致的panic  

##### Version 1.11.0
> 1. 合并大仓库
> 2. 修复丢失up主稿件问题

##### Version 1.10.10
> 1.fix 文章动态除0bug    

##### Version 1.10.9
> 1.fix 稿件状态判断    

##### Version 1.10.8
> 1.优化部分prom监控    

##### Version 1.10.7
> 1.更新获取up主最新投稿的接口    

##### Version 1.10.6
> 1.恢复up主过审稿件缓存    

##### Version 1.10.5
> 1.去掉up主过审稿件缓存    

##### Version 1.10.4
> 1.新增降级监控    

##### Version 1.10.3
> 1.文章动态支持动态不展示功能    

##### Version 1.10.2
> 1.fix文章动态翻页问题    

##### Version 1.10.1
> 1.支持文章动态    

##### Version 1.9.0
> 1.更新go-common、go-busines,以及更新关注接口    

##### Version 1.8.2
> 1.加回up主缓存    

##### Version 1.8.1
> 1.删除本地up主缓存    

##### Version 1.8.0
> 1.去掉番剧未读数API && 增加稿件转移rpc    

##### Version 1.7.4
> 1.多打一条日志    

##### Version 1.7.3
> 1.fixed fold bug    

