package vip

import (
	"context"
	"strconv"

	"go-common/app/interface/main/account/conf"
	"go-common/app/interface/main/account/dao/vip"
	"go-common/app/interface/main/account/model"
	mrl "go-common/app/service/main/relation/model"
	rlrpc "go-common/app/service/main/relation/rpc/client"
	resmdl "go-common/app/service/main/resource/model"
	rscrpc "go-common/app/service/main/resource/rpc/client"
	v1 "go-common/app/service/main/vip/api"
	vipmod "go-common/app/service/main/vip/model"
	viprpc "go-common/app/service/main/vip/rpc/client"
	"go-common/library/log"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// Service .
type Service struct {
	// conf
	c *conf.Config
	// http
	vipDao *vip.Dao
	// vip rpc
	vipRPC      *viprpc.Service
	relationRPC *rlrpc.Service
	resourceRPC *rscrpc.Service
	// vip service
	vipgRPC v1.VipClient
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:           c,
		vipDao:      vip.New(c),
		vipRPC:      viprpc.New(c.RPCClient2.Vip),
		relationRPC: rlrpc.New(c.RPCClient2.Relation),
		resourceRPC: rscrpc.New(c.RPCClient2.Resource),
	}
	vipgRPC, err := v1.NewClient(c.VipClient)
	if err != nil {
		panic(err)
	}
	s.vipgRPC = vipgRPC
	return
}

// OrderStatus .
func (s *Service) OrderStatus(c context.Context, arg *vipmod.ArgDialog) (res *vipmod.OrderResult, err error) {
	if res, err = s.vipRPC.OrderPayResult(c, arg); err != nil {
		return
	}
	if res == nil || res.Dialog == nil {
		log.Warn("s.vipRPC.OrderPayResult(%+v) get nil", arg)
		return
	}
	if res.Dialog.Follow {
		var (
			f    *mrl.Following
			ferr error
			ip   = metadata.String(c, metadata.RemoteIP)
			ar   = &mrl.ArgRelation{Mid: arg.Mid, Fid: s.c.Vipproperty.OfficialMid, RealIP: ip}
		)
		if f, ferr = s.relationRPC.Relation(c, ar); ferr != nil {
			log.Error("s.Relation(%+v) err(%v)", ar, ferr)
			return
		}
		if f == nil {
			return
		}
		// 如果已经关注就隐藏关注模块
		res.Dialog.Follow = !f.Following()
	}
	return
}

//Unfrozen user unfrozen vip
func (s *Service) Unfrozen(c context.Context, mid int64) (err error) {
	if err = s.vipRPC.Unfrozen(c, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//FrozenTime get frozen time
func (s *Service) FrozenTime(c context.Context, mid int64) (stime int64, err error) {
	if stime, err = s.vipRPC.SurplusFrozenTime(c, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// ResourceBanner .
func (s *Service) ResourceBanner(c context.Context, arg *model.ArgResource) (res []*resmdl.Banner, err error) {
	resID := ""
	if model.IsIPhone(arg.Plat) {
		resID = model.ResourceBannerIPhone
	}
	if model.IsAndroid(arg.Plat) {
		resID = model.ResourceBannerAndroid
	}
	if model.IsIPad(arg.Plat) {
		resID = model.ResourceBannerIPad
	}
	var argb = &resmdl.ArgBanner{
		Plat:    arg.Plat,
		Build:   arg.Build,
		MID:     arg.MID,
		ResIDs:  resID,
		Channel: "master",
		IP:      metadata.String(c, metadata.RemoteIP),
		Buvid:   arg.Buvid,
		Network: arg.Network,
		MobiApp: arg.MobiApp,
		Device:  arg.Device,
	}
	bs, err := s.resourceRPC.Banners(c, argb)
	if err != nil || bs == nil {
		log.Error("s.resourceRPC.Banners(%v) error(%+v) or bs is nil", argb, err)
		return
	}
	if len(bs.Banner) > 0 {
		var rid int
		rid, err = strconv.Atoi(resID)
		if err != nil {
			err = errors.WithStack(err)
			return
		}
		res = bs.Banner[rid]
	}
	return
}

// ResourceBuy .
func (s *Service) ResourceBuy(c context.Context, arg *model.ArgResource) (res []*resmdl.Banner, err error) {
	resID := ""
	if model.IsIPhone(arg.Plat) {
		resID = model.ResourceBuyIPhone
	}
	if model.IsAndroid(arg.Plat) {
		resID = model.ResourceBuyAndroid
	}
	if model.IsIPad(arg.Plat) {
		resID = model.ResourceBuyIPad
	}
	log.Info("ResourceBuy resID(%s) arg(%+v)", resID, arg)
	var argb = &resmdl.ArgBanner{
		Plat:    arg.Plat,
		Build:   arg.Build,
		MID:     arg.MID,
		ResIDs:  resID,
		Channel: "master",
		IP:      metadata.String(c, metadata.RemoteIP),
		Buvid:   arg.Buvid,
		Network: arg.Network,
		MobiApp: arg.MobiApp,
		Device:  arg.Device,
	}
	bs, err := s.resourceRPC.Banners(c, argb)
	if err != nil || bs == nil {
		log.Error("s.resourceRPC.Banners(%v) error(%+v) or bs is nil", argb, err)
		return
	}
	log.Info("s.resourceRPC.Banners(%+v)", bs.Banner)
	if len(bs.Banner) > 0 {
		var rid int
		rid, err = strconv.Atoi(resID)
		if err != nil {
			err = errors.WithStack(err)
			return
		}
		res = bs.Banner[rid]
	}
	return
}
