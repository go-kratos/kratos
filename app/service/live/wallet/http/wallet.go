package http

import (
	"strconv"

	"go-common/app/service/live/wallet/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func getBasicParam(c *bm.Context) *model.BasicParam {
	bp := new(model.BasicParam)
	var err error
	bp.TransactionId = c.Request.Form.Get("transaction_id")
	bp.BizCode = c.Request.Form.Get("biz_code")
	bp.Area, err = strconv.ParseInt(c.Request.Form.Get("area_id"), 10, 64)
	if err != nil {
		bp.Area = 0
	}
	bp.Source = c.Request.Form.Get("source")
	bp.BizSource = c.Request.Form.Get("biz_source")
	bp.MetaData = c.Request.Form.Get("metadata")
	bp.Reason, err = strconv.ParseInt(c.Request.Form.Get("biz_reason"), 10, 64)
	if err != nil {
		bp.Reason = 0
	}

	bp.Version, err = strconv.ParseInt(c.Request.Form.Get("version"), 10, 64)
	if err != nil {
		bp.Version = 0
	}

	return bp

}

func getWithMetal(c *bm.Context) (withMetal int, err error) {
	withMetalStr := c.Request.Form.Get("with_metal")
	if withMetalStr == "" {
		withMetal = 0
		return
	}
	// check params
	withMetal, err = strconv.Atoi(withMetalStr)
	if err != nil || (withMetal != 0 && withMetal != 1) {
		err = ecode.RequestErr
		return
	}
	return
}

func get(c *bm.Context) {
	uidStr := c.Request.Form.Get("uid")
	// check params
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil || uid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	withMetal, err := getWithMetal(c)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	bp := getBasicParam(c)
	platform := c.Request.Header.Get("platform")
	c.JSON(walletSvr.Get(c, bp, uid, platform, withMetal))
}

func delCache(c *bm.Context) {
	uidStr := c.Request.Form.Get("uid")
	// check params
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil || uid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	bp := getBasicParam(c)
	c.JSON(walletSvr.DelCache(c, bp, uid))
}

func getAll(c *bm.Context) {
	uidStr := c.Request.Form.Get("uid")
	// check params
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil || uid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	withMetal, err := getWithMetal(c)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	bp := getBasicParam(c)
	platform := c.Request.Header.Get("platform")
	c.JSON(walletSvr.GetAll(c, bp, uid, platform, withMetal))
}

func getTid(c *bm.Context) {
	typeStr := c.Request.Form.Get("type")
	// check params
	serviceType64, err := strconv.ParseInt(typeStr, 10, 64)
	serviceType := int32(serviceType64)
	if err != nil || serviceType < 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	params := c.Request.Form.Get("params")
	if params == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	bp := getBasicParam(c)
	c.JSON(walletSvr.GetTid(c, bp, 0, serviceType, params))
}

func recharge(c *bm.Context) {
	bp := getBasicParam(c)
	platform := c.Request.Header.Get("platform")

	arg := &model.RechargeOrPayForm{}
	if err := c.Bind(arg); err != nil {
		return
	}

	c.JSON(walletSvr.Recharge(c, bp, arg.Uid, platform, arg))
}

func modify(c *bm.Context) {
	bp := getBasicParam(c)
	platform := c.Request.Header.Get("platform")

	arg := &model.RechargeOrPayForm{}
	if err := c.Bind(arg); err != nil {
		return
	}

	c.JSON(walletSvr.Modify(c, bp, arg.Uid, platform, arg))
}

func pay(c *bm.Context) {
	bp := getBasicParam(c)
	platform := c.Request.Header.Get("platform")

	arg := &model.RechargeOrPayForm{}
	if err := c.Bind(arg); err != nil {
		return
	}

	var reason interface{}
	reasonFromHttp := c.Request.Form.Get("reason")
	if reasonFromHttp == "" {
		reason = nil
	} else {
		reason = reasonFromHttp
	}

	c.JSON(walletSvr.Pay(c, bp, arg.Uid, platform, arg, reason))
}

func exchange(c *bm.Context) {
	bp := getBasicParam(c)
	platform := c.Request.Header.Get("platform")

	arg := &model.ExchangeForm{}
	if err := c.Bind(arg); err != nil {
		return
	}

	c.JSON(walletSvr.Exchange(c, bp, arg.Uid, platform, arg))
}

func query(c *bm.Context) {
	bp := getBasicParam(c)
	platform := c.Request.Header.Get("platform")
	if platform == "" {
		platform = "pc"
	}

	tid := c.Request.Form.Get("transaction_id")
	if tid == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	uidStr := c.Request.Form.Get("uid")
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil || uid <= 0 {
		uid = 0
		return
	}

	c.JSON(walletSvr.Query(c, bp, uid, platform, tid))

}

func recordCoinStream(c *bm.Context) {
	bp := getBasicParam(c)
	platform := c.Request.Header.Get("platform")

	arg := &model.RecordCoinStreamForm{}
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(walletSvr.RecordCoinStream(c, bp, arg.Uid, platform, arg))
}
