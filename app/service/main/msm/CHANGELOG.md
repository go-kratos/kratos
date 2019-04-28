### msm-service
#### Version 1.7.5
1. ecode添加繁体信息接口
#### Version 1.7.4
1. ut补全

#### Version 1.7.3
1. 替换为verify.Verify

#### Version 1.7.1
>1. msm服务使用独立鉴权代码，防止鉴权回路.

#### Version 1.7.0
>1. 增加RPC鉴权接口
>2. 删除无用接口和代码
>3. 使用BM框架
>4. 迁移至business/service/main目录

#### Version 1.6.0
>1. ecode读取db2数据库

#### Version 1.5.0
>1. 修改了ecode的逻辑，分为2个接口codes兼容老的，codes/2支持新的

#### Version 1.4.0
>1. 增加limit接口，提供限流规则

#### Version 1.3.0
>1. 增加rules接口，databus配置规则全量获取

#### Version 1.2.0
>1. 优化ecode接口，增量拉取，用更新时间作为版本号

#### Version 1.1.0
>1. 增加获取codes接口

#### Version 1.0.0
>1. rpc服务修改节点 权重，分组信息
>2. rpc服务所有节点信息查询，节点删除
>3. 对配置中心进行broadCast调用所有节点，目前有push和setToken方法。





