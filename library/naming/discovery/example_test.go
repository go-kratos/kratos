package discovery_test

import (
	"context"
	"fmt"
	"go-common/library/naming"
	"go-common/library/naming/discovery"
	"time"
)

// this example creates a registry service to register instance info
// to discovery server.
// when the program is about to exit,registry.Cancel should be called.
func Example() {
	var c = &discovery.Config{
		Nodes:  []string{"api.bilibili.co"},
		Zone:   "sh001",
		Env:    "pre",
		Key:    "0c4b8fe3ff35a4b6",
		Secret: "b370880d1aca7d3a289b9b9a7f4d6812",
	}
	var ins = &naming.Instance{
		AppID: "main.arch.test2",
		Addrs: []string{
			"grpc://127.0.0.1:8080",
		},
		Version: "1",
		Metadata: map[string]string{
			"weight": "128",
			"color":  "blue",
		},
	}

	d := discovery.New(c)
	cacenl, err := d.Register(context.TODO(), ins)
	if err != nil {
		return
	}
	defer cacenl()
	//start to Serve
	time.Sleep(time.Second * 5)
}

// this example creates a discovery client to poll instances from discovery server.
func ExampleDiscovery() {
	d := discovery.Build("1231234")
	ch := d.Watch()
	for {
		<-ch
		ins, ok := d.Fetch(context.TODO())
		if ok {
			fmt.Println("new instances found:", ins)
		}
	}
}
