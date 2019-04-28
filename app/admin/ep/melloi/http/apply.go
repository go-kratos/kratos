package http

import (
	"go-common/app/admin/ep/melloi/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func queryApply(c *bm.Context) {
	qar := model.QueryApplyRequest{}
	if err := c.BindWith(&qar, binding.Form); err != nil {
		c.JSON(nil, err)
		return
	}
	if err := qar.Verify(); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.QueryApply(&qar))
}

func addApply(c *bm.Context) {
	apply := model.Apply{}
	if err := c.BindWith(&apply, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	cookie := c.Request.Header.Get("Cookie")
	c.JSON(nil, srv.AddApply(c, cookie, &apply))
}

func updateApply(c *bm.Context) {
	apply := model.Apply{}
	if err := c.BindWith(&apply, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	cookie := c.Request.Header.Get("Cookie")
	c.JSON(nil, srv.UpdateApply(cookie, &apply))
}

func delApply(c *bm.Context) {
	v := new(struct {
		ID int64 `form:"id"`
	})
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, srv.DeleteApply(v.ID))
}
