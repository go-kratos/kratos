package http

import (
	"github.com/pkg/errors"

	"go-common/app/service/openplatform/ticket-sales/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func syncOrder(c *bm.Context) {
	orderStatusMap := make(map[int64]string)
	orderStatusMap[model.OrderPaid] = "已付款"
	orderStatusMap[model.OrderRefunded] = "已退款"
	arg := new(model.DistOrderArg)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	_, ok := orderStatusMap[arg.Stat]
	if !ok {
		c.JSON("stat", ecode.TicketParamInvalid)
		return
	}
	if arg.Serial == "" {
		c.JSON("serial_num is empty", ecode.TicketParamInvalid)
		return
	}
	switch arg.Stat {
	case model.OrderPaid:
		arg.Stat = model.DistOrderNormal
		if arg.RefStat == model.OrderRefundPartly {
			arg.Stat = model.DistOrderRefunded
		}
	case model.OrderRefunded:
		arg.Stat = model.DistOrderRefunded
	}
	c.JSON(svc.SyncOrder(c, arg))
}

func getOrder(c *bm.Context) {
	arg := new(model.DistOrderGetArg)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}

	c.JSON(svc.GetOrder(c, arg.Oid))
}
