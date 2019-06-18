# 介绍

基于proto文件可以快速生成`bm`框架对应的代码，提前需要准备以下工作：

* 安装`kratos tool protoc`工具，请看[kratos工具](kratos-tool.md)
* 编写`proto`文件，示例可参考[kratos-demo内proto文件](https://github.com/bilibili/kratos-demo/blob/master/api/api.proto)

### kratos工具说明

`kratos tool protoc`工具可以生成`warden` `bm` `swagger`对应的代码和文档，想要单独生成`bm`代码只需加上`--bm`如：

```shell
# generate BM HTTP
kratos tool protoc --bm api.proto
```

### proto文件说明

请注意想要生成`bm`代码，需要特别在`proto`的`service`内指定`google.api.http`配置，如下：

```go
service Demo {
	rpc SayHello (HelloReq) returns (.google.protobuf.Empty);
	rpc SayHelloURL(HelloReq) returns (HelloResp) {
        option (google.api.http) = {     // 该配置指定SayHelloURL方法对应的url
            get:"/kratos-demo/say_hello" // 指定url和请求方式为GET
        };
    };
}
```

# 使用

建议在项目`api`目录下编写`proto`文件及生成对应的代码，可参考[kratos-demo内的api目录](https://github.com/bilibili/kratos-demo/tree/master/api)。

执行命令后生成的`api.bm.go`代码，注意其中的`type DemoBMServer interface`和`RegisterDemoBMServer`，其中：

* `DemoBMServer`接口，包含`proto`文件内配置了`google.api.http`选项的所有方法
* `RegisterDemoBMServer`方法提供注册`DemoBMServer`接口的实现对象，和`bm`的`Engine`用于注册路由
* `DemoBMServer`接口的实现，一般为`internal/service`内的业务逻辑代码，需要实现`DemoBMServer`接口

使用`RegisterDemoBMServer`示例代码请参考[kratos-demo内的http](https://github.com/bilibili/kratos-demo/blob/master/internal/server/http/server.go)内的如下代码：

```go
engine = bm.DefaultServer(hc.Server)
pb.RegisterDemoBMServer(engine, svc)
initRouter(engine)
```

`internal/service`内的`Service`结构实现了`DemoBMServer`接口可参考[kratos-demo内的service](https://github.com/bilibili/kratos-demo/blob/master/internal/service/service.go)内的如下代码：

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

# 文档

基于同一份`proto`文件还可以生成对应的`swagger`文档，运行命令如下：

```shell
# generate swagger
kratos tool protoc --swagger api.proto
```

该命令将生成对应的`swagger.json`文件，可用于`swagger`工具通过WEBUI的方式打开使用，可运行命令如下：

```shell
kratos tool swagger serve api/api.swagger.json
```

# 扩展阅读

[bm快速开始](blademaster-quickstart.md)  
[bm模块说明](blademaster-mod.md)  
[bm中间件](blademaster-mid.md)  

-------------

[文档目录树](summary.md)
