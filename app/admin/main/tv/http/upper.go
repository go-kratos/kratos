package http

import (
	"go-common/app/admin/main/tv/model"
	bm "go-common/library/net/http/blademaster"
)

func upAdd(c *bm.Context) {
	var (
		err error
	)
	v := new(struct {
		MIDs []int64 `form:"mids,split" validate:"required,min=1,dive,gt=0"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	c.JSON(tvSrv.AddMids(v.MIDs))
}

func upImport(c *bm.Context) {
	var (
		err error
	)
	v := new(struct {
		MIDs []int64 `form:"mids,split" validate:"required,min=1,dive,gt=0"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	c.JSON(tvSrv.ImportMids(v.MIDs))
}

func upDel(c *bm.Context) {
	var (
		err error
	)
	v := new(struct {
		MID int64 `form:"mid" validate:"required,min=1"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, tvSrv.DelMid(v.MID))
}

func upList(c *bm.Context) {
	v := new(struct {
		Order int    `form:"order" validate:"required,min=1,max=4"`
		Pn    int    `form:"pn"`
		Name  string `form:"name"`
		ID    int    `form:"id"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	if v.Pn == 0 {
		v.Pn = 1
	}
	c.JSON(tvSrv.UpList(v.Order, v.Pn, v.Name, v.ID))
}

func upcmsList(c *bm.Context) {
	v := new(model.ReqUpCms)
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(tvSrv.CmsList(c, v))
}

func upcmsAudit(c *bm.Context) {
	v := new(struct {
		Action string  `form:"action"`
		MIDs   []int64 `form:"mids,split" validate:"required,min=1,dive,gt=0"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(tvSrv.CmsAudit(c, v.MIDs, v.Action))
}

func upcmsEdit(c *bm.Context) {
	v := new(model.ReqUpEdit)
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, tvSrv.CmsEdit(c, v))
}
