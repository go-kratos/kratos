### member service

### Version 3.32.0
> 1. 官方认证增加最后提交审核时间
> 2. 官方认证是社会信用代码写入新表user_official_doc_addit中
> 3. 增加批量查询用户封禁详情接口
> 4. 删除 user_detail表相关代码

### Version 3.31.0
> 1. 官方认证增加官网地址与注册资金字段

### Version 3.30.0
> 1. 升级grpc 

### Version 3.29.0
> 1. 注册用户昵称先审后发

### Version 3.28.0
> 1. 查询用户 base 信息强制从主库查询
> 2. 增加实名认证非敏感信息的聚合接口
> 3. 增加通过身份证号查询用户 mid 的接口

### Version 3.26.0
> 1. 增加官方认证提交来源字段

### Version 3.25.0
> 1. 新增grpc接口 /BlockInfo 和 /BlockBatchInfo

### Version 3.24.2
> 1. RealnameCheck校验证件时统一转换为大写进行比较
> 2. 实名认证空数据

### Version 3.24.1
> 1. 优化芝麻认证

### Version 3.24.0
> 1. 增加芝麻实名认证

### Version 3.23.0
Added
- 新增grpc接口 /RealnameDetail

### Version 3.22.0
Changed
- 优化并新增 block source : bplus
- 优化并新增 bplus block notify msg

### Version 3.21.0
> 1. 优化 /realname/apply 接口
> 2. 优化 /realname/check 接口
> 3. 新增 /realname/tel/capture/check 接口，检查短信验证码

### Version 3.20.0
> 1. 增加 /realanme/check 接口，用户申诉验证实名证件

### Version 3.19.1
> 1. 优化block ut

### Version 3.19.0
> 1. remote ip from context metadata

### Version 3.18.4
> 1. fix block db tx

### Version 3.18.3
> 1. 增加 report manager 的初始化

### Version 3.18.2
> 1. 添加 block goRPC client 接口

### Version 3.18.1
> 1. 添加 block goRPC 接口

### Version 3.18.0
> 1. 嵌入 block service

### Version 3.17.1
> 1.添加member dao层test

### Version 3.17.1
> 1. 修正官方认证被取消时的清缓存通知

### Version 3.17.0
> 1. 删除 hbase
> 2. 删除无用配置
> 3. 修正 gorpc, bm 框架
> 4. 实名信息读主库
> 5. 规范缓存时间配置

### Version 3.16
> 1. 添加grpc支持

### Version 3.15.6
> 1. add notify sender mark.

###  Version 3.15.5
> 1. 再次修复节操参数校验

###  Version 3.15.4
> 1. 修复节操参数校验

###  Version 3.15.3
> 1. 增加节操修改rpc方法

###  Version 3.15.2
> 1. 修改拼写bug

###  Version 3.15.1
> 1. 删除老的参数绑定方式，改成新的参数绑定方式

#### Version 3.15.0
> 1. 改为 bm 内置参数解析工具

#### Version 3.14.4
> 1. 审核表增加extra 字段

#### Version 3.14.3
> 1. 修改/realname/adult 返回值

#### Version 3.14.2
> 1. 修改/realname/nonage --> /realname/adult

#### Version 3.14.1
> 1. 增加通过身份证号码，判断是否成年接口

#### Version 3.14.0
> 1. 根据节操行为日志撤回

#### Version 3.13.0
> 1. realname 提供实名详情接口

#### Version 3.12.4
> 1. 节操行为日志补充

#### Version 3.12.3
> 1. 经验日志输出 logid 字段

#### Version 3.12.2
> 1. 行为日志增加 log_id 字段

#### Version 3.12.1
> 1. 正常用户重复待审核状态归档

#### Version 3.12.0
> 1. 经验日志直接用行为日志

#### Version 3.11.1
> 1. 修改审核列表添加逻辑

#### Version 3.11.0
> 1. 节操日志进 report

#### Version 3.10.1
> 1. realname 移动端优化

#### Version 3.10.0
> 1. 去除 reload 配置文件

#### Version 3.9.2
> 1. 修复realname配置

#### Version 3.9.1
> 1. fix AddPropertyReview 

#### Version 3.9.0
> 1. 迁移到 bm

#### Version 3.8.1
> 1. 节操日志rowkey加入uuid防止日志覆写

#### Version 3.8.0
> 1. 允许设置为空签名

#### Version 3.7.0
> 1. 加入添加监控用户和签名审核功能

#### Version 3.6.9
> 1. setBase del cache.

#### Version 3.6.8
> 1. update recover date when moral decrease.

#### Version 3.6.7
> 1. fix moral notice.

#### Version 3.6.6
> 1. fix moral log status.

#### Version 3.6.5
> 1. fix moral change sql

