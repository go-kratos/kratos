package resource

import (
	"context"

	"go-common/app/interface/main/app-feed/conf"
	"go-common/app/service/main/resource/model"
	rscrpc "go-common/app/service/main/resource/rpc/client"
	"go-common/library/log"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

type Dao struct {
	c *conf.Config
	// rpc
	rscRPC *rscrpc.Service
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c: c,
		// rpc
		rscRPC: rscrpc.New(c.ResourceRPC),
	}
	return
}

func (d *Dao) Banner(c context.Context, plat int8, build int, mid int64, resIDs, channel, buvid, network, mobiApp, device string, isAd bool, openEvent, adExtra, hash string) (res map[int][]*model.Banner, version string, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &model.ArgBanner{Plat: plat, ResIDs: resIDs, Build: build, MID: mid, Channel: channel, IP: ip, Buvid: buvid, Network: network, MobiApp: mobiApp, Device: device, IsAd: isAd, OpenEvent: openEvent, AdExtra: adExtra, Version: hash}
	bs, err := d.rscRPC.Banners(c, arg)
	if err != nil {
		err = errors.Wrapf(err, "%v", arg)
		return
	}
	if bs != nil {
		res = bs.Banner
		version = bs.Version
	}
	return
}

// AbTest resource abtest
func (d *Dao) AbTest(ctx context.Context, groups string) (res map[string]*model.AbTest, err error) {
	arg := &model.ArgAbTest{
		Groups: groups,
	}
	if res, err = d.rscRPC.AbTest(ctx, arg); err != nil {
		log.Error("resource d.resRpc.AbTest error(%v)", err)
		return
	}
	return
}
