package push

import (
	"context"
	"fmt"
	"strings"
	"time"

	appres "go-common/app/interface/main/app-resource/api/v1"
	"go-common/library/log"
	"go-common/library/naming/discovery"
	"go-common/library/net/rpc/warden"
)

func errlog(step string, err error) {
	log.Error("CallPush Step [%s] Err [%v]", step, err)
}

// build grpc client for app-resource
func (d *Dao) grpcClient(addr string) (appres.AppResourceClient, error) {
	client := warden.NewClient(d.c.AppresClient)
	cc, err := client.Dial(context.Background(), addr)
	if err != nil {
		return nil, err
	}
	return appres.NewAppResourceClient(cc), nil
}

// CallRefresh picks ip addrs from discovery, and call grpc method
func (d *Dao) CallRefresh(ctx context.Context) (err error) {
	var (
		addrs    []string
		client   appres.AppResourceClient
		arg      = &appres.NoArgRequest{}
		succCall int
	)
	if addrs, err = d.pickAddrs(); err != nil {
		errlog("pickAddrs", err)
		return
	}
	// loop grpc call, ignore error
	for _, addr := range addrs {
		if client, err = d.grpcClient(addr); err != nil {
			errlog(fmt.Sprintf("grpcDial Addr [%s]", addr), err)
			continue
		}
		if _, err = client.ModuleUpdateCache(ctx, arg); err != nil {
			errlog(fmt.Sprintf("grpcCall Addr [%s] ", addr), err)
			continue
		}
		succCall++
	}
	if succCall == 0 { // addrs must be greater than 0, if succ call is zero, return error
		return fmt.Errorf("CallRefresh Addrs [%d], Zero Succ!", len(addrs))
	}
	log.Info("CallPush Refresh Succ [%d], Total [%d]", succCall, len(addrs))
	err = nil
	return
}

// pick all the app-resource grpc instances from discovery
func (d *Dao) pickAddrs() (grpcAddrs []string, err error) {
	dis := discovery.New(nil)
	defer dis.Close()
	b := dis.Build(d.c.Cfg.Grpc.ApiAppID)
	defer b.Close()
	ch := b.Watch()
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	select {
	case <-ch:
	case <-ctx.Done():
		err = fmt.Errorf("查找节点超时 请检查appid是否填写正确")
		return
	}
	ins, ok := b.Fetch(context.Background())
	if !ok {
		err = fmt.Errorf("discovery 拉取失败")
		return
	}
	for _, vs := range ins {
		for _, v := range vs {
			for _, addr := range v.Addrs {
				if strings.Contains(addr, "grpc://") {
					ip := strings.Replace(addr, "grpc://", "", -1)
					grpcAddrs = append(grpcAddrs, ip)
				}
			}
		}
	}
	if len(grpcAddrs) == 0 {
		err = fmt.Errorf("discovery 找不到服务节点")
	}
	return
}
