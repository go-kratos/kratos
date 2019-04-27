# 准备工作

推荐使用[kratos工具](kratos-tool.md)快速生成项目，如我们生成一个叫`kratos-demo`的项目。

生成目录结构如下：
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

# 扩展阅读

[bm模块说明](blademaster-mod.md) [bm中间件](blademaster-mid.md)  [bm基于pb生成](blademaster-pb.md)

-------------

[文档目录树](summary.md)
