package http

import (
	"go-common/app/admin/main/upload/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

// ping check server ok.
func ping(c *bm.Context) {
	if err := uaSvc.Ping(c); err != nil {
		c.Error = err
		c.AbortWithStatus(503)
	}
}

func add(c *bm.Context) {
	var err error
	ap := &model.AddParam{}
	if err = c.BindWith(ap, binding.FormPost); err != nil {
		return
	}

	c.JSON(nil, uaSvc.Add(c, ap))
}

func list(c *bm.Context) {
	var (
		err error
	)
	lp := &model.ListParam{}
	if err = c.Bind(lp); err != nil {
		return
	}
	c.JSON(uaSvc.List(c, lp))
}

func deleteFile(c *bm.Context) {
	var (
		err     error
		ok      bool
		adminID interface{}
	)
	dp := new(model.DeleteParam)
	if adminID, ok = c.Get("uid"); !ok {
		c.JSON(nil, ecode.UserNotExist)
		return
	}
	dp.AdminID = adminID.(int64)
	if err = c.Bind(dp); err != nil {
		return
	}
	c.JSON(nil, uaSvc.Delete(c, dp))
}

func deleteRawFile(c *bm.Context) {
	var (
		err error
	)
	dp := new(model.DeleteRawParam)
	if err = c.Bind(dp); err != nil {
		return
	}
	c.JSON(nil, uaSvc.DeleteRaw(c, dp))
}

func deleteFileV2(c *bm.Context) {
	var (
		adminID interface{}
		err     error
		ok      bool
	)
	dp := new(model.DeleteV2Param)
	if adminID, ok = c.Get("uid"); !ok {
		c.JSON(nil, ecode.UserNotExist)
		return
	}
	dp.AdminID = adminID.(int64)
	if err = c.Bind(dp); err != nil {
		return
	}
	if err = uaSvc.DeleteV2(c, dp); err != nil {
		log.Error("deleteFileV2 error(%v)", err)
	}
	c.JSON(nil, err)
}

func multiList(c *bm.Context) {
	var (
		err error
	)
	lp := &model.MultiListParam{}
	if err = c.Bind(lp); err != nil {
		return
	}
	c.JSON(uaSvc.MultiList(c, lp))
}
