# 安装失败，提示go mod 错误

执行
```shell
go get -u github.com/bilibili/kratos/tool/kratos
```
出现以下错误时
```shell
go: github.com/prometheus/client_model@v0.0.0-20190220174349-fd36f4220a90: parsing go.mod: missing module line
go: github.com/remyoudompheng/bigfft@v0.0.0-20190806203942-babf20351dd7e3ac320adedbbe5eb311aec8763c: parsing go.mod: missing module line
```
如果你使用了https://goproxy.io/ 代理,那你要使用其他代理来替换它，然后删除GOPATH目录下的mod缓存文件夹（`go clean --modcache`）,然后重新执行安装命令

代理列表

```
export GOPROXY=https://mirrors.aliyun.com/goproxy/
export GOPROXY=https://goproxy.cn/
export GOPROXY=https://goproxy.io/
```

