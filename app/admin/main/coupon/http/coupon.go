package http

import (
	"go-common/app/admin/main/coupon/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// batchadd add.
func batchadd(c *bm.Context) {
	var (
		err error
	)
	arg := new(model.ArgBatchInfo)
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	operator, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	b := new(model.CouponBatchInfo)
	b.AppID = arg.AppID
	b.Name = arg.Name
	if arg.MaxCount == 0 {
		b.MaxCount = -1
	} else {
		b.MaxCount = arg.MaxCount
	}
	if arg.LimitCount == 0 {
		b.LimitCount = -1
	} else {
		b.LimitCount = arg.LimitCount
	}
	b.StartTime = arg.StartTime
	b.ExpireTime = arg.ExpireTime
	b.Operator = operator.(string)
	if err = svc.AddBatchInfo(c, b); err != nil {
		log.Error("svc.AddBatchInfo(%v) err(%+v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func batchlist(c *bm.Context) {
	var (
		err error
		res []*model.CouponBatchResp
	)
	arg := new(model.ArgBatchList)
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	if res, err = svc.BatchList(c, arg); err != nil {
		log.Error("svc.BatchList(%v) err(%+v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(res, nil)
}

func allAppInfo(c *bm.Context) {
	c.JSON(svc.AllAppInfo(c), nil)
}

func salaryCoupon(c *bm.Context) {
	var err error
	arg := new(model.ArgSalaryCoupon)
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	if err = svc.SalaryCoupon(c, arg.Mid, arg.CouponType, arg.Count, arg.BranchToken); err != nil {
		log.Error("svc.SalaryCoupon(%v) err(%+v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func batchBlock(c *bm.Context) {
	var err error
	arg := new(model.ArgAllowance)
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	operator, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	c.JSON(nil, svc.UpdateBatchStatus(c, model.BatchStateBlock, operator.(string), arg.ID))
}

// allowanceUnBlock .
func batchUnBlock(c *bm.Context) {
	var err error
	arg := new(model.ArgAllowance)
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	operator, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	c.JSON(nil, svc.UpdateBatchStatus(c, model.BatchStateNormal, operator.(string), arg.ID))
}
