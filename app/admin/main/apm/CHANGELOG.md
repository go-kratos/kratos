### apm-admin
##### Version  1.23.7
> 1. 切换canal更新接口为/home/config/update

##### Version 1.23.6
> 1. 增加databus获取消费者地址接口  

##### Version 1.23.5
> 1. go-common/library包加入ut检验

##### Version 1.23.4
> 1. 去掉 ut_app log 中 has_ut 字段.

##### Version 1.23.3
> 1. 修改alarm接口中百分比字段的默认值为0
> 2. 修改创建group时的百分比字段默认值为80

##### Version 1.23.2
> 1. 优化和修正因增加library单测的ut逻辑

##### Version 1.23.1
> 1. 批量修改告警功能

##### Version 1.23.0
> 1.unit_test.sh 大重构.  
> 2.增加对整个项目覆盖率的统计.  
> 3.现在所有涉及到ut_pkganls表都有过滤项目逻辑(方便以后统计逻辑变更).

##### Version 1.22.1
> 1.修改gitface token

##### Version 1.22.0 
> 1.删除canal申请auth

##### Version 1.21.0
> 1.添加databus topic查询

##### Version 1.20.4
> 1.告警接口延迟

##### Version 1.20.3
> 1.去除部分冗余代码，加入group批量添加告警接口

##### Version 1.20.2
> 1.add unlock.

##### Version 1.20.1
> 1.增加canal申请审核微信通知
> 2.增加需求建议微信通知

##### Version 1.20.0
> 1.优化覆盖率计算方式.

##### Version 1.19.3
> 1.databus告警接口添加project

##### Version 1.19.2
> 1.ut.Upload 接口增加耗时日志.

##### Version 1.19.1
> 1.多语言错误码后台

##### Version 1.18.8
> 1.增加rank排名变化显示.  
> 2.增加ut_app owner更新逻辑

##### Version 1.18.7
> 1.upload date_file 不为空.

##### Version 1.18.6
> 1.ut上传接口增加ut_app表处理  

##### Version 1.18.5
> 1.增加redis配置  
> 2.ut 新增周五群发周刊消息  

##### Version 1.18.4
> 1.reply-feed透传 

##### Version 1.18.3
> 1.fix canal infoc字段

##### Version 1.18.2
> 1.fix canal 编辑project&leader字段被覆盖

##### Version 1.18.1
> 1.databus告警

##### Version 1.17.3
> 1.userAuth接口增加返回头像  
> 2.ut包覆盖率改为用包内文件计算

##### Version 1.17.2
> 1.增加dashboard项目负责人视图

##### Version 1.17.1
> 1.增加项目接入ut情况的统计

##### Version 1.16.13
> 1.修复utRank endtime格式

##### Version 1.16.12
> 1.添加monkey

##### Version 1.16.11
> 1.添加一个接口

##### Version 1.16.10
> 1.同步code码信息获取线上信息改为内部接口调用

##### Version 1.16.9
> 1.同步code码信息

##### Version 1.16.8
> 1.修复canal databus group重复返回错误提示

##### Version 1.16.7
> 1.优化UT列表和历史列表接口项目覆盖率计算方式  
> 2.修复UT发送微信消息超过2048字节被截断问题

##### Version 1.16.6
> 1.修复group_concat函数默认长度为1024导致数据被截断问题  
> 2.处理调用cmd.Start()导致产生僵尸进程的问题

##### Version 1.16.5
> 1.go command add cmd.Wait().  
> 2.commented activeWarning func code.


##### Version 1.16.4
> 1.ut增加原始覆盖数据文件处理  
> 2.ut企业微信和Git报告增加行数相关信息  
> 3.修复databus生产group申请冲突问题

##### Version 1.16.3
> 1.ecode增加告警级别和繁体信息字段

##### Version 1.16.2
> 1.重写monitor模块  
> 2.修复upload接口未关闭http链接问题

##### Version 1.16.1
> 1.修改canal scan databus&infoc为空返回null

##### Version 1.16.0
> 1.添加运维报警接口

