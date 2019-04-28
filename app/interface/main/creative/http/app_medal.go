package http

import (
	mdMdl "go-common/app/interface/main/creative/model/medal"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"

	"strconv"
)

func appMedalStatus(c *bm.Context) {
	// check user
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	st, err := mdSvc.Medal(c, mid)
	if err == ecode.Int(510002) {
		st = &mdMdl.Medal{
			UID:          strconv.FormatInt(mid, 10),
			LiveStatus:   "0",
			MasterStatus: "0",
			Status:       "0",
		}
		err = nil
	}
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(st, nil)
}

func appMedalCheck(c *bm.Context) {
	req := c.Request
	params := req.Form
	name := params.Get("medal_name")
	// check user
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	if name == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	valid, err := mdSvc.CheckMedal(c, mid, name)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"valid": valid,
	}, nil)
}

func appMedalOpen(c *bm.Context) {
	req := c.Request
	params := req.Form
	name := params.Get("medal_name")
	// check user
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	if name == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err := mdSvc.OpenMedal(c, mid, name)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func appMedalRename(c *bm.Context) {
	req := c.Request
	params := c.Request.Form
	name := params.Get("medal_name")
	cookie := req.Header.Get("cookie")
	ak := params.Get("access_key")
	// check user
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	if name == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	code := mdSvc.Rename(c, mid, name, ak, cookie)
	c.JSON(nil, code)
}
