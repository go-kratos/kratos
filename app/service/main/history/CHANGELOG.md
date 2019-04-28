# history-service

# v1.3.0
1. tidb client   

# v1.2.13
1. UserHistoriesReq Ps 1000  

# v1.2.12
1. 屏蔽掉 18446744073709551615 

# v1.2.11
1. 去掉 DeleteHistories-NothingFound  

# v1.2.10
1. DeleteHistories-NothingFound   

# v1.2.9
1. 翻页重复   

# v1.2.8
1. context.Background() 

# v1.2.7
1. grpc v1 

# v1.2.6
1. remove ping

# v1.2.5
1. 初始化grpc reply

# v1.2.4
1. 重新构建

# v1.2.1
1. 重新rebase master

# v1.2.0
1. merge增加kid 去重防止重复消费
2. 使用pipeline简化合并逻辑
3. 修复user接口panic bug

# v1.1.1
1. 修复panic bug

# v1.1.0
1. 异步清除播放历史
2. 分批删除数据

# v1.0.5
1. job写入数据库
# v1.0.4
1. 聚合批量写入数据库

# v1.0.3
1. 修改grpc package 名称
2. 去掉device验证

# v1.0.1
1. add ping redis

# v1.0.0
1. 上线功能播放历史
