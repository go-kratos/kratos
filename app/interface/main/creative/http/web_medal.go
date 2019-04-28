package http

import (
	"strconv"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func webMedalStatus(c *bm.Context) {
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	st, err := mdSvc.Medal(c, mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(st, nil)
}

func webMedalOpen(c *bm.Context) {
	params := c.Request.Form
	name := params.Get("medal_name")
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	if name == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mid, _ := midI.(int64)
	err := mdSvc.OpenMedal(c, mid, name)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func webRecentFans(c *bm.Context) {
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	fans, err := mdSvc.RecentFans(c, mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(fans, nil)
}

func webMedalCheck(c *bm.Context) {
	params := c.Request.Form
	name := params.Get("medal_name")
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	if name == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mid, _ := midI.(int64)
	valid, err := mdSvc.CheckMedal(c, mid, name)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"valid": valid,
	}, nil)
}

func webMedalRank(c *bm.Context) {
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	ranks, err := mdSvc.Rank(c, mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(ranks, nil)
}

func webMedalRename(c *bm.Context) {
	req := c.Request
	params := c.Request.Form
	name := params.Get("medal_name")
	cookie := req.Header.Get("cookie")
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	if name == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mid, _ := midI.(int64)
	code := mdSvc.Rename(c, mid, name, "", cookie)
	c.JSON(nil, code)
}

func webFansMedal(c *bm.Context) {
	params := c.Request.Form
	tmidStr := params.Get("tmid")
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	// check params
	var (
		tmid int64
		err  error
	)
	if tmidStr != "" {
		tmid, err = strconv.ParseInt(tmidStr, 10, 64)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	mid, _ := midI.(int64)
	if tmid > 0 && dataSvc.IsWhite(mid) {
		mid = tmid
	}
	data, err := dataSvc.UpFansMedal(c, mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}
