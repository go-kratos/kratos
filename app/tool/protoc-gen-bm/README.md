# Deprecated: 请大佬移步 [https://git.bilibili.co/platform/go-common/tree/master/app/tool/bmproto](https://git.bilibili.co/platform/go-common/tree/master/app/tool/bmproto) 生成的代码不一样但是功能一致, 可以无缝迁移!

# protoc-gen-bm

通过 protobuf 自动生成 bm server，统一数据传输模型.

### 安装

```bash
go install go-common/app/tool/protoc-gen-bm
```

### 如何生成 bm server 代码 （example）

```bash
# 生成 examples/helloworld bm server 代码
protoc -I. -Ithird_party/googleapis --bm_out=:. examples/helloworld/api/v1/helloworld.proto
# 使用 jsonpb
protoc -I. -Ithird_party/googleapis --bm_out=jsonpb=true:. examples/helloworld/api/v1/helloworld.proto
```

### 给 protobuf 添加 http 描述示例

```protobuf
service Messaging {
  rpc GetMessage(GetMessageRequest) returns (Message) {
    option (google.api.http) = {
        get: "/message";
    };
  }
}
message GetMessageRequest {
  string name = 1; // Mapped to URL path.
}
message Message {
  string text = 1; // The resource content.
}
```
http 到 gRPC 的映射如下

| HTTP                    | gRPC                      |
|-------------------------|---------------------------|
| `GET /message?name=foo` | `GetMessage(name: "foo")` |

对于 GET 之类不包含 Body 的请求，所有的 Request Message 将被映射到 URL Query 中。
对于嵌套的字段通过 `.` 分割 例如:

```protobuf
service Messaging {
  rpc GetMessage(GetMessageRequest) returns (Message) {
    option (google.api.http) = {
        get:"/messages"
    };
  }
}
message GetMessageRequest {
  message SubMessage {
    string subfield = 1;
  }
  string message_id = 1;
  int64 revision = 2;
  SubMessage sub = 3;
}
```
| HTTP                                                          | gRPC                                                                            |
|---------------------------------------------------------------|---------------------------------------------------------------------------------|
| `GET /messages?message_id=123456&revision=2&sub.subfield=foo` | `GetMessage(message_id: "123456" revision: 2 sub: SubMessage(subfield: "foo"))` |

对于包含 Body 的 http 请求，body 字段可以用来指定数据映射

```protobuf
service Messaging {
  rpc CreateMessage(CreateMessageRequest) returns (Message) {
    option (google.api.http) = {
      post: "/messages"
      body: "*"
    };
  }
}
message CreateMessageRequest {
  string message_id = 1; 
  string text = 2; 
}
```
| HTTP                                                         | gRPC                                               |
|--------------------------------------------------------------|----------------------------------------------------|
| `POST /v1/messages {"message_id": "123456", "text": "Hi!" }` | `CreateMessage(message_id: "123456", text: "Hi!")` |

### Response 响应
Json Response
```json
{
    "code": 0,
    "message": "this is message",
    "data": {// Response Message Marshal as JSON}
}
```
Protobuf Response 
```protobuf
message PB {
	int64 Code = 1;
	string Message = 2;
	uint64 TTL = 3;
	google.protobuf.Any Data = 4;
}
```

### 与 google.api.http 不一致的地方

- 因为 bm 不支持 restful 格式 API，所以不支持映射参数到 Path
- Response 多了层包装，而不是直接返回 message 主体，message 主体被放到 data 中

### 注意事项

#### bazel 编译问题

目前自动生成的 BUILD 的文件无法正常编译，需要手动修改，以 example/helloworld 项目为例，需要对 BUILD 文件进行以下修改

```diff
--- a/app/tool/protoc-gen-bm/examples/helloworld/api/v1/BUILD
+++ b/app/tool/protoc-gen-bm/examples/helloworld/api/v1/BUILD
@@ -13,8 +13,8 @@ load(
 proto_library(
     name = "v1_proto",
     srcs = ["helloworld.proto"],
-    tags = ["automanaged"],
-    deps = ["google/api/annotations.proto"],
+    tags = ["manual"],
+    deps = ["@go_googleapis//google/api:annotations_proto"],
 )

 go_proto_library(
@@ -22,14 +22,14 @@ go_proto_library(
     compilers = ["@io_bazel_rules_go//proto:go_grpc"],
     importpath = "go-common/app/tool/protoc-gen-bm/examples/helloworld/api/v1",
     proto = ":v1_proto",
-    tags = ["automanaged"],
-    deps = ["google/api/annotations.proto"],
+    tags = ["manual"],
+    deps = ["@go_googleapis//google/api:annotations_go_proto"],
 )

 go_library(
     name = "go_default_library",
     srcs = ["helloworld.pb.bm.go"],
-    embed = ["v1_go_proto"],
+    embed = [":v1_go_proto"],
     importpath = "go-common/app/tool/protoc-gen-bm/examples/helloworld/api/v1",
     tags = ["automanaged"],
     visibility = ["//visibility:public"],
```
恩、不解释了，应该都能看懂

#### 需要自定义 form tag

因为 blademaster 在 Bind parameters 时默认是大小写敏感的，而且没有自动转换驼峰与下划线，所以在定义 Proto message 时需要用 gogo 自定义一些 tag

```protobuf
message User {
    // mid
    int64 mid = 1 [(gogoproto.moretags) = "form:\"mid\" validate:\"required,min=1\""];
}
```

参考 https://github.com/gogo/protobuf/blob/master/extensions.md

### TODO

- [ ] 完善错误提示
- [ ] 更加完善的 google.api.http 支持

### Roadmap

- [x] 修改已有的工具实现读取 google.api.http option 生成 bm server
- [ ] 添加 bilibili.api.extra option，扩展 google.api.http

### 已知问题

* 复杂结构问题，由于 http from 表达能力有限，无法很好的描述有些复杂的 protobuf message，解决: 内嵌 protobuf 或者 使用 json？
* 老项目兼容问题，有些项目会返回类似 map[string]interface{} 之类的结构，无法再 protobuf，解决: 内嵌 json ?

### 参考文档
* [google.api.http 定义](https://github.com/googleapis/googleapis/blob/master/google/api/http.proto)
* [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)

### client 代码生成

#### iOS 代码生成

TODO

#### Android 代码生成

TODO

/cc @liugang  @zhoujiahui
