package http

import (
	"go-common/app/interface/main/player/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
)

func playurl(c *blademaster.Context) {
	var (
		mid int64
		err error
	)
	arg := new(model.PlayurlArg)
	if err = c.Bind(arg); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if arg.OType == "" {
		arg.OType = model.OtypeXML
	}
	if arg.OType != model.OtypeXML && arg.OType != model.OtypeJSON {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if arg.Player != 0 && arg.Player != 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if arg.Player == 1 && arg.OType == model.OtypeXML {
		log.Warn("playurl warn arg(%+v)", arg)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if arg.OType == model.OtypeJSON {
		c.JSON(playSvr.Playurl(c, mid, arg))
	} else {
		c.XML(playSvr.Playurl(c, mid, arg))
	}
}
