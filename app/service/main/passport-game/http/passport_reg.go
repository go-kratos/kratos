package http

import (
	"go-common/app/service/main/passport-game/model"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

func regV3(c *bm.Context) {
	var argRegV3 = new(model.ArgRegV3)
	err := c.Bind(argRegV3)
	if err != nil {
		return
	}

	var cookie = c.Request.Header.Get("Cookie")
	c.JSON(srv.RegV3(c, model.TdoRegV3{Arg: *argRegV3, IP: metadata.String(c, metadata.RemoteIP), Cookie: cookie}))
}

func regV2(c *bm.Context) {
	var argRegV2 = new(model.ArgRegV2)
	err := c.Bind(argRegV2)
	if err != nil {
		return
	}

	var cookie = c.Request.Header.Get("Cookie")
	c.JSON(srv.RegV2(c, model.TdoRegV2{Arg: *argRegV2, IP: metadata.String(c, metadata.RemoteIP), Cookie: cookie}))
}

func reg(c *bm.Context) {
	var argReg = new(model.ArgReg)
	err := c.Bind(argReg)
	if err != nil {
		return
	}

	var cookie = c.Request.Header.Get("Cookie")
	c.JSON(srv.Reg(c, model.TdoReg{Arg: *argReg, IP: metadata.String(c, metadata.RemoteIP), Cookie: cookie}))
}

func byTel(c *bm.Context) {
	var argByTel = new(model.ArgByTel)
	err := c.Bind(argByTel)
	if err != nil {
		return
	}

	var cookie = c.Request.Header.Get("Cookie")
	c.JSON(srv.ByTel(c, model.TdoByTel{Arg: *argByTel, IP: metadata.String(c, metadata.RemoteIP), Cookie: cookie}))
}

func captcha(c *bm.Context) {
	c.JSON(srv.Captcha(c, metadata.String(c, metadata.RemoteIP)))
}

func sendSms(c *bm.Context) {
	var sendSmsp = new(model.SendSms)
	berr := c.Bind(sendSmsp)
	if berr != nil {
		return
	}

	var cookie = c.Request.Header.Get("Cookie")
	err := srv.SendSms(c, model.TdoSendSms{Arg: *sendSmsp, IP: metadata.String(c, metadata.RemoteIP), Cookie: cookie})
	c.JSON(nil, err)
}
