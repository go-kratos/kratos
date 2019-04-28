package http

import (
	bm "go-common/library/net/http/blademaster"
)

// tagList get enabled tag list by business
func tagList(c *bm.Context) {
	tp := new(struct {
		Business int8 `form:"business" validate:"required"`
	})
	if err := c.Bind(tp); err != nil {
		return
	}
	c.JSON(wkfSvc.TagSlice(tp.Business), nil)
}

// tagList3 .
func tagList3(c *bm.Context) {
	tp := new(struct {
		BID int64 `form:"bid" validate:"required"`
		RID int64 `form:"rid"`
	})
	if err := c.Bind(tp); err != nil {
		return
	}
	c.JSON(wkfSvc.Tag3(tp.BID, tp.RID), nil)
}
