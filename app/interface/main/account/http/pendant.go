package http

import (
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"strconv"
)

func pendantAll(c *bm.Context) {
	mid, ok := c.Get("mid")
	//ip := c.RemoteIP()
	if !ok {
		c.JSON(nil, ecode.AccountNotLogin)
		return
	}
	c.JSON(usSvc.Group(c, mid.(int64)))
}

func pendantMy(c *bm.Context) {
	mid, ok := c.Get("mid")
	//ip := c.RemoteIP()
	if !ok {
		c.JSON(nil, ecode.AccountNotLogin)
		return
	}
	c.JSON(usSvc.My(c, mid.(int64)))
}

func pendantMyHistory(c *bm.Context) {
	//ip := c.RemoteIP()
	mid, ok := c.Get("mid")
	pageStr := c.Request.Form.Get("page")
	if !ok {
		c.JSON(nil, ecode.AccountNotLogin)
		return
	}
	page, _ := strconv.ParseInt(pageStr, 10, 64)
	c.JSON(usSvc.MyHistory(c, mid.(int64), page))
}

func pendantCurrent(c *bm.Context) {
	mid, ok := c.Get("mid")
	//ip := c.RemoteIP()
	if !ok {
		c.JSON(nil, ecode.AccountNotLogin)
		return
	}
	c.JSON(usSvc.Equipment(c, mid.(int64)))
}

func pendantEntry(c *bm.Context) {
	mid, ok := c.Get("mid")
	//ip := c.RemoteIP()
	if !ok {
		c.JSON(nil, ecode.AccountNotLogin)
		return
	}
	c.JSON(usSvc.GroupEntry(c, mid.(int64)))
}

func pendantSingle(c *bm.Context) {
	pidStr := c.Request.Form.Get("pid")
	//ip := c.RemoteIP()
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(usSvc.Pendant(c, pid))
}

func pendantVIP(c *bm.Context) {
	mid, ok := c.Get("mid")
	//ip := c.RemoteIP()
	if !ok {
		c.JSON(nil, ecode.AccountNotLogin)
		return
	}
	c.JSON(usSvc.GroupVIP(c, mid.(int64)))
}

func pendantCheckOrder(c *bm.Context) {
	mid, ok := c.Get("mid")
	//ip := c.RemoteIP()
	if !ok {
		c.JSON(nil, ecode.AccountNotLogin)
		return
	}
	orderID := c.Request.Form.Get("orderId")
	c.JSON(nil, usSvc.CheckOrder(c, mid.(int64), orderID))
}

func pendantVIPGet(c *bm.Context) {
	params := c.Request.Form
	mid, ok := c.Get("mid")
	//ip := c.RemoteIP()
	if !ok {
		c.JSON(nil, ecode.AccountNotLogin)
		return
	}
	pidStr := params.Get("pid")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if pid == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	activatedStr := params.Get("isActivated")
	activated, err := strconv.Atoi(activatedStr)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if activated == 0 {
		activated = 1
	} else {
		activated = 2
	}
	c.JSON(nil, usSvc.VipGet(c, mid.(int64), pid, int8(activated)))
}

func pendantOrder(c *bm.Context) {
	params := c.Request.Form
	mid, ok := c.Get("mid")
	//ip := c.RemoteIP()
	if !ok {
		c.JSON(nil, ecode.AccountNotLogin)
		return
	}
	pidStr := params.Get("pid")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	timeLengthStr := params.Get("timeLength")
	timeLength, err := strconv.ParseInt(timeLengthStr, 10, 64)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if pid <= 0 || timeLength <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var (
		moneyType    int8
		moneyTypeStr string
	)
	moneyTypeStr = params.Get("moneyType")
	switch moneyTypeStr {
	case "coin":
		moneyType = 0
	case "bcoin":
		moneyType = 1
	case "point":
		moneyType = 2
	default:
		c.JSON(nil, ecode.PendantPayTypeErr)
		return
	}
	c.JSON(usSvc.Order(c, mid.(int64), pid, timeLength, moneyType))
}
