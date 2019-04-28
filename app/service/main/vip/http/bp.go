package http

import (
	"go-common/app/service/main/vip/model"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func bpList(c *bm.Context) {
	var (
		err error
		res *model.BcoinSalaryResp
	)
	arg := new(struct {
		Mid int64 `form:"mid" validate:"required"`
	})
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	if res, err = vipSvc.BcoinGive(c, arg.Mid); err != nil {
		log.Error(" BcoinGive mid(%d), err(%+v)", arg.Mid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(res, nil)
}
