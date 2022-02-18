# Nacos Config

```go
import (
	kconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
)


sc := []constant.ServerConfig{
	*constant.NewServerConfig("127.0.0.1", 8848),
}

cc := &constant.ClientConfig{
	NamespaceId:         "public", //namespace id
	TimeoutMs:           5000,
	NotLoadCacheAtStart: true,
	LogDir:              "/tmp/nacos/log",
	CacheDir:            "/tmp/nacos/cache",
	RotateTime:          "1h",
	MaxAge:              3,
	LogLevel:            "debug",
}

// a more graceful way to create naming client
client, err := clients.NewConfigClient(
	vo.NacosClientParam{
		ClientConfig:  cc,
		ServerConfigs: sc,
	},
)
if err != nil {
	log.Panic(err)
}
```