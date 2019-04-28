package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	colapi "go-common/app/service/main/coupon/api"
	coumol "go-common/app/service/main/coupon/model"
	v1 "go-common/app/service/main/vip/api"
	"go-common/app/service/main/vip/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

var (
	_emptyVipPanelInfo = make([]*model.VipPanelInfo, 0)
)

func (s *Service) loadVipPriceConfig() {
	var (
		err  error
		c    = context.TODO()
		vpcs []*model.VipPriceConfig
		mvp  map[int64]*model.VipDPriceConfig
	)
	if vpcs, err = s.dao.VipPriceConfigs(c); err != nil {
		log.Error("s.dao.VipPriceConfigs error(%+v)", err)
		return
	}
	if mvp, err = s.dao.VipPriceDiscountConfigs(c); err != nil {
		log.Error("s.dao.VipPriceDiscountConfigs error(%+v)", err)
		return
	}
	vipConfSuitMax := make(map[int64]int8)
	vipPriceConfMap := make(map[int64]map[int8][]*model.VipPriceConfig)
	vipPriceMap := make(map[int64]*model.VipPriceConfig, len(vpcs))
	for _, v := range vpcs {
		if _, ok := vipPriceConfMap[v.Plat]; !ok {
			vipPriceConfMap[v.Plat] = make(map[int8][]*model.VipPriceConfig)
		}
		v.DoCheckRealPrice(mvp)
		v.DoTopSuitType()
		vipPriceConfMap[v.Plat][v.SuitType] = append(vipPriceConfMap[v.Plat][v.SuitType], v)
		if max, ok := vipConfSuitMax[v.Plat]; !ok {
			vipConfSuitMax[v.Plat] = v.SuitType
		} else {
			if v.SuitType > max {
				vipConfSuitMax[v.Plat] = v.SuitType
			}
		}
		vipPriceMap[v.ID] = v
	}
	s.vipPriceConf.SetPriceConfig(c, vipPriceConfMap)
	s.vipConfSuitMax = vipConfSuitMax
	s.vipPriceMap = vipPriceMap
}

// VipUserPanel .
func (s *Service) VipUserPanel(c context.Context,
	mid, plat int64,
	sortTp int8,
	build int64) (vps []*model.VipPanelInfo, err error) {
	defer func() {
		if len(vps) == 0 {
			log.Warn("mid(%d) plat(%d)  vip panel  empty , must check!", mid, plat)
		}
	}()
	vpc, ok := s.vipPriceConf.GetPriceConfig(plat)
	if !ok {
		vps = _emptyVipPanelInfo
		log.Warn("not plat(%d) in vip_price_conf_map(%+v)", plat, s.vipPriceConf)
		return
	}
	if _, ok := s.vipConfSuitMax[plat]; !ok {
		vps = _emptyVipPanelInfo
		log.Warn("not plat(%d) in vip_max_price_map(%+v)", plat, s.vipConfSuitMax)
		return
	}
	if vps, err = s.doVipPanelPrices(c, mid, vpc, s.vipConfSuitMax[plat], 0, build); err != nil {
		return
	}
	if len(vps) == 0 {
		log.Warn("len(vps) == 0 s.doVipPanelPrices(%d, s.vipConfSuitMax[plat](%+v),plat(%d))", mid, s.vipConfSuitMax[plat], plat)
		vps = _emptyVipPanelInfo
		return
	}
	s.sortVipConfig(vps, sortTp)
	var selectd = false
	for _, v := range vps {
		if v.Selected == model.PanelSelected && !selectd {
			selectd = true
			continue
		}
		v.Selected = model.PanelNotSelected
	}
	if len(vps) > 0 && !selectd {
		vps[0].Selected = int32(model.PanelSelected)
	}
	return
}

