# warden proto 自动生成工具

[参考文档](http://info.bilibili.co/display/documentation/gRPC+Golang+Warden+Gen)

# protoc.sh

platform/go-common 仓库 pb 文件生成工具

默认使用 gofast, 可以使用 PROTOC_GEN 环境变量指定.

默认 Proto Import Path 为 go-common 目录, go-common/vendor 目录, /usr/local/include 以及proto文件所在目录, 可以通过 PROTO_PATH 环境变量指定.

### 安装

```bash
# 如果你只有一个 GOPATH 的话
ln -s ${GOPATH}/src/go-common/app/tool/warden/protoc.sh /usr/local/bin/protoc.sh && chmod +x /usr/local/bin/protoc.sh
```

### 使用方法

```
> cd {proto 文件所在目录}
> protoc.sh
```

### TODO

- [ ] 支持更多系统
- [ ] 纠正 pb 文件中错误的 import
