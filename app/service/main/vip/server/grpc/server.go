// Package server generate by warden_gen
package server

import (
	"context"
	"encoding/json"
	"net"

	colv1 "go-common/app/service/main/coupon/api"
	col "go-common/app/service/main/coupon/model"
	v1 "go-common/app/service/main/vip/api"
	"go-common/app/service/main/vip/model"
	service "go-common/app/service/main/vip/service"
	"go-common/library/log"
	"go-common/library/net/rpc/warden"
)

// New VipInfo warden rpc server
func New(c *warden.ServerConfig, svr *service.Service) *warden.Server {
	ws := warden.NewServer(c)
	v1.RegisterVipServer(ws.Server(), &server{svr})
	ws, err := ws.Start()
	if err != nil {
		panic(err)
	}
	return ws
}

type server struct {
	svr *service.Service
}

var _ v1.VipServer = &server{}

// RegisterOpenID register open id.
func (s *server) RegisterOpenID(c context.Context, req *v1.RegisterOpenIDReq) (res *v1.RegisterOpenIDReply, err error) {
	var ro *model.RegisterOpenIDResp
	if ro, err = s.svr.RegisterOpenID(c, &model.ArgRegisterOpenID{
		AppID: req.AppId,
		Mid:   req.Mid,
	}); err != nil {
		return
	}
	return &v1.RegisterOpenIDReply{OpenId: ro.OpenID}, nil
}

// OpenBindByOutOpenID associate user bind by out_open_id [third -> bilibili].
func (s *server) OpenBindByOutOpenID(c context.Context, req *v1.OpenBindByOutOpenIDReq) (res *v1.OpenBindByOutOpenIDReply, err error) {
	if err = s.svr.OpenBindByOutOpenID(c, &model.ArgBind{
		AppID:     req.AppId,
		OpenID:    req.OpenId,
		OutOpenID: req.OutOpenId,
	}); err != nil {
		return
	}
	return &v1.OpenBindByOutOpenIDReply{}, nil
}

// UserInfoByOpenID get userinfo by open_id.
func (s *server) UserInfoByOpenID(c context.Context, req *v1.UserInfoByOpenIDReq) (res *v1.UserInfoByOpenIDReply, err error) {
	var u *model.UserInfoByOpenIDResp
	if u, err = s.svr.UserInfoByOpenID(c, &model.ArgUserInfoByOpenID{
		AppID:  req.AppId,
		OpenID: req.OpenId,
		IP:     req.Ip,
	}); err != nil {
		return
	}
	return &v1.UserInfoByOpenIDReply{
		Name:      u.Name,
		BindState: u.BindState,
		OutOpenId: u.OutOpenID,
	}, nil
}

// BindInfoByMid bind info by mid[bilibili->third].
func (s *server) BindInfoByMid(c context.Context, req *v1.BindInfoByMidReq) (res *v1.BindInfoByMidReply, err error) {
	var (
		b *model.BindInfo
		o *v1.BindOuter
		a *v1.Account
	)
	if b, err = s.svr.BindInfoByMid(c, &model.ArgBindInfo{
		AppID: req.AppId,
		Mid:   req.Mid,
	}); err != nil {
		return
	}
	if b.Account != nil {
		a = &v1.Account{
			Mid:  b.Account.Mid,
			Face: b.Account.Face,
			Name: b.Account.Name,
		}
	}
	if b.Outer != nil {
		o = &v1.BindOuter{
			Tel:       b.Outer.Tel,
			BindState: b.Outer.BindState,
		}
	}
	return &v1.BindInfoByMidReply{
		Account: a,
		Outer:   o,
	}, nil
}

// BilibiliPrizeGrant vip prize grant for third [third->bilibili].
func (s *server) BilibiliPrizeGrant(c context.Context, req *v1.BilibiliPrizeGrantReq) (res *v1.BilibiliPrizeGrantReply, err error) {
	var sr *col.SalaryCouponForThirdResp
	if sr, err = s.svr.BilibiliPrizeGrant(c, &model.ArgBilibiliPrizeGrant{
		AppID:    req.AppId,
		OpenID:   req.OpenId,
		PrizeKey: req.PrizeKey,
		UniqueNo: req.UniqueNo,
	}); err != nil {
		return
	}
	return &v1.BilibiliPrizeGrantReply{
		Amount:      sr.Amount,
		FullAmount:  sr.FullAmount,
		Description: sr.Description,
	}, nil
}