##### Version 1.15.15
> 1.ut rank增加抛物线时间衰减因子统计.

##### Version 1.15.14
> 1.ut曲线图修复按hour统计数据问题
> 2.补充canal dao层测试用例
> 3.canal修复project字段编辑被覆盖

##### Version 1.15.13
> 1.dashboard 个人排行榜增加排位统计.  
> 2.解除openView默认权限改为需授权状态

##### Version 1.15.12
> 1.排行榜威尔逊区间算法增加牛顿冷却因子.

##### Version 1.15.11
> 1.ut general commit排序功能修复  
> 2.ut upload接口增加author参数  
> 3.ut commit context with background.  
> 4.修正ut数据对比若干问题

##### Version 1.15.10
> 1.utRank排名算法优化  
> 2.ut下线发送Bottom 10微信群消息的功能

##### Version 1.15.9
> 1.修复pprof返回的图片地址参数错误的问题  
> 2.新增性能管理权限,去除/discovery权限验证  
> 3.修复ut微信报告中用户名取值问题

##### Version 1.15.8
> 1.ut dashboard 返回baseline  
> 2.ut增加发送企业微信消息功能

##### Version 1.15.7
> 1.ut dashboard 新增最近10条个人历史commit记录接口  
> 2.增加ut聚合commit概要信息  
> 3.ut shell脚本magic执行顺序修改及upload失败直接返回错误不做check

##### Version 1.15.6
> 1.优化utRank排名实现方式

##### Version 1.15.5
> 1.告警抓取的信息存储至db  
> 2.新增告警信息查询接口

##### Version 1.15.4
> 1.修复 format 时乱序的问题

##### Version 1.15.3
> 1.修复git评论ut简报包链接跳转问题  
> 2.修复ut检查接口commit_id不存在返回不报错的问题  
> 3.ut dashboard 新增个人全量排名

##### Version 1.15.2
> 1.优化dashboard时间顺序展示

##### Version 1.15.1
> 1.删除存储告警信息  
> 2.新增根据告警信息实时抓取内存图和性能图

##### Version 1.15.0
> 1.ut dashboard 增加个人质量聚合数据

##### Version 1.14.17
> 1.增加ut个人dashboard文件展示功能  
> 2.更改SAGAReaport功能自动调用gitlab接口添加留言  
> 3.修复dao层ut用例

##### Version 1.14.16
> 1.增加拉取gitlab头像的token

##### Version 1.14.15
> 1.单元测试tyrant接口改名为check

##### Version 1.14.14
> 1.修复canal monitor_period为空bug  
> 2.优化canal 历史代码

##### Version 1.14.13
> 1.添加rank排行头像

##### Version 1.14.12
> 1.添加存储告警数据接口

##### Version 1.14.11
> 1.修复canal编辑字段缺失

##### Version 1.14.10
> 1.修复monitor采集databus数据类型转换的错误bug

##### Version 1.14.9
> 1.单元测试merge/set接口改为接收webhook参数形式

##### Version 1.14.8
> 1.增加单元测试统计数据的排行utRank

##### Version 1.14.8
> 1.ut添加当前用户最近10次提交记录  
> 2.upload开启一个事物存储表数据

##### Version 1.14.7
> 1.增加一批用户默认权限

##### Version 1.14.6
> 1.notify添加机房选项

##### Version 1.14.5
> 1.canal databus信息增加从app旧表获取key&secret  
> 2.canal config信息增加monitor_period字段  
> 3.canal check master接口改为post  
> 4.canal DB增加infoc部分

##### Version 1.14.4
> 1.ut/check接口去掉增长率的判断

##### Version 1.14.3
> 1.新增监控rpc数据  
> 2.修复获取在线人数统计  
> 3.扩展monitor/prometheus 接口

##### Version 1.14.2
> 1.优化查询历史ut_detail sql

##### Version 1.14.1
> 1.增加pprof火焰图

##### Version 1.13.13
> 1.新增monitor监控数据查询接口

##### Version 1.13.12
> 1.ut/check 取历史最大值比较

