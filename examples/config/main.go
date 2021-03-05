package main

import (
	"flag"
	"log"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"gopkg.in/yaml.v2"
)

var flagconf string

func init() {
	flag.StringVar(&flagconf, "conf", "config.yaml", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
		config.WithDecoder(func(kv *config.KeyValue, v map[string]interface{}) error {
			return yaml.Unmarshal(kv.Value, v)
		}),
	)
	if err := c.Load(); err != nil {
		panic(err)
	}
	// struct
	var v struct {
		Serivce struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"service"`
	}
	if err := c.Scan(&v); err != nil {
		panic(err)
	}

	log.Printf("config: %+v", v)
	// key/value
	name, err := c.Value("service.name").String()
	if err != nil {
		panic(err)
	}
	log.Printf("service: %s", name)
	// watch key
	if err := c.Watch("service.name", func(key string, value config.Value) {
		log.Printf("config changed: %s = %v\n", key, value)
	}); err != nil {
		panic(err)
	}

	<-make(chan struct{})
}
