package service

import (
	"context"
	"go-common/app/common/openplatform/random"
	rpcV1 "go-common/app/service/openplatform/ticket-sales/api/grpc/v1"
	"go-common/app/service/openplatform/ticket-sales/model"
	"go-common/app/service/openplatform/ticket-sales/model/consts"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"math"
	"time"
)

//GetPromoOrder 获取拼团订单详情
func (s *Service) GetPromoOrder(c context.Context, orderID int64) (res *rpcV1.PromoOrder, err error) {
	var promoOrder *model.PromotionOrder
	if promoOrder, err = s.dao.PromoOrder(c, orderID); err != nil || promoOrder == nil {
		err = ecode.TicketPromotionOrderLost
		return
	}
	res = &rpcV1.PromoOrder{
		PromoID:  promoOrder.PromoID,
		GroupID:  promoOrder.GroupID,
		OrderID:  promoOrder.OrderID,
		IsMaster: promoOrder.IsMaster,
		UID:      promoOrder.UID,
		Status:   promoOrder.Status,
		Ctime:    int64(promoOrder.Ctime),
		Mtime:    int64(promoOrder.Mtime),
		SKUID:    promoOrder.SKUID,
	}
	return
}

//CheckCreateStatus 检测用户是否满足创单条件
func (s *Service) CheckCreateStatus(c context.Context, arg *rpcV1.CheckCreatePromoOrderRequest) (res *rpcV1.CheckCreatePromoOrderResponse, err error) {
	var promo *model.Promotion
	if promo, err = s.CheckPromoStatus(c, arg.PromoID, arg.SKUID); err != nil {
		return
	}
	var promoOrderInfo *model.PromotionOrder
	if promoOrderInfo, err = s.GetUserJoinPromoOrder(c, arg.PromoID, promo.SKUCount, arg.GroupID, arg.UID); err != nil && err != ecode.TicketPromotionRepeatJoin {
		return
	}
	res = &rpcV1.CheckCreatePromoOrderResponse{
		Amount:        promo.Amount,
		SKUID:         promo.SKUID,
		PrivSKUID:     promo.PrivSKUID,
		UsableCoupons: promo.UsableCoupons,
	}
	if err == ecode.TicketPromotionRepeatJoin {
		err = nil
		res.RepeatOrder = &rpcV1.RepeatOrder{
			OrderID:  promoOrderInfo.OrderID,
			IsMaster: promoOrderInfo.IsMaster,
			Status:   promoOrderInfo.Status,
		}
	}
	return
}

//GetUserJoinPromoOrder 获取用户正在参与该活动的订单
func (s *Service) GetUserJoinPromoOrder(c context.Context, promoID int64, skuCount int64, groupID int64, uid int64) (res *model.PromotionOrder, err error) {
	currentTime := time.Now().Unix()
	res = new(model.PromotionOrder)
	if groupID == 0 {
		//1.1找出该活动正在拼团中且未过期的拼团
		if _, err = s.dao.GetUserGroupDoing(c, promoID, uid, consts.GroupDoing); err != nil && err != sql.ErrNoRows {
			log.Warn(err.Error())
			return
		}
		if err == sql.ErrNoRows {
			//1.2该活动没有正在拼的团，找到用户发起未支付的拼团
			if res, err = s.dao.PromoOrderByStatus(c, promoID, groupID, uid, consts.PromoOrderUnpaid); err != nil && err != sql.ErrNoRows {
				log.Warn(err.Error())
				return
			}

			if err == sql.ErrNoRows {
				//未找到
				err = nil
				return
			}
			//找到，重复参加
			err = ecode.TicketPromotionRepeatJoin
			return
		}

		//1.3该活动有正在拼的团，查看该用户的拼团订单信息  to do
		if res, err = s.dao.PromoOrderByStatus(c, promoID, groupID, uid, consts.PromoOrderPaid); err != nil && err != sql.ErrNoRows {
			log.Warn(err.Error())
			return
		}

		err = ecode.TicketPromotionRepeatJoin
		return
	}
	var promoGroup *model.PromotionGroup
	//2.1检查该拼团的状态
	if promoGroup, err = s.dao.PromoGroup(c, groupID); err != nil || promoGroup == nil {
		err = ecode.TicketPromotionGroupLost
		return
	}
	if promoGroup.Status != consts.GroupDoing || promoGroup.ExpireAt < currentTime {
		err = ecode.TicketPromotionGroupLost
		return
	}

	if promoGroup.OrderCount >= skuCount {
		err = ecode.TicketPromotionGroupFull
		return
	}
	//2.2检查用户是否参与过该拼团
	if res, err = s.dao.PromoOrderDoing(c, promoID, groupID, uid); err != nil && err != sql.ErrNoRows {
		log.Warn(err.Error())
		return
	}

	if err == sql.ErrNoRows {
		err = nil
		return
	}

	err = ecode.TicketPromotionRepeatJoin
	return
}

