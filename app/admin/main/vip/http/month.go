package http

import (
	"go-common/app/admin/main/vip/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func monthList(c *bm.Context) {
	var (
		res []*model.VipMonth
		err error
	)
	if res, err = vipSvc.MonthList(c); err != nil {
		c.JSON(nil, err)
		return
	}
	page := &model.PageInfo{Count: len(res), Item: res}
	c.JSON(page, nil)
}

func monthEdit(c *bm.Context) {
	var (
		err error
		arg = new(model.ArgIDExtra)
	)
	username, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	if err = c.Bind(arg); err != nil {
		return
	}
	arg.Operator = username.(string)
	c.JSON(nil, vipSvc.MonthEdit(c, arg.ID, arg.Status, arg.Operator))
}

func priceList(c *bm.Context) {
	var (
		res []*model.VipMonthPrice
		err error
		arg = new(model.ArgID)
	)
	if err = c.Bind(arg); err != nil {
		return
	}
	if res, err = vipSvc.PriceList(c, arg.ID); err != nil {
		c.JSON(nil, err)
		return
	}
	page := &model.PageInfo{Count: len(res), Item: res}
	c.JSON(page, nil)
}

func priceAdd(c *bm.Context) {
	var (
		err error
		mp  = new(model.VipMonthPrice)
	)
	username, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	mp.Operator = username.(string)
	if err = c.Bind(mp); err != nil {
		return
	}
	c.JSON(nil, vipSvc.PriceAdd(c, mp))
}

func priceEdit(c *bm.Context) {
	var (
		err error
		mp  = new(model.VipMonthPrice)
	)
	username, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	mp.Operator = username.(string)
	if err = c.Bind(mp); err != nil {
		return
	}
	if mp.ID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, vipSvc.PriceEdit(c, mp))
}
