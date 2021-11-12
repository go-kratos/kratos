## Discovery Registry 

This module implements a `registry.Registrar` and `registry.Discovery` interface in kratos based `bilibili/discovery`.

### Quick Start

**_Register a service_**

```go
func main() {
	logger := log.NewStdLogger(os.Stdout)
	logger = log.With(logger, "service", "example.registry.discovery")

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
func main() {
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
		grpc.WithDiscovery(r),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	
	// request and log
}
```

### References 

- [bilibili/discovery](https://github.com/bilibili/discovery)

