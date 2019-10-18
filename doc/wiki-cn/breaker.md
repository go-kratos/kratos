## 熔断器/Breaker
熔断器是为了当依赖的服务已经出现故障时，主动阻止对依赖服务的请求。保证自身服务的正常运行不受依赖服务影响，防止雪崩效应。

## kratos内置breaker的组件
一般情况下直接使用kratos的组件时都自带了熔断逻辑，并且在提供了对应的breaker配置项。
目前在kratos内集成熔断器的组件有:
- RPC client: pkg/net/rpc/warden/client
- Mysql client：pkg/database/sql
- Tidb client：pkg/database/tidb
- Http client：pkg/net/http/blademaster

## 使用说明
```go
 //初始化熔断器组
 //一组熔断器公用同一个配置项，可从分组内取出单个熔断器使用。可用在比如mysql主从分离等场景。
 brkGroup := breaker.NewGroup(&breaker.Config{}) 
 //为每一个连接指定一个brekaker
 //此处假设一个客户端连接对象实例为conn
 //breakName定义熔断器名称 一般可以使用连接地址
 breakName = conn.Addr
 conn.breaker = brkGroup.Get(breakName)
 
 //在连接发出请求前判断熔断器状态
 if err = conn.breaker.Allow(); err != nil {
		return
  }
 
 //连接执行成功或失败将结果告知braker
 if(respErr != nil){
      conn.breaker.MarkFailed()
 }else{
      conn.breaker.MarkSuccess()
 }
 
```

## 配置说明
```go
type Config struct {
	SwitchOff bool // 熔断器开关,默认关 false.

	K float64  //触发熔断的错误率（K = 1 - 1/错误率）

	Window  xtime.Duration //统计桶窗口时间
	Bucket  int  //统计桶大小
	Request int64 //触发熔断的最少请求数量（请求少于该值时不会触发熔断）
}
```

