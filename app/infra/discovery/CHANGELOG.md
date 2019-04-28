### discovery
#### Version 1.8.4
> 1.disocoery不区分环境

#### Version 1.8.3
> 1.修复自发现按配置zone进行同步

#### Version 1.8.2
> 1.添加zones种子节点同步

#### Version 1.8.1
> 1.修复自发现zone过滤其他机房  

#### Version 1.8.0 
1. 迁移infra
#### Version 1.7.1
> 1.添加自己发现discovery节点

#### Version 1.7.0
> 1.refactor set 
> 2.set metadata 
> 3.remove color 

#### Version 1.6.3
> 1.修复同一个hostname多个长连接时串ch的bug  

#### Version 1.6.2
> 1.删除兼容的http rpc字段  

#### Version 1.6.1
> 1.增加polling查看当前正在poll的host  

#### Version 1.6.0
> 1.添加批量fetch接口
> 2.移动weight到metadata  

#### Version 1.5.0
> 1.优化修改app存储结构  
> 2.去除treeid兼容  

#### Version 1.4.4
1. 更新时返回全部zone  

#### Version 1.4.3
1. 修复broadcast 没下发zoneinstances  

#### Version 1.4.2
1. 修复latest_timestamp 更新  

#### Verson 1.4.1 
1. 区分zone返回instances  

#### Version 1.4.0
1. 支持返回多个zone实例  

#### Version 1.3.1
1. 修复同时注册更新时间相同导致的304，修改精度为纳秒 

#### Version 1.3.0
1. 支持多注册中心数据同步  

#### Version 1.2.0
1. 增加treeid  
2. 增加polls 批量订阅  
3. 增加chan连接池  
4. 使用 "部门.分组" 为key减小锁粒度  

#### Version 1.1.0
1. 删除replication多余参数。
2. 完善单测，mock http请求覆盖replication
3. poll新增host字段，每次poll请求结束后删除conn

#### Version 1.0.0
1. 支持结点同步注册，心跳，取消请求。
2. 支持服务自我保护，防止网络闪断。
3. 支持长轮询更新服务变化。
4. 支持心跳同步过程中根据dirtytime纠正同步数据。
