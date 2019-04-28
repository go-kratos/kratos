package service

import (
	"context"
	rpcV1 "go-common/app/service/openplatform/ticket-sales/api/grpc/v1"
	"go-common/app/service/openplatform/ticket-sales/model"
	"go-common/app/service/openplatform/ticket-sales/model/consts"
	"go-common/library/ecode"
	"time"
)

//GetPromo 获取拼团详情
func (s *Service) GetPromo(c context.Context, arg *rpcV1.PromoID) (res *rpcV1.Promo, err error) {
	var promo *model.Promotion
	if promo, err = s.dao.Promo(c, arg.PromoID); err != nil || promo == nil {
		err = ecode.TicketPromotionLost
		return
	}
	res = &rpcV1.Promo{
		PromoID:       promo.PromoID,
		Type:          promo.Type,
		ItemID:        promo.ItemID,
		SKUID:         promo.SKUID,
		Extra:         promo.Extra,
		ExpireSec:     promo.ExpireSec,
		SKUCount:      promo.SKUCount,
		Amount:        promo.Amount,
		BuyerCount:    promo.BuyerCount,
		BeginTime:     promo.BeginTime,
		EndTime:       promo.EndTime,
		Status:        promo.Status,
		Ctime:         int64(promo.Ctime),
		Mtime:         int64(promo.Mtime),
		PrivSKUID:     promo.PrivSKUID,
		UsableCoupons: promo.UsableCoupons,
	}
	return
}

//CreatePromo 创建拼团
func (s *Service) CreatePromo(c context.Context, arg *rpcV1.CreatePromoRequest) (res *rpcV1.PromoID, err error) {
	var promo *model.Promotion
	if promo, err = s.dao.Promo(c, arg.PromoID); err != nil {
		err = ecode.TicketUnKnown
		return
	}

	if promo != nil {
		err = ecode.PromoExists
		return
	}

	res = new(rpcV1.PromoID)
	promo = &model.Promotion{
		PromoID:       arg.PromoID,
		Type:          arg.Type,
		ItemID:        arg.ItemID,
		SKUID:         arg.SKUID,
		Extra:         arg.Extra,
		ExpireSec:     arg.ExpireSec,
		SKUCount:      arg.SKUCount,
		Amount:        arg.Amount,
		BuyerCount:    arg.BuyerCount,
		BeginTime:     arg.BeginTime,
		EndTime:       arg.EndTime,
		PrivSKUID:     arg.PrivSKUID,
		UsableCoupons: arg.UsableCoupons,
	}
	if res.PromoID, err = s.dao.CreatePromo(c, promo); err != nil {
		err = ecode.TicketUnKnown
		return
	}
	return
}

//HasPromoOfSKU 判断是否有相同sku已上架的活动
func (s *Service) HasPromoOfSKU(c context.Context, skuID int64, beginTime int64, endTime int64) (num int64, err error) {
	if num, err = s.dao.HasPromoOfSKU(c, skuID, beginTime, endTime); err != nil {
		err = ecode.TicketUnKnown
		return
	}
	return
}

//OperatePromo 修改活动状态
func (s *Service) OperatePromo(c context.Context, arg *rpcV1.OperatePromoRequest) (res *rpcV1.CommonResponse, err error) {
	var promo *model.Promotion
	if promo, err = s.dao.Promo(c, arg.PromoID); err != nil || promo == nil {
		err = ecode.TicketPromotionLost
		return
	}
	var fromStatus, toStatus int16

	if arg.OperateType == consts.DoUpShelf {
		var hasPromo int64
		if hasPromo, err = s.HasPromoOfSKU(c, promo.SKUID, promo.BeginTime, promo.EndTime); err != nil {
			return
		}
		if 0 != hasPromo {
			err = ecode.TicketPromoExistSameTime
			return
		}
		fromStatus = consts.PromoWaitShelf
		toStatus = consts.PromoUpShelf
	} else if arg.OperateType == consts.DoDelShelf {
		fromStatus = consts.PromoUpShelf
		toStatus = consts.PromoDelShelf
	} else {
		err = ecode.IllegalPromoOperate
		return
	}

	if fromStatus != promo.Status {
		err = ecode.PromoStatusChanged
		return
	}

	res = new(rpcV1.CommonResponse)
	if res.Res, err = s.dao.OperatePromo(c, arg.PromoID, fromStatus, toStatus); err != nil {
		err = ecode.TicketUnKnown
		return
	}

	s.dao.DelCachePromo(c, arg.PromoID)
	return
}

//CheckPromoStatus 检测拼团状态，skuID可为0，为0则不检测skuID是否一致
func (s *Service) CheckPromoStatus(c context.Context, promoID int64, skuID int64) (res *model.Promotion, err error) {
	if res, err = s.dao.Promo(c, promoID); err != nil || res == nil {
		err = ecode.TicketPromotionLost
		return
	}
	if skuID != 0 && res.SKUID != skuID {
		err = ecode.TicketPromotionLost
		return
	}

	currentTime := time.Now().Unix()
	if res.Status != consts.PromoUpShelf || res.BeginTime > currentTime || res.EndTime < currentTime {
		err = ecode.TicketPromotionEnd
		return
	}
	return
}

//EditPromo 编辑Promo  to do
func (s *Service) EditPromo(c context.Context, arg *rpcV1.EditPromoRequest) (res *rpcV1.CommonResponse, err error) {
	res = new(rpcV1.CommonResponse)
	var promo *model.Promotion
	currentTime := time.Now().Unix()
	if promo, err = s.dao.Promo(c, arg.PromoID); err != nil || promo == nil {
		err = ecode.TicketPromotionLost
		return
	}

	if promo.BeginTime < currentTime {
		switch promo.Status {
		case consts.PromoFinishShelf, consts.PromoDelShelf:
			err = ecode.PromoEditFieldNotALlowed
			return
		case consts.PromoWaitShelf:
			return
		case consts.PromoUpShelf:
			if promo.EndTime < currentTime {
				err = ecode.PromoEditFieldNotALlowed
				return
			}
			if promo.EndTime > currentTime && (arg.ExpireSec != 0 || arg.EndTime != 0 || arg.BeginTime != 0 || arg.SKUCount != 0 || arg.PrivSKUID != 0 || arg.Amount != 0 || arg.UsableCoupons == "") {
				err = ecode.PromoEditFieldNotALlowed
				return
			}
		}
	}

	if arg.ExpireSec != 0 {
		promo.ExpireSec = arg.ExpireSec
	}
	//if arg.SKUCount != 0 {
	//	promo.SKUCount = arg.SKUCount
	//}
	//if arg.BeginTime != 0 {
	//	promo.BeginTime = arg.BeginTime
	//}
	//if arg.EndTime != 0 {
	//	promo.EndTime = arg.EndTime
	//}
	//if arg.Amount != 0 {
	//	arg.Amount = promo.Amount
	//}
	//if arg.PrivSKUID != 0 {
	//	promo.PrivSKUID = arg.PrivSKUID
	//}
	//
	//if arg.UsableCoupons != "" {
	//	promo.UsableCoupons = arg.UsableCoupons
	//}

	if res.Res, err = s.dao.UpdatePromo(c, promo); err != nil {
		return
	}

	s.dao.DelCachePromo(c, arg.PromoID)

	return
}
