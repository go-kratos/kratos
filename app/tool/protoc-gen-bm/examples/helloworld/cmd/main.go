package main

import (
	"log"

	"go-common/app/tool/protoc-gen-bm/examples/helloworld/api"
	"go-common/app/tool/protoc-gen-bm/examples/helloworld/service"
	bm "go-common/library/net/http/blademaster"
)

func main() {
	engine := bm.NewServer(nil)
	// 注册 middleware 支持正则匹配
	engine.Inject("^/echo", func(c *bm.Context) {
		// do something
	})
	s := new(service.Service)
	v1.RegisterHelloBMServer(engine, s)
	if err := engine.Run("127.0.0.1:8000"); err != nil {
		log.Fatal(err)
	}
}
