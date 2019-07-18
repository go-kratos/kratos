package etcd

import (
	"context"
	"fmt"
	"github.com/bilibili/kratos/pkg/naming"
	"github.com/coreos/etcd/clientv3"
	"testing"
	"time"
)

func TestNew(t *testing.T) {

	config := &clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
		DialTimeout:time.Second*3,
	}
	builder,err := New(config)

	if(err != nil){
		fmt.Println("etcd 连接失败")
		return
	}
	app1 := builder.Build("app1")

	go func() {
		fmt.Printf("Watch \n")
		for {
			select {
			case <-app1.Watch():
				fmt.Printf("app1 节点发生变化 \n")
			}

		}

	}()
	time.Sleep(time.Second)

	app1Cancel, err := builder.Register(context.Background(), &naming.Instance{
		AppID:    "app1",
		Hostname: "h1",
		Zone:     "z1",
	})

	fmt.Println(2222)
	app2Cancel, err := builder.Register(context.Background(), &naming.Instance{
		AppID:    "app2",
		Hostname: "h5",
		Zone:     "z3",
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(1111)
	app2 := builder.Build("app2")
	fmt.Println(1111)

	go func() {
		fmt.Println("节点列表")
		for {
			fmt.Printf("app1: ")
			r1, _ := app1.Fetch(context.Background())
			if r1 != nil {
				for z, ins := range r1.Instances {
					fmt.Printf("zone: %s :", z)
					for _, in := range ins {
						fmt.Printf("app: %s host %s \n", in.AppID, in.Hostname)
					}
				}
			}
			fmt.Printf("app2: ")
			r2, _ := app2.Fetch(context.Background())
			if r2 != nil {
				for z, ins := range r2.Instances {
					fmt.Printf("zone: %s :", z)
					for _, in := range ins {
						fmt.Printf("app: %s host %s \n", in.AppID, in.Hostname)
					}
				}
			}
			time.Sleep(time.Second)
		}
	}()

	fmt.Println(1111)
	time.Sleep(time.Second * 5)
	fmt.Println("取消app1")
	app1Cancel()

	time.Sleep(time.Second * 10)
	fmt.Println("取消app2")
	app2Cancel()

	fmt.Println(1111)
	time.Sleep(30*time.Second)
}
