# 介绍

kratos包含了一批好用的工具集，比如项目一键生成、基于proto生成http&grpc代码，生成缓存回源代码，生成memcache执行代码，生成swagger文档等。

# 获取工具

执行以下命令，即可快速安装好`kratos`工具
```shell
go get -u github.com/bilibili/kratos/tool/kratos
```

那么接下来让我们快速开始熟悉工具的用法~

# kratos本体

`kratos`是所有工具集的本体，就像`go`一样，拥有执行各种子工具的能力，如`go build`和`go tool`。先让我们看看`-h`的输出：

```
NAME:
   kratos - kratos tool

USAGE:
   kratos [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
     new, n        create new project
     build, b      kratos build
     run, r        kratos run
     tool, t       kratos tool
     version, v    kratos version
     self-upgrade  kratos self-upgrade
     help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

可以看到`kratos`有如：`new` `build` `run` `tool`等在内的COMMANDS，那么接下来一一演示如何使用。

# kratos new

`kratos new`是快速创建一个项目的命令，执行如下：

```shell
kratos new kratos-demo
```

即可快速在当前目录生成一个叫`kratos-demo`的项目。此外还支持指定owner和path，如下：

```shell
kratos new kratos-demo -o YourName -d YourPath
```

注意，`kratos new`默认是不会生成通过protobuf定义的`grpc`和`bm`示例代码的，如需生成请加`--proto`，如下：

```shell
kratos new kratos-demo -o YourName -d YourPath --proto
```

特别注意，如果不是macos系统，生成的示例项目`api`目录下的`proto`文件并不会自动生成对应的`.pb.go`和`.bm.go`文件，需要参考以下说明进行生成。

[protoc说明](protoc.md)

# kratos build&run

`kratos build`和`kratos run`是`go build`和`go run`的封装，并无特别用途。

# kratos tool

`kratos tool`是基于proto生成http&grpc代码，生成缓存回源代码，生成memcache执行代码，生成swagger文档等工具集，先看下`kratos tool`的执行效果：

```
swagger(已安装): swagger api文档 Author(goswagger.io) [2019/05/05]
protoc(已安装): 快速方便生成pb.go和bm.go的protoc封装，windows、Linux请先安装protoc工具 Author(kratos) [2019/05/04]
kratos(已安装): Kratos工具集本体 Author(kratos) [2019/04/02]

安装工具: kratos tool install demo
执行工具: kratos tool demo
安装全部工具: kratos tool install all

详细文档： https://github.com/bilibili/kratos/blob/master/doc/wiki-cn/kratos-tool.md
```

***小小说明：如未安装工具，第一次运行也可自动安装，不需要特别执行install***

目前已经集成的工具有：
* `protoc`用于快速生成`*.pb.go`和`*.bm.go`文件，但目前不支持windows，Linux也需要先自己安装`protoc`工具。
* `swagger`用于显示自动生成的BM API接口文档，通过 `kratos tool swagger serve api/api.swagger.json` 可以访问到文档。
* TODOs...

### kratos tool protoc

该命令运行没其他参数，直接`kratos tool protoc`运行即可。但使用前需特别说明：

* 该工具在Windows/Linux下运行，需提前安装好`protoc`工具

该工具实际是一段`shell`脚本，其中自动将`protoc`命令进行了拼接，识别了需要`*.proto`的目录和当前目录下的`proto`文件，最终会拼接为如下命令进行执行：

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

大家也可以参考该命令进行`proto`生成`*.pb.go`和`*.bm.go`文件，也可以参考[protoc说明](protoc.md)。

### Tool examples
```shell
// new a project
cd $GOPATH
kratos new kratos-demo -o Tinker --proto

// build & run
cd kratos-demo
kratos run

// swagger docs
kratos tool swagger serve kratos-demo/api/api.swagger.json

// generate proto
cd kratos-demo/api
kratos tool protoc api.proto
```

-------------

[文档目录树](summary.md)
