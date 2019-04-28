package http

import (
	"go-common/app/admin/ep/merlin/model"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func queryImage(c *bm.Context) {
	c.JSON(svc.QueryImage())
}

func addImage(c *bm.Context) {
	var (
		image = &model.Image{}
		err   error
	)
	if err = c.BindWith(image, binding.JSON); err != nil {
		return
	}
	c.JSON(nil, svc.AddImage(image))
}

func updateImage(c *bm.Context) {
	var (
		image = &model.Image{}
		err   error
	)
	if err = c.BindWith(image, binding.JSON); err != nil {
		return
	}
	c.JSON(nil, svc.UpdateImage(image))
}

func delImage(c *bm.Context) {
	var (
		image = &model.Image{}
		err   error
	)
	if err = c.BindWith(image, binding.JSON); err != nil {
		return
	}
	c.JSON(nil, svc.DeleteImage(image.ID))
}
