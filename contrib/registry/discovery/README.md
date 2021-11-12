## Discovery Registry 

This module implements a `registry.Registrar` and `registry.Discovery` interface in kratos based `bilibili/discovery`.

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/go-kratos/kratos/contrib/registry/discovery/v2)

### Quick Start

**_Register a service_**

```go
import (
	"github.com/go-kratos/kratos/contrib/registry/discovery/v2"
)

func main() {
	logger := log.NewStdLogger(os.Stdout)
	logger = log.With(logger, "service", "example.registry.discovery")
	
	// initialize a registry
	r := discovery.New(&discovery.Config{
		Nodes:  []string{"0.0.0.0:7171"},
		Env:    "dev",
		Region: "sh1",
		Zone:   "zone1",
		Host:   "hostname",
	}, logger)

	// construct srv instance
	// ...
	
	app := kratos.New(
		kratos.Name("helloworld"),
		kratos.Server(
			httpSrv,
			grpcSrv,
		),
		kratos.Metadata(map[string]string{"color": "gray"}),
		// use Registrar
		kratos.Registrar(r),
	)
	
	if err := app.Run(); err != nil {
		log.NewHelper(logger).Fatal(err)
	}	
}
```

**_Discover a service_**

```go
import (
	"github.com/go-kratos/kratos/contrib/registry/discovery/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

func main() {
	// initialize a discovery
	r := discovery.New(&discovery.Config{
		Nodes:  []string{"0.0.0.0:7171"},
		Env:    "dev",
		Region: "sh1",
		Zone:   "zone1",
		Host:   "localhost",
	}, nil)

	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///appid"),
		// use discovery
		grpc.WithDiscovery(r),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	
	// request and log
}
```

### Config explain

```go
type Config struct {
	Nodes  []string // discovery nodes address
	Region string   // region of the service, sh
	Zone   string   // zone of region, sh001
	Env    string   // env of service, dev, prod and etc
	Host   string   // hostname of service
}
```

### References 

- [bilibili/discovery](https://github.com/bilibili/discovery)

