### v1.1.8
1. try fix vip info reply panic
2. 添加guard client
3. 添加获取某个主播全量守护数量接口
4. 添加multi守护接口 
5. 用户经验接口替换
6. 替换用户经验接口

### v1.1.7
1. 优化房管admin change message

### v1.1.6
1. xuser vip获取db记录时优化sql: no rows in result set error
2. 迁移某个主播最近开通的总督弹窗接口

### v1.1.5
1. 大航海读接口迁移
2. 异步更新经验

### v1.1.4
1. 添加新接口判断是否房管不返回用户信息

### v1.1.3
1.兼容vip cache中vip/svip可能是int/string的问题

### v1.1.2
1.修复vip 从db/cache获取数据时可能是过期数据导致vip/svip状态返回不正确的问题

### v1.1.1
1. 添加房管不再发送私信

### v1.1.0
1. 5.36房管 & 房管全量接口 gRPC api
2. 添加主播等级颜色

### v1.0.1
1. 购买姥爷更新流水表来源字段
2. 添加用户经验接口字段

### v1.0.0
1. 直播vip buy & vip info grpc api
2. 经验
3. 购买大航海