// VipUserPrice 获取单个价格配置 by month.
func (s *Service) VipUserPrice(c context.Context,
	mid int64,
	month int16,
	plat int64,
	subType int8,
	ignoreAutoRenewStatus int8,
	build int64) (vp *model.VipPanelInfo, err error) {
	var (
		ok              bool
		vipConfSuitMax  int8
		vps             []*model.VipPanelInfo
		vipPriceConf    map[int8][]*model.VipPriceConfig
		vipPriceConfMap = make(map[int8][]*model.VipPriceConfig)
	)
	defer func() {
		if vp == nil {
			log.Warn("mid(%d) plat(%d) month(%d) subType(%d) vip panel price empty , must check!", mid, plat, month, subType)
		}
	}()
	if vipPriceConf, ok = s.vipPriceConf.GetPriceConfig(plat); !ok {
		log.Warn("not plat(%d) in vip_price_conf_map(%+v)", plat, vipPriceConf)
		return
	}
	if vipConfSuitMax, ok = s.vipConfSuitMax[plat]; !ok {
		log.Warn("not plat(%d) in vip_max_price_map(%+v)", plat, vipConfSuitMax)
		return
	}
	for suit, items := range vipPriceConf {
		for _, vp := range items {
			if vp.Month != month || vp.SubType != subType {
				continue
			}
			if suit > vipConfSuitMax {
				vipConfSuitMax = suit
			}
			vipPriceConfMap[suit] = append(vipPriceConf[suit], vp)
		}
	}
	if len(vipPriceConfMap) == 0 {
		return
	}
	if vps, err = s.doVipPanelPrices(c, mid, vipPriceConfMap, vipConfSuitMax, ignoreAutoRenewStatus, build); err != nil {
		return
	}
	for _, v := range vps {
		if int16(v.Month) == month && int8(v.SubType) == subType {
			vp = v
			return
		}
	}
	return
}

// doVipPanelPrices 根据用户制作适用于他的价格列表.
func (s *Service) doVipPanelPrices(c context.Context,
	mid int64,
	vipPriceConfMap map[int8][]*model.VipPriceConfig,
	vipConfSuitMax int8,
	ignoreAutoRenewStatus int8,
	build int64) (vps []*model.VipPanelInfo, err error) {
	var (
		isVip      bool
		vipPayType int8
		vipInfo    *model.VipInfoResp
		mpo        map[string]struct{}
	)
	if vipInfo, err = s.ByMid(c, mid); err != nil {
		return
	}
	if vipInfo.VipType > model.NotVip {
		isVip = true
	}
	if vipConfSuitMax > model.NewSubVIP {
		if mpo, err = s.dao.VipPayOrderSuccs(c, mid); err != nil {
			return
		}
	}
	vipMonthMap := make(map[int8]map[int16]*model.VipPriceConfig)
	// 根据适用用户配置、续期类型、购买时间、是否vip来制作vip面板内容
	for suit := range vipPriceConfMap {
		log.Info("vipPriceConfMap suit(%d)", suit)
		switch suit {
		case model.OldPackVIP:
			if mpo == nil {
				continue
			}
			for _, vpc := range vipPriceConfMap[suit] {
				if _, ok := mpo[vpc.DoSubMonthKey()]; ok {
					s.makeVipConfig(vipMonthMap, vpc, suit, build)
				}
			}
		case model.NewPackVIP:
			for _, vpc := range vipPriceConfMap[suit] {
				if _, ok := mpo[vpc.DoSubMonthKey()]; !ok {
					s.makeVipConfig(vipMonthMap, vpc, suit, build)
				}
			}
		case model.OldSubVIP:
			for _, vpc := range vipPriceConfMap[suit] {
				if vipInfo.AutoRenewed == model.IsAutoRenewed {
					s.makeVipConfig(vipMonthMap, vpc, suit, build)
				}
			}
		case model.NewSubVIP:
			for _, vpc := range vipPriceConfMap[suit] {
				if vipInfo.AutoRenewed != model.IsAutoRenewed {
					s.makeVipConfig(vipMonthMap, vpc, suit, build)
				}
			}
		case model.OldVIP:
			for _, vpc := range vipPriceConfMap[suit] {
				if isVip {
					s.makeVipConfig(vipMonthMap, vpc, suit, build)
				}
			}
		case model.NewVIP:
			for _, vpc := range vipPriceConfMap[suit] {
				if !isVip {
					s.makeVipConfig(vipMonthMap, vpc, suit, build)
				}
			}
		case model.AllUser:
			for _, vpc := range vipPriceConfMap[suit] {
				s.makeVipConfig(vipMonthMap, vpc, suit, build)
			}
		}
	}
	if vipInfo != nil {
		vipPayType = int8(vipInfo.PayType)
	}
	if ignoreAutoRenewStatus == 1 {
		vipPayType = model.General
	}
	for subVip := range vipMonthMap {
		if subVipConf, ok := vipMonthMap[subVip]; ok {
			for k, v := range subVipConf {
				vp := &model.VipPanelInfo{
					Month:    int32(k),
					PdName:   v.PdName,
					PdID:     v.PdID,
					SubType:  int32(v.SubType),
					SuitType: int32(v.SuitType),
					OPrice:   v.OPrice,
					DPrice:   v.DPrice,
					DRate:    v.Superscript,
					Remark:   v.Remark,
					Id:       v.ID,
					Selected: v.Selected,
				}
				if vp.DRate == "" {
					vp.DRate = v.FormatRate()
				}
				switch vipPayType {
				case model.General:
					vps = append(vps, vp)
				case model.AutoRenew:
					if subVip == model.General {
						vps = append(vps, vp)
					}
					// else if subVip == model.AutoRenew && vipInfo.PayChannelID != model.IapPayChannelID && vipInfo.VipStatus == model.Expire {
					// 	// 非IAP续费 即使is_auto_renew=1 & vip_status=0时自动续期套餐还是要露出
					// 	vps = append(vps, vp)
					// }
				}
			}
		}
	}
	return
}

