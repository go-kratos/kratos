##### ep-melloi

##### Version 1.9.7
1.业务断言支持
2.用户上传脚本默认生成场景类型报告
3.增加测试报告的开始和结束时间
4.修复场景脚本无法修改成功问题
5.紧急修复sceneId bug
6.key-value 容错
7.fmt

##### Version 1.9.6
1.压测审批通知增加依赖服务
2.修改服务bug

##### Version 1.9.5
1. 压测时间段走审批

##### Version 1.9.4
1.gRPC 自助压测页面支持保存

##### Version 1.9.3
1.支持permit配置超时
2.配合caster 新增容器修改

##### Version 1.9.2
1.压测执行代码重构
2.重新合master

##### Version 1.9.1
1.上传脚本压测功能
2.添加cookie

##### Version 1.9.0
1.压测熔断提示后端
2.script 局部更新

##### Version 1.8.10
1.删除咸鱼容器

##### Version 1.8.9
1.新增定时器功能

##### Version 1.8.8
1.场景压测自助压测及我的脚本中，均支持熔断成功率的配置

##### Version 1.8.7
1.remove mail password
2.alter wechat notify content

##### Version 1.8.6
1.对加压逻辑改造，压测停止后，没添加的容器的不再添加
2.优化平均响应时间统计算法

##### Version 1.8.5
1.场景脚本排序逻辑调整
2.添加ttp、grpc以及场景脚本的删除功能

##### Version 1.8.4
1.运行接口增加cookie

##### Version 1.8.3
1.微信通知增加依赖服务通知
2.微信群通知增加服务依赖
3.压测容器详情

##### Version 1.8.2
1.微信通知增加依赖服务通知

##### Version 1.8.1
1.压测接口url请求参数包含json格式时，自动encode
2.修改脚本后保存时断言、数据文件删不掉的问题修复

##### Version 1.8.0
1. 增加 ping 功能和显示
2. 修复keepAlive 修改失败问题
3. 日志查询优化

##### Version 1.7.9
1. 添加grpc debug 标签

##### Version 1.7.8
1.强制删除容器修改为双层删除逻辑，自身接口删除失败会调用paas接口来删除

##### Version 1.7.7
1. 多断言
2. 修复报告不更新问题
3. 删除main二进制文件

##### Version 1.7.6
1.首页bugfix&数据清洗

##### Version 1.7.5
1.场景压测支持绑定非默认172.22.22.222的其他ip

##### Version 1.7.4
1.修改grpc returnType 包含. 的问题

##### Version 1.7.3
1.场景复制功能

##### Version 1.7.2
1.强制删除容器，不依赖caster平台

##### Version 1.7.1
1.melloi 服务树增加nil判断

##### Version 1.7.0
1.melloi首页统计、场景压测增加UT

##### Version 1.6.9
1.melloi部分代码增加UT
2.更新db.sql
3.report 更新策略优化

#####  Version 1.6.8
1. 下载功能

#####  Version 1.6.7
1. 对获取不到物理机ip 的容器 添加补偿机制
2. 解除查询job 限制

#####  Version 1.6.6
1. http 压测增加301，302 code

#### Version 1.6.5
1. 修复场景压测post接口编辑时不展示header和form的问题
2. 修复不同分组添加接口跑到同一组的问题
3. 修复场景接口 json 解析失败
4. 修复jmeter报的 json 错误


#####  Version 1.6.4
1. url 解析错误提示

#### Version 1.6.3
1. 选择已有接口支持选择到服务树

##### Version 1.6.2
1.  场景脚本批量执行

##### Version 1.6.1
1.  grpc 脚本详情支持白名单用户查看

#### Version 1.6.0
1. 场景名称支持输入中文，脚本路径命名调整

#### Version 1.5.9
1. 获取压测job容器所在的物理机ip并落地

##### Version 1.5.8
1. 批量增加容器
2. headers 和 argument-param 初始化

##### Version 1.5.7
1. 修复reportSummary 状态问题

##### Version 1.5.6
1. 获取压测job容器所在的物理机ip并落地

##### Version 1.5.5
1. 压测报告支持百分位显示

##### Version 1.5.4
1.获取压测job容器所在的物理机ip并落地

#### Version 1.5.3
1. 修复只显示queryFree 10 个的问题
2. 修改 sign
3. 增加 sign

##### Version 1.5.2
1.修复场景分组问题
2.hostinfo配置化

##### Version 1.5.1
1.域名支持绑定用户指定的非172.16.1.1的host

#### Version 1.5.0
1. 定义 script.TestType

##### Version 1.4.9
1.支持grpc复制功能

