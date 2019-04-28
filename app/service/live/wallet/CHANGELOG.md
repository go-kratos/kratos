# v2.0.5
1. 避免当开启事务失败返回的指针没有初始化panic

# v2.0.4
1. 兑换硬币如果硬币接口失败还是认为兑换已经成功且设置返回避免panic

# v2.0.3
1. 兑换硬币如果硬币接口失败还是认为兑换已经成功

# v2.0.2
1. 取数据没取到不写db
2. 新建一个goroutine发消息队列

# v2.0.1
1. cache出错不回源db
2. 记流水时修改cnt字段

# v2.0.0
1. 事务化
2. 双余额

# v1.0.11
1. 修改retry为0

# v1.0.10
1. 修改retry为1

# v1.0.9
1. add test

# v1.0.8
1. add version support

# v1.0.7
1. add cache version support

# v1.0.6
1. add bi support
2. add cost_base support

# v1.0.5
1. build bazel

# v1.0.4
1. add pub wallet change info

# v1.0.3
1. add reason for pay metal
2. fix err bug

# v1.0.2
1. fix coin stream index can not fix

# v1.0.1
1. 增加delCache接口
2. recordCoinStream 增加写入数据库的字段
3. move dir

# v1.0.0
1. 新增
