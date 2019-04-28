### Version 2.18.1
> 1. 修正批量获取用户佩戴勋章时的闭包问题

### Verson 2.18.0
> 1. 佩戴勋章区分来源

### Verson 2.17.2
> 1. 修复大会员无法佩戴不在背包里的大会员挂件

### Version 2.17.1
> 1. 记录邀请码购买时的 IPv6 地址
> 2. 佩戴挂件时不再检查是否为 VIP

### Version 2.16.3
> 1. 修复批量获取佩戴挂件时的空挂件现象
> 2. 修复获取单个用户佩戴挂件时的空挂件现象

### Version 2.16.0
> 1. 变更vip挂件装备逻辑  
> 2. vip 挂件直落装备随意切换在vip有效期内  

### Version 2.15.0
> 1. 在购买邀请码时添加扣除硬币时遗漏的IP参数

### Version 2.14.0
> 1. dao层的ut使用新规范完善  

### Version 2.13.1
> 1. 删除block-service import   

### Version 2.13.0
> 1. 新增挂件红点提醒   
> 2. context.todo -> context.backround  

### Version 2.12.3
> 1. fix medal popup bug 

### Version 2.12.2
> 1. blockinfo use memberRPC 

### Version 2.12.1
> 1. fix 分享group rpc 接口  

### Version 2.12.0
> 1. 新增分享group rpc 接口 

### Version 2.11.1
> 1. fix message

### Version 2.11.0
> 1. 改为 metadata RemoteIP

### Version 2.10.4
> 1. 迁移identify-->verify

### Version 2.10.3
> 1.use account notify

### Version 2.10.2
> 1.change medal owners mc key 

### Version 2.10.1
> 1.ignore paymentSystem err status callback  

### Version 2.10.0
> 1.usersuit load gid ref pid    

### Version 2.9.3
> 1.usersuit notify  

### Version 2.9.2
> 1.fix medal user  

### Version 2.9.1
> 1.refined medal all
  
### Version 2.9.0
> 1.add mid pkg and equip local cache

### Version 2.8.6
> 1.fix redis expire time

### Version 2.8.5
> 1.add medal multi activated rpc api  

### Version 2.8.4
> 1.add medal activated rpc  

### Version 2.8.3
> 1.refined all info  

### Version 2.8.1
> 1.edit allinfo

### Version 2.8.0
> 1.incr chan size  

### Version 2.7.0
> 1.refined model and pb  

### Version 2.6.8
> 1.refined code  

### Version 2.6.7
> 1.fix level  

### Version 2.6.6
> 1.refined medal logic  

### Version 2.6.5
> 1.refined medal activated logic  

### Version 2.6.4
> 1.add medal activated info api  

### Version 2.6.3
> 1.medal info add level  

### Version 2.6.2
> 1.update sql

### Version 2.6.1
> 1.fix 勋章cond字段bug

### Version 2.6.0
> 1.勋章重构 

### Version 2.5.0
> 1.bm

### Version 2.4.3
> 1.remove useless code

### Version 2.4.1
> 1.add register

### Version 2.4.0
> 1.add mutiple grant pandent rpc  

### Version 2.4.0 
> 1.usersuit medal 勋章重构 

### Version 2.3.4
> 1.默认返回的节操7000  

### Version 2.3.3
> 1.delete vipinfo proto   

### Version 2.3.2
> 1.add equip rpc   

### Version 2.3.1
> 1. 给member-service的rpc接口加default  

### Version 2.3.0
> 1.usersuit 去除member-service的myinfo接口  

### Version 2.2.13
> 1.fix panic bug  

#### Version 2.2.12
> 1.fix pendant not online display

#### Version 2.2.11
> 1.pendant and group local cache rebuild  

#### Version 2.2.10
> 1.修改邀请码接口 member-service  myInfo +  coin数据

#### Version 2.2.9
> 1.add empty cache for Equipment

#### Version 2.2.8
> 1.fix reids flush issue

#### Version 2.2.7
> 1.fix pendant purge cache

#### Version 2.2.6 
> 1.fix redis conn.close

#### Version 2.2.5 
> 1.背包缓存清理改为同步执行

#### Version 2.2.4 
> 1.装备挂件db写操作，背包信息从cache改为从db中取 

#### Version 2.2.3 
> 1.fix normal vip equip failed      

#### Version 2.2.2 
> 1.修复equipment单个接口nil      

#### Version 2.2.1 
> 1.equipment rpc 批量查询改为100.

#### Version 2.2.0 
> 1.增加 equipment 批量rpc接口  

#### Version 2.1.3
> 1.fix pendant expire bug

#### Version 2.1.2
> 1.fix panic

#### Version 2.1.1
> 1.pay

#### Version 2.1.0
> 1.change part ut to convey

> 2.remove unused conf item

#### Version 2.0.1
> 1.添加购买历史

#### Version 2.0.0
> 1.初始化挂件信息

#### Version 1.0.0
> 1.添加邀请码功能：包括购买邀请码、使用邀请码、获取邀请码状态
