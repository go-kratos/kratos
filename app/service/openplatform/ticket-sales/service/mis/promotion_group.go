package mis

import (
	"context"
	rpcV1 "go-common/app/service/openplatform/ticket-sales/api/grpc/v1"
	"go-common/app/service/openplatform/ticket-sales/model"
	"go-common/library/ecode"
)

//GetGroupOrdersMis 获取拼团订单
func (m *Mis) GetGroupOrdersMis(c context.Context, arg *rpcV1.GetGroupOrdersMisRequest) (res *rpcV1.GetGroupOrdersMisResponse, err error) {
	var groupID int64
	if arg.GroupID != 0 {
		groupID = arg.GroupID
	} else if arg.OrderID != 0 {
		var promoOrder *model.PromotionOrder
		if promoOrder, err = m.dao.PromoOrder(c, arg.OrderID); err != nil || promoOrder == nil {
			err = ecode.TicketPromotionOrderLost
			return
		}
		groupID = promoOrder.GroupID
	}

	if groupID == 0 {
		err = ecode.TicketPromotionGroupLost
		return
	}

	var orders []*model.PromotionOrder
	if orders, err = m.dao.PromoOrders(c, groupID); err != nil || orders == nil {
		err = ecode.TicketPromotionOrderLost
		return
	}

	res = new(rpcV1.GetGroupOrdersMisResponse)
	for _, order := range orders {
		tempOrder, _ := m.dao.Orders(c, &model.OrderMainQuerier{OrderID: []int64{order.OrderID}})
		r := &rpcV1.PromoOrderMis{
			PromoID:  order.PromoID,
			GroupID:  order.GroupID,
			OrderID:  order.OrderID,
			IsMaster: order.IsMaster,
			UID:      order.UID,
			Status:   order.Status,
			Ctime:    int64(order.Ctime),
			PayTime:  tempOrder[0].PayTime,
			SKUID:    order.SKUID,
		}
		res.Orders = append(res.Orders, r)
	}
	return
}
