#### account-service
`account-service`主要是为主站各个业务提供一个大的cache和rpc接口。
* 聚合了账号多个业务(`passport` `account` `big` `pay` `relation-service`)的缓存

* 提供了rpc服务

1. 为没有rpc服务的业务(`passport` `account` `big` `pay`)封装rpc接口
2. 已rpc服务的业务(`relation-service`)提供接口封装，透传rpc服务

##### 依赖环境
Go 1.7.5或更高版本

##### API文档
TODO example code api
