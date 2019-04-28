# push-interface

### v2.3.3
1. 支持漫画

### v2.3.1
1. 支持国际版上报

### v2.3.0
1. using push grpc

### v2.2.2
1. 更换极光callback url
2. revert v2.2.1

### v2.2.1
1. 修复iOS上报时对设备的判断

### v2.2.0
1. 增加接收旧版本APP上报token的接口

### v2.1.2
1. 根据platform判断使用何种推送SDK改为匹配前缀方式

### v2.1.1
1. 支持极光批量回执

### v2.1.0
1. 添加极光回执

### v2.0.2
1. 小米callback加barStatus

### v2.0.1
1. use bm auth middleware

### v2.0.0
1. 去除推送能力

### v1.9.1
1. report接口token为空时返回成功

### v1.9.0
1. remove RPC service

### v1.8.16
1. iOS支持image

### v1.8.15
1. 使用 go-common/env

### v1.8.14
1. 迁移model至push-service

### v1.8.13
1. 更改jobName生成规则

### v1.8.12
1. 小米接入VIP线路

### v1.8.11
1. 第三方依赖的http的prom错误码写成err

### v1.8.10
1. 加内网设置用户开关接口

### v1.8.9
1. 获取用户开关配置接口返回值调整

### v1.8.8
1. 全量推送建任务移到 admin 项目 

### v1.8.7
1. fix bug

### v1.8.6
1. add setSetting rpc method

### v1.8.5
1. add ecode

### v1.8.4
1. 添加用户上报稿件和直播的开关设置

### v1.8.3
1. business add silent time & push count limit
2. add platform midValid into progress

### v1.8.2
1. add task接口加签名
2. 回执时读取report时包含已删除的信息

### v1.8.1
1. 增加任务时加 uuid 验证

### v1.8.0
1. 接入极光
2. 回调中清理华为token

### v1.7.5
1. 回调中加入品牌

### v1.7.4
1. 在任务信息中加入 job name

### v1.7.3
1. 统计中加品牌计数

### v1.7.2
1. 推送消息过滤

### v1.7.1
1. 完善后台token推送

### v1.7.0
1. 支持后台任务推送

### v1.6.0
1. 上报接口记录Android设备信息(brand/model/os version)
2. 去掉获取小米推送结果

### v1.5.16
1. 华为流控后稍后重试

### v1.5.15
1. 华为流控判断
2. http参数错误日志级别改成warn 

### v1.5.14
1. 上报接口支持oppo

### v1.5.13
1. aps中tid换成device token

### v1.5.12
1. get buvid from header at click callback

### v1.5.11
1. fix ios space scheme

### v1.5.10
1. 上报缓存中添加token对应的id

### v1.5.9
1. fix service close chan

### v1.5.8
1. 改oppo为单个token推送

### v1.5.7
1. 支持全量推送

### v1.5.5
1. 优化oppo推送

### v1.5.4
1. 消息默认不响铃、不振动

### v1.5.3
1. http接口异步处理时加返回值

### v1.5.2
1. single push接口改成POST

### v1.5.1
1. 加 revover task
2. 添加任务时加 job name

### v1.5.0
1. 接入bm

### v1.4.11
1. midValid值更精确

### v1.4.10
1. 实现oppo callback
2. decode点击回执中的token

### v1.4.9
1. 查上报缓存失败不回源

### v1.4.8
1. 添加用户上报缓存RPC

### v1.4.7
1. 修复callback platform错误

### v1.4.6
1. 增加添加上报缓存RPC

### v1.4.5
1. 优化callback extra字段

### v1.4.4
1. 将送达回执的token状态写入DB

### v1.4.3
1. 优化推送计数代码

### v1.4.2
1. 对接小米卸载token接口

### v1.4.1
1. 接小米送达回执

### v1.4.0
1. 处理送达和点击回执，入库

### v1.3.10
1. 客户端点击回执
2. 程序停止时结束当前在执行的任务再退出

### v1.3.9
1. add task fix sql

### v1.3.8
1. 新建task支持group

### v1.3.7
1. 拉取小米推送结果的天数可配置
2. 去除单推测试接口mid白名单

### v1.3.6
1. 小米regid回调
2. 优化获取小米推送结果的方式

### v1.3.5
1. 实现华为回调接口

### v1.3.4
1. 修复小米老的108位token的上报问题

### v1.3.3
1. 优化新增上报缓存

### v1.3.2
1. 添加推送任务时支持定时

### v1.3.1
1. 优化批量读取上报数据

### v1.3.0
1. 接OPPO
2. 支持按token批量推送

### v1.2.20
1. fix upload name

### v1.2.13
1. 创作姬加入H5推送协议
2. 小米推送结果拉取时间改成7天

### v1.2.12
1. 优化的推送结果统计

### v1.2.11
1. 添加华为透传消息

### v1.2.10
1. 优化华为授权

### v1.2.9
1. 华为无效token的处理
2. add task支持版本判断
3. 优化推送结果统计

### v1.2.8
1. 推送接口支持pass_through参数
2. 接入华为推送

### v1.2.7
1. 调整prom

### v1.2.6
1. 补注释，处理golint报错的问题

### v1.2.5
1. 获取小米无效token加日志记录调用详情

### v1.2.4
1. 没有token的mid用小米推

### v1.2.3
1. add apns-collapse-id

### v1.2.2
1. 加入按token推送的接口供测试使用
2. 调整声音和振动http参数的校验方式
3. 改变 temp task id 长度为8位

### v1.2.1
1. 修复免打扰时间计算的bug

### v1.2.0
1. pull mipush result

### v1.1.2
1. mipush add jobkey

### v1.1.1
1. 轮询token进行无效删除
2. apns换超时机制

### v1.1.0
1. iOS push加context withtimeout
2. 加入一个利用缓存中已有mid批量测试的接口
3. 免打扰时间段内不处理推送任务
4. 修复sql in查询bug (xstr.JoinInts 不能作为占位符的值)

### v1.0.2
1. push/single接口调整参数校验
2. apns是否走代理改成配置项
3. 推送时把mid和token打印出来，方便项目上线前期排查问题
4. 直播iPad支持推送协议

### v1.0.1
1. device_token为空不报-400
2. 加创作姬推送协议

### v1.0.0
1. 项目初始化
