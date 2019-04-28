package http

import (
	"go-common/app/interface/main/answer/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func pendantRec(c *bm.Context) {
	arg := new(model.ReqPendant)
	mid, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.AccountNotLogin)
		return
	}
	if err := c.Bind(arg); err != nil {
		return
	}
	arg.MID = mid.(int64)
	c.JSON(nil, svc.PendantRec(c, arg))
}