// BilibiliVipGrant bilibili associate vip grant [third -> bilibili]
func (s *server) BilibiliVipGrant(c context.Context, req *v1.BilibiliVipGrantReq) (res *v1.BilibiliVipGrantReply, err error) {
	if err = s.svr.BilibiliVipGrant(c, &model.ArgBilibiliVipGrant{
		AppID:      req.AppId,
		OpenID:     req.OpenId,
		OutOpenID:  req.OutOpenId,
		OutOrderNO: req.OutOrderNo,
		Duration:   req.Duration,
	}); err != nil {
		return
	}
	return &v1.BilibiliVipGrantReply{}, nil
}

// CreateAssociateOrder create associate order.
func (s *server) CreateAssociateOrder(c context.Context, req *v1.CreateAssociateOrderReq) (res *v1.CreateAssociateOrderReply, err error) {
	var cr *model.CreateOrderRet
	if cr, err = s.svr.CreateAssociateOrder(c, &model.ArgCreateOrder2{
		Mid:         req.Mid,
		Month:       req.Month,
		Platform:    req.Platform,
		MobiApp:     req.MobiApp,
		Device:      req.Device,
		AppID:       req.AppId,
		AppSubID:    req.AppSubId,
		OrderType:   int8(req.OrderType),
		Dtype:       int8(req.Dtype),
		ReturnURL:   req.ReturnUrl,
		CouponToken: req.CouponToken,
		Bmid:        req.Bmid,
		PanelType:   req.PanelType,
		Build:       req.Build,
		IP:          net.ParseIP(req.IP),
	}); err != nil {
		return
	}
	marshal, err := json.Marshal(cr.PayParam)
	if err != nil {
		log.Error("json.Marshal(%+v) err(%+v)", cr.PayParam, err)
		return &v1.CreateAssociateOrderReply{}, err
	}
	return &v1.CreateAssociateOrderReply{PayParam: string(marshal)}, nil
}

// AssociatePanel associate panel.
func (s *server) AssociatePanel(c context.Context, req *v1.AssociatePanelReq) (res *v1.AssociatePanelReply, err error) {
	var (
		pl   []*model.AssociatePanelInfo
		list = []*v1.AssociatePanelInfo{}
	)
	if pl, err = s.svr.AssociatePanel(c, &model.ArgAssociatePanel{
		Mid:       req.Mid,
		SortTP:    int8(req.SortTp),
		IP:        req.IP,
		MobiApp:   req.MobiApp,
		Device:    req.Device,
		Platform:  req.Platform,
		PanelType: req.PanelType,
		Build:     req.Build,
	}); err != nil {
		return
	}
	for _, v := range pl {
		list = append(list, &v1.AssociatePanelInfo{
			Id:            v.ID,
			Month:         v.Month,
			ProductName:   v.PdName,
			ProductId:     v.PdID,
			SubType:       v.SubType,
			SuitType:      v.SuitType,
			OriginalPrice: v.OPrice,
			DiscountPrice: v.DPrice,
			DiscountRate:  v.DRate,
			Remark:        v.Remark,
			Selected:      v.Selected,
			PayState:      int32(v.PayState),
			PayMessage:    v.PayMessage,
		})
	}
	return &v1.AssociatePanelReply{List: list}, nil
}

// OpenAuthCallBack third open call back.
func (s *server) OpenAuthCallBack(c context.Context, req *v1.OpenAuthCallBackReq) (res *v1.OpenAuthCallBackReply, err error) {
	if err = s.svr.OpenAuthCallBack(c, &model.ArgOpenAuthCallBack{
		Mid:       req.Mid,
		AppID:     req.AppId,
		ThirdCode: req.ThirdCode,
	}); err != nil {
		return
	}
	return &v1.OpenAuthCallBackReply{}, nil
}

// EleRedPackages red packages.
func (s *server) EleRedPackages(c context.Context, req *v1.EleRedPackagesReq) (res *v1.EleRedPackagesReply, err error) {
	var data []*model.EleRedPackagesResp
	if data, err = s.svr.EleRedPackages(c); err != nil {
		return
	}
	list := []*v1.ModelEleRedPackage{}
	for _, v := range data {
		list = append(list, &v1.ModelEleRedPackage{
			Name:         v.Name,
			Amount:       v.Amount,
			SumCondition: v.SumCondition,
		})
	}
	return &v1.EleRedPackagesReply{List: list}, nil
}

