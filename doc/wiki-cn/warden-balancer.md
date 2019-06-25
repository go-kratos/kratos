# Warden Balancer

## 介绍
grpc-go内置了round-robin轮询，但由于自带的轮询算法不支持权重，也不支持color筛选等需求，故需要重新实现一个负载均衡算法。

## WRR (Weighted Round Robin)
该算法在加权轮询法基础上增加了动态调节权重值，用户可以在为每一个节点先配置一个初始的权重分，之后算法会根据节点cpu、延迟、服务端错误率、客户端错误率动态打分，在将打分乘用户自定义的初始权重分得到最后的权重值。

## P2C (Pick of two choices)
本算法通过随机选择两个node选择优胜者来避免羊群效应，并通过ewma尽量获取服务端的实时状态。

服务端：
服务端获取最近500ms内的CPU使用率（需要将cgroup设置的限制考虑进去，并除于CPU核心数），并将CPU使用率乘与1000后塞入每次grpc请求中的的Trailer中夹带返回：
cpu_usage
uint64 encoded with string	
cpu_usage : 1000

客户端：
主要参数：
* server_cpu：通过每次请求中服务端塞在trailer中的cpu_usage拿到服务端最近500ms内的cpu使用率
* inflight：当前客户端正在发送并等待response的请求数（pending request）
* latency: 加权移动平均算法计算出的接口延迟
* client_success:加权移动平均算法计算出的请求成功率（只记录grpc内部错误，比如context deadline）

目前客户端，已经默认使用p2c负载均衡算法`grpc.WithBalancerName(p2c.Name)`：
```go
// NewClient returns a new blank Client instance with a default client interceptor.
// opt can be used to add grpc dial options.
func NewClient(conf *ClientConfig, opt ...grpc.DialOption) *Client {
	c := new(Client)
	if err := c.SetConfig(conf); err != nil {
		panic(err)
	}
	c.UseOpt(grpc.WithBalancerName(p2c.Name))
	c.UseOpt(opt...)
	c.Use(c.recovery(), clientLogging(), c.handle())
	return c
}
```

-------------

[文档目录树](summary.md)
