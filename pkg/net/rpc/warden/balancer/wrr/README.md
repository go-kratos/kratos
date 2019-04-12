#### business/warden/balancer/wrr

##### 项目简介

warden 的 weighted round robin负载均衡模块，主要用于为每个RPC请求返回一个Server节点以供调用

##### 编译环境

- **请只用 Golang v1.9.x 以上版本编译执行**

##### 依赖包

- [grpc](google.golang.org/grpc)