package vip

import (
	"context"

	"go-common/app/interface/main/account/model"
	vipv1 "go-common/app/service/main/vip/api"
	vipml "go-common/app/service/main/vip/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"

	"github.com/pkg/errors"
)

// VipPanel .
func (s *Service) VipPanel(c context.Context, mid int64, a *model.VipPanelRes) (res *vipml.VipPirceResp, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	res, err = s.vipRPC.VipPanelInfo2(c, &vipml.ArgPanel{Mid: mid, SortTp: a.SortTP, IP: ip, Device: a.Device, MobiApp: a.MobiApp, Platform: a.Platform, PanelType: a.PanelType, SubType: a.SubType, Month: a.Month, Build: a.Build})
	return
}

// VipPanel5 .
func (s *Service) VipPanel5(c context.Context, mid int64, a *model.VipPanelRes) (res *vipml.VipPirceResp5, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	res, err = s.vipRPC.VipPanelInfo5(c, &vipml.ArgPanel{Mid: mid, SortTp: a.SortTP, IP: ip, Device: a.Device, MobiApp: a.MobiApp, Platform: a.Platform, PanelType: a.PanelType, SubType: a.SubType, Month: a.Month, Build: a.Build})
	return
}

// VipPanelV2 vip panel v2.
func (s *Service) VipPanelV2(c context.Context, a *model.ArgVipPanel) (res *model.VipPanelResp, err error) {
	var (
		g  errgroup.Group
		p  *vipml.VipPirceResp5
		ts []*vipml.TipsResp
	)
	res = new(model.VipPanelResp)
	g.Go(func() (err error) {
		if p, err = s.vipRPC.VipPanelInfo5(c, &vipml.ArgPanel{
			Mid:       a.Mid,
			SortTp:    a.SortTP,
			IP:        a.IP,
			Device:    a.Device,
			MobiApp:   a.MobiApp,
			Platform:  a.Platform,
			PanelType: a.PanelType,
			Build:     a.Build,
		}); err != nil || p == nil {
			log.Error("s.vipRPC.VipPanelInfo2(%+v) error(%v)", a, err)
			return
		}
		res.Vps = p.Vps
		res.CodeSwitch = p.CodeSwitch
		res.GiveSwitch = p.GiveSwitch
		res.Privileges = p.Privileges
		return
	})
	g.Go(func() (err error) {
		if ts, err = s.vipRPC.Tips(c, &vipml.ArgTips{
			Version:  a.Build,
			Platform: a.Platform,
			Position: vipml.PanelPosition,
		}); err != nil {
			log.Error("s.vipRPC.Tips(%+v) error(%v)", a, err)
		}
		if len(ts) == 0 {
			return
		}
		res.TipInfo = ts[0]
		return
	})
	g.Go(func() (err error) {
		if res.UserInfo, err = s.vipRPC.PanelExplain(c, &vipml.ArgPanelExplain{
			Mid: a.Mid,
		}); err != nil {
			log.Error("s.vipRPC.PanelExplain(%+v) error(%v)", a, err)
		}
		return
	})
	if err = g.Wait(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// VipPanelV8 vip panel v8
func (s *Service) VipPanelV8(c context.Context, a *model.ArgVipPanel) (res *model.VipPanelV8Resp, err error) {
	var (
		g  errgroup.Group
		p  *vipml.VipPirceResp5
		ts []*vipml.TipsResp
	)
	res = new(model.VipPanelV8Resp)
	g.Go(func() (err error) {
		if p, err = s.vipRPC.VipPanelInfo5(c, &vipml.ArgPanel{
			Mid:       a.Mid,
			SortTp:    a.SortTP,
			IP:        a.IP,
			Device:    a.Device,
			MobiApp:   a.MobiApp,
			Platform:  a.Platform,
			PanelType: a.PanelType,
			Build:     a.Build,
		}); err != nil || p == nil {
			log.Error("s.vipRPC.VipPanelInfo2(%+v) error(%v)", a, err)
			return
		}
		res.Vps = p.Vps
		res.CodeSwitch = p.CodeSwitch
		res.GiveSwitch = p.GiveSwitch
		res.Privileges = p.Privileges
		res.CouponInfo = p.CouponInfo
		res.CouponSwith = p.CouponSwith
		return
	})
	g.Go(func() (err error) {
		if ts, err = s.vipRPC.Tips(c, &vipml.ArgTips{
			Version:  a.Build,
			Platform: a.Platform,
			Position: vipml.PanelPosition,
		}); err != nil {
			log.Error("s.vipRPC.Tips(%+v) error(%v)", a, err)
		}
		if len(ts) == 0 {
			return
		}
		res.TipInfo = ts[0]
		return
	})
	g.Go(func() (err error) {
		if res.UserInfo, err = s.vipRPC.PanelExplain(c, &vipml.ArgPanelExplain{
			Mid: a.Mid,
		}); err != nil {
			log.Error("s.vipRPC.PanelExplain(%+v) error(%v)", a, err)
		}
		return
	})
	g.Go(func() (err error) {
		if res.AssociateVips, err = s.vipRPC.AssociateVips(c, &vipml.ArgAssociateVip{
			Platform: a.Platform,
			Device:   a.Device,
			MobiApp:  a.MobiApp,
		}); err != nil {
			log.Error("s.vipRPC.AssociateVips(%+v) error(%v)", a, err)
		}
		return
	})
	if err = g.Wait(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// VipPanelV9 vip panel v9
func (s *Service) VipPanelV9(c context.Context, a *model.ArgVipPanel) (res *model.VipPanelRespV9, err error) {
	var (
		p  *vipv1.VipUserPanelReply
		ts []*vipml.TipsResp
	)
	eg, ec := errgroup.WithContext(c)
	res = new(model.VipPanelRespV9)
	eg.Go(func() (err error) {
		if p, err = s.vipgRPC.VipUserPanel(ec, &vipv1.VipUserPanelReq{
			Mid:       a.Mid,
			SortTp:    int32(a.SortTP),
			Ip:        a.IP,
			Device:    a.Device,
			MobiApp:   a.MobiApp,
			Platform:  a.Platform,
			PanelType: a.PanelType,
			Build:     a.Build,
		}); err != nil || p == nil {
			log.Error("s.vipRPC.VipPanelInfo2(%+v) error(%v)", a, err)
			return
		}
		res.Vps = p.PriceList
		res.CodeSwitch = p.CodeSwitch
		res.GiveSwitch = p.GiveSwitch
		res.Privileges = p.Privileges
		res.Coupon = p.Coupon
		res.CouponSwith = p.CouponSwitch
		return
	})
	eg.Go(func() (err error) {
		if ts, err = s.vipRPC.Tips(ec, &vipml.ArgTips{
			Version:  a.Build,
			Platform: a.Platform,
			Position: vipml.PanelPosition,
		}); err != nil {
			log.Error("s.vipRPC.Tips(%+v) error(%v)", a, err)
		}
		if len(ts) == 0 {
			return
		}
		res.TipInfo = ts[0]
		return
	})
	eg.Go(func() (err error) {
		if res.UserInfo, err = s.vipRPC.PanelExplain(ec, &vipml.ArgPanelExplain{
			Mid: a.Mid,
		}); err != nil {
			log.Error("s.vipRPC.PanelExplain(%+v) error(%v)", a, err)
		}
		return
	})
	eg.Go(func() (err error) {
		if res.AssociateVips, err = s.vipRPC.AssociateVips(ec, &vipml.ArgAssociateVip{
			Platform: a.Platform,
			Device:   a.Device,
			MobiApp:  a.MobiApp,
		}); err != nil {
			log.Error("s.vipRPC.AssociateVips(%+v) error(%v)", a, err)
		}
		return
	})
	if err = eg.Wait(); err != nil {
		err = errors.WithStack(err)
	}
	return
}