// EleSpecailFoods specail foods.
func (s *server) EleSpecailFoods(c context.Context, req *v1.EleSpecailFoodsReq) (res *v1.EleSpecailFoodsReply, err error) {
	var data []*model.EleSpecailFoodsResp
	if data, err = s.svr.EleSpecailFoods(c); err != nil {
		return
	}
	list := []*v1.ModelEleSpecailFoods{}
	for _, v := range data {
		list = append(list, &v1.ModelEleSpecailFoods{
			RestaurantName: v.RestaurantName,
			FoodName:       v.FoodName,
			FoodUrl:        v.FoodURL,
			Discount:       v.Discount,
			Amount:         v.Amount,
			OriginalAmount: v.OriginalAmount,
			RatingPoint:    v.RatingPoint,
		})
	}
	return &v1.EleSpecailFoodsReply{List: list}, nil
}

// EleVipGrant vip grant.
func (s *server) EleVipGrant(c context.Context, req *v1.EleVipGrantReq) (res *v1.EleVipGrantReply, err error) {
	if err = s.svr.EleVipGrant(c, &model.ArgEleVipGrant{OrderNO: req.OrderNo}); err != nil {
		return
	}
	return &v1.EleVipGrantReply{}, nil
}

// CouponBySuitID get coupon by mid and suit info.
func (s *server) CouponBySuitID(c context.Context, req *v1.CouponBySuitIDReq) (res *v1.CouponBySuitIDReply, err error) {
	var data *colv1.UsableAllowanceCouponV2Reply
	res = new(v1.CouponBySuitIDReply)
	if data, err = s.svr.CouponBySuitIDV2(c, req); err != nil {
		return
	}
	if data == nil {
		return
	}
	res.CouponTip = data.CouponTip
	if data.CouponInfo == nil {
		return
	}
	res.CouponInfo = &v1.ModelCouponAllowancePanelInfo{
		CouponToken:         data.CouponInfo.CouponToken,
		CouponAmount:        data.CouponInfo.CouponAmount,
		State:               data.CouponInfo.State,
		FullAmount:          data.CouponInfo.FullAmount,
		FullLimitExplain:    data.CouponInfo.FullLimitExplain,
		ScopeExplain:        data.CouponInfo.ScopeExplain,
		CouponDiscountPrice: data.CouponInfo.CouponDiscountPrice,
		StartTime:           data.CouponInfo.StartTime,
		ExpireTime:          data.CouponInfo.ExpireTime,
		Selected:            data.CouponInfo.Selected,
		DisablesExplains:    data.CouponInfo.DisablesExplains,
		OrderNo:             data.CouponInfo.OrderNo,
		Name:                data.CouponInfo.Name,
		Usable:              data.CouponInfo.Usable,
	}
	return
}

// VipUserPanel vip user panel
func (s *server) VipUserPanel(c context.Context, req *v1.VipUserPanelReq) (res *v1.VipUserPanelReply, err error) {
	var data *model.VipPirceRespV9
	res = new(v1.VipUserPanelReply)
	if data, err = s.svr.VipUserPanelV9(c, req); err != nil {
		return
	}
	if data == nil {
		return
	}
	res.CouponSwitch = int32(data.CodeSwitch)
	res.CodeSwitch = int32(data.CodeSwitch)
	res.GiveSwitch = int32(data.GiveSwitch)
	priceList := []*v1.ModelVipPanelInfo{}
	for _, v := range data.Vps {
		priceList = append(priceList, &v1.ModelVipPanelInfo{
			Month:         v.Month,
			ProductName:   v.PdName,
			ProductId:     v.PdID,
			SubType:       v.SubType,
			SuitType:      v.SuitType,
			OriginalPrice: v.OPrice,
			DiscountPrice: v.DPrice,
			DiscountRate:  v.DRate,
			Remark:        v.Remark,
			Selected:      v.Selected,
			Id:            v.Id,
			Type:          v.Type,
		})
	}
	res.PriceList = priceList
	if data.Coupon != nil {
		res.Coupon = &v1.CouponBySuitIDReply{
			CouponTip: data.Coupon.CouponTip,
		}
		if data.Coupon.CouponInfo != nil {
			res.Coupon.CouponInfo = &v1.ModelCouponAllowancePanelInfo{
				CouponToken:         data.Coupon.CouponInfo.CouponToken,
				CouponAmount:        data.Coupon.CouponInfo.CouponAmount,
				State:               data.Coupon.CouponInfo.State,
				FullAmount:          data.Coupon.CouponInfo.FullAmount,
				FullLimitExplain:    data.Coupon.CouponInfo.FullLimitExplain,
				ScopeExplain:        data.Coupon.CouponInfo.ScopeExplain,
				CouponDiscountPrice: data.Coupon.CouponInfo.CouponDiscountPrice,
				StartTime:           data.Coupon.CouponInfo.StartTime,
				ExpireTime:          data.Coupon.CouponInfo.ExpireTime,
				Selected:            data.Coupon.CouponInfo.Selected,
				DisablesExplains:    data.Coupon.CouponInfo.DisablesExplains,
				OrderNo:             data.Coupon.CouponInfo.OrderNo,
				Name:                data.Coupon.CouponInfo.Name,
				Usable:              data.Coupon.CouponInfo.Usable,
			}
		}
	}
	privileges := map[int32]*v1.ModelPrivilegeResp{}
	for k, v := range data.Privileges {
		list := []*v1.ModelPrivilege{}
		for _, p := range v.List {
			list = append(list, &v1.ModelPrivilege{
				Name:    p.Name,
				IconUrl: p.IconURL,
				Type:    int32(p.Type),
			})
		}
		privileges[int32(k)] = &v1.ModelPrivilegeResp{
			Title: v.Title,
			List:  list,
		}
	}
	res.Privileges = privileges
	return
}

