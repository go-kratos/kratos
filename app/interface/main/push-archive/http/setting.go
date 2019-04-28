package http

import (
	"go-common/app/interface/main/push-archive/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"strconv"
)

func getMID(c *bm.Context) (mid int64) {
	midi, _ := c.Get("mid")
	if midi != nil {
		mid = midi.(int64)
	}
	return
}

func setting(c *bm.Context) {
	mid := getMID(c)
	c.JSON(pushSrv.Setting(c, mid))
}

func setSetting(c *bm.Context) {
	mid := getMID(c)
	tp, _ := strconv.Atoi(c.Request.Form.Get("type"))
	if tp <= 0 {
		log.Error("type(%d) is wrong", tp)
		c.JSON(nil, ecode.RequestErr)
		return
	}

	st := &model.Setting{Type: tp}
	c.JSON(nil, pushSrv.SetSetting(c, mid, st))
}
