package http

import (
	"go-common/app/admin/main/activity/model"
	bm "go-common/library/net/http/blademaster"
)

func archives(c *bm.Context) {
	p := &model.ArchiveParam{}
	if err := c.Bind(p); err != nil {
		return
	}
	c.JSON(actSrv.Archives(c, p))
}
