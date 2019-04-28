package http

import (
	"mime/multipart"

	"go-common/app/admin/main/coupon/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func viewBatchAdd(c *bm.Context) {
	var err error
	arg := new(model.ArgCouponViewBatch)
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	operator, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	arg.Operator = operator.(string)
	c.JSON(nil, svc.CouponViewBatchAdd(c, arg))
}

func viewBatchSave(c *bm.Context) {
	var err error
	arg := new(model.ArgCouponViewBatch)
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	operator, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	arg.Operator = operator.(string)
	c.JSON(nil, svc.CouponViewbatchSave(c, arg))
}

func viewBlock(c *bm.Context) {
	var err error
	arg := new(struct {
		Mid         int64  `form:"mid" validate:"required"`
		CouponToken string `form:"coupon_token" validate:"required"`
	})
	if err = c.Bind(arg); err != nil {
		log.Error("c.bind err(%+v)", err)
		return
	}
	_, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	c.JSON(nil, svc.CouponViewBlock(c, arg.Mid, arg.CouponToken))
}

func viewUnblock(c *bm.Context) {
	var err error
	arg := new(struct {
		Mid         int64  `form:"mid" validate:"required"`
		CouponToken string `form:"coupon_token" validate:"required"`
	})
	if err = c.Bind(arg); err != nil {
		log.Error("c.bind err(%+v)", err)
		return
	}
	_, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	c.JSON(nil, svc.CouponViewUnblock(c, arg.Mid, arg.CouponToken))
}

func viewList(c *bm.Context) {
	var err error
	arg := new(model.ArgSearchCouponView)
	if err = c.Bind(arg); err != nil {
		log.Error("c.bind err(%+v)", err)
		return
	}
	res, count, err := svc.CouponViewList(c, arg)
	info := new(model.PageInfo)
	info.CurrentPage = arg.PN
	info.Count = int(count)
	info.Item = res
	c.JSON(info, err)
}

func salaryView(c *bm.Context) {
	var (
		f   multipart.File
		h   *multipart.FileHeader
		err error
	)
	arg := new(model.ArgAllowanceSalary)
	if err = c.BindWith(arg, binding.FormMultipart); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(arg.Mids) <= 0 {
		f, h, err = c.Request.FormFile("file")
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	c.JSON(svc.CouponViewSalary(c, f, h, arg.Mids, arg.BatchToken))
}
