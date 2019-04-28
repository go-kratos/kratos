package service

// import (
// 	"context"
// 	"fmt"
// 	"net/url"
// 	"time"

// 	"go-common/app/job/main/ugcpay/model"
// 	"go-common/library/log"
// )

// func (s *Service) repairOrderUser() {
// 	var (
// 		ctx       = context.Background()
// 		beginTime = time.Date(2018, time.October, 1, 0, 0, 0, 0, time.Local)
// 		endTime   = time.Now()
// 		state     = "paid"
// 		limit     = 1000
// 	)

// 	maxID, orderList, err := s.dao.OrderList(ctx, beginTime, endTime, state, 0, limit)
// 	if err != nil {
// 		log.Error("%+v", err)
// 		return
// 	}
// 	log.Info("repaireOrderUser got list: %d, maxID: %d", len(orderList), maxID)
// 	for _, order := range orderList {
// 		// 修复 pay_time
// 		order.PayTime = order.CTime
// 		// 修复 real_fee
// 		asset, err := s.dao.Asset(ctx, order.OID, order.OType, order.Currency)
// 		if err != nil {
// 			log.Error("order:%+v, err: %+v", order, err)
// 			continue
// 		}
// 		order.RealFee = asset.Price
// 		// 修复 pay_id
// 		if len(order.PayID) < 2 {
// 			order.PayID, err = s.tradeQuery(ctx, order.OrderID)
// 			if err != nil {
// 				log.Error("order:%+v, err: %+v", order, err)
// 				continue
// 			}
// 		}
// 		if _, err = s.dao.UpdateOrder(ctx, order); err != nil {
// 			log.Error("order:%+v, err: %+v", order, err)
// 			err = nil
// 		}
// 		log.Info("repaireOrderUser success, order: %+v", order)
// 	}
// }

// func (s *Service) tradeQuery(ctx context.Context, orderID string) (payID string, err error) {
// 	// 2. 从支付平台获取订单状态
// 	var (
// 		params    url.Values
// 		jsonData  string
// 		orders    map[string][]*model.PayOrder
// 		payOrders []*model.PayOrder
// 		ok        bool
// 	)
// 	params = s.pay.Query(orderID)
// 	if err = s.pay.Sign(params); err != nil {
// 		return
// 	}
// 	if jsonData, err = s.pay.ToJSON(params); err != nil {
// 		return
// 	}
// 	if orders, err = s.dao.PayQuery(ctx, jsonData); err != nil {
// 		return
// 	}
// 	if payOrders, ok = orders[orderID]; !ok || len(payOrders) == 0 {
// 		log.Info("tradeQuery from pay platform not found order: %s", orderID)
// 		return
// 	}
// 	for _, po := range payOrders {
// 		if po.TXID > 0 {
// 			log.Info("tradeQuery order: %s, txID: %d", orderID, po.TXID)
// 			payID = fmt.Sprintf("%d", po.TXID)
// 			return
// 		}
// 	}
// 	return
// }
