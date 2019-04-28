package http

import (
	"go-common/app/admin/main/tv/model"
	bm "go-common/library/net/http/blademaster"
)

func mangoList(c *bm.Context) {
	c.JSON(tvSrv.MangoList(c))
}

func mangoAdd(c *bm.Context) {
	param := new(struct {
		IDs   []int64 `form:"rids,split" validate:"required,min=1,dive,gt=0"`
		RType int     `form:"rtype" validate:"required,min=1,max=2"` // 1=pgc, 2=ugc
	})
	if err := c.Bind(param); err != nil {
		return
	}
	c.JSON(tvSrv.MangoAdd(c, param.RType, param.IDs))
}

func mangoEdit(c *bm.Context) {
	param := new(model.ReqMangoEdit)
	if err := c.Bind(param); err != nil {
		return
	}
	c.JSON(nil, tvSrv.MangoEdit(c, param))
}

func mangoDel(c *bm.Context) {
	param := new(struct {
		ID int64 `form:"id" validate:"required,min=1,gt=0"`
	})
	if err := c.Bind(param); err != nil {
		return
	}
	c.JSON(nil, tvSrv.MangoDel(c, param.ID))
}

func mangoPub(c *bm.Context) {
	param := new(struct {
		IDs []int64 `form:"ids,split" validate:"required,min=1,dive,gt=0"`
	})
	if err := c.Bind(param); err != nil {
		return
	}
	c.JSON(nil, tvSrv.MangoPub(c, param.IDs))
}
