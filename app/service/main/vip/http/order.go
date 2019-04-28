package http

import (
	"fmt"

	"go-common/app/service/main/vip/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

func status(c *bm.Context) {
	var (
		o   *model.OrderInfo
		vip *model.VipInfo
		err error
	)
	arg := new(struct {
		OrderNo string `form:"order_no" validate:"required"`
	})
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	if o, err = vipSvc.OrderInfo(c, arg.OrderNo); err != nil {
		log.Error("vipSvc.OrderInfo(%s)  err(%+v)", arg.OrderNo, err)
		c.JSON(nil, err)
		return
	}
	if o == nil {
		c.JSON(nil, ecode.VipOrderNotFoundErr)
		return
	}

	res := new(struct {
		OrderNo string              `json:"order_no"`
		Status  int8                `json:"status"`
		Message *model.OrderMessage `json:"message"`
	})

	res.OrderNo = o.OrderNo
	res.Status = o.Status

	if o.Status == model.SUCCESS {
		if vip, err = vipSvc.VipInfo(c, o.Mid); err != nil {
			c.JSON(nil, err)
			return
		}
		message := new(model.OrderMessage)
		message.RightButton = "知道了"
		message.Title = "开通成功"
		message.Content = fmt.Sprintf("你已成功开通%d个月大会员，目前有效期%s", o.BuyMonths, vip.VipOverdueTime.Time().Format("2006-01-02"))
		res.Message = message
	} else if o.Status == model.FAILED {
		message := new(model.OrderMessage)
		message.RightButton = "知道了"
		message.Title = "支付失败"
		message.Content = fmt.Sprintf("订单号:%s \nUID:%d \n支付失败了，试试重新购买吧。", o.OrderNo, vip.Mid)
		res.Message = message
	}

	c.JSON(res, nil)
}

func orders(c *bm.Context) {
	var (
		orders []*model.PayOrderResp
		total  int64
		err    error
	)
	arg := new(struct {
		Mid int64 `form:"mid" validate:"required,min=1,gte=1"`
		Ps  int   `form:"ps" default:"20" validate:"min=0,max=50"`
		Pn  int   `form:"pn" default:"1"`
	})
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	if orders, total, err = vipSvc.OrderList(c, arg.Mid, arg.Pn, arg.Ps); err != nil {
		log.Error("vipSvc.OrderList(%d) err(%+v)", arg.Mid, err)
		c.JSON(nil, err)
		return
	}
	res := new(struct {
		Data  []*model.PayOrderResp `json:"list"`
		Total int64                 `json:"total"`
	})
	res.Data = orders
	res.Total = total
	c.JSON(orders, nil)
}

func createOrder(c *bm.Context) {
	var (
		err error
		pp  map[string]interface{}
	)
	arg := new(model.ArgCreateOrder)
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	if pp, err = vipSvc.CreateOrder(c, arg, metadata.String(c, metadata.RemoteIP)); err != nil {
		log.Error("vipSvc.CreateOrder(%d) err(%+v)", arg.Mid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(pp, nil)
}

func pannelInfoNew(c *bm.Context) {
	var (
		pi  *model.PannelInfo
		err error
	)
	arg := new(model.ArgPannel)
	if err = c.Bind(arg); err != nil {
		log.Error("pannelInfoNew(%+v)", err)
		return
	}
	if pi, err = vipSvc.PannelInfoNew(c, arg.Mid, arg); err != nil {
		log.Error("vipSvc.PannelInfoNew(%d) err(%+v)", arg.Mid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(pi, nil)
}

func createOldOrder(c *bm.Context) {
	var (
		err error
	)
	arg := new(model.ArgOldPayOrder)
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	c.JSON(nil, vipSvc.CreateOldOrder(c, arg))
}

func orderMng(c *bm.Context) {
	var (
		err   error
		order *model.OrderMng
	)
	arg := new(struct {
		Mid int64 `form:"mid" validate:"required,min=1,gte=1"`
	})
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	if order, err = vipSvc.OrderMng(c, arg.Mid); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(order, err)
}

func rescision(c *bm.Context) {
	var (
		err error
	)

	arg := new(struct {
		Mid        int64 `form:"mid" validate:"required,mid=1,gte=1"`
		DeviceType int32 `form:"deviceType" validate:"required"`
	})

	if err = c.Bind(arg); err != nil {
		return
	}

	err = vipSvc.Rescision(c, arg.Mid, arg.DeviceType)
	c.JSON(nil, err)

}
