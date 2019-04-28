package http

import (
	"strconv"
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

//func getUpStatus(c *bm.Context) {
//	params := c.Request.Form
//	midStr := params.Get("mid")
//	mid, err := strconv.ParseInt(midStr, 10, 64)
//	if err != nil {
//		log.Error("strconv.ParseInt mid(%s) error(%v)", midStr, err)
//		c.JSON(nil, ecode.RequestErr)
//		return
//	}
//	if mid <= 0 {
//		log.Error("http.getUpStatus mid(%d) <= 0", mid)
//		c.JSON(nil, ecode.RequestErr)
//		return
//	}
//
//	data, err := svc.GetUpStatus(c, mid)
//	if err != nil {
//		log.Error("svc.GetUpStatus mid(%v), error(%v)", err)
//		c.JSON(nil, err)
//		return
//	}
//	c.JSON(data, nil)
//}

func getUpStatus(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	data, err := svc.GetUpStatus(c, mid, ip)
	if err != nil {
		log.Error("svc.GetUpStatus mid(%v), error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func join(c *bm.Context) {
	params := c.Request.Form
	midStr := params.Get("mid")
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt mid(%s) error(%v)", midStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if mid <= 0 {
		log.Error("http.getUpStatus mid(%d) <= 0", mid)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	accountType, err := strconv.Atoi(params.Get("account_type"))
	if err != nil {
		log.Error("strconv.Atoi account_type(%s) error(%v)", accountType, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	signType, err := strconv.Atoi(params.Get("sign_type"))
	if err != nil {
		log.Error("strconv.Atoi sign_type(%s) error(%v)", signType, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}

	err = svc.JoinAv(c, accountType, mid, signType)
	if err != nil {
		log.Error("svc.AddUp accountType(%d) mid(%d) signType(%d) error(%v)", accountType, mid, signType, err)
	}
	c.JSON(nil, err)
}

func upCharge(c *bm.Context) {
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	data, err := svc.GetUpCharge(c, mid, time.Now())
	if err != nil {
		log.Error("growup svc.GetUpCharge error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func quit(c *bm.Context) {
	params := c.Request.Form
	midStr := params.Get("mid")
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt mid(%s) error(%v)", midStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if mid <= 0 {
		log.Error("svc.quit mid(%d) <= 0", mid)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	reason := params.Get("reason")
	err = svc.Quit(c, mid, reason)
	if err != nil {
		log.Error("svc.Quit mid(%d) error(%v)", mid, err)
	}
	c.JSON(nil, err)
}

func quit1(c *bm.Context) {
	v := new(struct {
		Reason string `form:"reason"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	err := svc.Quit(c, mid, v.Reason)
	if err != nil {
		log.Error("svc.Quit mid(%d) error(%v)", mid, err)
	}
	c.JSON(nil, err)
}

func banner(c *bm.Context) {
	data, err := svc.GetBanner(c)
	if err != nil {
		log.Error("svc.Banner error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func upBill(c *bm.Context) {
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)

	data, err := svc.UpBill(c, mid)
	if err != nil {
		log.Error("svc.UpBill mid(%d) error(%v)", mid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func upYear(c *bm.Context) {
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	data, err := svc.UpYear(c, mid)
	if err != nil {
		log.Error("svc.UpYear mid(%d) error(%v)", mid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}
