package vip

import (
	"context"
	"encoding/json"

	"go-common/app/interface/main/account/model"
	vol "go-common/app/service/main/vip/model"

	v1 "go-common/app/service/main/vip/api"
)

// BindInfoByMid bind info by mid[bilibili->third].
func (s *Service) BindInfoByMid(c context.Context, a *model.ArgBindInfo) (res *v1.BindInfoByMidReply, err error) {
	return s.vipgRPC.BindInfoByMid(c, &v1.BindInfoByMidReq{
		Mid:   a.Mid,
		AppId: a.AppID,
	})
}

// CreateAssociateOrder create associate order.
func (s *Service) CreateAssociateOrder(c context.Context, req *model.ArgCreateAssociateOrder) (res map[string]interface{}, err error) {
	var p *v1.CreateAssociateOrderReply
	if p, err = s.vipgRPC.CreateAssociateOrder(c, &v1.CreateAssociateOrderReq{
		Mid:         req.Mid,
		Month:       req.Month,
		Platform:    req.Platform,
		MobiApp:     req.MobiApp,
		Device:      req.Device,
		AppId:       req.AppID,
		AppSubId:    req.AppSubID,
		OrderType:   int32(req.OrderType),
		Dtype:       int32(req.Dtype),
		ReturnUrl:   req.ReturnURL,
		CouponToken: req.CouponToken,
		Bmid:        req.Bmid,
		PanelType:   req.PanelType,
		Build:       req.Build,
		IP:          req.IP,
	}); err != nil {
		return
	}
	json.Unmarshal([]byte(p.PayParam), &res)
	return
}

// AssociatePanel associate panel.
func (s *Service) AssociatePanel(c context.Context, req *vol.ArgAssociatePanel) (res []*v1.AssociatePanelInfo, err error) {
	var reply *v1.AssociatePanelReply
	if reply, err = s.vipgRPC.AssociatePanel(c, &v1.AssociatePanelReq{
		Mid:       req.Mid,
		SortTp:    int32(req.SortTP),
		IP:        req.IP,
		MobiApp:   req.MobiApp,
		Device:    req.Device,
		Platform:  req.Platform,
		PanelType: req.PanelType,
		Build:     req.Build,
	}); err != nil {
		return
	}
	res = reply.List
	return
}

// EleRedPackages ele red packages.
func (s *Service) EleRedPackages(c context.Context) (res []*v1.ModelEleRedPackage, err error) {
	var reply *v1.EleRedPackagesReply
	if reply, err = s.vipgRPC.EleRedPackages(c, &v1.EleRedPackagesReq{}); err != nil {
		return
	}
	res = reply.List
	return
}

// EleSpecailFoods ele speacail foods.
func (s *Service) EleSpecailFoods(c context.Context) (res []*v1.ModelEleSpecailFoods, err error) {
	var reply *v1.EleSpecailFoodsReply
	if reply, err = s.vipgRPC.EleSpecailFoods(c, &v1.EleSpecailFoodsReq{}); err != nil {
		return
	}
	res = reply.List
	return
}
