package resource

import (
	"context"

	"go-common/app/interface/main/app-resource/conf"
	resource "go-common/app/service/main/resource/model"
	resrpc "go-common/app/service/main/resource/rpc/client"
	"go-common/library/log"
)

type Dao struct {
	c *conf.Config
	// rpc
	resRpc *resrpc.Service
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c: c,
		// rpc
		resRpc: resrpc.New(c.ResourceRPC),
	}
	return
}

// ResSideBar resource ressidebar
func (d *Dao) ResSideBar(ctx context.Context) (res *resource.SideBars, err error) {
	if res, err = d.resRpc.SideBars(ctx); err != nil {
		log.Error("resource d.resRpc.SideBars error(%v)", err)
		return
	}
	return
}

// AbTest resource abtest
func (d *Dao) AbTest(ctx context.Context, groups string) (res map[string]*resource.AbTest, err error) {
	arg := &resource.ArgAbTest{
		Groups: groups,
	}
	if res, err = d.resRpc.AbTest(ctx, arg); err != nil {
		log.Error("resource d.resRpc.AbTest error(%v)", err)
		return
	}
	return
}
