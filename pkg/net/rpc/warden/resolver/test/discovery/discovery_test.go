package discovery

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go-common/library/conf/env"
	"go-common/library/naming"
	"go-common/library/naming/discovery"
	"go-common/library/net/rpc/warden/resolver"

	gr "google.golang.org/grpc/resolver"
)

const (
	appID = "test.server"
)

type ClientConn struct {
	addrs []gr.Address
}

func (c *ClientConn) UpdateState(gr.State) {

}

func (c *ClientConn) NewAddress(addresses []gr.Address) {
	fmt.Println("get addr:", addresses)
	c.addrs = addresses

}

func (c *ClientConn) NewServiceConfig(serviceConfig string) {

}

func init() {
	ctx := context.Background()
	env.DeployEnv = "uat"
	env.Hostname = "host1"
	dis := discovery.New(&discovery.Config{
		Env:    "uat",
		Host:   "host1",
		Zone:   "sh001",
		Key:    "1",
		Secret: "1",
	})
	dis.Register(ctx, &naming.Instance{
		AppID: appID,
		Addrs: []string{
			"grpc://127.0.0.1:9000",
		},
		Version: "1",
	})
	env.Hostname = "host2"
	dis2 := discovery.New(&discovery.Config{
		Env:    "uat",
		Host:   "host2",
		Zone:   "sh001",
		Key:    "1",
		Secret: "1",
	})
	dis2.Register(ctx, &naming.Instance{
		AppID: appID,
		Addrs: []string{
			"grpc://127.0.0.2:9000",
		},
		Version: "1",
	})

}
func Test_Disocvery(t *testing.T) {
	dis := discovery.New(&discovery.Config{
		Env:    "uat",
		Host:   "host1",
		Zone:   "sh001",
		Key:    "1",
		Secret: "1",
	})
	b := resolver.Builder{dis}
	cc := &ClientConn{}
	r, err := b.Build(gr.Target{"discovery", "d", appID}, cc, gr.BuildOption{})
	if err != nil {
		t.Fatalf("b.Build error(%v)", err)
	}
	r.ResolveNow(gr.ResolveNowOption{})
	time.Sleep(time.Second)
	if len(cc.addrs) != 2 {
		t.Fatalf("get addrs expected 2,but only get %d", len(cc.addrs))
	}
}
