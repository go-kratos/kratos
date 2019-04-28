# 目录

<!-- vim-markdown-toc GitLab -->

* [简介](#简介)
        * [HTTP访问](#http访问)
        * [grpc](#grpc)
* [安装](#安装)
* [用法](#用法)
    * [生成结果](#生成结果)
* [命名规范](#命名规范)
    * [proto包名与版本号](#proto包名与版本号)
    * [生成的go文件的包名](#生成的go文件的包名)
    * [多个proto文件](#多个proto文件)
* [其他特性](#其他特性)
    * [添加http框架的Middleware](#添加http框架的middleware)
    * [自定义Url或者指定http方法为post](#自定义url或者指定http方法为post)
    * [生成service模板](#生成service模板)
    * [form tag和json tag](#form-tag和json-tag)
    * [指定输入参数的约束条件](#指定输入参数的约束条件)
    * [同步Markdown文档到live-doc仓库（直播Only）](#同步markdown文档到live-doc仓库直播only)
    * [支持json做为输入](#支持json做为输入)
* [直播部门老的用法](#直播部门老的用法)
    * [兼容直播服务列表（按照discovery id）](#兼容直播服务列表按照discovery-id)

<!-- vim-markdown-toc -->

## 简介
根据protobuf文件，生成grpc和blademaster框架http代码及文档


```protobuf
syntax = "proto3";
package department.app;

option go_package = "api";

import "github.com/gogo/protobuf/gogoproto/gogo.proto";

service Greeter{
    // api 标题
    // api 说明
    rpc SayHello(HelloRequest) returns (HelloResponse);
}

message HelloRequest {
    // 请求参数说明
    string param1 = 1 [(gogoproto.moretags) = 'form:"param1"']; 
}

message HelloResponse {
    // 返回字段说明
    string ret_string = 1 [(gogoproto.jsontag) = 'ret_string'];
}
```

#### HTTP访问

```
GET  /department.app.Greeter/SayHello?param1=p1
响应
{
    "code": 0,
    "message": "ok",
    "data": {
        "ret_string": "anything"
    }
}
```
#### grpc
路径  /department.app.Greeter/SayHello

## 安装

```shell
go install go-common/app/tool/bmproto/...
```

## 用法

- cd 项目目录
- 在 api目录下新建api.proto文件(参见上面的例子) 例如 api/api.proto
- 运行 bmgen（在项目的任意位置）
- 创建 internal/service/greeter.go(这是写业务代码的地方)（也可以通过bmgen -t直接生成）


```go
import pb "path-to-project/api"
....

// 实现 pb.GreeterBMServer 和 grpc的 pb.GreeterServer 
type GreeterService struct {
    
}


func (s *GreeterService) SayHello(ctx context.Context, req *pb.SayHelloRequest)
(resp *pb.SayHelloResp, err error) {
    
}

```

- 在server/http.go 初始化代码(一般是`route`方法)里加入代码


```go
import pb "path-to-project/api"
import svc "path-to-project/internal/service"
......
pb.RegisterGreeterBMServer(engine, &svc.GreeterService{})
```

- 如果是grpc 在server/grpc/server.go 初始化里面加入代码

`pb.RegisterGreeterServer(grpcServer, &svc.GreeterService{})`

- 启动服务
- 访问接口 `curl 127.0.0.1:8000/department.app.Greeter/SayHello` (默认路由规则为 `/package.service/method`)

### 生成结果

```
project-
|------|--internal/service/greeter.go (使用bmgen -t 会生成，如果proto新增加方法，会自动往这里面添加模板代码）
       |--api/
           |--api.greeter.md (HTTP API文档)
           |--api.bm.go       
           |--api.pb.go
           |--api.proto           
```

## 命名规范
### proto包名与版本号
- DISCOVERY_ID 或者 DISCOVERY_ID.v*
- DISCOVERY_ID的构成为 `部门.服务` 并且去掉中划线
- 第一个版本不用加版本号，从第二个版本加v2
- **示例** 部门 department 服务 hello-world， 则 package为`department.helloworld`, 文件目录为api/
- 第二个版本package `department.helloworld.v2", 目录为api/v2
### 生成的go文件的包名
- golang一般原则上保持包名和目录名一致
- proto 可以指定`option go_package = "xxx"; `

比如对于api/api.proto `option go_package = "api"; `

对于api/v2/api.proto `option go_package = "v2"; `
### 多个proto文件
一个文件夹下面可以有多个proto文件，但是要满足以下约束

- 同目录下的proto package 一致
- message service 等定义不能重复（因为是在统一package下面）


## 其他特性
### 添加http框架的Middleware

在RegisterXXBMServer之前加入代码


```
bm.Inject(pb.PathGreeterSayHello, middleware1, middleware2)
```

### 自定义Url或者指定http方法为post


```protobuf
.....
package department.app;
....
import "google/api/annotations.proto";
....
service Greeter{
    rpc SayHelloCustomUrl(HelloRequest) returns (HelloResponse) {
        option (google.api.http) = {
            get:"/say_hello" // GET /say_hello
        };
    };
    
    rpc SayHelloPost(HelloRequest) returns (HelloResponse) {
        option (google.api.http) = {
            post:"" // POST /department.app.Greeter/SayHelloPost
        };
    };
}
```

### 生成service模板
`bmgen -t` 生成service模板代码在 internal/service/serviceName.go

### form tag和json tag
```
对于HTTP接口
现在请求字段需要加上form tag以解析请求参数，
响应参数需要加上json tag 以避免 字段为0或者空字符串时不显示，
这两个tag都建议和字段名保持一致
现在是必须加，将来考虑维护一个自己的proto仓库，以移除这个多余的tag
```


### 指定输入参数的约束条件

```protobuf
...
import "github.com/gogo/protobuf/gogoproto/gogo.proto"
...
message Request {
    int param1 = 1 [(gogoproto.moretags) = 'validate:"required"')]; // 参数必传，不能等于0
}
```

### 同步Markdown文档到live-doc仓库（直播Only）
`bmgen -l`

### 支持json做为输入

```
curl 127.0.0.1:8000/department.app.Greeter/SayHello -H "Content-Type: application/json" -d "{"param1":"p1"}" -X POST
```



## 直播部门老的用法
**对于以下"兼容直播服务列表中的服务"有效**


- URL：/xlive/项目名/v*/service开头小写/method
- 注册路由：使用RegisterXXXService而不是 RegisterXXBMServer
- middleware：不支持RegisterXXXMiddleware 而是 使用注解


```go
api/api.proto
service Greeter {
    // `method:"POST"` // 表示请求方法为POST
    // `midware:"user"`
    rpc SayHello(A) returns (B);
}

// server/http/http.go
import bm "go-common/library/net/http/blademaster"
....
userAuthMiddleware := xxxxx
pb.RegisterXXService(e, svc, map[string]bm.HandlerFunc{"user":userAuthMiddleware})
```


- 注解，在方法或者字段上方的注视添加和 go的tag格式一样的注解，实现一定的功能

  注解列表：

  | key          | 位置                  | 说明                                                         |
  | ------------ | --------------------- | ------------------------------------------------------------ |
  | midware      | rpc method上方        | midware:"auth,verify" 中间件，auth 是验证登录态，verify是校验签名， |
  | method       | rpc method上方        | method:"POST" 指定http请求方法                               |
  | mock         | 响应message的字段上方 | mock:"mockdata" mock数据，生成文档的时候有用                 |
  | internal     | 不建议继续使用        | 不建议继续使用                                               |
  | dynamic      | 不建议继续使用        | 不建议继续使用                                               |
  | dynamic_resp | 不建议继续使用        | 不建议继续使用                                               |



### 兼容直播服务列表（按照discovery id）
-  "live.webucenter"
-  "live.webroom"
-  "live.appucenter"
-  "live.appblink"
-  "live.approom"
-  "live.appinterface"
-  "live.liveadmin"
-  "live.resource"
-  "live.livedemo"
-  "live.lotteryinterface"


