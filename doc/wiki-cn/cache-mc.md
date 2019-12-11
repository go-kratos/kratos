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

创建项目成功后，进入项目中的configs目录，打开memcache.toml，我们可以看到：

```toml
[Client]
	name = "demo"
	proto = "tcp"
	addr = "127.0.0.1:11211"
	active = 50
	idle = 10
	dialTimeout = "100ms"
	readTimeout = "200ms"
	writeTimeout = "300ms"
    idleTimeout = "80s"
```
在该配置文件中我们可以配置memcache的连接方式proto、连接地址addr、连接池的闲置连接数idle、最大连接数active以及各类超时。

## 初始化

进入项目的internal/dao目录，打开mc.go，其中：

```go
var cfg struct {
    Client *memcache.Config
}
checkErr(paladin.Get("memcache.toml").UnmarshalTOML(&mc))
```
使用paladin配置管理工具将上文中的memcache.toml中的配置解析为我们需要使用的配置。

```go
// dao dao.
type dao struct {
	mc          *memcache.Memcache
	mcExpire    int32
}
```

在dao的主结构提中定义了memcache的连接池对象和过期时间。

```go
d = &dao{
    // memcache
    mc:       memcache.New(mc.Demo),
    mcExpire: int32(time.Duration(mc.DemoExpire) / time.Second),
}
```

使用kratos/pkg/cache/memcache包的New方法进行连接池对象的初始化，需要传入上文解析的配置。

## Ping

```go
// Ping ping the resource.
func (d *dao) Ping(ctx context.Context) (err error) {
	return d.pingMC(ctx)
}

func (d *dao) pingMC(ctx context.Context) (err error) {
	if err = d.mc.Set(ctx, &memcache.Item{Key: "ping", Value: []byte("pong"), Expiration: 0}); err != nil {
		log.Error("conn.Set(PING) error(%v)", err)
	}
	return
}
```

生成的dao层模板中自带了memcache相关的ping方法，用于为负载均衡服务的健康监测提供依据，详见[blademaster](blademaster-quickstart.md)。

## 关闭

```go
// Close close the resource.
func (d *Dao) Close() {
	d.mc.Close()
}
```

在关闭dao层时，通过调用memcache连接池对象的Close方法，我们可以关闭该连接池，从而释放相关资源。

# 常用方法

推荐使用[memcache代码生成器](kratos-genmc.md)帮助我们生成memcache操作的相关代码。

以下我们来逐一解析以下kratos/pkg/cache/memcache包中提供的常用方法。

## 单个查询

```go
// CacheDemo get data from mc
func (d *Dao) CacheDemo(c context.Context, id int64) (res *Demo, err error) {
	key := demoKey(id)
	res = &Demo{}
	if err = d.mc.Get(c, key).Scan(res); err != nil {
		res = nil
		if err == memcache.ErrNotFound {
			err = nil
		}
	}
	if err != nil {
		prom.BusinessErrCount.Incr("mc:CacheDemo")
		log.Errorv(c, log.KV("CacheDemo", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}
```

如上为代码生成器生成的进行单个查询的代码，使用到mc.Get(c,key)方法获得返回值，再使用scan方法将memcache的返回值转换为golang中的类型（如string，bool, 结构体等）。

## 批量查询使用

```go
replies, err := d.mc.GetMulti(c, keys)
for _, key := range replies.Keys() {
    v := &Demo{}
    err = replies.Scan(key, v)
}
```

如上为代码生成器生成的进行批量查询的代码片段，这里使用到mc.GetMulti(c,keys)方法获得返回值，与单个查询类似地，我们需要再使用scan方法将memcache的返回值转换为我们定义的结构体。

## 设置KV

```go
// AddCacheDemo Set data to mc
func (d *Dao) AddCacheDemo(c context.Context, id int64, val *Demo) (err error) {
	if val == nil {
		return
	}
	key := demoKey(id)
	item := &memcache.Item{Key: key, Object: val, Expiration: d.demoExpire, Flags: memcache.FlagJSON | memcache.FlagGzip}
	if err = d.mc.Set(c, item); err != nil {
		prom.BusinessErrCount.Incr("mc:AddCacheDemo")
		log.Errorv(c, log.KV("AddCacheDemo", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}
```

如上为代码生成器生成的添加结构体进入memcache的代码，这里需要使用到的是mc.Set方法进行设置。
这里使用的item为memcache.Item结构体，包含key, value, 超时时间（秒）, Flags。

### Flags


上文添加结构体进入memcache中，使用到的flags为：memcache.FlagJSON | memcache.FlagGzip代表着：使用json作为编码方式，gzip作为压缩方式。

Flags的相关常量在kratos/pkg/cache/memcache包中进行定义，包含编码方式如gob, json, protobuf，和压缩方式gzip。

```go
const(
	// Flag, 15(encoding) bit+ 17(compress) bit

	// FlagRAW default flag.
	FlagRAW = uint32(0)
	// FlagGOB gob encoding.
	FlagGOB = uint32(1) << 0
	// FlagJSON json encoding.
	FlagJSON = uint32(1) << 1
	// FlagProtobuf protobuf
	FlagProtobuf = uint32(1) << 2
	// FlagGzip gzip compress.
	FlagGzip = uint32(1) << 15
)
```

## 删除KV

```go
// DelCacheDemo delete data from mc
func (d *Dao) DelCacheDemo(c context.Context, id int64) (err error) {
	key := demoKey(id)
	if err = d.mc.Delete(c, key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		prom.BusinessErrCount.Incr("mc:DelCacheDemo")
		log.Errorv(c, log.KV("DelCacheDemo", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}
```
如上为代码生成器生成的从memcache中删除KV的代码，这里需要使用到的是mc.Delete方法。
和查询时类似地，当memcache中不存在参数中的key时，会返回error为memcache.ErrNotFound。如果不需要处理这种error，可以参考上述代码将返回出去的error置为nil。

# 扩展阅读

[memcache代码生成器](kratos-genmc.md)  
[redis模块说明](cache-redis.md)  

-------------

[文档目录树](summary.md)