##### Version 1.13.11
> 1.canal 增加databus auth表判断

##### Version 1.14.0
> 1.重构 ut 相关所有接口.  

##### Version 1.13.10
> 1.ut/upload 强制更新默认值

##### Version 1.13.9
> 1.ut/check 添加日志  
> 2.增加通过率判断

##### Version 1.13.8
> 1.ut/check 去除多余判断

##### Version 1.13.7
> 1.add user 需求module

##### Version 1.13.6
> 1.ut/check接口修改默认值

##### Version 1.13.5
> 1.ut单元测试添加达标检测接口

##### Version 1.13.4
> 1.open鉴权管理透传

##### Version 1.13.3
> 1.discovery列表接口变更

##### Version 1.13.2
> dapper 添加独立 Host

##### Version 1.13.1
> 1.添加monitor监控定时任务  
> 2.bfs upload接口返回错误信息

##### Version 1.13.0
> 1.增加monitor数据监控收集模块

##### Version 1.12.0
> 1.增加需求与建议模块

##### Version 1.11.8
> 1.增加ut baseline接口参数配置化  
> 2.canal审核接口增加force强制推送  
> 3.ut列表返回三个月之内数据  
> 4.utlist,historycommit增加mtime字段
##### Version 1.11.7
> 1.ut detail 接口pkg数据去重  
> 2.修改upload返回sven地址  
> 3.增加unittest.sh对pkg下是否存在go文件的判断
##### Version 1.11.6
> 1.ut增加history commit及bfs proxy接口
##### Version 1.11.5
> 1.upload接口增加TestResult结果过滤
##### Version 1.11.4
> 1.canal process 增加对master_info 表中addr的检验
##### Version 1.11.3
> 1.上传bfs接口删除二级目录/ut/
##### Version 1.11.2
> 1.上传bfs接口添加二级目录/ut/
##### Version 1.11.1
> 1.heap接口创建文件改成串行
##### Version 1.11.0
> 1.bm bind返回修改  
> 2.pprof heap添加独立接口

##### Version 1.10.11
> 1.canal修改config配置参数,server_id,active
##### Version 1.10.10
> 1.pprof以主机名为纬度存图
##### Version 1.10.9
> 1.notify申请支持重命名功能  
> 2.对pprof增加对内存的数据抓取
##### Version 1.10.8
> 1.增加canal申请对user/pwd进行checkmaster检查.  
> 2.canal config table字段增加primarykey,omitfield
##### Version 1.10.7
> 1.修复notify编辑的错误  
> 2.修复读取svg图的权限错误  
> 3.修复用户模块和权限编辑功能  
> 4.修复一些bm框架接收数据的问题
##### Version 1.10.6
> 1.修复required 错误
##### Version 1.10.5
> 1.增加canalsvenCo配置文件信息
##### Version 1.10.4
> 1.pprof命令拼错的bug修改
##### Version 1.10.3
> 1.parse修改为bm框架的方法
##### Version 1.10.2
> 1.pprof的go路径配置
##### Version 1.10.0
> 1.增加cacheview权限 
> 2.抽离permit  

##### Version 1.9.1
> 1.pprof功能更新，加入cpu和内存图  
> 2.平台管理透传

##### Version 1.9.0
> 1.新增ut模块，用户查询saga执行的单测结果  
> 2.新增upload文件上传接口，提供给saga调用，将单测执行的结果存储至bfs

##### Version 1.8.3
> 1.增加token参数调用配置中心canal接口

##### Version 1.8.2
> 1.fix config update bug

##### Version 1.8.1
> 1.pprof功能
##### Version 1.8.0
> 1.canal 新增申请并同步至配置中心接口 apply/config  
> 2.canal scan接口增加管理员编辑获取user/pwd  
> 3.canal list接口增加status/operator返回字段  
> 4.canal 新增返回所有已存在canal addr信息接口addrall  
> 5.canal 新增canal编辑申请接口config/edit  
> 6.canal 新增apply申请列表 apply  
> 7.canal 新增process审核接口 process
##### Version 1.7.5
> 1.topic修改集群时在新的kafka自动创建topic

