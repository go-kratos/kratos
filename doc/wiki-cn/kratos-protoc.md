### kratos tool protoc

```
// generate all
kratos tool protoc api.proto
// generate gRPC
kratos tool protoc --grpc api.proto
// generate BM HTTP
kratos tool protoc --bm api.proto
// generate swagger
kratos tool protoc --swagger api.proto
```
执行对应生成 `api.pb.go/api.bm.go/api.swagger.json` 源文档。

> 该工具在Windows/Linux下运行，需提前安装好 protobuf 工具

该工具实际是一段`shell`脚本，其中自动将`protoc`命令进行了拼接，识别了需要的`*.proto`文件和当前目录下的`proto`文件，最终会拼接为如下命令进行执行：

```shell
export $KRATOS_HOME = kratos路径
export $KRATOS_DEMO = 项目路径

// 生成：api.pb.go
protoc -I$GOPATH/src:$KRATOS_HOME/third_party:$KRATOS_DEMO/api --gofast_out=plugins=grpc:$KRATOS_DEMO/api $KRATOS_DEMO/api/api.proto

// 生成：api.bm.go
protoc -I$GOPATH/src:$KRATOS_HOME/third_party:$KRATOS_DEMO/api --bm_out=$KRATOS_DEMO/api $KRATOS_DEMO/api/api.proto

// 生成：api.swagger.json
protoc -I$GOPATH/src:$KRATOS_HOME/third_party:$KRATOS_DEMO/api --bswagger_out=$KRATOS_DEMO/api $KRATOS_DEMO/api/api.proto
```

大家也可以参考该命令进行`proto`生成，也可以参考 [protobuf](https://github.com/google/protobuf) 官方参数。


-------------

[文档目录树](summary.md)
