#### dynamic-service

##### 项目简介
> 1.提供稿件动态服务

##### 编译环境
> 请使用golang v1.7.x以上版本编译执行。  

##### 依赖包 
> 1.公共包go-common  

##### 编译执行
> 在主目录执行go build。   
> 编译后可执行 ./cmd -conf dynamic-service-test.toml 使用项目本地配置文件启动服务。  
> 也可执行 ./dynamic-service  -conf_appid=dynamic-service -conf_version=shsb-server-1  -conf_host=config.bilibili.co -conf_path=/data/conf/dynamic-service -conf_env=10 -conf_token=I4oCdH5cWAfP8wHqjEnbOu0qnB1miyEN 使用配置中心测试环境配置启动服务，如无法启动，可检查token是否正确。  

##### 特别说明  
> http接口文档可参考 http://info.bilibili.co/pages/viewpage.action?pageId=2493593