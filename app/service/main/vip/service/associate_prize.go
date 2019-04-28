package service

import (
	"context"

	col "go-common/app/service/main/coupon/model"
	"go-common/app/service/main/vip/model"
	"go-common/library/ecode"
)

// ThirdPrizeGrant associate prize grant[bilibili->third].
func (s *Service) ThirdPrizeGrant(c context.Context, a *model.ArgThirdPrizeGrant) (err error) {
	var (
		b  *model.OpenBindInfo
		oi *model.OpenInfo
	)
	if b, err = s.dao.BindInfoByMid(c, a.Mid, a.AppID); err != nil {
		return
	}
	if b == nil {
		return ecode.VipActivityAccountNotSupport
	}
	if oi, err = s.dao.OpenInfoByMid(c, a.Mid, a.AppID); err != nil {
		return
	}
	if oi == nil {
		return ecode.VipAssociateOpenIDNotExsitErr
	}
	if a.PrizeType == model.AssociatePrizeTypeCode {
		if a.PrizeKey == 0 {
			err = ecode.VipAssociatePrizeKeyErr
			return
		}
		// grant code.
		if err = s.ResourceBatchOpenVip(c, &model.ArgUseBatch{
			Mid:     a.Mid,
			OrderNo: a.UniqueNo,
			Remark:  a.Remark,
			Appkey:  a.Appkey,
			BatchID: a.PrizeKey,
		}); err != nil {
			return
		}
	}
	if a.PrizeType == model.AssociatePrizeTypeEleBag {
		// ele prize grant
		var data []*model.EleReceivePrizesResp
		if data, err = s.dao.EleUnionReceivePrizes(c, &model.ArgEleReceivePrizes{
			BliOpenID:   oi.OpenID,
			ElemeOpenID: b.OutOpenID,
			SourceID:    a.UniqueNo,
		}); err != nil {
			return
		}
		if len(data) < 0 {
			err = ecode.VipEleUnionReceivePrizesErr
			return
		}
	}
	return
}

// BilibiliPrizeGrant vip prize grant for third [third->bilibili].
func (s *Service) BilibiliPrizeGrant(c context.Context, a *model.ArgBilibiliPrizeGrant) (res *col.SalaryCouponForThirdResp, err error) {
	var (
		oi         *model.OpenInfo
		batchToken string
	)
	if batchToken = s.c.AssociateConf.BilibiliPrizeGrantKeyMap[a.PrizeKey]; batchToken == "" {
		return nil, ecode.VipAssociatePrizeKeyErr
	}
	if oi, err = s.dao.RawOpenInfoByOpenID(c, a.OpenID, a.AppID); err != nil {
		return
	}
	if oi == nil {
		return nil, ecode.VipAssociateOpenIDNotExsitErr
	}
	if res, err = s.couRPC.SalaryCouponForThird(context.Background(), &col.ArgSalaryCoupon{
		Mid:        oi.Mid,
		CouponType: col.CouponAllowance,
		Origin:     col.AllowanceBusinessReceive,
		Count:      1,
		BatchToken: batchToken,
		AppID:      a.AppID, //联合会员APPID
		UniqueNo:   a.UniqueNo,
	}); err != nil {
		return
	}
	return
}
