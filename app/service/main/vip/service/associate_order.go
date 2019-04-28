package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/service/main/vip/model"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

// CreateAssociateOrder create associate order.
func (s *Service) CreateAssociateOrder(c context.Context, a *model.ArgCreateOrder2) (r *model.CreateOrderRet, err error) {
	var (
		oi      *model.OpenInfo
		o       *model.OpenBindInfo
		eleType int32
		order   *model.PayOrder
	)
	if o, err = s.dao.RawBindInfoByMid(c, a.Mid, a.AppID); err != nil {
		return
	}
	if o == nil {
		return nil, ecode.VipActivityAccountNotSupport
	}
	if oi, err = s.dao.OpenInfoByMid(c, a.Mid, a.AppID); err != nil {
		return
	}
	if oi == nil {
		return nil, ecode.VipAssociateOpenIDNotExsitErr
	}
	// 商品限定
	if a.Bmid > 0 {
		// not suuport give friend.
		return nil, ecode.VipEleUnionBuyGiveFirendErr
	}
	if a.AppID == model.EleAppID && a.PanelType != "ele" {
		return nil, ecode.VipEleUnionBuyCanProductErr
	}
	switch a.Month {
	case model.UnionOneMonth:
		eleType = model.EleMonthVip
	case model.UnionOneYear:
		eleType = model.EleYearVip
	default:
		err = ecode.VipEleUnionBuyCanProductErr
		return
	}
	// ele 风控检查
	if _, err = s.dao.EleCanPurchase(c, &model.ArgEleCanPurchase{
		ElemeOpenID: o.OutOpenID,
		BliOpenID:   oi.OpenID,
		UserIP:      IPStr(a.IP),
		VipType:     eleType,
	}); err != nil {
		return
	}
	if r, order, err = s.CreateOrder2(c, a); err != nil {
		return
	}
	//TODO Deprecated add old order.
	if err = s.dao.AddOldPayOrder(c, &model.VipOldPayOrder{
		OrderNo:     order.OrderNo,
		AppID:       order.AppID,
		Mid:         order.Mid,
		BuyMonths:   order.BuyMonths,
		Money:       order.Money,
		Status:      order.Status,
		Ver:         order.Ver,
		Platform:    order.Platform,
		AppSubID:    order.AppSubID,
		Bmid:        order.ToMid,
		OrderType:   order.OrderType,
		CouponMoney: order.CouponMoney,
		PID:         order.PID,
		UserIP:      order.UserIP,
	}); err != nil {
		return
	}
	//TODO Deprecated add old recharge order.
	err = s.dao.AddOldRechargeOrder(c, &model.VipOldRechargeOrder{
		AppID:        order.AppID,
		PayMid:       order.Mid,
		OrderNo:      s.orderID(),
		RechargeBp:   order.Money,
		PayOrderNO:   order.OrderNo,
		Status:       order.Status,
		Remark:       "",
		Ver:          1,
		ThirdTradeNO: "",
	})
	return
}

// BilibiliVipGrant bilibili associate vip grant [third -> bilibili]
func (s *Service) BilibiliVipGrant(c context.Context, a *model.ArgBilibiliVipGrant) (err error) {
	var (
		oi         *model.OpenInfo
		ob         *model.OpenBindInfo
		count      int64
		limitCount int64
		days       int64
		grantC     *model.VipAssociateGrantCount
	)
	if limitCount = s.c.AssociateConf.GrantDurationMap[fmt.Sprintf("%d", a.Duration)]; limitCount == 0 {
		return ecode.VipAssociateGrantDurationErr
	}
	if days = model.EleGrantVipDays[a.Duration]; days == 0 {
		return ecode.VipAssociateGrantDayErr
	}
	if oi, err = s.dao.RawOpenInfoByOpenID(c, a.OpenID, a.AppID); err != nil {
		return
	}
	if oi == nil {
		return ecode.VipAssociateOpenIDNotExsitErr
	}
	if ob, err = s.dao.RawBindInfoByMid(c, oi.Mid, a.AppID); err != nil {
		return
	}
	if ob == nil || ob.OutOpenID != a.OutOpenID {
		return ecode.VipActivityAccountNotSupport
	}
	if count, err = s.dao.CountGrantOrderByOutTradeNo(c, a.OutOrderNO, a.AppID); err != nil {
		return
	}
	if count > 0 {
		// return ecode.VipAssociateGrantOutTradeNoExsitErr
		return
	}
	if grantC, err = s.associateGrantCount(c, oi.Mid, a.AppID, a.Duration); err != nil {
		return
	}
	if grantC.CurrentCount >= limitCount {
		return ecode.VipAssociateGrantLimitErr
	}
	if err = s.addAssociateGrant(c, &model.VipOrderAssociateGrant{
		AppID:      a.AppID,
		Mid:        oi.Mid,
		OutOpenID:  ob.OutOpenID,
		OutTradeNO: a.OutOrderNO,
		Ctime:      xtime.Time(time.Now().Unix()),
		Months:     a.Duration,
		State:      model.AssociateGrantStateHadGrant,
	}, grantC, limitCount, days); err != nil {
		return
	}
	// update bind state had pay.
	if err = s.UpdateBindState(c, &model.OpenBindInfo{
		Mid:   oi.Mid,
		State: model.AssociateBindStatePurchased,
		AppID: a.AppID,
		Ver:   ob.Ver,
	}); err != nil {
		return
	}
	// salary mail coupon
	if err1 := s.MailCouponCodeCreate(c, ob.Mid); err1 != nil {
		log.Error("s.MailCouponCodeCreate[ele->bilibili](%d) err(%+v)", ob.Mid, err1)
		return
	}
	return
}

