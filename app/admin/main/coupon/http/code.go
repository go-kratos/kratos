package http

import (
	"fmt"

	"go-common/app/admin/main/coupon/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

func exportCode(c *bm.Context) {
	var (
		err   error
		codes []string
	)
	arg := new(model.ArgCouponCode)
	if err = c.Bind(arg); err != nil {
		return
	}
	if codes, err = svc.ExportCode(c, arg); err != nil {
		c.JSON(nil, err)
		return
	}
	writer := c.Writer
	header := writer.Header()
	header.Add("Content-disposition", "attachment; filename="+fmt.Sprintf("%v", arg.BatchToken)+".txt")
	header.Add("Content-Type", "application/x-download;charset=utf-8")
	for _, v := range codes {
		writer.Write([]byte(fmt.Sprintf("%v\r\n", v)))
	}
}

func codePage(c *bm.Context) {
	var err error
	arg := new(model.ArgCouponCode)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(svc.CodePage(c, arg))
}

func codeBlock(c *bm.Context) {
	var err error
	arg := new(model.ArgCouponCode)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, svc.CodeBlock(c, arg))
}

func codeUnBlock(c *bm.Context) {
	var err error
	arg := new(model.ArgCouponCode)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, svc.CodeUnBlock(c, arg))
}

func codeAddBatch(c *bm.Context) {
	var (
		err   error
		token string
	)
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
	if arg.MaxCount == 0 || arg.MaxCount > model.BatchCodeMaxCount {
		c.JSON(nil, ecode.CouponCodeMaxLimitErr)
		return
	}
	b.MaxCount = arg.MaxCount
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
	b.CouponType = model.CouponAllowanceCode
	if token, err = svc.AddAllowanceBatchInfo(c, b); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, svc.InitCodes(c, token))
}

func codeBatchModify(c *bm.Context) {
	var err error
	arg := new(model.ArgAllowanceBatchInfoModify)
	if err = c.Bind(arg); err != nil {
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
	c.JSON(nil, svc.UpdateCodeBatchInfo(c, b))
}

func codeBatchList(c *bm.Context) {
	var err error
	arg := new(model.ArgBatchList)
	if err = c.Bind(arg); err != nil {
		return
	}
	arg.Type = model.CouponAllowanceCode
	c.JSON(svc.BatchList(c, arg))
}

func codeBatchBlock(c *bm.Context) {
	var err error
	arg := new(model.ArgAllowance)
	if err = c.Bind(arg); err != nil {
		return
	}
	operator, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	c.JSON(nil, svc.UpdateBatchStatus(c, model.BatchStateBlock, operator.(string), arg.ID))
}

func codeBatchUnBlock(c *bm.Context) {
	var err error
	arg := new(model.ArgAllowance)
	if err = c.Bind(arg); err != nil {
		return
	}
	operator, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	c.JSON(nil, svc.UpdateBatchStatus(c, model.BatchStateNormal, operator.(string), arg.ID))
}
