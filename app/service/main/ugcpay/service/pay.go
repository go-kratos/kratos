package service

import (
	"context"

	"go-common/app/service/main/ugcpay/model"
	"go-common/library/log"
)

func (s *Service) createPayData(ctx context.Context, order *model.Order) (data string, err error) {
	if order == nil {
		return
	}
	title := s.payTitle(order.OID, order.Fee)
	payParam := s.pay.Create(order.OrderID, order.OID, order.Fee, s.pay.DeviceType(order.Platform), s.pay.ServiceType(order.Platform), order.MID, title)
	if err = s.pay.Sign(payParam); err != nil {
		return
	}
	return s.pay.ToJSON(payParam)
}

func (s *Service) refundOrder(ctx context.Context, order *model.Order, desc string) (err error) {
	if order == nil {
		return
	}
	fn := func() (err error) {
		refundParam := s.pay.Refund(order.PayID, order.Fee, desc)
		if err = s.pay.Sign(refundParam); err != nil {
			return
		}
		refundJSON, err := s.pay.ToJSON(refundParam)
		if err != nil {
			return
		}
		return s.dao.PayRefund(ctx, refundJSON)
	}
	return s.tryHard(fn, "payRefund", 2)
}

func (s *Service) cancelOrder(ctx context.Context, orderID string) (err error) {
	cancelParam := s.pay.Cancel(orderID)
	if err = s.pay.Sign(cancelParam); err != nil {
		return
	}
	cancelJSON, err := s.pay.ToJSON(cancelParam)
	if err != nil {
		return
	}
	return s.dao.PayCancel(ctx, cancelJSON)
}

func (s *Service) tryHard(fn func() error, fname string, tryTimes int) (err error) {
	for tryTimes > 0 {
		tryTimes--
		if err = fn(); err != nil {
			log.Error("func: %s, err: %+v, try times left: %d", fname, err, tryTimes)
		} else {
			log.Info("func: %s, run success, try times left: %d", fname, tryTimes)
			break
		}
	}
	return
}
