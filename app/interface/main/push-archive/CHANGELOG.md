# push-archive

### v1.6.0
1. 特殊关注和普通关注用不同的业务组推

### v1.5.0
1. using push grpc

### v1.4.3
1. 调整abtest策略叠加

### v1.4.2
1. prodSwitch

### v1.4.1
1. 优化abtest方式

### v1.4.0
1. abtest

### v1.3.11
1. hbase v2
2. 频率控制全部有配置控制

### v1.3.10
1. order配置可用分组和分组优先级，order元素=分组key
2. active配置活跃过滤开关，无值为关，有值为开且为默认活跃时间
3. fangroup.hitby指定分组规则：default/hbase, default=全部命中,hbase=走hbase表过滤


### v1.3.9
1. 移除up主维度半小时限制策略
2. 移除粉丝维度，每个up主限制策略
3. 批量id范围去删除统计数据，避免context超时
4. 基础库identify改成auth+verify，http.serverconfig改成bladermaster.serverconfig,bm改为newserver+start模式

### v1.3.8
1. xints迁移到model

### v1.3.7
1. 移除hbase的ping

### v1.3.6
1. 新推送接口必需提供uuid，避免由于网络问题而重试产生的重复推送

### v1.3.5
1. 推送接口修改为/push-strategy/task/add

### v1.3.4
1. 迁移http到blademaster

### v1.3.3
1. 迁移push model至push-service

### v1.3.2
1. rpc双写推送开关数据到推送平台

### v1.3.1
1. 推送平台接口支持签名认证

### v1.3.0
1. 迁移至interface/main目录下

### v1.2.24
1. 使用account-service v7

### v1.2.23
1. 支持默认活跃时间过滤

### v1.2.22
1. fix 推送名单过滤前后变量共享导致的过滤无效问题

### v1.2.21
1. 过滤不在活跃时间内的粉丝推送

### v1.2.20
1. 特殊关注粉丝免限制条件: up同时在最近常看的列表中
2. 普通关注粉丝分组优先级策略

### v1.2.19
1. 特殊关注粉丝和普通关注粉丝的推送次数分开限制

### v1.2.18
1. 限制特殊关注的推送频率: cd时间内总推送次数 + cd时间内每个粉丝关注的每个up主的推送次数

### v1.2.17
1. pgc稿件，屏蔽非特殊关注粉丝的推送

### v1.2.16
1. 统计数据redis缓存获取间隔为100ms, 降低获取 qps

### v1.2.15
1. 推送禁止时间分段设置

### v1.2.14
1. 推送开关由缓存查询，改为db查询, antispam限制查询频率

### v1.2.13
1. 统计数据db del panic fix

### v1.2.12
1. 统计数据redis cache消耗间隔 1ms

### v1.2.11
1. 统计数据缓存到redis，fix channel panic

### v1.2.10
1. 文案从统一共享，改成各组配置

### v1.2.9
1. 统计数据落库bilibili_push_archive.push_statistics
2. sms报警改为wechat报警
3. fangroup增加name=(hbasetable/special),用户prominfo统计
4. 推送接口的group参数截短,比如:ai:pushlist_offline_up截短为offline

### v1.2.8
1. 添加过审稿件统计、过审稿件的粉丝统计、特殊关注粉丝命中统计、特殊关注粉丝实际推送统计

### v1.2.7
1. 禁止推送时间内消费databus消息，使得非禁止时间推送的都是非禁止时间内更新的稿件

### v1.2.6
1. 普关粉丝未设置推送开关时，默认推送，用户推送开关分批加载

### v1.2.5
1. 实验组hbase查询改为rowkey get形式，移除之前的rowfilter形式（经常超时）

### v1.2.4
1. 实验组取值比例连续性移除

### v1.2.3
1. 实验组取值比例有起始值

### v1.2.2
1. 统计名字重复fix

### v1.2.1
1. 普通关注分3实验组

### v1.1.10
1. account rpc Info3 改成 info2

### v1.1.9
1. 修改稿件推送标题

### v1.1.8
1. 去除upper主白名单

### v1.1.7
1. 请求push接口加上版本限制

### v1.1.6
1. 修复 canal relation tag

### v1.1.5
1. up主过审推送频率限制存redis
2. 调账号RPC查up主名称

### v1.1.4
1. 修复稿件过审判断逻辑

### v1.1.3
1. 稿件过审后加up主推送白名单
2. 替换prom为go-common中对象

### v1.1.2
1. 控制databus消费速率

### v1.1.1
1. dao中加入hbase的close和ping

### v1.1.0
1. 稿件推送

### v1.0.0
1. 用户稿件推送设置的上传下载接口
