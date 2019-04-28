package service

import (
	"context"
	"strings"
	"time"

	"go-common/app/admin/main/vip/model"
	vipmol "go-common/app/service/main/vip/model"

	"go-common/library/ecode"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

//OrderList order list.
func (s *Service) OrderList(c context.Context, arg *model.ArgPayOrder) (res []*model.PayOrder, count int64, err error) {
	if count, err = s.dao.OrderCount(c, arg); err != nil {
		err = errors.WithStack(err)
		return
	}
	if res, err = s.dao.OrderList(c, arg); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//Refund refund order.
func (s *Service) Refund(c context.Context, orderNo, username string, refundAmount float64) (err error) {
	var (
		order *model.PayOrder
	)
	if order, err = s.dao.SelOrder(c, orderNo); err != nil {
		err = errors.WithStack(err)
		return
	}
	if order == nil {
		err = ecode.VipOrderNoErr
		return
	}
	if order.Status != model.SUCCESS {
		err = ecode.VipOrderStatusErr
		return
	}
	if order.ToMid > 0 {
		err = ecode.VipOrderToMidErr
		return
	}
	if order.AppID == vipmol.EleAppID {
		return ecode.VipEleOrderCanNotReFundErr
	}
	if order.Money < refundAmount {
		err = ecode.VipOrderMoneyErr
		return
	}
	if order.Money < refundAmount+order.RefundAmount {
		err = ecode.VipOrderMoneyErr
		return
	}
	var timefmt = time.Now().Format("20060102150405")

	if len(order.ThirdTradeNo) <= 0 || order.ThirdTradeNo[0:len(timefmt)] == timefmt {
		err = ecode.VipOldOrderErr
		return
	}
	refundID := strings.Replace(uuid.New().String(), "-", "", -1)[4:20]
	if err = s.dao.PayRefund(c, order, refundAmount, refundID); err != nil {
		err = errors.WithStack(err)
		return
	}
	olog := new(model.PayOrderLog)
	olog.OrderNo = order.OrderNo
	olog.Mid = order.Mid
	olog.RefundID = refundID
	olog.Status = model.REFUNDING
	olog.Operator = username
	olog.RefundAmount = refundAmount
	if err = s.dao.AddPayOrderLog(c, olog); err != nil {
		err = errors.WithStack(err)
	}
	return
}
