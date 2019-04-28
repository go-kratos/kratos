package http

import (
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func porderCfgList(c *bm.Context) {
	cfgList, err := vdpSvc.PorderCfgList(c)
	if err != nil {
		log.Error("vdpSvc.porderCfgList() error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(cfgList, nil)
}

func porderArcList(c *bm.Context) {
	params := c.Request.Form
	begin := params.Get("begin")
	end := params.Get("end")
	data, err := vdpSvc.PorderArcList(c, begin, end)
	if err != nil {
		log.Error("vdpSvc.PorderArcList() error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}
