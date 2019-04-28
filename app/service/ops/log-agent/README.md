# log-agent

##### 项目简介
> 1. 公司级日志收集组件，收集应用日志，然后发送给lancer
> 2. 支持本地缓存，确保数据可靠传输
> 3. 暴露http api，支持日志流式查看。
> 4. 支持日志采样
> 5. app级别日志流实时监控

##### 编译环境
> 1. 请使用golang v1.7.x以上版本编译执行。

##### 依赖包
> 1. 依赖github.com/prometheus/client_golang/prometheus

##### 编译执行
> 1. 编译后在物理机上启动log-agent即可。
> 2. 基于rms进程agent的管理与升级

##### 使用方式
> 1. [日志规范](http://info.bilibili.co/pages/viewpage.action?pageId=3674729)
> 2. [日志接入流程](http://info.bilibili.co/pages/viewpage.action?pageId=3678680)
