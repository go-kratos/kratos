#### database/sql

##### 项目简介
MySQL数据库驱动，进行封装加入了链路追踪和统计。

如果需要SQL级别的超时管理 可以在业务代码里面使用context.WithDeadline实现 推荐超时配置放到application.toml里面 方便热加载

##### 依赖包
1. [Go-MySQL-Driver](https://github.com/go-sql-driver/mysql)
