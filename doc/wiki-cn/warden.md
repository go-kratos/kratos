# 背景

我们需要统一的rpc服务，经过选型讨论决定直接使用成熟的跨语言的gRPC。

# 概览

* 不改gRPC源码，基于接口进行包装集成trace、log、prom等组件
* 打通自有服务注册发现系统[discovery](https://github.com/bilibili/discovery)
* 实现更平滑可靠的负载均衡算法
  
# 拦截器

gRPC暴露了两个拦截器接口，分别是：

* `grpc.UnaryServerInterceptor`服务端拦截器
* `grpc.UnaryClientInterceptor`客户端拦截器
  
基于两个拦截器可以针对性的定制公共模块的封装代码，比如`warden/logging.go`是通用日志逻辑。

[warden拦截器](warden-mid.md)

# 服务发现

gRPC暴露了服务发现的接口`resolver.Resolver`，`warden/resolver/resolver.go`实现了该接口，并基于了`pkg/naming/naming.go`内的`Resolver`接口进行`Fetch``Watch`等操作。

`pkg/naming/discovery/discovery.go`内实现了`pkg/naming/naming.go`内的`Resolver`接口，使用[discovery](https://github.com/bilibili/discovery)来进行服务发现。

注意：`pkg/naming/naming.go`内的`Resolver`接口是`kratos`的一层封装，暴露的接口主要：

* 相对原生`resolver.Resolver`内`ResolveNow`更友好的方法`Fetch``Watch`
* 统一应用的实例信息结构体`naming.Instance`

想要用非[discovery](https://github.com/bilibili/discovery)的请参考下面文档进行开发。

[warden服务发现](warden-resolver.md)

# 负载均衡

实现了`wrr`和`p2c`两种算法，默认使用`p2c`。

[warden负载均衡](warden-balancer.md)

# 扩展阅读

[warden快速开始](warden-quickstart.md) [warden拦截器](warden-mid.md) [warden负载均衡](warden-balancer.md) [warden基于pb生成](warden-pb.md) [warden服务发现](warden-resolver.md)

-------------

[文档目录树](summary.md)
