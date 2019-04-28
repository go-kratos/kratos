# push-job

### v2.0.5
1. waitgroup

### v2.0.4
1. 修改定时delete的任务

### v2.0.3
1. fix close addTaskCh

### v2.0.2
1. mid文件目录增加目录层级，防止同一个目录下文件数过多

### v2.0.1
1. 优化全量推送建任务速度

### v2.0.0
1. using push grpc
2. remove abtest code

### v1.8.0
1. 删除老上报

### v1.7.3
1. fix 老上报 HD 归类错误

### v1.7.2
1. stop abtest

### v1.7.1
1. 优化生成abtest池子

### v1.7.0
1. abtest

### v1.6.0
1. 支持图片推送字段

### v1.5.0
1. 对接数据平台用户画像

### v1.4.2
1. 调整callback写入速率

### v1.4.1
1. 优化-批量写入token缓存

### v1.4.0
1. 刷新缓存的时候，增加token级别的缓存

### v1.3.0
1. 推送服务切换到push-service

### v1.2.1
1. 使用 go-common/env

### v1.2.0
1. 接bm

### v1.1.9
1. 迁移model至push-service

### v1.1.8
1. 定期删除 task

### v1.1.7
1. 修复 write closed chan

### v1.1.6
1. remove consumer uninstall mi token

### v1.1.5
1. 刷新token后释放内存

### v1.1.4
1. 项目迁移到main目录下

### v1.1.3
1. 上报增加设备信息
2. 去掉获取小米推送结果

### v1.1.2
1. 定时删除callback数据
2. 定时刷新上报缓存

### v1.1.1
1. 更改报警方式为企业微信

### v1.1.0
1. add push callback

### v1.0.15
1. 不外理新版本的老的上报数据
2. 更改prom为go-common中的对象

### v1.0.14
1. fix kafka consume

### v1.0.13
1. fix kafka consume

### v1.0.12
1. fix kafka commit

### v1.0.11
1. 优化 kafka 消费的写法

### v1.0.10
1. fix kafka consumer

### v1.0.9
1. pull push result

### v1.0.4
1. 升级 trace / http client auto sign

### v1.0.0
1. 项目初始化
