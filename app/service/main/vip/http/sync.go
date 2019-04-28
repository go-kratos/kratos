package http

import (
	"go-common/app/service/main/vip/model"
	bm "go-common/library/net/http/blademaster"
)

func syncUser(c *bm.Context) {
	var err error
	user := new(model.VipUserInfo)
	if err = c.Bind(user); err != nil {
		return
	}
	vipSvc.SyncUser(c, user)
	c.JSON(nil, nil)
}
