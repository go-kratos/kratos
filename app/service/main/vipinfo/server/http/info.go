package http

import (
	"go-common/app/service/main/vipinfo/model"
	bm "go-common/library/net/http/blademaster"
)

func info(c *bm.Context) {
	var err error
	arg := new(model.ArgMid)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(srv.Info(c, arg.Mid))
}

func infos(c *bm.Context) {
	var err error
	arg := new(model.ArgMids)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(srv.Infos(c, arg.Mids))
}
