#### database/tidb

##### 项目简介
TiDB数据库驱动 对mysql驱动进行封装

##### 功能
1. 支持discovery服务发现 多节点直连
2. 支持通过lvs单一地址连接
3. 支持prepare绑定多个节点
4. 支持动态增减节点负载均衡
5. 日志区分运行节点

##### 编译环境
> 请只用golang v1.8.x以上版本编译执行。

##### 依赖包
> 1.[Go-MySQL-Driver](https://github.com/go-sql-driver/mysql)