func (s *Service) addAssociateGrant(c context.Context, oa *model.VipOrderAssociateGrant, grantC *model.VipAssociateGrantCount, limitCount int64, days int64) (err error) {
	var (
		hv  *model.OldHandlerVip
		aff int64
	)
	oldtx, err := s.dao.OldStartTx(c)
	if err != nil {
		return
	}
	defer func() {
		if err == nil {
			if err = oldtx.Commit(); err != nil {
				oldtx.Rollback()
				return
			}
			s.dao.DelVipInfoCache(context.Background(), oa.Mid)
			s.cache(func() {
				s.oldCleanCacheAndNotify(context.Background(), hv, "")
			})
			s.asyncBcoin(func() {
				s.OldProcesserHandler(context.Background(), hv, "")
			})
		} else {
			oldtx.Rollback()
		}
	}()
	if aff, err = s.dao.TxInsertAssociateGrantOrder(oldtx, oa); err != nil {
		return
	}
	if aff != 1 {
		return ecode.VipAssociateGrantOutTradeNoExsitErr
	}
	if hv, err = s.OldUpdateVipWithHistory(c, oldtx, &model.VipChangeBo{
		ChangeType: model.ChangeTypeSystem,
		Remark:     model.ElemeGrantRemark,
		RelationID: oa.OutTradeNO,
		Days:       days,
		Mid:        oa.Mid,
	}); err != nil {
		err = errors.WithStack(err)
		return
	}
	// limit count
	if err = s.dao.UpdateAssociateGrantCount(c, grantC); err != nil {
		return
	}
	return
}

// EleVipGrant ele vip grant [bilibili -> third].
func (s *Service) EleVipGrant(c context.Context, a *model.ArgEleVipGrant) (err error) {
	var (
		act     *model.VipOrderActivityRecord
		oi      *model.OpenInfo
		o       *model.VipPayOrderOld
		ob      *model.OpenBindInfo
		eleType int32
		aff     int64
	)
	if act, err = s.dao.ActivityOrder(c, a.OrderNO); err != nil {
		return
	}
	if act == nil {
		return ecode.VipActivityOrderNotFoundErr
	}
	if act.AssociateState == model.AssociateGrantStateHadGrant {
		return
	}
	if act.PanelType != model.PanelTypeEle {
		return ecode.VipPriceNotEleProductErr
	}
	if o, err = s.dao.SelOldPayOrder(c, a.OrderNO); err != nil {
		return
	}
	if o == nil {
		return ecode.VipOrderInfoNotFoundErr
	}
	if o.Status != model.SUCCESS {
		return ecode.VipOrderStatusNotSuccessErr
	}
	// check open_id
	if oi, err = s.dao.OpenInfoByMid(c, o.Mid, o.AppID); err != nil {
		return
	}
	if oi == nil {
		return ecode.VipAssociateOpenIDNotExsitErr
	}
	if ob, err = s.dao.RawBindInfoByMid(c, o.Mid, o.AppID); err != nil {
		return
	}
	if ob == nil {
		return ecode.VipActivityAccountNotSupport
	}
	switch int32(o.BuyMonths) {
	case model.UnionOneMonth:
		eleType = model.EleMonthVip
	case model.UnionOneYear:
		eleType = model.EleYearVip
	default:
		err = ecode.VipEleUnionBuyCanProductErr
		return
	}
	if _, err = s.dao.EleBindUnion(c, &model.ArgEleBindUnion{
		ElemeOpenID: ob.OutOpenID,
		BliOpenID:   oi.OpenID,
		SourceID:    o.OrderNo,
		UserIP:      IPStr(o.UserIP),
		VipType:     eleType,
	}); err != nil {
		return
	}
	// update act order.
	if aff, err = s.dao.UpdateActivityState(c, model.AssociateGrantStateHadGrant, a.OrderNO); err != nil {
		return
	}
	if aff == 1 {
		// update bind state had pay.
		if err = s.UpdateBindState(c, &model.OpenBindInfo{
			Mid:   oi.Mid,
			State: model.AssociateBindStatePurchased,
			AppID: ob.AppID,
			Ver:   ob.Ver,
		}); err != nil {
			return
		}
		// salary mail coupon
		if err1 := s.MailCouponCodeCreate(c, o.Mid); err1 != nil {
			log.Error("s.MailCouponCodeCreate[bilibili->ele](%d) err(%+v)", o.Mid, err1)
			return
		}
	}
	return
}

// associateGrantCount associate grant limit.
func (s *Service) associateGrantCount(c context.Context, mid int64, appID int64, months int32) (res *model.VipAssociateGrantCount, err error) {
	if res, err = s.dao.AssociateGrantCountInfo(c, mid, appID, months); err != nil {
		return
	}
	if res != nil {
		return
	}
	// init.
	res = &model.VipAssociateGrantCount{
		AppID:        appID,
		Mid:          mid,
		Months:       months,
		CurrentCount: 0,
	}
	if err = s.dao.AddAssociateGrantCount(c, res); err != nil {
		return
	}
	return
}
