运营后台搜索项目-Job
===============

v1.5.2
1. activity_all 

v1.5.1
1. esports_fav_all 

v1.5.0
1. TODO 去掉data_fields的_id和_mtime，改为读新的配置字段
2. TODO bulk优化

v1.4.4
1. 支持workflow特殊字段
2. 行为日志支持配置

v1.4.3
1. fix slice bounds out of range

v1.4.2
1. 回滚数据字符处理`

v1.4.1
1. extra支持多字段

v1.4.0
1. 支持多个sliceField

v1.3.9
1. 修复app_databus全量不加offset问题

v1.3.8
1. 废除databus_index_id

v1.3.7
1. 去除弹幕老逻辑

v1.3.6
1. dm_date bug修复
2. 打印es写入时间

v1.3.5
1. app_multiple_databus更加通用

v1.3.4
1. 日志只发送评论

v1.3.3
1. 行为日志infoc打印log

v1.3.2
1. dmreport预发上线
2. dataExtra新增条件过滤

v1.3.1
1.fix log bug

v1.3.0
1. update infoc sdk

v1.2.0
1.日志平台支持数组
2.workflow新索引修改

v1.1.9
1. 迁移bm框架
2. 日志平台支持多个集群
3. workflow新索引上线

v1.1.8
1. 修复amd的indexname的问题

v1.1.7
1. TODO: 合并BulkDatabusData和BulkDBData，所有数据全部移到model层处理完，不再对数据额外处理，直接循环bulk

v1.1.6
1. TODO: 全量url新增参数index_version，导数据到一个新版本索引（当不一定有别名时也支持）

v1.1.5
1. IndexAliasPrefix索引别名支持
2. business配置，去除手动写businessPool
3. 修复indexName bug

v1.1.4
1. data_fields改成json格式，兼容db和databus
2. extra_data兼容db和databus
3. 新增base.go，兼容自定义包
4. dao方法和变量对外开放

v1.1.3
1. 弹幕监控和历史上线

v1.1.2
1. 释放dataMap

v1.1.1
1. 增量前移time
2. remove无用代码

v1.1.0
1. 支持recover
2. 支持数组型字段配置
3. 简化dtb

v1.0.9
1. single、multiple的配置化支持
2. 优化attrs
3. extra_data跨库支持

v1.0.8
1. bug修复

v1.0.7
1. workflow_group_common
2. workflow_chall_common

v1.0.6
1. blocked_case增加databus消息聚合量

v1.0.5
1. 增加workflow_feedback

v1.0.4
1. 修改Unmarshal date bug

v1.0.3
1. 风纪委重构
2. 预发和上线

v1.0.2
1. 修改blocked的commit逻辑
2. 预发和上线

v1.0.1
1. 增加配置index(bool), 判断全量 or 增量
2. 预发和上线

v1.0.0
1. 风纪委项目上线(blocked)
2. 预发和上线
