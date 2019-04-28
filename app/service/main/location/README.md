#### location-service

##### 项目简介
> 1.提供IP信息查询的接口  
> 2.提供各种不同应用场景的，关于IP与稿件地区限制规则的查询接口  

##### 编译环境
> 请只用golang v1.7.x以上版本编译执行。  

##### 依赖包
> 1.公共包go-common  

##### 编译执行
> 在主目录执行go build。   
> 编译后可执行 ./location -conf location-example.toml 使用项目本地配置文件启动服务。  
> 也可执行 ./location -conf_appid=location-service -conf_version=v2.1.0 -conf_host=172.16.33.134:9011 -conf_path=/data/conf/location-service -conf_env=10 -conf_token=J8MnvX8woVZmdwQ78HIFe26QOgleuxd4 使用配置中心测试环境配置启动服务，如无法启动，可检查token是否正确。  

##### RPC测试
> 具体的测试内容可修改rpc/rpc_test.go文件。  
> 在rpc目录执行go test测试rpc接口。  

##### 特别说明
> 1.model目录可能会被其他项目引用，请谨慎请改并通知各方。  
> 2.http接口文档可参考 http://syncsvn.bilibili.co/platform/doc/blob/master/api/location-service/v3.md  