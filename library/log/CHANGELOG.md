### go-common/log

### v1.16.1
> 1.修复文件 handler 不能删除已有老文件问题

### v1.16
> 1. 增加是否批量写日志的判断

### v1.15
> 1.log 如果已定义了 source 字段就不再重新获取
> 2.合并了部分代码
> 3.添加了 log.no-agent flag 可以强制关闭 log-agent

### v1.15
> 1.修复向log agent批量写日志的bug

### v1.14
> 1. 实现文件日志，移除 log4go 依赖

### v1.13.1
> 1.修复infov，参数丢失...问题

### v1.13.2
> 1.优化log agent的性能

### v1.13.1
> 1.infoc support mirror request

### v1.13.0
> 1.support mirror request

### v1.12.2
> 1.fix logw bug and add test

### v1.12.1
> 1.add infoc write timeout  

### v1.12.0
> 1.add log doc  

### v1.11.0
> 1.use library/conf/dsn parse  

### v1.10.1
> 1.修复pattern中获取当前行信息的错误，之前的设置在go 1.9中会获取到错误的行，在go 1.10中是正确的(ノへ￣、)。

### v1.10.0
> 1.log error report to prometheus  

### v1.9.0
> 1.log dsn

### v1.8.4
> 1.优化文件日志输出内容

### v1.8.4
> 1.library/log enhancement  

### v1.8.4
> 1.infoc新增超过最大重试次数的日志

### v1.8.3
> 1.修改report包，从log协议改成databus

### v1.8.2
> 1.fixed funcname  

### v1.8.1
> 1.新增report包，支持上报行为日志

### v1.8.0
> 1. 优化log.D pool

#### v1.7.1
> 1. agent log enhance   

#### v1.7.0
> 1. update infoc sdk    

#### v1.6.3
> 1. add zone info

#### v1.6.2
> 1. update verbose doc and stdout log

#### v1.6.1
> 1. 更改默认日志发送等待时间

#### v1.6.0
> 1. add stdout log handler

#### v1.5.3
> 1. close log nil check

#### v1.5.2
> 1. 优先使用Caster环境变量

#### v1.5.1
> 1. 支持Caster环境变量

#### v1.5.0
> 1. agent批量写日志

#### v1.4.0
> 1. 移除log.XXXContext()

#### v1.3.0
> 1. 添加verbose log

#### v1.2.3
> 1. 修复agent退出，连接未关闭；
> 2. 修复conn重连，导致的饥饿无法退出；

#### v1.2.2
> 1. 完善net/http, net/rpc日志 

#### v1.2.1
> 1. 修复log handler未初始化panic 

#### v1.2.0
> 1. 结构化日志 

#### v1.1.0
> 1. 剔除elk, synclog  

#### v1.0.0
> 1. 初始化项目，更新依赖
