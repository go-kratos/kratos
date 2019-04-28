package http

import (
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func webNsMd5(c *bm.Context) {
	params := c.Request.Form
	nidStr := params.Get("nid")
	nid, err := strconv.ParseInt(nidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", nidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	appkey := params.Get("appkey")
	if appkey != "bilibili" {
		log.Error("(%s) error(%v)", appkey, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	md5 := params.Get("md5")
	if len(md5) != 32 {
		log.Error("len(%s) (%d) error(%v)", md5, len(md5), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	companyID := 2
	err = arcSvc.AddNetSafeMd5(c, nid, md5)
	if err != nil {
		c.JSON(nil, err)
	}
	c.JSONMap(map[string]interface{}{
		"nid":       nid,
		"md5":       md5,
		"cid":       companyID,
		"companyId": companyID,
		"response":  "ok",
	}, nil)
}
