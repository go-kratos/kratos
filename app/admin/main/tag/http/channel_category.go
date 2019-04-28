package http

import (
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func categoryList(c *bm.Context) {
	c.JSON(svc.ChannelCategories(c))
}

func categoryAdd(c *bm.Context) {
	var (
		err   error
		param = new(struct {
			Name      string `form:"name" validate:"required"`
			INTShield int32  `form:"int_shield" validate:"gte=0,lte=1"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if param.Name, err = svc.CheckChannelCategory(param.Name); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, svc.AddChannelCategory(c, param.Name, param.INTShield))
}

func categoryDelete(c *bm.Context) {
	var (
		err   error
		param = new(struct {
			ID   int64  `form:"id"`
			Name string `form:"name"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if param.ID <= 0 && param.Name == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if param.Name != "" {
		if param.Name, err = svc.CheckChannelCategory(param.Name); err != nil {
			c.JSON(nil, err)
			return
		}
	}
	c.JSON(nil, svc.DeleteChannelCategory(c, param.ID, param.Name))
}

func categorySort(c *bm.Context) {
	var (
		err   error
		param = new(struct {
			IDs []int64 `form:"ids,split" validate:"required,min=1,dive,gt=0"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	c.JSON(nil, svc.SortChannelCategory(c, param.IDs))
}

func categoryShieldINT(c *bm.Context) {
	var (
		err   error
		param = new(struct {
			ID    int64 `form:"id" validate:"required,gt=0"`
			State int32 `form:"state" validate:"gte=0,lte=1"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	_, uname := managerInfo(c)
	c.JSON(nil, svc.CategoryShieldINT(c, param.ID, param.State, uname))
}
