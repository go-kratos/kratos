package service

import (
	"context"
	"encoding/json"
	"net/url"

	"go-common/app/service/openplatform/ticket-sales/api/grpc/v1"
	"go-common/app/service/openplatform/ticket-sales/model"
	"go-common/app/service/openplatform/ticket-sales/model/consts"
	"go-common/library/ecode"
	"go-common/library/log"
)

// PayNotify 支付回调通知接口
func (s *Service) PayNotify(c context.Context, req v1.PayNotifyRequest) (resp v1.PayNotifyResponse, err error) {
	log.Info("s.PayNotify(%s)", req.String())
	if req.MsgID == "" || req.MsgContent == "" {
		err = ecode.ParamInvalid
		return
	}

	msgContent, err := url.QueryUnescape(req.MsgContent)
	if err != nil {
		log.Error("s.PayNotify() MsgContent=%s error(%v)", req.MsgContent, err)
		return
	}

	mc := model.MsgContent{}
	if err = json.Unmarshal([]byte(msgContent), &mc); err != nil {
		log.Error("s.PayNotify() msgContent=%s error(%v)", msgContent, err)
		return
	}

	// 不是测试模式要验证签名
	if !req.TestMode && !mc.ValidSign() {
		err = ecode.SignCheckErr
		return
	}

	if mc["payStatus"] == consts.PayStatusSuccess {
		var charge *model.Charge
		if charge, err = mc.ToCharge(); err != nil {
			log.Error("s.PayNotify() mc.ToCharge(%v) error(%v)", mc, err)
			return
		}

		s.chargeCallback(c, charge, req.TestMode)
	}
	return
}

// chargeCallback 支付回调处理
func (s *Service) chargeCallback(c context.Context, charge *model.Charge, testMode bool) (err error) {
	orders, err := s.dao.RawOrders(c, &model.OrderMainQuerier{OrderID: []int64{charge.OrderID}})
	if err != nil {
		log.Error("s.chargeCallback() s.dao.RawUserOrders(%d) error(%v)", charge.OrderID, err)
		return
	}
	if len(orders) == 0 {
		log.Error("s.chargeCallback() s.dao.RawUserOrders() 订单 ID %d 未找到", charge.OrderID)
		return
	}
	order := orders[0]

	// 测试模式却不是测试项目
	if testMode && !s.isTestProject(order.ItemID) {
		err = ecode.ParamInvalid
		log.Warn("testMode %t but not test project", testMode)
		return
	}

	if order.PayMoney != charge.Amount {
		log.Warn("订单金额不匹配 order: %d, charge: %d", order.PayMoney, charge.Amount)
		return
	}

	return
}

func (s *Service) isTestProject(itemID int64) bool {
	for _, ID := range s.c.TestProject.IDs {
		if ID == itemID {
			return true
		}
	}
	return false
}