// sortVipConfig 排序月份和续期
func (s *Service) sortVipConfig(vps []*model.VipPanelInfo, sortTp int8) {
	sort.Slice(vps, func(i int, j int) bool {
		if sortTp == model.PanelMonthDESC {
			return vps[i].Month > vps[j].Month
		}
		return vps[i].Month < vps[j].Month
	})
	sort.Slice(vps, func(i int, j int) bool {
		return vps[i].SubType > vps[j].SubType
	})
}

// makeVipConfig 此接口用于vip面板数据合并
func (s *Service) makeVipConfig(vipMonthMap map[int8]map[int16]*model.VipPriceConfig,
	vpc *model.VipPriceConfig,
	localSuit int8,
	build int64) {
	var (
		exists bool
		vm     *model.VipPriceConfig
	)
	// build check
	if !vpc.FilterBuild(build) {
		return
	}
	// 初始化map
	if _, ok := vipMonthMap[vpc.SubType]; !ok {
		vipMonthMap[vpc.SubType] = make(map[int16]*model.VipPriceConfig)
	}
	// 默认对象初始化
	if vm, exists = vipMonthMap[vpc.SubType][vpc.Month]; !exists {
		vipMonthMap[vpc.SubType][vpc.Month] = vpc
		return
	}
	// 内存中的适用人群从大到小写入到内存
	if vm.SuitType > localSuit {
		return
	}
	// 不能跨适用对象组写数据
	if vm.TopSuitType != vpc.TopSuitType {
		vipMonthMap[vpc.SubType][vpc.Month] = vpc
	}
}

