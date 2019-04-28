package service

import (
	"context"
	rpcV1 "go-common/app/service/openplatform/ticket-sales/api/grpc/v1"
	"go-common/app/service/openplatform/ticket-sales/model"
	"go-common/app/service/openplatform/ticket-sales/model/consts"
	"go-common/library/ecode"
	"time"
)

//GetPromoGroupInfo 获取拼团详情，同时检测当前团状态
func (s *Service) GetPromoGroupInfo(c context.Context, arg *rpcV1.GetPromoGroupInfoRequest) (res *rpcV1.GetPromoGroupInfoResponse, err error) {
	var promoOrder *model.PromotionOrder
	if promoOrder, err = s.dao.PromoOrder(c, arg.OrderID); err != nil || promoOrder == nil {
		err = ecode.TicketPromotionOrderLost
		return
	}
	var promoGroup *model.PromotionGroup
	if promoGroup, err = s.dao.PromoGroup(c, promoOrder.GroupID); err != nil || promoGroup == nil {
		err = ecode.TicketPromotionGroupLost
		return
	}
	var promo *model.Promotion
	if promo, err = s.dao.Promo(c, promoOrder.PromoID); err != nil || promo == nil {
		err = ecode.TicketPromotionLost
		return
	}
	if promoGroup.Status == consts.GroupDoing && promoGroup.ExpireAt <= time.Now().Unix() {
		promoGroup.Status = consts.GroupFailed
	}
	res = &rpcV1.GetPromoGroupInfoResponse{
		PromoID:    promo.PromoID,
		SKUCount:   promo.SKUCount,
		Amount:     promo.Amount,
		GroupID:    promoGroup.GroupID,
		OrderCount: promoGroup.OrderCount,
		ExpireAt:   promoGroup.ExpireAt,
		Status:     promoGroup.Status,
		Ctime:      int64(promoGroup.Ctime),
	}
	return
}

//GroupFailed 拼团失败
func (s *Service) GroupFailed(c context.Context, arg *rpcV1.GroupFailedRequest) (res *rpcV1.GroupID, err error) {
	if _, err = s.dao.UpdateGroupStatusAndOrderCount(c, arg.GroupID, arg.CancelNum, consts.GroupDoing, consts.GroupFailed); err != nil {
		err = ecode.TicketUpdatePromoGroupFail
		return
	}
	s.dao.DelCachePromoGroup(c, arg.GroupID)
	s.dao.DelCachePromoOrders(c, arg.GroupID)
	res = &rpcV1.GroupID{GroupID: arg.GroupID}
	return
}
