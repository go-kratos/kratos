package http

import (
	"go-common/app/admin/main/tv/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func addTime(c *bm.Context) {
	tm := new(model.UgcTime)
	if err := c.Bind(tm); err != nil {
		return
	}
	c.JSON(nil, tvSrv.AddUgcTm(tm))
}

func editTime(c *bm.Context) {
	tm := new(model.EditUgcTime)
	if err := c.Bind(tm); err != nil {
		return
	}
	c.JSON(nil, tvSrv.EditUgcTm(tm))
}

func actLabels(c *bm.Context) {
	param := new(struct {
		IDs    []int64 `form:"ids,split" validate:"required,min=1,dive,gt=0"`
		Action string  `form:"action" validate:"required"` // 0 = hide, 1 = recover
	})
	if err := c.Bind(param); err != nil {
		return
	}
	c.JSON(nil, tvSrv.ActLabels(param.IDs, atoi(param.Action)))
}

func delTmLabels(c *bm.Context) {
	param := new(struct {
		IDs []int64 `form:"ids,split" validate:"required,min=1,dive,gt=0"`
	})
	if err := c.Bind(param); err != nil {
		return
	}
	c.JSON(nil, tvSrv.DelLabels(param.IDs))
}

func ugcLabels(c *bm.Context) {
	req := new(model.ReqLabel)
	if err := c.Bind(req); err != nil {
		return
	}
	if req.Param != model.ParamUgctime && req.Param != model.ParamTypeid {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(tvSrv.PickLabels(req, model.UgcLabel))
}

func pgcLabels(c *bm.Context) {
	req := new(model.ReqLabel)
	if err := c.Bind(req); err != nil {
		return
	}
	c.JSON(tvSrv.PickLabels(req, model.PgcLabel))
}

func pgcLblTps(c *bm.Context) {
	param := new(struct {
		Category int `form:"category" validate:"required,min=1,gt=0"`
	})
	if err := c.Bind(param); err != nil {
		return
	}
	c.JSON(tvSrv.LabelTp(param.Category))
}

func editLabel(c *bm.Context) {
	param := new(struct {
		ID   int64  `form:"id" validate:"required"`
		Name string `form:"name" validate:"required"`
	})
	if err := c.Bind(param); err != nil {
		return
	}
	c.JSON(nil, tvSrv.EditLabel(param.ID, param.Name))
}

func pubLabel(c *bm.Context) {
	param := new(struct {
		IDs []int64 `form:"ids,split" validate:"required,min=1,dive,gt=0"`
	})
	if err := c.Bind(param); err != nil {
		return
	}
	c.JSON(nil, tvSrv.PubLabel(param.IDs))
}
