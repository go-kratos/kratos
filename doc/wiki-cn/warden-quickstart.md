# 准备工作

推荐使用[kratos工具](kratos-tool.md)快速生成带`grpc`的项目，如我们生成一个叫`kratos-demo`的项目。

```
kratos new kratos-demo --proto
```

# pb文件

创建项目成功后，进入`api`目录下可以看到`api.proto`和`api.pb.go`和`generate.go`文件，其中：

* `api.proto`是gRPC server的描述文件
* `api.pb.go`是基于`api.proto`生成的代码文件
* `generate.go`是用于`kratos tool protoc`执行`go generate`进行代码生成的临时文件

接下来可以将以上三个文件全部删除或者保留`generate.go`，之后编写自己的proto文件，确认proto无误后，进行代码生成：

* 可直接执行`kratos tool protoc`，该命令会调用protoc工具生成`.pb.go`文件
* 如`generate.go`没删除，也可以执行`go generate`命令，将调用`kratos tool protoc`工具进行代码生成

[kratos工具请看](kratos-tool.md)

### 如没看kprotoc文档，请看下面这段话

`kratos tool protoc`用于快速生成`pb.go`文件，但目前windows和Linux需要先自己安装`protoc`工具，具体请看[protoc说明](protoc.md)。

# 注册server

进入`internal/server/grpc`目录打开`server.go`文件，可以看到以下代码，只需要替换以下注释内容就可以启动一个gRPC服务。

```go
package grpc

import (
	pb "kratos-demo/api"
	"kratos-demo/internal/service"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/net/rpc/warden"
)

// New new a grpc server.
func New(svc *service.Service) *warden.Server {
	var rc struct {
		Server *warden.ServerConfig
	}
	if err := paladin.Get("grpc.toml").UnmarshalTOML(&rc); err != nil {
		if err != paladin.ErrNotExist {
			panic(err)
		}
	}
	ws := warden.NewServer(rc.Server)
	// 注意替换这里：
	// RegisterDemoServer方法是在"api"目录下代码生成的
	// 对应proto文件内自定义的service名字，请使用正确方法名替换
	pb.RegisterDemoServer(ws.Server(), svc)
	ws, err := ws.Start()
	if err != nil {
		panic(err)
	}
	return ws
}
```

### 注册注意

```go
// SayHello grpc demo func.
func (s *Service) SayHello(ctx context.Context, req *pb.HelloReq) (reply *empty.Empty, err error) {
	reply = new(empty.Empty)
	fmt.Printf("hello %s", req.Name)
	return
}
```

请进入`internal/service`内找到`SayHello`方法，注意方法的入参和出参，都是按照gRPC的方法声明对应的：

* 第一个参数必须是`context.Context`，第二个必须是proto内定义的`message`对应生成的结构体
* 第一个返回值必须是proto内定义的`message`对应生成的结构体，第二个参数必须是`error`
* 在http框架bm中，如果共用proto文件生成bm代码，那么也可以直接使用该service方法

建议service严格按照此格式声明方法使其能够在bm和warden内共用。

# client调用

请进入`internal/dao`方法内，一般对资源的处理都会在这一层封装。  
对于`client`端，前提必须有对应`proto`文件生成的代码，那么有两种选择：

* 拷贝proto文件到自己项目下并且执行代码生成
* 直接import服务端的api package

> 这也是业务代码我们加了一层`internal`的关系，服务对外暴露的只有接口

不管哪一种方式，以下初始化gRPC client的代码建议伴随生成的代码存放在统一目录下：

```go
package dao

import (
	"context"

	"github.com/bilibili/kratos/pkg/net/rpc/warden"

	"google.golang.org/grpc"
)

// target server addrs.
const target = "direct://default/127.0.0.1:9000,127.0.0.1:9091" // NOTE: example

// NewClient new member grpc client
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (DemoClient, error) {
	client := warden.NewClient(cfg, opts...)
	conn, err := client.Dial(context.Background(), target)
	if err != nil {
		return nil, err
	}
	// 注意替换这里：
	// NewDemoClient方法是在"api"目录下代码生成的
	// 对应proto文件内自定义的service名字，请使用正确方法名替换
	return NewDemoClient(conn), nil
}
```

其中，`target`为gRPC用于服务发现的目标，使用标准url资源格式提供给resolver用于服务发现。`warden`默认使用`direct`直连方式，直接与`server`端进行连接。如果在使用其他服务发现组件请看[warden服务发现](warden-resolver.md)。

有了初始化`Client`的代码，我们的`Dao`对象即可进行初始化和使用，以下以直接import服务端api包为例：

```go
package dao

import(
	demoapi "kratos-demo/api"
	grpcempty "github.com/golang/protobuf/ptypes/empty"
	"github.com/bilibili/kratos/pkg/net/rpc/warden"

	"github.com/pkg/errors"
)

type Dao struct{
	demoClient demoapi.DemoClient
}

// New account dao.
func New() (d *Dao) {
	cfg := &warden.ClientConfig{}
	paladin.Get("grpc.toml").UnmarshalTOML(cfg)
	d = &Dao{}
	var err error
	if d.demoClient, err = demoapi.NewClient(cfg); err != nil {
		panic(err)
	}
	return
}

// SayHello say hello.
func (d *Dao) SayHello(c context.Context, req *demoapi.HelloReq) (resp *grpcempty.Empty, err error) {
	if resp, err = d.demoClient.SayHello(c, req); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}
```

如此在`internal/service`层就可以进行资源的方法调用。

# 扩展阅读

[warden拦截器](warden-mid.md)  
[warden基于pb生成](warden-pb.md)  
[warden服务发现](warden-resolver.md)  
[warden负载均衡](warden-balancer.md) 

-------------

[文档目录树](summary.md)