#### Version 3.6.4
> 1. 增加节操修改和批量修改接口

#### Version 3.6.3
> 1. restruct proto file

#### Version 3.6.2
> 1. optimize realname

#### Version 3.6.1
> 1. expLog limit 7 day

#### Version 3.6.0
> 1. fix official api and moralLog limit 7 day

#### Version 3.5.9
> 1. fix official api

#### Version 3.5.8
> 1. official submit

#### Version 3.5.7
> 1. add realname api

#### Version 3.5.6
> 1. 使用单个httpserver

#### Version 3.5.5
> 1. add base/set interface.

#### Version 3.5.4 
> 1. fix moral protobuf.

#### Version 3.5.3 
> 1. fix moral default moral.

#### Version 3.5.2
> 1. fix moral read.

#### Version 3.5.1
> 1. member个人信息修改删除自己缓存和通知业务方缓存
> 2. 删除没有使用的代码
> 3. 节操加缓存

#### Version 3.5.0
> 1. remove notify/close logic to account.  
> 2. super super clean code.  

#### Version 3.4.0
> 1. 增加member 个人信息修改接口

#### Version 3.3.0
> 1. 官方认证读新库

#### Version 3.2.4
> 1. fix moral.  

#### Version 3.2.3
> 1. remove add&set exp check&init logic.  
> 2. super clean code.

#### Version 3.2.2
> 1. add identify apply status check for notive v2.

#### Version 3.2.1
> 1. add moral&log interface.  
> 2. add moral&log rpc client.

#### Version 3.2.0
> 1. add NickUpdated and SetNickUpdated rpc interface.

用户账号会员系统的服务层。
#### Version 3.1.0
> 1. 添加实名认证通知
> 2. 添加实名认证RPC API

#### Version 3.0.0
> 1. update path

#### Version 2.10.3
> 1. using parent context in errgroup directly.

#### Version 2.10.2
> 1. add BaseExp rpc interface.

#### Version 2.10.1
> 1. remove login watch shareClick  
> 2. move project to main path  
> 3. clean code  

#### Version 2.10.0
> 1. 删除无用info card等  
> 2. 移除对account-service的依赖  
> 3. watch_ac typo

#### Version 2.9.6
> 1. 用户经验日志加 json tag

#### Version 2.9.5
> 1. 修复 exps 缓存出错

#### Version 2.9.4
> 1. MyInfo 继续使用 member 库中的经验值
> 2. remove statsd

#### Version 2.9.1-ooc
> 1. 因数据同步问题，MyInfo 直接用 account-java 的经验信息

#### Version 2.9.1
> 1. 修改经验接口参数ts名称改为ptime

#### Version 2.9.0
> 1. 提供base,detail,batchbase 的rpc接口
> 2. detail接口增加level信息

#### Version 2.8.1
> 1. 提供等级批量接口
> 2. 修改经验缓存为mc 

#### Version 2.7.1
> 1.add rpc service

#### Version 2.7.0
> 1.补充单元测试  

#### Version 2.6.2
> 1.hbase scan日志修复
#### Version 2.6.1
> 1.修复空经验值  

#### Version 2.6.0
> 1.pb格式兼容  
> 2.exp 经验值转发  

#### Version 2.5.0
> 1.经验重构  

#### Version 2.4.3
> 1.无头像返回默认头像

#### Version 2.4.0
> 1.会员信息基础接口上线

#### Version 2.3.2
> 1.删除upInfo接口逻辑

#### Version 2.3.1
> 1.修复移动端关闭实名认证、异地登录提醒接口关闭异地登录提醒不成功

#### Version 2.3.0
> 1.update http client conf and usage
 
#### Version 2.2.0
> 1.移动端获取及关闭异地登录、实名制提醒
 
#### Version 2.1.0
> 1.myinfo接口调用account-java的获取用户信息myinfo内网接口

#### Version 2.0.0
> 1.增加info card等信息

#### Version 1.4.2
> 1.fix upinfo cache error

#### Version 1.4.1
> 1.添加upinfo容错信息

#### Version 1.4.0
> 1.新增up infos: pic and blink

#### Version 1.3.0
> 1.剔除答题逻辑

#### Version 1.2.0
> 1.实名制提示
> 2.获取UP主权限

#### Version 1.1.0
> 1.merge member-service into Kratos

#### Version 1.0.7
> 1.排行榜昵称

#### Version 1.0.6
> 1.答题转正重试

#### Version 1.0.5
> 1.新增异地登录功能

#### Version 1.0.4
> 1.修复h5不显示答题

#### Version 1.0.3
> 1.答题图片bfs

#### Version 1.0.2
> 1.接入prom

#### Version 1.0.1
> 1.修复封禁时长

#### Version 1.0.0
> 1.基础api
