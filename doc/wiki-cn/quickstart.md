# 快速开始

快速使用kratos项目，可以使用`kratos`工具，如下：

```shell
go get -u github.com/bilibili/kratos/tool/kratos
cd $GOPATH/src
kratos new kratos-demo
```

根据提示可以快速创建项目，如[kratos-demo](https://github.com/bilibili/kratos-demo)就是通过工具创建生成。目录结构如下：

```
├── CHANGELOG.md 
├── OWNERS
├── README.md
├── api                     # api目录为对外保留的proto文件及生成的pb.go文件
│   ├── api.bm.go
│   ├── api.pb.go           # 通过go generate生成的pb.go文件
│   ├── api.proto
│   └── client.go
├── cmd
│   └── main.go             # cmd目录为main所在
├── configs                 # configs为配置文件目录
│   ├── application.toml    # 应用的自定义配置文件，可能是一些业务开关如：useABtest = true
│   ├── db.toml             # db相关配置
│   ├── grpc.toml           # grpc相关配置
│   ├── http.toml           # http相关配置
│   ├── memcache.toml       # memcache相关配置
│   └── redis.toml          # redis相关配置
├── go.mod
├── go.sum
└── internal                # internal为项目内部包，包括以下目录：
│   ├── dao                 # dao层，用于数据库、cache、MQ、依赖某业务grpc|http等资源访问
│   │   ├── dao.bts.go
│   │   ├── dao.go
│   │   ├── db.go
│   │   ├── mc.cache.go
│   │   ├── mc.go
│   │   └── redis.go
│   ├── di                  # 依赖注入层 采用wire静态分析依赖
│   │   ├── app.go
│   │   ├── wire.go         # wire 声明
│   │   └── wire_gen.go     # go generate 生成的代码
│   ├── model               # model层，用于声明业务结构体
│   │   └── model.go
│   ├── server              # server层，用于初始化grpc和http server
│   │   ├── grpc            # grpc层，用于初始化grpc server和定义method
│   │   │   └── server.go
│   │   └── http            # http层，用于初始化http server和声明handler
│   │       └── server.go
│   └── service             # service层，用于业务逻辑处理，且为方便http和grpc共用方法，建议入参和出参保持grpc风格，且使用pb文件生成代码
│       └── service.go
└── test                    # 测试资源层 用于存放测试相关资源数据 如docker-compose配置 数据库初始化语句等
    └── docker-compose.yaml
```

生成后可直接运行如下：

```shell
cd kratos-demo/cmd
go build
./cmd -conf ../configs
```

打开浏览器访问：[http://localhost:8000/kratos-demo/start](http://localhost:8000/kratos-demo/start)，你会看到输出了`Golang 大法好 ！！！`

[kratos工具](kratos-tool.md)

-------------

[文档目录树](summary.md)
