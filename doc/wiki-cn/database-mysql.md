# 准备工作

推荐使用[kratos工具](kratos-tool.md)快速生成项目，如我们生成一个叫`kratos-demo`的项目。目录结构如下：

```
├── CHANGELOG.md
├── OWNERS
├── README.md
├── api
│   ├── api.bm.go
│   ├── api.pb.go
│   ├── api.proto
│   └── client.go
├── cmd
│   ├── cmd
│   └── main.go
├── configs
│   ├── application.toml
│   ├── db.toml
│   ├── grpc.toml
│   ├── http.toml
│   ├── memcache.toml
│   └── redis.toml
├── go.mod
├── go.sum
├── internal
│   ├── dao
│   │   ├── dao.bts.go
│   │   ├── dao.go
│   │   ├── db.go
│   │   ├── mc.cache.go
│   │   ├── mc.go
│   │   └── redis.go
│   ├── di
│   │   ├── app.go
│   │   ├── wire.go
│   │   └── wire_gen.go
│   ├── model
│   │   └── model.go
│   ├── server
│   │   ├── grpc
│   │   │   └── server.go
│   │   └── http
│   │       └── server.go
│   └── service
│       └── service.go
└── test
    └── docker-compose.yaml
```

# 开始使用

## 配置

创建项目成功后，进入项目中的configs目录，mysql.toml，我们可以看到：

```toml
[demo]
	addr = "127.0.0.1:3306"
	dsn = "{user}:{password}@tcp(127.0.0.1:3306)/{database}?timeout=1s&readTimeout=1s&writeTimeout=1s&parseTime=true&loc=Local&charset=utf8mb4,utf8"
	readDSN = ["{user}:{password}@tcp(127.0.0.2:3306)/{database}?timeout=1s&readTimeout=1s&writeTimeout=1s&parseTime=true&loc=Local&charset=utf8mb4,utf8","{user}:{password}@tcp(127.0.0.3:3306)/{database}?timeout=1s&readTimeout=1s&writeTimeout=1s&parseTime=true&loc=Local&charset=utf8,utf8mb4"]
	active = 20
	idle = 10
	idleTimeout ="4h"
	queryTimeout = "200ms"
	execTimeout = "300ms"
	tranTimeout = "400ms"
```

在该配置文件中我们可以配置mysql的读和写的dsn、连接地址addr、连接池的闲置连接数idle、最大连接数active以及各类超时。

如果配置了readDSN，在进行读操作的时候会优先使用readDSN的连接。

## 初始化

进入项目的internal/dao目录，打开db.go，其中：

```go
var cfg struct {
    Client *sql.Config
}
checkErr(paladin.Get("db.toml").UnmarshalTOML(&dc))
```
使用paladin配置管理工具将上文中的db.toml中的配置解析为我们需要使用db的相关配置。

```go
// Dao dao.
type Dao struct {
	db          *sql.DB
}
```

在dao的主结构提中定义了mysql的连接池对象。

```go
d = &dao{
    db: sql.NewMySQL(dc.Demo),
}
```

使用kratos/pkg/database/sql包的NewMySQL方法进行连接池对象的初始化，需要传入上文解析的配置。

## Ping

```go
// Ping ping the resource.
func (d *dao) Ping(ctx context.Context) (err error) {
	return d.db.Ping(ctx)
}
```

生成的dao层模板中自带了mysql相关的ping方法，用于为负载均衡服务的健康监测提供依据，详见[blademaster](blademaster-quickstart.md)。

## 关闭

```go
// Close close the resource.
func (d *dao) Close() {
	d.db.Close()
}
```

在关闭dao层时，通过调用mysql连接池对象的Close方法，我们可以关闭该连接池，从而释放相关资源。

# 常用方法

## 单个查询 

```go
// GetDemo 用户角色
func (d *dao) GetDemo(c context.Context, did int64) (demo int8, err error) {
	err = d.db.QueryRow(c, _getDemoSQL, did).Scan(&demo)
	if err != nil && err != sql.ErrNoRows {
		log.Error("d.GetDemo.Query error(%v)", err)
		return
	}
	return demo, nil
}
```

db.QueryRow方法用于返回最多一条记录的查询，在QueryRow方法后使用Scan方法即可将mysql的返回值转换为Golang的数据类型。

当mysql查询不到对应数据时，会返回sql.ErrNoRows，如果不需处理，可以参考如上代码忽略此error。

## 批量查询

```go
// ResourceLogs ResourceLogs.
func (d *dao) GetDemos(c context.Context, dids []int64) (demos []int8, err error) {
	rows, err := d.db.Query(c, _getDemosSQL, dids)
	if err != nil {
		log.Error("query  error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var tmpD int8
		if err = rows.Scan(&tmpD); err != nil {
			log.Error("scan demo log error(%v)", err)
			return
		}
		demos = append(demos, tmpD)
	}
	return
}
```

db.Query方法一般用于批量查询的场景，返回*sql.Rows和error信息。
我们可以使用rows.Next()方法获得下一行的返回结果，并且配合使用rows.Scan()方法将该结果转换为Golang的数据类型。当没有下一行时，rows.Next方法将返回false，此时循环结束。

注意，在使用完毕rows对象后，需要调用rows.Close方法关闭连接，释放相关资源。

## 执行语句 

```go
// DemoExec exec
func (d *Dao) DemoExec(c context.Context, id int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _demoUpdateSQL, id)
	if err != nil {
		log.Error("db.DemoExec.Exec(%s) error(%v)", _demoUpdateSQL, err)
		return
	}
	return res.RowsAffected()
}
```

执行UPDATE/DELETE/INSERT语句时，使用db.Exec方法进行语句执行，返回*sql.Result和error信息：

```go

// A Result summarizes an executed SQL command.
type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}
```

Result接口支持获取影响行数和LastInsertId（一般用于获取Insert语句插入数据库后的主键ID）


## 事务

kratos/pkg/database/sql包支持事务操作，具体操作示例如下：

开启一个事务：

```go
tx := d.db.Begin()
if err = tx.Error; err != nil {
    log.Error("db begin transcation failed, err=%+v", err)
    return
}
```

在事务中执行语句：

```go
res, err := tx.Exec(_demoSQL, did)
if err != nil {
    return
}
rows := res.RowsAffected()
```

提交事务：

```go
if err = tx.Commit().Error; err!=nil{
    log.Error("db commit transcation failed, err=%+v", err)
}
```

回滚事务：

```go
if err = tx.Rollback().Error; err!=nil{
    log.Error("db rollback failed, err=%+v", rollbackErr)
}
```

# 扩展阅读

[tidb模块说明](database-tidb.md)
[hbase模块说明](database-hbase.md)

-------------

[文档目录树](summary.md)
