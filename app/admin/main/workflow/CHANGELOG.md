#### 工单流审核后台接口

##### Version 1.6.6
> 1.account服务切到grpc

##### Version 1.6.5
> 1.account服务切回到gorpc

##### Version 1.6.4
> 1.group列表支持是否优先匹配关键字
> 2.支持搜索举报对象的第一个用户举报tag id

##### Version 1.6.3
> 1.接入配置中心sdk paladin

##### Version 1.6.2
> 1.v1举报接口下线

##### Version 1.6.1
> 1.评论举报增加火鸟来源, 站内信推送过滤

##### Version 1.6.0
> 1.长短评接入v2

##### Version 1.5.2
> 1.封禁发送站内信

##### Version 1.5.1
> 1.优化用户tag统计

##### Version 1.5.0
> 1.稿件投诉迁移v3

##### Version 1.4.5
> 1.站内信title长度截断

##### Version 1.4.4
> 1.v3/group/set 操作校验 rid, 有效操作同步写rid

##### Version 1.4.3
> 1.评论业务站内信内容优化 

##### Version 1.4.2
> 1.增加封禁类型（恶意冒充他人）

##### Version 1.4.1
> 1.封禁操作上报小黑屋

##### Version 1.4.0
> 1.对接评论举报

##### Version 1.3.7
> 1.修复manager_v4请求修改申诉状态

##### Version 1.3.6
> 1.修复站内信推送内容过长

##### Version 1.3.5
> 1.修复es排序参数

##### Version 1.3.4
> 1.工单列表增加rid二次验证

##### Version 1.3.3
> 2.优化日志上报

##### Version 1.3.2
> 1.字幕举报对接

##### Version 1.3.1
> 1.修复oid溢出问题
> 2.tid参数允许0值

##### Version 1.3.0
> 1.评论举报对接  
> 2.model代码优化   
> 3.ctx取admin_name上报行为日志    
> 4.不再依赖旧的es appid  
> 5.新的es appid迁移databus

##### Version 1.2.9
> 1.去掉操作日志写db

##### Version 1.2.8
> 1.工单列表meta_data字段通用化

##### Version 1.2.7
> 1.account-service 切换rpc请求

##### Version 1.2.6
> 1.全局通过 gid 查询 business object 
> 2.操作日志查询行为日志  
> 3.修复反馈详情列表    
> 4.优化操作日志查询    
> 5.移除本地tag依赖,迁移到manager-admin

##### Version 1.2.5
> 1.修复errgroup内存泄露  
> 2.优化es查询工单详情  
> 3.修复工单列表日志显示     
> 4.修复log,tag列表 
> 5.es client 配置化
> 6.请求列表后访问db check工单状态并过滤

##### Version 1.2.4
> 1.修复logproc

##### Version 1.2.3
> 1.优化调用账号api的错误日志

##### Version 1.2.2
> 1.修复参数错误

##### Version 1.2.1
> 1.优化日志上报

##### Version 1.2.0
> 1.显示用户 special tag, 粉丝数  
> 2.支持分词查询回复内容, 搜索受理人 username  
> 3.修复点评展示的 typename        
> 4.消息推送记录行为日志

##### Version 1.1.1
> 1.不依赖 php  
> 2.修复api权限验证   
> 3.修复列表翻页和排序   
> 4.老api下线   
> 5.log 接入行为日志  
> 6.修复点评展示的 typename     

##### Version 1.1.0
> 1.工作台 v1.1 基础功能   
> 2.稿件申诉回复进入工作台     
> 3.引入 redis       
> 4.规范路由    
> 5.部分api不再依赖php    
> 6.规范api返回格式   
> 7.修复延时导致回复状态错误     
> 8.修复工作台分页展示   
> 9.前端与 manager 解耦,直接访问 workflow    
> 10.workflow tag默认状态为开启   

##### Version 1.0.1
> 1.读取新的状态字段    
> 2.复合的状态字段不返回到前端   
> 3.修复对group操作的状态双写     
> 4.移除双写操作，只使用复合的状态字段   
> 5.优化group操作，操作group时异步对challenge操作     
> 6.修复审核操作日志丢失    

##### Version 1.0.0
> 1.增加工作台功能     
> 2.修复工作台工单展示       
> 3.接入memcache  
> 4.处理工单同时更新group表  
> 5.整合修复项目路由    
> 6.修复状态双写  

##### Version 0.2.1
> 1.增加小黑屋申诉业务

##### Version 0.2.0
> 1.投诉通用化   
> 2.给出申诉的反馈处理时间

##### Version 0.1.0
> 1.添加申诉反馈单的功能

##### Version 0.0.2
> 1.添加申诉执法单功能

##### Version 0.0.1
> 1.工单流项目初版