//CreatePromoOrder 创建活动订单
func (s *Service) CreatePromoOrder(c context.Context, arg *rpcV1.CreatePromoOrderRequest) (res *rpcV1.OrderID, err error) {
	var (
		number int64
		status = consts.PromoOrderUnpaid
	)
	if arg.PayMoney == 0 {
		status = consts.PromoOrderPaid
	}
	if arg.GroupID == 0 {
		//团长订单
		if status == consts.PromoOrderUnpaid {
			//非0元单
			if number, err = s.dao.AddPromoOrder(c, arg.PromoID, 0, arg.OrderID, 1, arg.UID, status, arg.PromoSKUID, arg.Ctime); err != nil || number == 0 {
				err = ecode.TicketAddPromoOrderFail
				return
			}
			res = &rpcV1.OrderID{OrderID: arg.OrderID}
			return
		}

		//0元单
		var conn *sql.Tx
		if conn, err = s.dao.BeginTx(c); err != nil {
			err = ecode.TicketUnKnown
			return
		}
		var groupID = random.Uniqid(15)
		if number, err = s.dao.TxAddPromoOrder(c, conn, arg.PromoID, groupID, arg.OrderID, 1, arg.UID, status, arg.PromoSKUID, arg.Ctime); err != nil || number == 0 {
			err = ecode.TicketAddPromoOrderFail
			conn.Rollback()
			return
		}
		var promo *model.Promotion
		if promo, err = s.dao.Promo(c, arg.PromoID); err != nil || promo == nil {
			err = ecode.TicketPromotionLost
			conn.Rollback()
			return
		}
		currentTime := time.Now().Unix()
		expireAt := int64(math.Max(math.Min(float64(promo.EndTime), float64(promo.ExpireSec+currentTime)), float64(currentTime)))
		if number, err = s.dao.TxAddPromoGroup(c, conn, arg.PromoID, groupID, arg.UID, 1, status, expireAt); err != nil || number == 0 {
			err = ecode.TicketAddPromoGroupFail
			conn.Rollback()
			return
		}
		conn.Commit()
		res = &rpcV1.OrderID{OrderID: arg.OrderID}
		return
	}

	//非团长订单
	var conn *sql.Tx
	if conn, err = s.dao.BeginTx(c); err != nil {
		err = ecode.TicketUnKnown
		return
	}
	if number, err = s.dao.TxAddPromoOrder(c, conn, arg.PromoID, arg.GroupID, arg.OrderID, 0, arg.UID, status, arg.PromoSKUID, arg.Ctime); err != nil || number == 0 {
		err = ecode.TicketAddPromoOrderFail
		conn.Rollback()
		return
	}
	var promo *model.Promotion
	if promo, err = s.dao.Promo(c, arg.PromoID); err != nil || promo == nil {
		err = ecode.TicketPromotionLost
		conn.Rollback()
		return
	}
	if number, err = s.dao.TxUpdateGroupOrderCount(c, conn, 1, arg.GroupID, promo.SKUCount); err != nil || number == 0 {
		err = ecode.TicketUpdatePromoGroupFail
		conn.Rollback()
		return
	}
	conn.Commit()
	s.dao.DelCachePromoGroup(c, arg.GroupID)
	s.dao.DelCachePromoOrders(c, arg.GroupID)
	res = &rpcV1.OrderID{OrderID: arg.OrderID}
	return
}

