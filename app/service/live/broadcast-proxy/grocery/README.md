# sven.go

可能是最最精简的配置中心 SDK 了，简单接入。


# 集成
```go
package main

import (
	"nano-repo/grocery"
	"log"
)

func main() {
	sven, err := grocery.NewSvenClient("13586", "sh001", "dev",
		"docker-1", "7c41388b593d562120bec1bcb355e538")
	if err != nil {
		panic(err)
	}
	//Get the latest configuration with Config method anytime
	c := sven.Config()
	log.Printf("Initial version:%d", c.Version)
	log.Printf("Initial config :%v", c.Config)

	go func(){
	    //Get configuration change event with ConfigNotify method
		for config := range sven.ConfigNotify() {
			log.Printf("New version:%d", config.Version)
			log.Printf("New config: %v", config.Config)
		}
	}()

	go func(){
		for e := range sven.LogNotify() {
		     log.Printf("Sven log return, level:%v, message:%v", e.Level, e.Message)
		}
	}()
	quit := make(chan struct{})
	<- quit
}

```