//VipUserPanelV4 vip user panel v4
func (s *Service) VipUserPanelV4(c context.Context, arg *model.ArgPanel) (res *model.VipPirceResp, err error) {
	var (
		vps                          []*model.VipPanelInfo
		cp                           *coumol.CouponAllowancePanelInfo
		pirce                        float64
		sid                          int64
		exist                        int8
		prodLimMonth, prodLimRenewal int8
	)
	res = new(model.VipPirceResp)
	if vps, err = s.VipUserPanel(c, arg.Mid, s.GetPlatID(c, arg.Platform, arg.PanelType, arg.MobiApp, arg.Device, arg.Build), arg.SortTp, arg.Build); err != nil {
		return
	}
	if len(vps) <= 0 {
		return
	}
	if arg.Month != 0 {
		for _, v := range vps {
			v.Selected = 0
			if arg.Month == v.Month && arg.SubType == v.SubType {
				v.Selected = model.PanelSelected
				break
			}
		}
	}
	for _, v := range vps {
		if v.Selected == model.PanelSelected {
			pirce = v.DPrice
			sid = v.Id
			prodLimMonth = int8(v.Month)
			prodLimRenewal = model.MapProdLlimRenewal[int8(v.SubType)]
			break
		}
	}
	// privilege
	if res.Privileges, err = s.PrivilegesBySid(c, &model.ArgPrivilegeBySid{
		Sid:      sid,
		Platform: arg.Platform,
	}); err != nil {
		return
	}
	res.Vps = vps
	res.CodeSwitch = s.c.Property.CodeSwitch
	res.GiveSwitch = s.c.Property.GiveSwitch
	if s.c.Property.AllowanceSwitch == model.SwitchClose {
		return
	}
	res.CouponSwith = s.c.Property.AllowanceSwitch
	if pirce <= 0 {
		return
	}
	if cp, err = s.couRPC.UsableAllowanceCoupon(c, &coumol.ArgAllowanceCoupon{Mid: arg.Mid, Pirce: pirce, Platform: int(s.GetPlatID(c, arg.Platform, arg.PanelType, arg.MobiApp, arg.Device, arg.Build)), ProdLimMonth: prodLimMonth, ProdLimRenewal: prodLimRenewal}); err != nil {
		err = errors.WithStack(err)
		return
	}
	if cp != nil && cp.Amount > 0 {
		res.CouponInfo = cp
		exist = 1
	}
	if exist == 0 {
		for _, v := range res.Vps {
			colres, err1 := s.couRPC.AllowanceCouponPanel(c, &coumol.ArgAllowanceCoupon{
				Mid:      arg.Mid,
				Pirce:    v.DPrice,
				Platform: int(s.GetPlatID(c, arg.Platform, arg.PanelType, arg.MobiApp, arg.Device, arg.Build)),
			})
			if err1 != nil {
				log.Error("allowance coupon panel error(%+v)", err1)
			}
			if colres != nil && (len(colres.Usables) > 0 || len(colres.Using) > 0) {
				exist = 1
				break
			}
		}
	}
	res.ExistCoupon = exist
	return
}

//VipUserPanelV5 vip user panel v5
func (s *Service) VipUserPanelV5(c context.Context, arg *model.ArgPanel) (res *model.VipPirceResp5, err error) {
	var (
		vps                          []*model.VipPanelInfo
		cp                           *coumol.CouponAllowancePanelInfo
		pirce                        float64
		pmap                         = map[int8]*model.PrivilegesResp{}
		prodLimMonth, prodLimRenewal int8
	)
	res = new(model.VipPirceResp5)
	if vps, err = s.VipUserPanel(c, arg.Mid, s.GetPlatID(c, arg.Platform, arg.PanelType, arg.MobiApp, arg.Device, arg.Build), arg.SortTp, arg.Build); err != nil {
		return
	}
	if len(vps) <= 0 {
		return
	}
	if arg.Month != 0 {
		for _, v := range vps {
			v.Selected = 0
			if arg.Month == v.Month && arg.SubType == v.SubType {
				v.Selected = model.PanelSelected
				break
			}
		}
	}
	for _, v := range vps {
		if v.Selected == model.PanelSelected {
			pirce = v.DPrice
			prodLimMonth = int8(v.Month)
			prodLimRenewal = model.MapProdLlimRenewal[int8(v.SubType)]
		}
		switch {
		case v.Month >= _annualMonth:
			pmap[model.OnlyAnnualPrivilege] = new(model.PrivilegesResp)
			v.Type = int32(model.OnlyAnnualPrivilege)
		default:
			pmap[model.AllPrivilege] = new(model.PrivilegesResp)
			v.Type = int32(model.AllPrivilege)
		}
	}
	res.Vps = vps
	// 安卓国际版不增加 代金券,权益信息返回 激活码,好友赠送入口关闭
	if arg.MobiApp == "android_i" {
		return
	}
	// privilege
	for k := range pmap {
		pmap[k], _ = s.PrivilegesList(c, k, arg.Lang)
	}
	res.Privileges = pmap
	res.CodeSwitch = s.c.Property.CodeSwitch
	res.GiveSwitch = s.c.Property.GiveSwitch
	// iOS platform not support coupon
	if s.c.Property.AllowanceSwitch == model.SwitchClose || arg.Platform == "ios" {
		return
	}
	res.CouponSwith = s.c.Property.AllowanceSwitch
	if pirce <= 0 {
		return
	}
	if cp, err = s.couRPC.UsableAllowanceCoupon(c, &coumol.ArgAllowanceCoupon{Mid: arg.Mid, Pirce: pirce, Platform: int(s.GetPlatID(c, arg.Platform, arg.PanelType, arg.MobiApp, arg.Device, arg.Build)), ProdLimMonth: prodLimMonth, ProdLimRenewal: prodLimRenewal}); err != nil {
		err = errors.WithStack(err)
		return
	}
	if cp != nil && cp.Amount > 0 {
		res.CouponInfo = cp
	}
	return
}

