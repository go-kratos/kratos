package main

import (
	"log"

	"github.com/go-kratos/kratos/contrib/config/apollo/v2"
	"github.com/go-kratos/kratos/v2/config"
)

func main() {
	c := config.New(
		config.WithSource(
			apollo.NewSource(
				apollo.WithAppID("kratos"),
				apollo.WithCluster("dev"),
				apollo.WithEndpoint("http://localhost:8080"),
				apollo.WithNamespace("application"),
				apollo.WithEnableBackup(),
				apollo.WithSecret("895da1a174934ababb1b1223f5620a45"),
			),
		),
	)
	if err := c.Load(); err != nil {
		panic(err)
	}
	// Get a value associated with the key
	name, err := c.Value("name").String()
	if err != nil {
		panic(err)
	}
	log.Printf("service: %s", name)

	// Defines the config JSON Field
	var v struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}

	// Unmarshal the config to struct
	if err = c.Scan(&v); err != nil {
		panic(err)
	}
	log.Printf("config: %+v", v)

	// Get a value associated with the key
	name, err = c.Value("name").String()
	if err != nil {
		panic(err)
	}
	log.Printf("service: %s", name)

	// watch key
	if err = c.Watch("name", func(key string, value config.Value) {
		n, e := value.String()
		if e != nil {
			panic(e)
		}
		log.Printf("config changed: %s = %s\n", key, n)
	}); err != nil {
		panic(err)
	}

	<-make(chan struct{})
}
