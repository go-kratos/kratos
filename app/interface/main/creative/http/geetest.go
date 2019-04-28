package http

import (
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"strconv"
)

func gtPreProcessAdd(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	process, err := gtSvc.PreProcessAdd(c, mid, ip, "web", 1)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(process, nil)
}

func gtPreProcess(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	process, err := gtSvc.PreProcess(c, mid, ip, "web", 1)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(process, nil)
}

func gtValidate(c *bm.Context) {
	req := c.Request
	ip := metadata.String(c, metadata.RemoteIP)
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	challenge := req.Form.Get("geetest_challenge")
	validate := req.Form.Get("geetest_validate")
	seccode := req.Form.Get("geetest_seccode")
	success := req.Form.Get("geetest_success")
	successi, err := strconv.Atoi(success)
	if err != nil {
		successi = 1
	}
	status := gtSvc.Validate(c, challenge, validate, seccode, "web", ip, successi, mid)
	if !status {
		c.JSON(nil, ecode.CreativeGeetestErr)
	}
	c.JSON(nil, nil)
}