// EleVipGrant get welfare list.
func (s *server) WelfareList(c context.Context, req *v1.WelfareReq) (res *v1.WelfareReply, err error) {
	var (
		data  []*model.WelfareListResp
		count int64
		list  = []*v1.WelfareListDetail{}
	)
	if data, count, err = s.svr.WelfareList(c, &model.ArgWelfareList{
		Tid:       req.Tid,
		Recommend: req.Recommend,
		Ps:        req.Ps,
		Pn:        req.Pn,
	}); err != nil {
		return
	}
	for _, v := range data {
		list = append(list, &v1.WelfareListDetail{
			Id:          v.ID,
			Name:        v.Name,
			HomepageUri: v.HomepageUri,
			BackdropUri: v.BackdropUri,
			Tid:         v.Tid,
			Rank:        v.Rank,
		})
	}
	return &v1.WelfareReply{Count: count, List: list}, nil
}

// WelfareTypeList get welfare type list.
func (s *server) WelfareTypeList(c context.Context, req *v1.WelfareTypeReq) (res *v1.WelfareTypeReply, err error) {
	var (
		data []*model.WelfareTypeListResp
		list = []*v1.WelfareTypeListDetail{}
	)
	if data, err = s.svr.WelfareTypeList(c); err != nil {
		return
	}
	for _, v := range data {
		list = append(list, &v1.WelfareTypeListDetail{
			Id:   v.ID,
			Name: v.Name,
		})
	}
	return &v1.WelfareTypeReply{List: list}, nil
}

// WelfareInfo get welfare info.
func (s *server) WelfareInfo(c context.Context, req *v1.WelfareInfoReq) (res *v1.WelfareInfoReply, err error) {
	var data *model.WelfareInfoResp
	if data, err = s.svr.WelfareInfo(c, &model.ArgWelfareInfo{ID: req.Id, MID: req.Mid}); err != nil {
		return
	}

	return &v1.WelfareInfoReply{
		Id:          data.ID,
		Name:        data.Name,
		Desc:        data.Desc,
		HomepageUri: data.HomepageUri,
		BackdropUri: data.BackdropUri,
		Finished:    data.Finished,
		Received:    data.Received,
		VipType:     data.VipType,
		Stime:       int64(data.Stime),
		Etime:       int64(data.Etime),
	}, nil
}

// WelfareReceive receive welfare.
func (s *server) WelfareReceive(c context.Context, req *v1.WelfareReceiveReq) (res *v1.WelfareReceiveReply, err error) {
	if err = s.svr.WelfareReceive(c, &model.ArgWelfareReceive{Wid: req.Wid, Mid: req.Mid}); err != nil {
		return
	}
	return &v1.WelfareReceiveReply{}, nil
}

// MyWelfare get my welfare
func (s *server) MyWelfare(c context.Context, req *v1.MyWelfareReq) (res *v1.MyWelfareReply, err error) {
	var (
		data []*model.MyWelfareResp
		list = []*v1.MyWelfareDetail{}
	)
	if data, err = s.svr.MyWelfare(c, req.Mid); err != nil {
		return
	}
	for _, v := range data {
		list = append(list, &v1.MyWelfareDetail{
			Wid:        v.Wid,
			Name:       v.Name,
			Desc:       v.Desc,
			UsageForm:  v.UsageForm,
			ReceiveUri: v.ReceiveUri,
			Code:       v.Code,
			Expired:    v.Expired,
			Stime:      int64(v.Stime),
			Etime:      int64(v.Etime),
		})
	}
	return &v1.MyWelfareReply{List: list}, nil
}
