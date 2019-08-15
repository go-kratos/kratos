# 背景

基于bm的handler机制，可以自定义很多middleware(中间件)进行通用的业务处理，比如用户登录鉴权。接下来就以鉴权为例，说明middleware的写法和用法。

# 写自己的中间件

middleware本质上就是一个handler，接口和方法声明如下代码：

```go
// Handler responds to an HTTP request.
type Handler interface {
	ServeHTTP(c *Context)
}

// HandlerFunc http request handler function.
type HandlerFunc func(*Context)

// ServeHTTP calls f(ctx).
func (f HandlerFunc) ServeHTTP(c *Context) {
	f(c)
}
```

1. 实现了`Handler`接口，可以作为engine的全局中间件使用：`engine.Use(YourHandler)`
2. 声明为`HandlerFunc`方法，可以作为engine的全局中间件使用：`engine.UseFunc(YourHandlerFunc)`，也可以作为router的局部中间件使用：`e.GET("/path", YourHandlerFunc)`

简单示例代码如下：

```go
type Demo struct {
	Key   string
	Value   string
}
// ServeHTTP implements from Handler interface
func (d *Demo) ServeHTTP(ctx *bm.Context) {
	ctx.Set(d.Key, d.Value)
}

e := bm.DefaultServer(nil)
d := &Demo{}

// Handler使用如下：
e.Use(d)

// HandlerFunc使用如下：
e.UseFunc(d.ServeHTTP)
e.GET("/path", d.ServeHTTP)

// 或者只有方法
myHandler := func(ctx *bm.Context) {
    // some code
}
e.UseFunc(myHandler)
e.GET("/path", myHandler)
```

# 全局中间件

在blademaster的`server.go`代码中，有以下代码：

```go
func DefaultServer(conf *ServerConfig) *Engine {
	engine := NewServer(conf)
	engine.Use(Recovery(), Trace(), Logger())
	return engine
}
```

会默认创建一个`bm engine`，并注册`Recovery(), Trace(), Logger()`三个middlerware用于全局handler处理，优先级从前到后。如果想要将自定义的middleware注册进全局，可以继续调用Use方法如下：

```go
engine.Use(YourMiddleware())
```

此方法会将`YourMiddleware`追加到已有的全局middleware后执行。如果需要全部自定义全局执行的middleware，可以使用`NewServer`方法创建一个无middleware的engine对象，然后使用`engine.Use/UseFunc`进行注册。

# 局部中间件

先来看一段鉴权伪代码示例([auth示例代码位置](https://github.com/bilibili/kratos/tree/master/example/blademaster/middleware/auth))：

```go
func Example() {
	myHandler := func(ctx *bm.Context) {
		mid := metadata.Int64(ctx, metadata.Mid)
		ctx.JSON(fmt.Sprintf("%d", mid), nil)
	}

	authn := auth.New(&auth.Config{DisableCSRF: false})

	e := bm.DefaultServer(nil)

	// "/user"接口必须保证登录用户才能访问，那么我们加入"auth.User"来确保用户鉴权通过，才能进入myHandler进行业务逻辑处理
	e.GET("/user", authn.User, myHandler)
	// "/guest"接口访客用户就可以访问，但如果登录用户我们需要知道mid，那么我们加入"auth.Guest"来尝试鉴权获取mid，但肯定会继续执行myHandler进行业务逻辑处理
	e.GET("/guest", authn.Guest, myHandler)

    // "/owner"开头的所有接口，都需要进行登录鉴权才可以被访问，那可以创建一个group并加入"authn.User"
	o := e.Group("/owner", authn.User)
	o.GET("/info", myHandler) // 该group创建的router不需要再显示的加入"authn.User"
	o.POST("/modify", myHandler) // 该group创建的router不需要再显示的加入"authn.User"

	e.Start()
}
```

# 内置中间件

## Recovery

代码位于`pkg/net/http/blademaster/recovery.go`内，用于recovery panic。会被`DefaultServer`默认注册，建议使用`NewServer`的话也将其作为首个中间件注册。

## Trace

代码位于`pkg/net/http/blademaster/trace.go`内，用于trace设置，并且实现了`net/http/httptrace`的接口，能够收集官方库内的调用栈详情。会被`DefaultServer`默认注册，建议使用`NewServer`的话也将其作为第二个中间件注册。

## Logger

代码位于`pkg/net/http/blademaster/logger.go`内，用于请求日志记录。会被`DefaultServer`默认注册，建议使用`NewServer`的话也将其作为第三个中间件注册。

## CSRF

代码位于`pkg/net/http/blademaster/csrf.go`内，用于防跨站请求。如要使用如下：

```go
e := bm.DefaultServer(nil)
// 挂载自适应限流中间件到 bm engine，使用默认配置
csrf := bm.CSRF([]string{"bilibili.com"}, []string{"/a/api"})
e.Use(csrf)
// 或者
e.GET("/api", csrf, myHandler)
```

## CORS

代码位于`pkg/net/http/blademaster/cors.go`内，用于跨域允许请求。请注意该：
1. 使用该中间件进行全局注册后，可"省略"单独为`OPTIONS`请求注册路由，如示例一。
2. 使用该中间单独为某路由注册，需要为该路由再注册一个`OPTIONS`方法的同路径路由，如示例二。

示例一：
```go
e := bm.DefaultServer(nil)
// 挂载自适应限流中间件到 bm engine，使用默认配置
cors := bm.CORS([]string{"github.com"})
e.Use(cors)
// 该路由可以默认针对 OPTIONS /api 的跨域请求支持
e.POST("/api", myHandler)
```

示例二：
```go
e := bm.DefaultServer(nil)
// 挂载自适应限流中间件到 bm engine，使用默认配置
cors := bm.CORS([]string{"github.com"})
// e.Use(cors) 不进行全局注册
e.OPTIONS("/api", cors, myHandler) // 需要单独为/api进行OPTIONS方法注册
e.POST("/api", cors, myHandler)
```

## 自适应限流

更多关于自适应限流的信息可参考：[kratos 自适应限流](/doc/wiki-cn/ratelimit.md)。如要使用如下：

```go
e := bm.DefaultServer(nil)
// 挂载自适应限流中间件到 bm engine，使用默认配置
limiter := bm.NewRateLimiter(nil)
e.Use(limiter.Limit())
// 或者
e.GET("/api", csrf, myHandler)
```

# 扩展阅读

[bm快速开始](blademaster-quickstart.md)   
[bm模块说明](blademaster-mod.md)  
[bm基于pb生成](blademaster-pb.md)  

-------------

[文档目录树](summary.md)
