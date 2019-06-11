# 准备工作

推荐使用[kratos工具](kratos-tool.md)快速生成项目，如我们生成一个叫`kratos-demo`的项目。目录结构如下：

```
├── CHANGELOG.md
├── CONTRIBUTORS.md
├── LICENSE
├── README.md
├── cmd
│   ├── cmd
│   └── main.go
├── configs
│   ├── application.toml
│   ├── grpc.toml
│   ├── http.toml
│   ├── log.toml
│   ├── memcache.toml
│   ├── mysql.toml
│   └── redis.toml
├── go.mod
├── go.sum
└── internal
    ├── dao
    │   └── dao.go
    ├── model
    │   └── model.go
    ├── server
    │   └── http
    │       └── http.go
    └── service
        └── service.go
```

# 路由

创建项目成功后，进入`internal/server/http`目录下，打开`http.go`文件，其中有默认生成的`blademaster`模板。其中：

```go
engine = bm.DefaultServer(hc.Server)
initRouter(engine)
if err := engine.Start(); err != nil {
    panic(err)
}
```

是bm默认创建的`engine`及启动代码，我们看`initRouter`初始化路由方法，默认实现了：

```go
func initRouter(e *bm.Engine) {
	e.Ping(ping) // engine自带的"/ping"接口，用于负载均衡检测服务健康状态
	g := e.Group("/kratos-demo") // e.Group 创建一组 "/kratos-demo" 起始的路由组
	{
		g.GET("/start", howToStart) // g.GET 创建一个 "kratos-demo/start" 的路由，默认处理Handler为howToStart方法
	}
}
```

bm的handler方法，结构如下：

```go
func howToStart(c *bm.Context) // handler方法默认传入bm的Context对象
```

### Ping

engine自带Ping方法，用于设置`/ping`路由的handler，该路由统一提供于负载均衡服务做健康检测。服务是否健康，可自定义`ping handler`进行逻辑判断，如检测DB是否正常等。

```go
func ping(c *bm.Context) {
    if some DB check not ok {
        c.AbortWithStatus(503)
    }
}
```

# 默认路由

默认路由有：

* /metrics 用于prometheus信息采集
* /metadata 可以查看所有注册的路由信息

查看加载的所有路由信息：

```shell
curl 'http://127.0.0.1:8000/metadata'
```

输出：

```json
{
    "code": 0,
    "message": "0",
    "ttl": 1,
    "data": {
        "/kratos-demo/start": {
            "method": "GET"
        },
        "/metadata": {
            "method": "GET"
        },
        "/metrics": {
            "method": "GET"
        },
        "/ping": {
            "method": "GET"
        }
    }
}
```

# 性能分析

启动时默认监听了`2333`端口用于`pprof`信息采集，如：

```shell
go tool pprof http://127.0.0.1:8000/debug/pprof/profile
```

改变端口可以使用flag，如：`-http.perf=tcp://0.0.0.0:12333`

# 扩展阅读

[bm模块说明](blademaster-mod.md)  
[bm中间件](blademaster-mid.md)  
[bm基于pb生成](blademaster-pb.md)  

-------------

[文档目录树](summary.md)
