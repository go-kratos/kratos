package http

import (
	"go-common/app/admin/ep/merlin/model"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func genMachinesV2(c *bm.Context) {
	var (
		gmr      = &model.GenMachinesRequest{}
		err      error
		username string
	)

	if username, err = getUsername(c); err != nil {
		return
	}

	if err = c.BindWith(gmr, binding.JSON); err != nil {
		return
	}

	if err = gmr.Verify(); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, svc.GenMachinesV2(c, gmr, username))
}
