package http

import (
	"net"

	"go-common/app/service/main/vip/model"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

func createOrder2(c *bm.Context) {
	var (
		err error
		r   *model.CreateOrderRet
	)
	arg := new(model.ArgCreateOrder2)
	if err = c.Bind(arg); err != nil {
		return
	}
	arg.IP = net.ParseIP(metadata.String(c, metadata.RemoteIP))
	r, _, err = vipSvc.CreateOrder2(c, arg)
	c.JSON(r, err)
}

func createQrCodeOrder(c *bm.Context) {
	var err error
	arg := new(model.ArgCreateOrder2)
	if err = c.Bind(arg); err != nil {
		return
	}
	arg.IP = net.ParseIP(metadata.String(c, metadata.RemoteIP))
	c.JSON(vipSvc.CreateQrCodeOrder(c, arg))
}

func grantAssociateVip(c *bm.Context) {
	var err error
	arg := new(model.ArgEleVipGrant)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, vipSvc.EleVipGrant(c, arg))
}
