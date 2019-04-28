package service

import (
	"context"

	"go-common/app/interface/main/ugcpay/model"
)

// TradeCreate 订单创建
func (s *Service) TradeCreate(ctx context.Context, mid int64, platform string, oid int64, otype string, currency string) (orderID string, payData string, err error) {
	if platform == "" {
		platform = "web"
	}
	orderID, payData, err = s.dao.TradeCreate(ctx, platform, mid, oid, otype, currency)
	return
}

// TradeQuery 订单查询
func (s *Service) TradeQuery(ctx context.Context, orderID string) (order *model.TradeOrder, err error) {
	order, err = s.dao.TradeQuery(ctx, orderID)
	return
}

// TradeConfirm 订单二次确认
func (s *Service) TradeConfirm(ctx context.Context, orderID string) (order *model.TradeOrder, err error) {
	order, err = s.dao.TradeConfirm(ctx, orderID)
	return
}

// TradeCancel 订单取消
func (s *Service) TradeCancel(ctx context.Context, orderID string) (err error) {
	err = s.dao.TradeCancel(ctx, orderID)
	return
}
