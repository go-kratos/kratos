# protoc

`protobuf`是Google官方出品的一种轻便高效的结构化数据存储格式，可以用于结构化数据串行化，或者说序列化。它很适合做数据存储或 RPC 数据交换格式。可用于通讯协议、数据存储等领域的语言无关、平台无关、可扩展的序列化结构数据格式。

使用`protobuf`，需要先书写`.proto`文件，然后编译该文件。编译`proto`文件则需要使用到官方的`protoc`工具，安装文档请参看：[google官方protoc工具](https://github.com/protocolbuffers/protobuf#protocol-compiler-installation)。

注意：`protoc`是用于编辑`proto`文件的工具，它并不具备生成对应语言代码的能力，所以正常都是`protoc`配合对应语言的代码生成工具来使用，如Go语言的[gogo protobuf](https://github.com/gogo/protobuf)，请先点击按文档说明安装。

安装好对应工具后，我们可以进入`api`目录，执行如下命令：

```shell
export $KRATOS_HOME = kratos路径
export $KRATOS_DEMO = 项目路径

// 生成：api.pb.go
protoc -I$GOPATH/src:$KRATOS_HOME/tool/protobuf/pkg/extensions:$KRATOS_DEMO/api --gogofast_out=plugins=grpc:$KRATOS_DEMO/api $KRATOS_DEMO/api/api.proto

// 生成：api.bm.go
protoc -I$GOPATH/src:$KRATOS_HOME/tool/protobuf/pkg/extensions:$KRATOS_DEMO/api --bm_out=$KRATOS_DEMO/api $KRATOS_DEMO/api/api.proto

// 生成：api.swagger.json
protoc -I$GOPATH/src:$KRATOS_HOME/tool/protobuf/pkg/extensions:$KRATOS_DEMO/api --bswagger_out=$KRATOS_DEMO/api $KRATOS_DEMO/api/api.proto
```

请注意替换`/Users/felix/work/go/src`目录为你本地开发环境对应GOPATH目录，其中`--gogofast_out`意味着告诉`protoc`工具需要使用`gogo protobuf`的工具生成代码。

-------------

[文档目录树](summary.md)