##### Version 1.4.8
1.解析grpc 支持 java_package

##### Version 1.4.7
1.首页压测次数统计增加压测类型维度

##### Version 1.4.6
1.首页优化，支持grpc、场景相关数据统计

#### Version 1.4.5
1. defer close file

#### Version 1.4.4
1. 删除容器时更新测试报告状态

#### Version 1.4.3
1. 优化grpc执行代码

#### Version 1.4.2
1.修改grpc通知 压测人BUG

##### Version 1.4.1
1.调整grpc解析，将 method aa_bb 转成 aaBb

####Version 1.4.0
1. 增加http 短连接配置

##### Version 1.3.9
1.添加压测配置查询接口
2.修复添加场景时草稿箱状态错误问题

##### Version 1.3.8
1. 场景模块上传数据文件改造
2. 批量停止压测优化

##### Version 1.3.7
1.支持场景上传数据文件

##### Version 1.3.6
1.场景压测删除草稿、清空草稿箱功能

##### Version 1.3.5
1.output_params 初始化
2.停止所有容器接口

##### Version 1.3.4
1.GRPC 更新bug修复，支持上传新的文件

##### Version 1.3.3
1.支持压测熔断参数用户可配置
2.重构部分代码

##### Version 1.3.2
1.GRPC压测支持参数化
2.GRPC压测支持路径

##### Version 1.3.1
1.场景压测中编辑接口并保存后执行顺序问题修复

##### Version 1.3.0
1.修复生成多余的文件夹 和 jmx 文件
2.修复post 请求的 body 丢失问题

##### Version 1.2.9
1.修复场景压测未绑定容器host问题

##### Version 1.2.8
1. 场景 post 请求增加3个模板
2. 解决 参数依赖模块写入jmx 文件发生 html 转义问题
3. 修复修复输出参数兼容性引起的排序问题

##### Version 1.2.7
1.修复 quick-start 批量选择接口自动生成场景多余的 scene.jmx 文件问题
2.修复 script-list 批量选择接口生成场景的 script 表无 sceneId 问题

##### Version 1.2.6
1.接口输出参数新老数据兼容性处理
2.TestType配置化
3.可选参数列表bugfix

##### Version 1.2.5
1.修复单接口修改脚本的报错问题

##### Version 1.2.4
1.修复单接口调试、执行压测的502报错问题

##### Version 1.2.3
1.GRPC优化，支持多个proto文件编译
2.GRPC 编译&执行分离

##### Version 1.2.2
1.场景压测输出参数结构调整
2.可用参数列表接口调整
3.排序逻辑调整

##### Version 1.2.1
1.修复拿服务树节点的bug

##### Version 1.2.0
1.接口之间多个参数依赖
2.Cpu 核心数配置化

##### Version 1.1.9
1.场景压测预览图接口开发，支持前端层级关系图模式
2.从草稿箱选择一个场景后，再选择已有接口，接口报错问题修复
3.部分接口查询条件添加有效数据的过滤条件
4.场景压测相关的bug修复

##### Version 1.1.8
1.修复根据id 查询报告bug

##### Version 1.1.7
1.修复脚本登录状态 和 异步压测脚本无法修改问题
2.压测快照增加 binary 信息和展示
3.场景脚本 debug 线程数优化
4.修复新用户登录404 bug

##### Version 1.1.6
1.整拖拽功能底层逻辑，配合前端
2.选择已有接口逻辑调整
3.场景压测增加预览图接口，支持层级关系图模式
4.一些接口的bug修复

##### Version 1.1.5
1.grpc 场景脚本查询增加白名单
2.压测时间段校验忽略白名单
3.修复 grpc 压测微信通知报告地址bug
4.脚本快照增加 post 请求的 multipart/form-data 类型
5.优化场景脚本生成逻辑

##### Version 1.1.4
1.场景压测调试
2.异步压测脚本修改
3.post IO 流请求脚本新增和修改

##### Version 1.1.3
1.将jmx模板路径放到配置文件

##### Version 1.1.2
1. 新增脚本详情页压测入口
2. http 异步压测脚本生成
3. grpc 异步压测脚本生成

##### Version 1.1.1
1. 修复程序占用大量内存和 cpu 问题

##### Version 1.1
1. 脚本复制，脚本修改
2. 压测微信通知，错误熔断

##### Version 1.0.0
1 自助压测，填写接口、QPS、压测时长，快速开始压测
2 自测报告，查询压测中以及压测完成的报告
3 手工压测，填写压测申请，由ep部门辅助进行压测
4 脚本列表，已创建成功的压测脚本，可选择重复执行
5 集群资源，查询压测机器资源利用情况