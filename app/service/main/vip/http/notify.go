package http

import (
	"encoding/json"
	"go-common/app/service/main/vip/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

const (
	_payNotifySuccess = "SUCCESS"
	_payNotifyFail    = "FAIL"
)

func notify(c *bm.Context) {
	d := new(model.PayCallBackResult)
	if err := c.Bind(d); err != nil {
		log.Error("pr.Bind err(%+v)", err)
		return
	}
	if d.TradeNO == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if d.OutTradeNO == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if d.TradeStatus != model.TradeSuccess {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err := vipSvc.PayNotify(c, d); err != nil {
		log.Error("s.PayNotify  err(%+v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func notify2(c *bm.Context) {
	var (
		err error
	)
	arg := new(struct {
		MsgContent string `form:"msgContent" validate:"required"`
	})
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		c.Writer.Write([]byte(_payNotifyFail))
		return
	}
	p := &model.PayNotifyContent{}
	if err = json.Unmarshal([]byte(arg.MsgContent), p); err != nil {
		c.Writer.Write([]byte(_payNotifyFail))
		return
	}
	if err = vipSvc.PayNotify2(c, p); err != nil {
		log.Error("s.PayNotify2 err(%+v)", err)
		if err == ecode.VipOrderAlreadyHandlerErr {
			c.Writer.Write([]byte(_payNotifySuccess))
			return
		}
		c.Writer.Write([]byte(_payNotifyFail))
		return
	}
	c.Writer.Write([]byte(_payNotifySuccess))
}

func signNotify(c *bm.Context) {
	var (
		err error
	)
	arg := new(struct {
		MsgContent string `form:"msgContent" validate:"required"`
	})
	if err = c.Bind(arg); err != nil {
		return
	}
	p := new(model.PaySignNotify)
	if err = json.Unmarshal([]byte(arg.MsgContent), p); err != nil {
		c.Writer.Write([]byte(_payNotifyFail))
		return
	}
	if err = vipSvc.PaySignNotify(c, p); err != nil {
		log.Error("vip.paySignNotify(%+v) error(%+v)", p, err)
		c.Writer.Write([]byte(_payNotifyFail))
		return
	}
	c.Writer.Write([]byte(_payNotifySuccess))
}

func refundOrderNotify(c *bm.Context) {
	var (
		err error
	)
	arg := new(struct {
		MsgContent string `form:"msgContent" validate:"required"`
	})
	if err = c.Bind(arg); err != nil {
		return
	}
	p := new(model.PayRefundNotify)
	log.Info("refun order notify params:%+v", arg.MsgContent)
	if err = json.Unmarshal([]byte(arg.MsgContent), p); err != nil {
		log.Error("error(%+v)", err)
		c.Writer.Write([]byte(_payNotifyFail))
		return
	}
	if err = vipSvc.RefundNotify(c, p); err != nil {
		log.Error("vip.refundNotify(%+v)  error(%+v)", arg.MsgContent, err)
		c.Writer.Write([]byte(_payNotifyFail))
		return
	}
	c.Writer.Write([]byte(_payNotifySuccess))
}