//VipUserPanelV9 vip user panel v9
func (s *Service) VipUserPanelV9(c context.Context, arg *v1.VipUserPanelReq) (res *model.VipPirceRespV9, err error) {
	var (
		vps    []*model.VipPanelInfo
		pmap   = map[int8]*model.PrivilegesResp{}
		coupon *colapi.UsableAllowanceCouponV2Reply
	)
	res = new(model.VipPirceRespV9)
	platID := s.GetPlatID(c, arg.Platform, arg.PanelType, arg.MobiApp, arg.Device, arg.Build)
	if vps, err = s.VipUserPanel(c, arg.Mid, platID, int8(arg.SortTp), arg.Build); err != nil {
		return
	}
	if len(vps) <= 0 {
		return
	}
	if arg.Month != 0 {
		for _, v := range vps {
			v.Selected = 0
			if arg.Month == v.Month && arg.SubType == v.SubType {
				v.Selected = model.PanelSelected
				break
			}
		}
	}
	for _, v := range vps {
		switch {
		case v.Month >= _annualMonth:
			pmap[model.OnlyAnnualPrivilege] = new(model.PrivilegesResp)
			v.Type = int32(model.OnlyAnnualPrivilege)
		default:
			pmap[model.AllPrivilege] = new(model.PrivilegesResp)
			v.Type = int32(model.AllPrivilege)
		}
	}
	res.Vps = vps
	// 安卓国际版不增加 代金券,权益信息返回 激活码,好友赠送入口关闭
	if arg.MobiApp == "android_i" {
		return
	}
	// privilege
	for k := range pmap {
		pmap[k], _ = s.PrivilegesList(c, k, arg.Lang)
	}
	res.Privileges = pmap
	res.CodeSwitch = s.c.Property.CodeSwitch
	res.GiveSwitch = s.c.Property.GiveSwitch
	// iOS platform not support coupon
	if s.c.Property.AllowanceSwitch == model.SwitchClose || arg.Platform == "ios" {
		return
	}
	res.CouponSwith = s.c.Property.AllowanceSwitch
	if coupon, err = s.bestCoupon(c, arg.Mid, int64(platID), vps); err != nil {
		return
	}
	res.Coupon = coupon
	return
}

// VipPrice get vip pirce.
func (s *Service) VipPrice(c context.Context,
	mid int64,
	month int16,
	plat int64,
	subType int8,
	token string,
	pstr string,
	build int64) (res *model.VipPirce, err error) {
	var (
		vp *model.VipPanelInfo
		cp *coumol.CouponAllowanceInfo
	)
	res = new(model.VipPirce)
	if vp, err = s.VipUserPrice(c, mid, month, plat, subType, 0, build); err != nil {
		return
	}
	res.Panel = vp
	if token == "" || vp == nil {
		return
	}
	if cp, err = s.couRPC.JudgeCouponUsable(c, &coumol.ArgJuageUsable{
		Mid:            mid,
		Pirce:          vp.DPrice,
		CouponToken:    token,
		Platform:       int(plat),
		ProdLimMonth:   int8(vp.Month),
		ProdLimRenewal: model.MapProdLlimRenewal[int8(vp.SubType)],
	}); err != nil {
		return
	}
	res.Coupon = cp
	return
}

