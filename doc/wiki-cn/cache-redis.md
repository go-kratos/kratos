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

创建项目成功后，进入项目中的configs目录，打开redis.toml，我们可以看到：

```toml
[Client]
	name = "kratos-demo"
	proto = "tcp"
	addr = "127.0.0.1:6389"
	idle = 10
	active = 10
	dialTimeout = "1s"
	readTimeout = "1s"
	writeTimeout = "1s"
	idleTimeout = "10s"
```

在该配置文件中我们可以配置redis的连接方式proto、连接地址addr、连接池的闲置连接数idle、最大连接数active以及各类超时。

## 初始化

进入项目的internal/dao目录，打开redis.go，其中：

```go
var cfg struct {
    Client *memcache.Config
}
checkErr(paladin.Get("redis.toml").UnmarshalTOML(&rc))
```
使用paladin配置管理工具将上文中的redis.toml中的配置解析为我们需要使用的配置。

```go
// Dao dao.
type Dao struct {
	redis       *redis.Pool
	redisExpire int32
}
```

在dao的主结构提中定义了redis的连接池对象和过期时间。

```go
d = &dao{
    // redis
    redis:       redis.NewPool(rc.Demo),
    redisExpire: int32(time.Duration(rc.DemoExpire) / time.Second),
}
```

使用kratos/pkg/cache/redis包的NewPool方法进行连接池对象的初始化，需要传入上文解析的配置。

## Ping

```go
// Ping ping the resource.
func (d *dao) Ping(ctx context.Context) (err error) {
	return d.pingRedis(ctx)
}

func (d *dao) pingRedis(ctx context.Context) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	if _, err = conn.Do("SET", "ping", "pong"); err != nil {
		log.Error("conn.Set(PING) error(%v)", err)
	}
	return
}
```

生成的dao层模板中自带了redis相关的ping方法，用于为负载均衡服务的健康监测提供依据，详见[blademaster](blademaster-quickstart.md)。

## 关闭

```go
// Close close the resource.
func (d *Dao) Close() {
	d.redis.Close()
}
```

在关闭dao层时，通过调用redis连接池对象的Close方法，我们可以关闭该连接池，从而释放相关资源。

# 常用方法

## 发送单个命令 Do

```go
// DemoIncrby .
func (d *dao) DemoIncrby(c context.Context, pid int) (err error) {
	cacheKey := keyDemo(pid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("INCRBY", cacheKey, 1); err != nil {
		log.Error("DemoIncrby conn.Do(INCRBY) key(%s) error(%v)", cacheKey, err)
	}
	return
}
```
如上为向redis server发送单个命令的用法示意。这里需要使用redis连接池的Get方法获取一个redis连接conn，再使用conn.Do方法即可发送一条指令。
注意，在使用该连接完毕后，需要使用conn.Close方法将该连接关闭。

## 批量发送命令 Pipeline

kratos/pkg/cache/redis包除了支持发送单个命令，也支持批量发送命令（redis pipeline)，比如：

```go
// DemoIncrbys .
func (d *dao) DemoIncrbys(c context.Context, pid int) (err error) {
	cacheKey := keyDemo(pid)
	conn := d.redis.Get(c)
	defer conn.Close()
    if err = conn.Send("INCRBY", cacheKey, 1); err != nil {
        return 
    }
    if err = conn.Send("EXPIRE", cacheKey, d.redisExpire); err != nil {
        return
    }
    if err = conn.Flush(); err != nil {
        log.Error("conn.Flush error(%v)", err)
        return
    }
    for i := 0; i < 2; i++ {
        if _, err = conn.Receive(); err != nil {
            log.Error("conn.Receive error(%v)", err)
            return
        }
    }
    return
}
```

和发送单个命令类似地，这里需要使用redis连接池的Get方法获取一个redis连接conn，在使用该连接完毕后，需要使用conn.Close方法将该连接关闭。

这里使用conn.Send方法将命令写入客户端的buffer（缓冲区）中，使用conn.Flush将客户端的缓冲区内的命令打包发送到redis server。redis server按顺序返回的reply可以使用conn.Receive方法进行接收和处理。


## 返回值转换 

kratos/pkg/cache/redis包中也提供了Scan方法将redis server的返回值转换为golang类型。

除此之外，kratos/pkg/cache/redis包提供了大量返回值转换的快捷方式:

### 单个查询 

单个查询可以使用redis.Uint64/Int64/Float64/Int/String/Bool/Bytes进行返回值的转换，比如：

```go
// GetDemo get
func (d *Dao) GetDemo(ctx context.Context, key string) (string, error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	return redis.String(conn.Do("GET", key))
}
```

### 批量查询

批量查询时候，可以使用redis.Int64s,Ints,Strings,ByteSlices方法转换如MGET，HMGET，ZRANGE，SMEMBERS等命令的返回值。
还可以使用StringMap, IntMap, Int64Map方法转换HGETALL命令的返回值，比如：

```go
// HGETALLDemo get 
func (d *Dao) HGETALLDemo(c context.Context, pid int64) (res map[string]int64, err error) {
	var (
		key  = keyDemo(pid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if res, err = redis.Int64Map(conn.Do("HGETALL", key)); err != nil {
		log.Error("HGETALL %v failed error(%v)", key, err)
	}
	return
}
```

# 扩展阅读

[memcache模块说明](cache-mc.md)  

-------------

[文档目录树](summary.md)
