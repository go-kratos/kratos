package http

import (
	"go-common/app/admin/main/vip/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// tips info.
func tips(c *bm.Context) {
	var (
		res []*model.Tips
		err error
	)
	arg := new(struct {
		State    int8 `form:"state"`
		Platform int8 `form:"platform"`
		Position int8 `form:"position"`
	})
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	if res, err = vipSvc.TipList(c, arg.Platform, arg.State, arg.Position); err != nil {
		log.Error("vipSvc.TipList(%+v) err(%+v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(res, nil)
}

// tips info.
func tipbyid(c *bm.Context) {
	var (
		res *model.Tips
		err error
	)
	arg := new(struct {
		ID int64 `form:"id" validate:"required,min=1,gte=1"`
	})
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	if res, err = vipSvc.TipByID(c, arg.ID); err != nil {
		log.Error("vipSvc.TipByID(%d) err(%+v)", arg.ID, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(res, nil)
}

// tips add.
func tipadd(c *bm.Context) {
	var (
		err error
		arg = new(model.Tips)
	)
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	opI, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	arg.Operator = opI.(string)
	if err = vipSvc.AddTip(c, arg); err != nil {
		log.Error("vipSvc.AddTip(%v) err(%+v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// tips update.
func tipupdate(c *bm.Context) {
	var (
		err error
		arg = new(model.Tips)
	)
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	opI, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	arg.Operator = opI.(string)
	if err = vipSvc.TipUpdate(c, arg); err != nil {
		log.Error("vipSvc.TipUpdate(%v) err(%+v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// tips delete.
func tipdelete(c *bm.Context) {
	var (
		err error
	)
	arg := new(struct {
		ID int64 `form:"id" validate:"required,min=1,gte=1"`
	})
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	opI, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	if err = vipSvc.DeleteTip(c, arg.ID, opI.(string)); err != nil {
		log.Error("vipSvc.DeleteTip(%d) err(%+v)", arg.ID, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// tips expire.
func tipexpire(c *bm.Context) {
	var (
		err error
	)
	arg := new(struct {
		ID int64 `form:"id" validate:"required,min=1,gte=1"`
	})
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	opI, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	if err = vipSvc.ExpireTip(c, arg.ID, opI.(string)); err != nil {
		log.Error("vipSvc.ExpireTip(%d) err(%+v)", arg.ID, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}
