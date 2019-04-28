package http

import (
	"go-common/app/admin/main/upload/model"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func addDir(c *bm.Context) {
	var err error
	adp := &model.AddDirParam{}
	if err = c.BindWith(adp, binding.FormPost); err != nil {
		return
	}

	c.JSON(nil, uaSvc.AddDir(c, adp))
}
