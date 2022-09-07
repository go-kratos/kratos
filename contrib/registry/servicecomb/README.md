# Servicecomb Registry

## example
### server
```go
package main

import (
	"log"

	"github.com/go-chassis/sc-client"
	"github.com/go-kratos/kratos/contrib/registry/servicecomb/v2"
	"github.com/go-kratos/kratos/v2"
)

func main() {
	c, err := sc.NewClient(sc.Options{
		Endpoints: []string{"127.0.0.1:30100"},
	})
	if err != nil {
		log.Panic(err)
	}
	r := servicecomb.NewRegistry(c)
	app := kratos.New(
		kratos.Name("helloServicecomb"),
		kratos.Registrar(r),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

```
### client
```go
package main

import (
	"context"
	"log"

	"github.com/go-chassis/sc-client"
	"github.com/go-kratos/kratos/contrib/registry/servicecomb/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

func main() {
	c, err := sc.NewClient(sc.Options{
		Endpoints: []string{"127.0.0.1:30100"},
	})
	if err != nil {
		log.Panic(err)
	}
	r := servicecomb.NewRegistry(c)
	ctx := context.Background()
	conn, err := grpc.DialInsecure(
		ctx,
		grpc.WithEndpoint("discovery:///helloServicecomb"),
		grpc.WithDiscovery(r),
	)
	if err != nil {
		return
	}
	defer conn.Close()
}

```