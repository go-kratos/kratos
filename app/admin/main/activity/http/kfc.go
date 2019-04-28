package http

import (
	kfcmdl "go-common/app/admin/main/activity/model/kfc"
	bm "go-common/library/net/http/blademaster"
)

func kfcList(c *bm.Context) {
	arg := new(kfcmdl.ListParams)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(kfcSrv.List(c, arg))
}
