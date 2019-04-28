package http

import (
	"strconv"

	"go-common/app/admin/main/vip/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func drawback(c *bm.Context) {
	var (
		mid      int64
		params   = c.Request.Form
		err      error
		remark   = params.Get("remark")
		username string
		days     int
	)
	if nameInter, ok := c.Get("username"); ok {
		username = nameInter.(string)
	}
	midStr := params.Get("mid")
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	daysStr := params.Get("days")
	if days, err = strconv.Atoi(daysStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = vipSvc.Drawback(c, days, mid, username, remark); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func historyList(c *bm.Context) {
	var (
		err   error
		count int
		res   []*model.VipChangeHistory
		req   = new(model.UserChangeHistoryReq)
	)

	if err = c.Bind(req); err != nil {
		return
	}
	if res, count, err = vipSvc.HistoryPage(c, req); err != nil {
		c.JSON(nil, err)
		return
	}
	page := &model.PageInfo{Count: count, Item: res}
	c.JSON(page, nil)
}

func vipInfo(c *bm.Context) {
	arg := new(struct {
		Mid int64 `form:"mid"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(vipSvc.VipInfo(c, arg.Mid))
}
