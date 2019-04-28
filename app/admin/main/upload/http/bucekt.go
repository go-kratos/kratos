package http

import (
	"net/http"

	"go-common/app/admin/main/upload/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
	"go-common/library/net/http/blademaster/render"
)

func addBucket(c *bm.Context) {
	var err error
	abp := new(model.AddBucketParam)
	if err = c.BindWith(abp, binding.FormPost); err != nil {
		return
	}
	if len(abp.KeyID) != 16 {
		c.Render(http.StatusOK, render.JSON{
			Code:    ecode.RequestErr.Code(),
			Message: "key_id should length 16",
			Data:    nil,
		})
		c.Abort()
		return
	}
	if len(abp.KeySecret) != 30 {
		c.Render(http.StatusOK, render.JSON{
			Code:    ecode.RequestErr.Code(),
			Message: "key_secret should length 30",
			Data:    nil,
		})
		c.Abort()
		return
	}

	c.JSON(nil, uaSvc.AddBucket(c, abp))
}

func listBucket(c *bm.Context) {
	var err error
	lbp := new(model.ListBucketParam)
	if err = c.Bind(lbp); err != nil {
		return
	}
	c.JSON(uaSvc.ListBucket(c, lbp))
}

func listPublicBucket(c *bm.Context) {
	var err error
	lbp := new(model.ListBucketParam)
	if err = c.Bind(lbp); err != nil {
		return
	}
	c.JSON(uaSvc.ListPublicBucket(c, lbp))
}

func detailBucket(c *bm.Context) {
	var err error
	dbp := new(struct {
		Bucket string `json:"bucket" form:"bucket" validate:"required"`
	})
	if err = c.Bind(dbp); err != nil {
		return
	}
	c.JSON(uaSvc.DetailBucket(c, dbp.Bucket))
}
