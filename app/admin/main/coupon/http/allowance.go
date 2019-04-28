package http

import (
	"bytes"
	"io"
	"mime/multipart"

	"go-common/app/admin/main/coupon/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
	"go-common/library/xstr"
)

// batchadd add.
func allowanceBatchadd(c *bm.Context) {
	var err error
	arg := new(model.ArgAllowanceBatchInfo)
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	operator, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	if _, ok := model.ProdLimMonthMap[arg.ProdLimMonth]; !ok {
		c.JSON(nil, ecode.RequestErr)
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
	b.Amount = arg.Amount
	b.FullAmount = arg.FullAmount
	b.ExpireDay = arg.ExpireDay
	b.PlatformLimit = xstr.JoinInts(arg.PlatformLimit)
	b.ProdLimMonth = arg.ProdLimMonth
	b.ProdLimRenewal = arg.ProdLimRenewal
	b.CouponType = model.CouponAllowance
	_, err = svc.AddAllowanceBatchInfo(c, b)
	c.JSON(nil, err)
}

// batchadd allowance modify.
func allowanceBatchModify(c *bm.Context) {
	var err error
	arg := new(model.ArgAllowanceBatchInfoModify)
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
	b.Operator = operator.(string)
	b.PlatformLimit = xstr.JoinInts(arg.PlatformLimit)
	b.ProdLimMonth = arg.ProdLimMonth
	b.ProdLimRenewal = arg.ProdLimRenewal
	b.ID = arg.ID
	c.JSON(nil, svc.UpdateAllowanceBatchInfo(c, b))
}

// allowanceBlock .
func allowanceBatchBlock(c *bm.Context) {
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
func allowanceBatchUnBlock(c *bm.Context) {
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

func allowanceSalary(c *bm.Context) {
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
	c.JSON(svc.AllowanceSalary(c, f, h, arg.Mids, arg.BatchToken, arg.MsgType))
}

func batchInfo(c *bm.Context) {
	var (
		err error
	)
	arg := new(model.ArgAllowanceInfo)
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	_, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	c.JSON(svc.BatchInfo(c, arg.BatchToken))
}

// func allowancePage(c *bm.Context) {
// 	var err error
// 	arg := new(model.ArgAllowanceSearch)
// 	if err = c.Bind(arg); err != nil {
// 		log.Error("c.Bind err(%+v)", err)
// 		return
// 	}
// 	_, ok := c.Get("username")
// 	if !ok {
// 		c.JSON(nil, ecode.AccessDenied)
// 		return
// 	}
// 	c.JSON(svc.AllowancePage(c, arg))
// }

func allowanceList(c *bm.Context) {
	var err error
	arg := new(model.ArgAllowanceSearch)
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	_, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	c.JSON(svc.AllowanceList(c, arg))
}

func allowanceBlock(c *bm.Context) {
	var (
		err error
	)
	arg := new(model.ArgAllowanceState)
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	_, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	c.JSON(nil, svc.UpdateAllowanceState(c, arg.Mid, model.Block, arg.CouponToken))
}

func allowanceUnBlock(c *bm.Context) {
	var err error
	arg := new(model.ArgAllowanceState)
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	_, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	c.JSON(nil, svc.UpdateAllowanceState(c, arg.Mid, model.NotUsed, arg.CouponToken))
}

func batchSalaryCoupon(c *bm.Context) {
	var err error
	req := new(model.ArgBatchSalaryCoupon)
	if err = c.Bind(req); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	_, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	c.JSON(nil, svc.ActivitySalaryCoupon(c, req))
}

func uploadFile(c *bm.Context) {
	var (
		f   multipart.File
		err error
	)
	_, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	arg := new(model.ArgUploadFile)
	if err = c.BindWith(arg, binding.FormMultipart); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	f, _, err = c.Request.FormFile("file")
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	defer f.Close()
	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, f); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, svc.OutFile(c, buf.Bytes(), arg.FileURL))
}
