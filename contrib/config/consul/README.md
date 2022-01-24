# Consul Config

```go
import (
    "github.com/go-kratos/kratos/contrib/config/consul/v2"
    "github.com/hashicorp/consul/api"
)
func main() {

    consulClient, err := api.NewClient(&api.Config{
    Address: "127.0.0.1:8500",
    })
    if err != nil {
        panic(err)
    }
    cs, err := consul.New(consulClient, consul.WithPath("app/cart/configs/"))
    //consul中需要标注文件后缀，kratos读取配置需要适配文件后缀
    //The file suffix needs to be marked, and kratos needs to adapt the file suffix to read the configuration.
    if err != nil {
        panic(err)
    }
    c := config.New(config.WithSource(cs))
}
```