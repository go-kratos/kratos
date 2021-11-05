# Nacos Registry

```go
import (
		"github.com/go-kratos/kratos/v2"
		"github.com/go-kratos/kratos/v2/transport/grpc"
		
		"github.com/nacos-group/nacos-sdk-go/clients"
		"github.com/nacos-group/nacos-sdk-go/common/constant"
    	"github.com/nacos-group/nacos-sdk-go/vo"
)

sc := []constant.ServerConfig{
	*constant.NewServerConfig("127.0.0.1", 8848),
}

cc := constant.ClientConfig{
	NamespaceId:         "public",
	TimeoutMs:           5000,
}

client, err := clients.NewNamingClient(
	vo.NacosClientParam{
		ClientConfig:  &cc,
		ServerConfigs: sc,
	},
)

if err != nil {
	log.Panic(err)
}

r := nacos.New(client)

// server
app := kratos.New(
	kratos.Name("helloworld"),
	kratos.Registrar(r),
)
if err := app.Run(); err != nil {
	log.Fatal(err)
}

// client
conn, err := grpc.DialInsecure(
	context.Background(),
	grpc.WithEndpoint("discovery:///helloworld"),
	grpc.WithDiscovery(r),
)
```