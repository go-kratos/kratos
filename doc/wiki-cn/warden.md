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

`warden`默认使用`direct`方式直连，正常线上都会使用第三方服务注册与发现中间件，`warden`内包含了[discovery](https://github.com/bilibili/discovery)的逻辑实现，想使用如`etcd`、`zookeeper`等也可以，都请看下面文档。

[warden服务发现](warden-resolver.md)

# 负载均衡

实现了`wrr`和`p2c`两种算法，默认使用`p2c`。

[warden负载均衡](warden-balancer.md)

# 扩展阅读

[warden快速开始](warden-quickstart.md) [warden拦截器](warden-mid.md) [warden负载均衡](warden-balancer.md) [warden基于pb生成](warden-pb.md) [warden服务发现](warden-resolver.md)

-------------

[文档目录树](summary.md)
