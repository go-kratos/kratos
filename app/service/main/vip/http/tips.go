package http

import (
	"go-common/app/service/main/vip/model"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// tips info.
func tips(c *bm.Context) {
	var (
		res []*model.TipsResp
		err error
		arg = new(model.ArgTips)
	)
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	if res, err = vipSvc.Tips(c, arg); err != nil {
		log.Error("vipSvc.Tips(%+v) err(%+v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(res, nil)
}
