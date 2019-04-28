package resource

import (
	"context"

	"go-common/app/interface/main/creative/conf"
	resmdl "go-common/app/service/main/resource/model"
	resrpc "go-common/app/service/main/resource/rpc/client"
	"go-common/library/log"
	"strconv"
)

// Dao str
type Dao struct {
	c *conf.Config
	// rpc
	resRPC *resrpc.Service
}

// New str
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:      c,
		resRPC: resrpc.New(c.ResourceRPC),
	}
	return
}

// Banner get search banner
func (d *Dao) Banner(c context.Context, mobiApp, device, network, channel, ip, buvid, adExtra, resIDStr string, build int, plat int8, mid int64, isAd bool) (res map[int][]*resmdl.Banner, err error) {
	var bs *resmdl.Banners
	arg := &resmdl.ArgBanner{
		MobiApp: mobiApp,
		Device:  device,
		Network: network,
		Channel: channel,
		IP:      ip,
		Buvid:   buvid,
		AdExtra: adExtra,
		ResIDs:  resIDStr,
		Build:   build,
		Plat:    plat,
		MID:     mid,
		IsAd:    isAd,
	}
	if bs, err = d.resRPC.Banners(c, arg); err != nil || bs == nil {
		log.Error("d.resRPC.Banners(%v) error(%v) or bs is nil", arg, err)
		return
	}
	if bs == nil {
		return
	}
	if len(bs.Banner) > 0 {
		res = bs.Banner
	}
	return
}

// SimpleResource simple resource
func (d *Dao) SimpleResource(c context.Context, resID int) (res *resmdl.Resource, err error) {
	arg := &resmdl.ArgRes{ResID: resID}
	if res, err = d.resRPC.Resource(c, arg); err != nil || res == nil {
		log.Error("d.resRPC.Resource(%v) error(%v) or bs is nil", arg, err)
		return
	}
	return
}

// Resource get resource
func (d *Dao) Resource(c context.Context, resID int) (aidMap map[int64]struct{}, err error) {
	var rs *resmdl.Resource
	arg := &resmdl.ArgRes{ResID: resID}
	if rs, err = d.resRPC.Resource(c, arg); err != nil || rs == nil {
		log.Error("d.resRPC.Resource(%v) error(%v) or bs is nil", arg, err)
		return
	}
	aidMap = make(map[int64]struct{})
	for _, ass := range rs.Assignments {
		aid, _ := strconv.ParseInt(ass.URL, 10, 64)
		if aid > 0 {
			aidMap[aid] = struct{}{}
		}
	}
	return
}
