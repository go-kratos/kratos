package http

import (
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func nsMd5(c *bm.Context) {
	params := c.Request.Form
	nidStr := params.Get("nid")
	nid, err := strconv.ParseInt(nidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", nidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	md5 := params.Get("md5")
	if len(md5) != 32 {
		log.Error("strconv.ParseInt(%s) error(%v)", md5, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = vdpSvc.AddNetSafeMd5(c, nid, md5)
	if err != nil {
		log.Error(" vdpSvc.AddNetSafeMd5(%d) error(%v)|nid(%v)|md5(%v)", nid, md5, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(map[string]interface{}{
		"nid": nid,
	}, nil)
}
