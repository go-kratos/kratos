# database/sql

## 背景
数据库驱动，进行封装加入了熔断、链路追踪和统计，以及链路超时。  
通常数据模块都写在`internal/dao`目录中，并提供对应的数据访问接口。

## MySQL
MySQL数据库驱动，支持读写分离、context、timeout、trace和统计功能，以及错误熔断防止数据库雪崩。  
[mysql client](database-mysql.md)  
[mysql client orm](database-mysql-orm.md)

## HBase
HBase客户端，支持trace、slowlog和统计功能。  
[hbase client](database-hbase.md)

## TiDB
TiDB客户端，支持服务发现和熔断功能。  
[tidb client](database-tidb.md)

-------------

[文档目录树](summary.md)