//CancelOrder 取消订单
func (s *Service) CancelOrder(c context.Context, arg *rpcV1.OrderID) (res *rpcV1.OrderID, err error) {
	var promoOrder *model.PromotionOrder
	if promoOrder, err = s.dao.PromoOrder(c, arg.OrderID); err != nil || promoOrder == nil {
		err = ecode.TicketPromotionOrderLost
		return
	}
	var conn *sql.Tx
	if conn, err = s.dao.BeginTx(c); err != nil {
		err = ecode.TicketUnKnown
		return
	}
	var number int64
	if number, err = s.dao.TxUpdatePromoOrderStatus(c, conn, arg.OrderID, consts.PromoOrderCancel); err != nil {
		err = ecode.TicketUpdatePromoOrderFail
		conn.Rollback()
		return
	}
	if number == 0 {
		conn.Rollback()
		res = &rpcV1.OrderID{OrderID: 0}
		return
	}
	if promoOrder.GroupID == 0 {
		conn.Commit()
		s.dao.DelCachePromoOrder(c, arg.OrderID)
		res = arg
		return
	}
	if number, err = s.dao.TxUpdateGroupOrderCount(c, conn, -1, promoOrder.GroupID, 100000); err != nil || number == 0 {
		err = ecode.TicketUpdatePromoGroupFail
		conn.Rollback()
		return
	}
	conn.Commit()
	s.dao.DelCachePromoOrder(c, arg.OrderID)
	s.dao.DelCachePromoGroup(c, promoOrder.GroupID)
	s.dao.DelCachePromoOrders(c, promoOrder.GroupID)
	res = arg
	return
}

// PromoPayNotify 支付通知
func (s *Service) PromoPayNotify(c context.Context, arg *rpcV1.OrderID) (res *rpcV1.OrderID, err error) {
	var promoOrder *model.PromotionOrder
	if promoOrder, err = s.dao.PromoOrder(c, arg.OrderID); err != nil || promoOrder == nil {
		err = ecode.TicketPromotionOrderLost
		return
	}

	var promo *model.Promotion
	if promo, err = s.dao.Promo(c, promoOrder.PromoID); err != nil || promo == nil {
		err = ecode.TicketPromotionLost
		return
	}
	if promo.Status != consts.PromoUpShelf {
		err = ecode.TicketPromotionEnd
	}
	if promoOrder.GroupID != 0 {
		//非团长订单
		if _, err = s.dao.UpdatePromoOrderStatus(c, arg.OrderID, consts.PromoOrderPaid); err != nil {
			err = ecode.TicketUpdatePromoOrderFail
			return
		}
		s.dao.DelCachePromoOrder(c, arg.OrderID)
		s.dao.DelCachePromoOrders(c, promoOrder.GroupID)
		res = arg
		return
	}

	//团长订单
	var (
		conn   *sql.Tx
		number int64
	)
	if conn, err = s.dao.BeginTx(c); err != nil {
		err = ecode.TicketUnKnown
		return
	}
	currentTime := time.Now().Unix()
	groupID := random.Uniqid(15)
	expireAt := int64(math.Max(math.Min(float64(promo.EndTime), float64(promo.ExpireSec+currentTime)), float64(currentTime)))
	if number, err = s.dao.TxAddPromoGroup(c, conn, promoOrder.PromoID, groupID, promoOrder.UID, 1, consts.GroupDoing, expireAt); err != nil || number == 0 {
		err = ecode.TicketAddPromoGroupFail
		conn.Rollback()
		return
	}
	if number, err = s.dao.TxUpdatePromoOrderGroupIDAndStatus(c, conn, arg.OrderID, groupID, consts.PromoOrderPaid); err != nil || number == 0 {
		err = ecode.TicketUpdatePromoOrderFail
		conn.Rollback()
		return
	}
	conn.Commit()
	s.dao.DelCachePromoOrder(c, arg.OrderID)
	res = arg
	return
}