// VipPriceV2 get vip price v2
func (s *Service) VipPriceV2(c context.Context, a *model.ArgPriceV2) (res *model.VipPirce, err error) {
	var (
		vp   *model.VipPanelInfo
		cp   *coumol.CouponAllowanceInfo
		plat = s.GetPlatID(c, a.Platform, a.PanelType, a.MobiApp, a.Device, a.Build)
	)
	res = new(model.VipPirce)
	if vp, err = s.VipUserPrice(c, a.Mid, a.Month, plat, a.SubType, 0, a.Build); err != nil {
		return
	}
	res.Panel = vp
	if a.Token == "" || vp == nil {
		return
	}
	if cp, err = s.couRPC.JudgeCouponUsable(c, &coumol.ArgJuageUsable{
		Mid:            a.Mid,
		Pirce:          vp.DPrice,
		CouponToken:    a.Token,
		Platform:       int(plat),
		ProdLimMonth:   int8(a.Month),
		ProdLimRenewal: model.MapProdLlimRenewal[a.SubType],
	}); err != nil {
		return
	}
	res.Coupon = cp
	return
}

// VipPanelExplain vip panel explain.
func (s *Service) VipPanelExplain(c context.Context, a *model.ArgPanelExplain) (res *model.VipPanelExplain, err error) {
	var v *model.VipInfo
	res = new(model.VipPanelExplain)
	res.BackgroundURL = s.c.Property.PanelBgURL
	if a.Mid == 0 {
		res.Explain = model.UserNotLoginExplain
		return
	}
	if v, err = s.VipInfo(c, a.Mid); err != nil {
		return
	}
	if v.VipStatus == model.NotVip {
		if v.VipOverdueTime.Time().Unix() == 0 {
			res.Explain = model.NotVipExplain
			return
		}
		res.Explain = model.ExpireVipExplain
		return
	}
	d := (v.VipOverdueTime.Time().Unix() - time.Now().Unix()) / int64(_daysecond)
	switch {
	case d <= 0:
		res.Explain = model.ExpireVipExplain
		return
	case d > _willexpiredays || v.VipPayType != model.NormalPay:
		res.Explain = fmt.Sprintf(model.YYYYDDVipExplain, v.VipOverdueTime.Time().Format(_yyyymmdd))
		return
	default:
		res.Explain = fmt.Sprintf(model.WillExplainVipExplain, d)
	}
	return
}

// PriceByProductID get price by product id
func (s *Service) PriceByProductID(c context.Context, productID string) (res *model.VipPriceConfig, err error) {
	var vds []*model.VipDPriceConfig
	if vds, err = s.dao.VipPriceDiscountByProductID(c, productID); err != nil {
		return
	}
	for _, v := range vds {
		if res = s.vipPriceMap[v.ID]; res != nil {
			res.PdID = v.PdID
			res.DPrice = v.DPrice
			res.Remark = v.Remark
			break
		}
	}
	if res != nil {
		return
	}
	if res, err = s.dao.VipPriceByProductID(c, productID); err != nil {
		return
	}
	if res != nil {
		res.DPrice = res.OPrice
	}
	return
}

func (s *Service) setPanelTypeByID(c context.Context, p *model.VipPriceConfig) (err error) {
	var conf *model.ConfPlatform
	if conf, err = s.PlatformByID(c, int64(p.Plat)); err != nil {
		return
	}
	if conf == nil {
		return ecode.VipPlatformByIDNotFoundErr
	}
	p.PanelType = conf.PanelType
	return
}

// VipPriceByID vip price by id.
func (s *Service) VipPriceByID(c context.Context, a *model.ArgVipPriceByID) (vpc *model.VipPriceConfig, err error) {
	if vpc, err = s.dao.VipPriceByID(c, a.ID); err != nil {
		return
	}
	if vpc == nil {
		err = ecode.VipPriceInfoNotFoundErr
		return
	}
	err = s.setPanelTypeByID(c, vpc)
	return
}
