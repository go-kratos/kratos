# Etcd Config

```go
import (
	"log"

	clientv3 "go.etcd.io/etcd/client/v3"

	cfg "github.com/go-kratos/kratos/contrib/config/etcd/v3"
	"github.com/go-kratos/kratos/v3/config"
)

// create an etcd client
client, err := clientv3.New(clientv3.Config{
    Endpoints:   []string{"127.0.0.1:2379"},
    DialTimeout: time.Second,
})
if err != nil {
    log.Fatal(err)
}

// configure the source, "path" is required
source, err := cfg.New(client, cfg.WithPath("/app-config"), cfg.WithPrefix(true))
if err != nil {
    log.Fatalln(err)
}

// create a config instance with source
c := config.New(config.WithSource(source))
defer c.Close()

// load sources before get
if err := c.Load(); err != nil {
    log.Fatalln(err)
}

// acquire config value
foo, err := c.Value("/app-config").String()
if err != nil {
    log.Fatalln(err)
}

log.Println(foo)
```
