package main

import (
	"fmt"
	"log"

	_ "github.com/SeeMusic/kratos/v2/encoding/json"
	_ "github.com/SeeMusic/kratos/v2/encoding/yaml"

	"github.com/SeeMusic/kratos/contrib/config/apollo/v2"
	"github.com/SeeMusic/kratos/v2/config"
)

type bootstrap struct {
	Application struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"application"`

	Event struct {
		Key   string   `json:"key"`
		Array []string `json:"array"`
	} `json:"event"`

	Demo struct {
		Deep struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"deep"`
	} `json:"demo"`
}

func main() {
	c := config.New(
		config.WithSource(
			apollo.NewSource(
				apollo.WithAppID("kratos"),
				apollo.WithCluster("dev"),
				apollo.WithEndpoint("http://localhost:8080"),
				apollo.WithNamespace("application,event.yaml,demo.json"),
				apollo.WithEnableBackup(),
				apollo.WithSecret("ad75b33c77ae4b9c9626d969c44f41ee"),
			),
		),
	)
	var bc bootstrap
	if err := c.Load(); err != nil {
		panic(err)
	}

	scan(c, &bc)

	value(c, "application")
	value(c, "application.name")
	value(c, "event.array")
	value(c, "demo.deep")

	watch(c, "application")
	<-make(chan struct{})
}

func scan(c config.Config, bc *bootstrap) {
	err := c.Scan(bc)
	fmt.Printf("=========== scan result =============\n")
	fmt.Printf("err: %v\n", err)
	fmt.Printf("cfg: %+v\n\n", bc)
}

func value(c config.Config, key string) {
	fmt.Printf("=========== value result =============\n")
	v := c.Value(key).Load()
	fmt.Printf("key=%s, load: %+v\n\n", key, v)
}

func watch(c config.Config, key string) {
	if err := c.Watch(key, func(key string, value config.Value) {
		log.Printf("config(key=%s) changed: %s\n", key, value.Load())
	}); err != nil {
		panic(err)
	}
}
