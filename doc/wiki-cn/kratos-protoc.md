### kratos tool protoc

```shell
# generate all
kratos tool protoc api.proto
# generate gRPC
kratos tool protoc --grpc api.proto
# generate BM HTTP
kratos tool protoc --bm api.proto
# generate ecode
kratos tool protoc --ecode api.proto
# generate swagger
kratos tool protoc --swagger api.proto
```

执行生成如 `api.pb.go/api.bm.go/api.swagger.json/api.ecode.go` 的对应文件，需要注意的是：`ecode`生成有固定规则，需要首先是`enum`类型，且`enum`名字要以`ErrCode`结尾，如`enum UserErrCode`。详情可见：[example](https://github.com/bilibili/kratos/tree/master/example/protobuf)

> 该工具在Windows/Linux下运行，需提前安装好 [protobuf](https://github.com/google/protobuf) 工具

`kratos tool protoc`本质上是拼接好了`protoc`命令然后进行执行，在执行时会打印出对应执行的`protoc`命令，如下可见：

```shell
protoc --proto_path=$GOPATH --proto_path=$GOPATH/github.com/bilibili/kratos/third_party --proto_path=. --bm_out=:. api.proto
protoc --proto_path=$GOPATH --proto_path=$GOPATH/github.com/bilibili/kratos/third_party --proto_path=. --gofast_out=plugins=grpc:. api.proto
protoc --proto_path=$GOPATH --proto_path=$GOPATH/github.com/bilibili/kratos/third_party --proto_path=. --bswagger_out=:. api.proto
protoc --proto_path=$GOPATH --proto_path=$GOPATH/github.com/bilibili/kratos/third_party --proto_path=. --ecode_out=:. api.proto
```

-------------

[文档目录树](summary.md)
