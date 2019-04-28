# 运营后台搜索项目-后台


v2.1.10
1. 支持combo not

v2.1.9
1. fix archive_check time

v2.1.8
1. 添加es集群ower接口

v.2.1.2
1. 用户行为日志索引多集群支持

v.2.1.1
1. SetSniff false while ops log cluster

v2.1.0
1. 删除无用的v2接口， tag_update/pgc/account

v2.0.2
1. update 大数字转int64

v1.9.5
1. 支持随机种子random seed
2. IgnoreUnavailable & AllowNoIndices

v1.9.4
1. 支持like level middle
2. 优化middle

v1.9.3
1. sven搜索目录修改

v1.9.2
1. 日志显示ip 排除ip为空字符串

v1.9.1
1. 日志显示ip 数据统计

v1.9.0
1. 删除无用的v2接口, tag/blocked/dm/vip/pgc

v1.8.9
1. sven接口格式修改

v1.8.8
1. 创作中心fix

v1.8.7
1. 新增es统计接口
2. ping改成异步
3. 新增创作中心稿件接口

v1.8.6
1. 行为日志支持like查询

v1.8.5
1. 稿件title查询集群修改

v1.8.4
1. 稿件title查询

v1.8.3
1. 支持自由组合过滤combo

v.1.8.2 
1. fix workflow param

v1.8.0
1. remove music search api

v1.7.2
1. QueryBasic变清真
2. 支持query mode nested
3. 电竞日历优化

v1.7.1
1. fix app struct sleep type, int => float64

v1.7.0
1. 对接sven

v1.6.4
1. 换成bm的binding
2. 定制化搜索时，支持调用QueryBasic

v1.6.3
1. 支持like level
2. 支持enhanced
3. 查询结果使用QueryResult，支持slice或map的情况
4. like过滤掉特殊字符
5. 删除Appid
6. 支持upsert

v1.6.2
1. 规规一波

v1.6.1
1. workflow返回oid

v1.5.9
1. search/query去掉appid参数

v1.5.8
1. 支持queryExtra，即查询体基础上自定义
2. 去除queryConf自定义部分的business

v1.5.7
1. search/query接口去掉校验business

v1.4.4
1. 更新es包。支持Collapse。

v1.4.3
1. 从db读取配置

v1.4.2
1. 弹幕举报group接口

v1.4.1
1. 查询体支持or和not
2. 修复日志panic bug

v1.4.0
1. 新增通用查询体Query
2. 支持接口Debug，可分别在dsl前后debug

v1.3.9
1. workflow接口

v1.3.8
1. workflow上线

v1.3.7
1. 日志平台默认只查两个索引、优化

v1.3.6
1. 日志支持数组

v1.3.5
1. 修复日志bug

v1.3.4
1. 修复archive check
2. 精准搜索

v1.3.3
1. workflow |改&

v1.3.2
1. workflow 新增接口

v1.3.1
1. log group by数量fix

v1.3.0
1. log group by
2. workflow fix

v1.2.9
1. 增加owner

v1.2.8
1. dm增加参数

v1.2.7
1. log配置化

v1.2.6
1. log支持位，配置化

v1.2.5
1. log新增不分表索引

v1.2.4
1. 去除account校验

v1.2.3
1. 修复log鉴权

v1.2.2
1. 修改日志集群

v1.2.1
1. copyright修改匹配度

v1.2.0
1. copyright修改查询方式

v1.1.9
1. copyright上线

v1.1.8
1. 日志平台修改权限点
2. 新增copyright接口

v1.1.7
1. 日志平台新增权限点
2. pgc接口
3. vip新增索引字段

v1.1.6
1. 增加打分

v1.1.5
1. 修复匹配度百分号位置

v1.1.4
1. 弹幕关键字搜索降低匹配模糊度

v1.1.3
1. 弹幕监控上线

v1.1.2
1. 移除stats

v1.1.1
1. 添加workflow接口参数
2. 合并后台管理接口

v1.1.0
1. 修改集群地址

v1.0.9
1. 弹幕接口上线

v1.0.8
1. 增加搜索条件

v1.0.7
1. workflow_group_common
2. workflow_chall_common

v1.0.6
1. search-interface迁到search-admin
2. 增加关键字高亮

v1.0.5
1. 解决冲突

v1.0.4
1. 修改接口参数

v1.0.3
1. blocked接口重构

v1.0.3
1. workflow模糊匹配精确度提高到80%

v1.0.2
1. 风纪委上线

v1.0.1
1. 修改tag_rounds参数

v1.0.0
1. 初始化项目，更新依赖
2. 预发和上线