//CheckIssue 出票检查
func (s *Service) CheckIssue(c context.Context, arg *rpcV1.OrderID) (res *rpcV1.CheckIssueResponse, err error) {
	var promoOrder *model.PromotionOrder
	if promoOrder, err = s.dao.PromoOrder(c, arg.OrderID); err != nil || promoOrder == nil {
		err = ecode.TicketPromotionOrderLost
		return
	}
	if promoOrder.IsMaster == 0 || promoOrder.Status != consts.PromoOrderPaid {
		err = ecode.TicketPromoOrderTypeErr
		return
	}
	var promoGroup *model.PromotionGroup
	if promoGroup, err = s.dao.PromoGroup(c, promoOrder.GroupID); err != nil || promoGroup == nil {
		err = ecode.TicketPromotionGroupLost
		return
	}
	if promoGroup.Status != consts.GroupDoing {
		err = ecode.TicketPromoGroupStatusErr
		return
	}
	var groupOrders []*model.PromotionOrder
	if groupOrders, err = s.dao.GroupOrdersByStatus(c, promoOrder.GroupID, consts.PromoOrderPaid); err != nil {
		err = ecode.TicketUnKnown
		return
	}
	var promo *model.Promotion
	if promo, err = s.dao.RawPromo(c, promoGroup.PromoID); err != nil {
		err = ecode.TicketPromotionLost
		return
	}
	stock := int64(len(groupOrders))
	if stock < promo.SKUCount {
		err = ecode.TicketPromotionGroupNotFull
		return
	}
	res = new(rpcV1.CheckIssueResponse)
	res.PromoID = promoOrder.PromoID
	res.GroupID = promoOrder.GroupID
	for _, v := range groupOrders {
		temp := &rpcV1.OrderID{OrderID: v.OrderID}
		res.PaidOrders = append(res.PaidOrders, temp)
	}
	return
}

//FinishIssue 完成出票
func (s *Service) FinishIssue(c context.Context, arg *rpcV1.FinishIssueRequest) (res *rpcV1.GroupID, err error) {
	var (
		number int64
		promo  *model.Promotion
		conn   *sql.Tx
	)
	if promo, err = s.dao.Promo(c, arg.PromoID); err != nil || promo == nil {
		err = ecode.TicketPromotionLost
		return
	}
	if conn, err = s.dao.BeginTx(c); err != nil {
		err = ecode.TicketUnKnown
		return
	}
	if number, err = s.dao.TxUpdateGroupStatus(c, conn, arg.GroupID, consts.GroupDoing, consts.GroupSuccess); err != nil {
		conn.Rollback()
		err = ecode.TicketUpdatePromoGroupFail
		return
	}
	if number == 0 {
		res = &rpcV1.GroupID{GroupID: arg.GroupID}
		return
	}
	if number, err = s.dao.TxUpdatePromoBuyerCount(c, conn, arg.GroupID, promo.SKUCount); err != nil || number == 0 {
		conn.Rollback()
		err = ecode.TicketUpdatePromoGroupFail
		return
	}
	conn.Commit()
	s.dao.DelCachePromo(c, arg.PromoID)
	s.dao.DelCachePromoGroup(c, arg.GroupID)
	s.dao.DelCachePromoOrders(c, arg.GroupID)
	res = &rpcV1.GroupID{GroupID: arg.GroupID}
	return
}

// PromoRefundNotify 退款通知
func (s *Service) PromoRefundNotify(c context.Context, arg *rpcV1.OrderID) (res *rpcV1.OrderID, err error) {
	var promoOrder *model.PromotionOrder
	if promoOrder, err = s.dao.PromoOrder(c, arg.OrderID); err != nil || promoOrder == nil {
		err = ecode.TicketPromotionOrderLost
		return
	}

	if _, err = s.dao.UpdatePromoOrderStatus(c, arg.OrderID, consts.PromoOrderRefund); err != nil {
		err = ecode.TicketUpdatePromoOrderFail
		return
	}
	s.dao.DelCachePromoGroup(c, promoOrder.GroupID)
	s.dao.DelCachePromoOrder(c, promoOrder.OrderID)
	s.dao.DelCachePromoOrders(c, promoOrder.GroupID)
	res = &rpcV1.OrderID{OrderID: promoOrder.OrderID}
	return
}
