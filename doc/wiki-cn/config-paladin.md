# Paladin SDK

## 配置模块化
进行配置的模块化是为了更好地管理配置，尽可能避免由修改配置带来的失误。  
在配置种类里，可以看到其实 环境配置 和 应用配置 已经由平台进行管理化。  
我们通常业务里只用配置 业务配置 和 在线配置 就可以了，之前我们大部分都是单个文件配置，而为了更好管理我们需要按类型进行拆分配置文件。  

例如：

| 名称 | 说明 |
|:------|:------|
| application.toml | 在线配置 |
| mysql.toml | 业务db配置 |
| hbase.toml | 业务hbase配置 |
| memcache.toml | 业务mc配置 |
| redis.toml | 业务redis配置 |
| http.toml | 业务http client/server/auth配置 |
| grpc.toml | 业务grpc client/server配置 |

## 使用方式

paladin 是一个config SDK客户端，包括了remote、file、mock几个抽象功能，方便使用本地文件或者远程配置中心，并且集成了对象自动reload功能。

### 远程配置中心
可以通过环境变量注入，例如：APP_ID/DEPLOY_ENV/ZONE/HOSTNAME，然后通过paladin实现远程配置中心SDK进行配合使用。

### 指定本地文件：
```shell
./cmd -conf=/data/conf/app/demo.toml
# or multi file
./cmd -conf=/data/conf/app/
```

### mock配置文件
```go
func TestMain(t *testing.M) {
    mock := make(map[string]string])
    mock["application.toml"] = `
        demoSwitch = false
        demoNum = 100
        demoAPI = "xxx"
    `
    paladin.DefaultClient = paladin.NewMock(mock)
}
```

### example main
```go
// main.go
func main() {
	flag.Parse()
    // 初始化paladin
    if err := paladin.Init(); err != nil {
        panic(err)
    }
    log.Init(nil) // debug flag: log.dir={path}
    defer log.Close()
}
```

### example HTTP/gRPC
```toml
# http.toml
[server]
    addr = "0.0.0.0:9000"
    timeout = "1s"
  
```

```go
// server.go
func NewServer() {
	// 默认配置用nil，这时读取HTTP/gRPC构架中的flag或者环境变量（可能是docker注入的环境变量，默认端口：8000/9000）
	engine := bm.DefaultServer(nil)

	// 除非自己要替换了配置，用http.toml
	var bc struct {
		Server *bm.ServerConfig
	}
	if err := paladin.Get("http.toml").UnmarshalTOML(&bc); err != nil {
		// 不存在时，将会为nil使用默认配置
		if err != paladin.ErrNotExist {
			panic(err)
		}
	}
	engine := bm.DefaultServer(bc.Server)
}
```

### example Service（在线配置热加载配置）
```go
# service.go
type Service struct {
	ac *paladin.Map
}

func New() *Service {
	// paladin.Map 通过atomic.Value支持自动热加载
	var ac = new(paladin.TOML)
	if err := paladin.Watch("application.toml", ac); err != nil {
		panic(err)
	}
	s := &Service{
		ac: ac,
	}
	return s
}

func (s *Service) Test() {
	sw, err := s.ac.Get("switch").Bool()
	if err != nil {
		// TODO
	}
	
	// or use default value
	sw := paladin.Bool(s.ac.Get("switch"), false)
}
```