##### Version 1.7.7
> 1.修复事务的bug
##### Version 1.7.5
> 1.topic修改集群时在新的kafka自动创建topic

##### Version 1.7.4
> 1.databus加入按时间和位置设置offset  
> 2.修正groups接口operation字段的返回值  
> 3.修正创建topic时的一个kafka集群判断bug  
> 4.discovery的下拉菜单权限开放

##### Version 1.7.3
> 1.user申请时非线上环境过滤配置中心权限申请  
> 2.databus的Group申请，编辑，审核中支持新增和修改topic相关功能

##### Version 1.7.2
> 1.user新增权限申请修改接口

##### Version 1.7.0
> 1.user新增申请操作权限，审核权限等功能
##### Version 1.6.0
> 1.bm框架
##### Version 1.5.3
> 1.修复查看canal详情获取错误bug  
> 2.env环境读取本地环境变量  
> 3.过滤敏感信息

##### Version 1.5.2
> 1.discovery获取服务列表接口变更
##### Version 1.5.1
> 1.group差值partition排序  
> 2.discovery的polling接口聚合  
> 3.topic列表查询修改cluster非必须  
> 4.tree/nodes服务树权限去重
##### Version 1.5.0 
> 1.添加canal/delete接口，删除canal配置  
> 2.添加canal/scan接口，查询canal详情
##### Version 1.4.2 
> 1.fix notify bug 
##### Version 1.4.1
> 1.添加kafka创建topic  
> 2.修复申请消费修改project但group没重新生成的bug  
> 3.去cluster

##### Version 1.4.0
> 1.添加notify申请编辑  

##### Version 1.3.7
> 1.修正group申请的一系列bug

##### Version 1.3.6
> 1.修正databus申请有重复无法创建的bug  
> 2.修改集群的时候，同时也要修改老表的cluster

##### Version 1.3.5
> 1.sven权限变更

##### Version 1.3.4
> 1.databus group申请及审核功能

##### Version 1.3.3
> 1.databus group设为kafka最老记录  
> 2.添加所有group 分页diff接口

##### Version 1.3.2
> 1.log日志接入  
> 2.apm-admin迁到main下

##### Version 1.3.1
> 1.修复canal不能编辑leader和cluster的问题

##### Version 1.3.0
> 1.增加鉴权功能

##### Version 1.2.9
> 1.修复group的重命名权限

##### Version 1.2.8
> 1.databus的app/edit的bug修正

##### Version 1.2.7
> 1.databus的offset差值规则修改

##### Version 1.2.7
> 1.app鉴权接口

##### Version 1.2.6
> 1.databus,group修改增加备注字段  
> 2.user/auth接口返回格式变更  
> 3.修复新加group的告警默认为0的bug  
> 4.服务树接口更换  
> 5.discovery的nodes接口特殊规则

##### Version 1.2.5
> 1.databus告警,group展示告警规则,修复了alarm接口连删除数据也返回的错误

##### Version 1.2.4
> 1.databus，canal告警功能，codes改为ecode

##### Version 1.2.3
> 1.discovery透传接口，增加appids接口

##### Version 1.2.2
> 1.databus新增功能，删除，group重命名，offset等

##### Version 1.2.1
> 1.cannl更优雅的写法，auth方法变更

##### Version 1.2.0
> 1.增加配置中心菜单权限

##### Version 1.1.9
> 1.增加group判断是否已存在

##### Version 1.1.8
> 1.databus告警需要的节点接口

##### Version 1.1.7
> 1.项目组重复增加group2  

##### Version 1.1.6
> 1.修复服务树接口  

##### Version 1.1.5
> 1.新增canal增，查，改接口

##### Version 1.1.4
> 1.codes新接口

##### Version 1.1.3
> 1.添加dapper接口支持

##### Version 1.1.2
> 1.修复apm permit

##### Version 1.1.1
> 1.过滤group的projects

##### Version 1.1.0
> 1.用户权限限制

##### Version 1.0.0
> 1.用户管理  
> 2.databus管理     
