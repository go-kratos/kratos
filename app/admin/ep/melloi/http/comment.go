package http

import (
	"go-common/app/admin/ep/melloi/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func queryComment(c *bm.Context) {
	comment := model.Comment{}
	if err := c.BindWith(&comment, binding.Form); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.QueryComment(&comment))
}

func addComment(c *bm.Context) {
	comment := model.Comment{}
	if err := c.BindWith(&comment, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, srv.AddComment(&comment))
}

func updateComment(c *bm.Context) {
	comment := model.Comment{}
	if err := c.BindWith(&comment, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, srv.UpdateComment(&comment))
}

func deleteComment(c *bm.Context) {
	v := new(struct {
		ID int64 `form:"id"`
	})
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, srv.DeleteComment(v.ID))
}
