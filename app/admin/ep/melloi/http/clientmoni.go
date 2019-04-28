package http

import (
	"go-common/app/admin/ep/melloi/model"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func queryClientMoni(c *bm.Context) {

	var (
		cli    model.ClientMoni
		resMap = make(map[string]interface{})
	)
	if err := c.BindWith(&cli, binding.Form); err != nil {
		c.JSON(nil, err)
		return
	}
	clientMonis, err := srv.QueryClientMoni(&cli)
	if err != nil {
		log.Error("srv.QueryClientMoni err (%v)", err)
		return
	}
	resMap["clientMonis"] = clientMonis
	c.JSON(resMap, err)
}
