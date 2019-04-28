# v2.1.0
1. api接口优化

# v2.0.0
1. 去除业务逻辑，仅作为网关

# v1.0.5
1. if order.app_id == 19 不更新榜单

# v1.0.4
1. if order.app_id == 19 榜单隐藏

# v1.0.3
1. 修改排行榜排序算法，兼容线上

# v1.0.2
1. 修复 prep_rank 读取时反复刷新TTL的问题
2. 修复 prep_rank updateOrder 取值问题
3. 修复 up主榜单留言显示问题

# v1.0.1
1. 修复av/detail:total_count取值错误
2. 改用效率更高的json序列化库

# v1.0.0
1. 增加 rank & prep_rank 
2. 增加 rank & prep_rank 相互流转逻辑
3. 新增订阅 order binlog 并更新 prep_rank
2. 增加local cache功能，允许rank通过配置使用

# v0.1.0
1. 服务基建