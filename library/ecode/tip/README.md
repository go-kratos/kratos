# go-common/business/ecode

##### 项目简介
> 1. 提供所有请求的错误码及其错误信息管理，错误信息在管理平台配置，支持自动更新
> 2. 提供错误使用规范及文档，包含堆栈信息使用

##### 编译环境
> 1. 请只用golang v1.7.x以上版本编译执行。

##### 依赖包
> 1. 依赖github.com/pkg/errors，当前版本v0.8.0

##### 编译执行
> 1. 启动执行 ecode.Init(conf.Conf.Ecode)，初始化ecode 配置
> 2. 配置参考 http://info.bilibili.co/pages/viewpage.action?pageId=3684076 配置

##### 测试
> 1. 执行当前目录下所有测试文件，测试所有功能

##### 特别说明
> 1. common.go 里面保存所有业务的code码，当有新增加code码需求时请记得一定及时更新common,并在管理平台配置对应信息
> 2. 管理平台地址 http://apm-monitor.bilibili.co/#/codes/codeslist?name=all
> 3. 按部门给错误码分大段，部门内部按业务模块继续分段具体参考info地址: http://info.bilibili.co/pages/viewpage.action?pageId=5374316