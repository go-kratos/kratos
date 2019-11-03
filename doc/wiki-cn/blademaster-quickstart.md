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
		g.GET("/start", howToStart) // g.GET 创建一个 "kratos-demo/start" 的路由，使用GET方式请求，默认处理Handler为howToStart方法
		g.POST("start", howToStart) // g.POST 创建一个 "kratos-demo/start" 的路由，使用POST方式请求，默认处理Handler为howToStart方法
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

# 路径参数

使用方式如下：

```go
func initRouter(e *bm.Engine) {
	e.Ping(ping)
	g := e.Group("/kratos-demo")
	{
		g.GET("/start", howToStart)

		// 路径参数有两个特殊符号":"和"*"
		// ":" 跟在"/"后面为参数的key，匹配两个/中间的值 或 一个/到结尾(其中不再包含/)的值
		// "*" 跟在"/"后面为参数的key，匹配从 /*开始到结尾的所有值，所有*必须写在最后且无法多个

		// NOTE：这是不被允许的，会和 /start 冲突
		// g.GET("/:xxx")

		// NOTE: 可以拿到一个key为name的参数。注意只能匹配到/param1/felix，无法匹配/param1/felix/hao(该路径会404)
		g.GET("/param1/:name", pathParam)
		// NOTE: 可以拿到多个key参数。注意只能匹配到/param2/felix/hao/love，无法匹配/param2/felix或/param2/felix/hao
		g.GET("/param2/:name/:value/:felid", pathParam)
		// NOTE: 可以拿到一个key为name的参数 和 一个key为action的路径。
		// NOTE: 如/params3/felix/hello，action的值为"/hello"
		// NOTE: 如/params3/felix/hello/hi，action的值为"/hello/hi"
		// NOTE: 如/params3/felix/hello/hi/，action的值为"/hello/hi/"
		g.GET("/param3/:name/*action", pathParam)
	}
}

func pathParam(c *bm.Context) {
	name, _ := c.Params.Get("name")
	value, _ := c.Params.Get("value")
	felid, _ := c.Params.Get("felid")
	action, _ := c.Params.Get("action")
	path := c.RoutePath // NOTE: 获取注册的路由原始地址，如: /kratos-demo/param1/:name
	c.JSONMap(map[string]interface{}{
		"name":   name,
		"value":  value,
		"felid":  felid,
		"action": action,
		"path":   path,
	}, nil)
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
