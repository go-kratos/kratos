### config-service

#### Version 2.3.5
>1. 增加一个http接口，直接获取当前最新发布的内容

#### Version 2.3.4
>1. log日志修改

#### Version 2.3.3
>1. sql语句中force是关键字要加``

#### Version 2.3.2
>1. 配置中心最新规则：
规则1：当前版本大于0时优先级 单机强制 > 全局强制(当前版本小于等于最近的一次强制版本号时才会拉取,主要为了castr发版时能拉到最新配置) > 指定版本 > 当前发布最新版本
规则2：当前版本小于等于0时没有单机强制和全局面强制逻辑，优先级为  指定版本 > 当前最新发布版本

#### Version 2.3.1
>1. ut补全

#### Version 2.2.1
>1. 忽视appiont,执行强制更新功能

#### Version 2.1.4
>1. sdk 多处连接泄露修复

#### Version 2.1.3
>1. sdk 连接泄露bug修复

#### Version 2.1.2
>1. discovery注册

#### Version 2.1.1
>1. 限流
>2. 去tree_id并兼容老tree_id接口

#### Version 2.1.0
>1. bm架构

#### Version 2.0.4
>1. 增加register

#### Version 2.0.3
>1. fix sh001 limit and mv file to main

#### Version 2.0.2
>1. 移除statsd 模块

#### Version 2.0.2
>1. 版本推送不依赖redis

#### Version 2.0.1
>1. 修复file.so 接口，兼容新的和老的file.so逻辑

#### Version 2.0.0
>1. 配置中心版本v4,走新库新表
>2. 支持客户端增量更新

#### Version 1.6.2
>1. push 接口推送时候如果数据库内容为空，返回失败

#### Version 1.6.1
>1. check接口增加自定义参数上报

#### Version 1.6.0

>1. 增加批量添加配置接口
>2. 增加配置拷贝接口
>3. 增加修改版本下的所有配置接口，没有就新加，有就覆盖
>4. 增加返回未配置完成版本ID 列表

#### Version 1.5.0

>1. 增加配置文件命名空间，支持公共配置

#### Version 1.4.3

>1. 再次修复数据库连接bad connection 问题

#### Version 1.4.2

>1. 修复数据库连接bad connection 问题

#### Version 1.4.1

>1. bugfix 修复缓存文件key 错误

#### Version 1.4.0

>1. 接入普罗米修斯监控
>2. 更新go-common和go-business 为最新的
>3. file.so 接口添加日志返回，配置内容本地缓存

#### Version 1.3.0

>1. 增加builds，versions接口,分别获取所有构建版本和版本id。
>2. check接口支持appoint参数，返回appoint作为版本id。

#### Version 1.2.4

>1. 更改主机历史记录保存时间，超过3小时没有在线。在查询的时候删除。
>2. 主机超时时间延长5秒

#### Version 1.2.3

>1. file接口去掉hostName参数，version改为可选参数。更改获取逻辑，直接从数据库获取单个配置文件。

#### Version 1.2.2

>1. 将返回值由json改为文件内容

#### Version 1.2.1

>1. 增加获取单个配置接口file


#### Version 1.2.0
>1. 启动参数增加token字段，应用注册生成，用于应用和环境权限限制
>2. 应用启动配置文件本地map缓存一份，同时写入配置文件


#### Version 1.0.1
>1. APM 添加build version映射表，对应于config版本号
>2. 发布版本，只需要修改这个build version映射表的config版本号。
>3. 目前先去掉推送主机和主机标记功能，下个版本考虑。
>4. 程序编译时自己带build version，可以直接通过映射表，获取使用的配置版本。


#### Version 1.0.0
>1. 支持推APM推服务主机配置，生成本地缓存版本。
>2. 支持长轮询监听新配置文件更新，并推送版本号到客户端。
>3. 支持按版本号获取配置文件。
>4. php以进程方式启动，并监听配置文件更新。
>5. go以零配置启动，下载配置文件，并监听配置更新。





