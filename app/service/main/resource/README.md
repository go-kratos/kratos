#### resource-service

##### 项目简介
> 1.为移动端和WEB端提供广告相关接口  
> 2.提供在线热数据配置和加载  

##### 编译环境
> 请只用golang v1.8.x以上版本编译执行。  

##### 依赖包
> 1.公共包go-common  

##### 编译执行
> 在主目录执行go build。  
> 编译后可执行 ./resource -conf resource-service-test.toml 使用项目本地配置文件启动服务。  
> 也可执行 -conf_appid=resource-service -conf_version=shsb-server-1 -conf_host=config.bilibili.co -conf_path=/data/conf/resource-service -conf_env=10 -conf_token=3DTSmKyHEwN6eYGcKhlIyqIM60yyyxQD 使用配置中心测试环境配置启动服务，如无法启动，可检查token是否正确。  

##### RPC测试
> 具体的测试内容可修改rpc/rpc_test.go文件。  
> 在rpc目录执行go test测试rpc接口。  

##### 特别说明
> 1.model目录可能会被其他项目引用，请谨慎请改并通知各方。  
> 2.http接口文档可参考 http://info.bilibili.co/pages/viewpage.action?pageId=3690413  