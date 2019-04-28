#### web-goblin

##### 项目简介
> 1.提供web端附加功能接口

##### 编译环境
> 请使用golang v1.7.x以上版本编译执行。

##### 依赖包
> 1.公共包go-common

##### 编译执行
> 在主目录执行go build。
> 编译后可执行 ./cmd/cmd -conf web-goblin-test.toml 使用项目本地配置文件启动服务。
> 也可执行 ./web-goblin -conf_appid=web-goblin -conf_version=v2.1.0 -conf_host=172.16.33.134:9011 -conf_path=/data/conf/web-goblin -conf_env=10 -conf_token=*** 使用配置中心测试环境配置启动服务，如无法启动，可检查token是否正确。

##### 特别说明
> http接口文档可参考 http://info.bilibili.co/pages/viewpage.action?pageId=9843305
