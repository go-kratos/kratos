# 介绍

基于proto文件可以快速生成`warden`框架对应的代码，提前需要准备以下工作：

* 安装`kratos tool protoc`工具，请看[kratos工具](kratos-tool.md)
* 编写`proto`文件，示例可参考[kratos-demo内proto文件](https://github.com/bilibili/kratos-demo/blob/master/api/api.proto)

### kratos工具说明

`kratos tool protoc`工具可以生成`warden` `bm` `swagger`对应的代码和文档，想要单独生成`warden`代码只需加上`--grpc`如：

```shell
# generate gRPC
kratos tool protoc --grpc api.proto
```

# 使用

建议在项目`api`目录下编写`proto`文件及生成对应的代码，可参考[kratos-demo内的api目录](https://github.com/bilibili/kratos-demo/tree/master/api)。

执行命令后生成的`api.pb.go`代码，注意其中的`DemoClient`和`DemoServer`，其中：

* `DemoClient`接口为客户端调用接口，相对应的有`demoClient`结构体为其实现
* `DemoServer`接口为服务端接口声明，需要业务自己实现该接口的所有方法，`kratos`建议在`internal/service`目录下使用`Service`结构体实现

`internal/service`内的`Service`结构实现了`DemoServer`接口可参考[kratos-demo内的service](https://github.com/bilibili/kratos-demo/blob/master/internal/service/service.go)内的如下代码：

```go
// SayHelloURL bm demo func.
func (s *Service) SayHelloURL(ctx context.Context, req *pb.HelloReq) (reply *pb.HelloResp, err error) {
	reply = &pb.HelloResp{
		Content: "hello " + req.Name,
	}
	fmt.Printf("hello url %s", req.Name)
	return
}
```

更详细的客户端和服务端使用请看[warden快速开始](warden-quickstart.md)

# 扩展阅读

[warden快速开始](warden-quickstart.md) [warden拦截器](warden-mid.md) [warden负载均衡](warden-balancer.md) [warden服务发现](warden-resolver.md)

-------------

[文档目录树](summary.md)
