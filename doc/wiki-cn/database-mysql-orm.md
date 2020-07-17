# 准备工作

推荐使用[kratos工具](kratos-tool.md)快速生成项目，如我们生成一个叫`kratos-demo`的项目。目录结构如下：

```
├── CHANGELOG.md
├── OWNERS
├── README.md
├── api
│   ├── api.bm.go
│   ├── api.pb.go
│   ├── api.proto
│   └── client.go
├── cmd
│   ├── cmd
│   └── main.go
├── configs
│   ├── application.toml
│   ├── db.toml
│   ├── grpc.toml
│   ├── http.toml
│   ├── memcache.toml
│   └── redis.toml
├── go.mod
├── go.sum
├── internal
│   ├── dao
│   │   ├── dao.bts.go
│   │   ├── dao.go
│   │   ├── db.go
│   │   ├── mc.cache.go
│   │   ├── mc.go
│   │   └── redis.go
│   ├── di
│   │   ├── app.go
│   │   ├── wire.go
│   │   └── wire_gen.go
│   ├── model
│   │   └── model.go
│   ├── server
│   │   ├── grpc
│   │   │   └── server.go
│   │   └── http
│   │       └── server.go
│   └── service
│       └── service.go
└── test
    └── docker-compose.yaml
```

# 开始使用

## 配置

创建项目成功后，进入项目中的configs目录，mysql.toml，更改为：

```toml
[demo]
   dsn = "{user}:{password}@tcp(127.0.0.1:3306)/{database}?timeout=1s&readTimeout=1s&writeTimeout=1s&parseTime=true&loc=Local&charset=utf8mb4,utf8"
   active = 20
   idle = 10
```

在该配置文件中我们不在区分读写dsn，可以设置为读写分离地址

## 初始化

进入项目的internal/dao目录，打开db.go，其中：

```go
type DB struct {
    DSN    string 
    Active int
    Idle   int 
} 
var cfg struct {
    DB *DB
}
checkErr(paladin.Get("db.toml").UnmarshalTOML(&dc))
```
使用paladin配置管理工具将上文中的db.toml中的配置解析为我们需要使用db的相关配置。



采用gorm1.9.14版本、并集成trace链路追踪，与原有区别多传一个ctx参数
```go
import "github.com/go-kratos/kratos/pkg/database/sql/gorm"

db, _ := gorm.Open("mysql", "{user}:{password}@tcp(127.0.0.1:3306)/{database}?timeout=1s&readTimeout=1s&writeTimeout=1s&parseTime=true&loc=Local&charset=utf8mb4,utf8")
db.DB().SetMaxIdleConns(idle)
db.DB().SetMaxOpenConns(active)

// 查询
db.Where("id = ?", 1).First(ctx, &user)
// 新增
db.Create(ctx, &User{ID: 1, Name: "hello"})
// 删除
db.Delete(ctx, &user)
// 更改
db.Model(&user).Updates(ctx, User{Name: "hello", Age: 18})

...

```

# TODO：补充常用方法

# 扩展阅读

[tidb模块说明](database-tidb.md)
[hbase模块说明](database-hbase.md)

-------------

[文档目录树](summary.md)