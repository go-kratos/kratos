package http

import (
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"strconv"
)

func appArticleList(c *bm.Context) {
	params := c.Request.Form
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	sortStr := params.Get("sort")
	groupStr := params.Get("group")
	categoryStr := params.Get("category")
	ip := metadata.String(c, metadata.RemoteIP)
	// check user
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	pn, err := strconv.Atoi(pnStr)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.Atoi(psStr)
	if err != nil || ps <= 10 {
		ps = 20
	}
	sort, err := strconv.Atoi(sortStr)
	if err != nil || sort < 0 {
		sort = 0
	}
	group, err := strconv.Atoi(groupStr)
	if err != nil || group < 0 {
		group = 0
	}
	category, err := strconv.Atoi(categoryStr)
	if err != nil || category < 0 {
		category = 0
	}
	arts, err := artSvc.Articles(c, mid, int(pn), int(ps), sort, group, category, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(arts, nil)
}
