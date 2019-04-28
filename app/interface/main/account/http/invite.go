package http

import (
	"strconv"

	"go-common/app/interface/main/account/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// stat get invite stat.
func inviteStat(c *bm.Context) {
	mid, _ := c.Get("mid")
	var err error
	var stat *model.RichInviteStat
	if stat, err = usSvc.Stat(c, mid.(int64)); err != nil {
		log.Error("memberService.Stat(%d) error(%v)", mid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(stat, nil)
}

// buy buy invite code.
func buy(c *bm.Context) {
	mid, _ := c.Get("mid")
	var err error
	var num int64
	req := c.Request.Form
	numStr := req.Get("num")
	if num, err = strconv.ParseInt(numStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	geeParam := new(model.GeeCheckRequest)
	if err = c.Bind(geeParam); err != nil {
		return
	}
	geeParam.MID = mid.(int64)
	if isPass := geetestSvr.Validate(c, geeParam); !isPass {
		c.JSON(nil, ecode.CaptchaErr)
		return
	}
	var invs []*model.RichInvite
	if invs, err = usSvc.Buy(c, mid.(int64), num); err != nil {
		log.Error("memberService.Buy(%d, %d) error(%v)", mid, num, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(invs, nil)
}

// apply apply invite code.
func apply(c *bm.Context) {
	mid, _ := c.Get("mid")
	var err error
	var num int64
	code := c.Request.Form.Get("invite_code")
	if code == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = usSvc.Apply(c, mid.(int64), code, c.Request.Header.Get("Cookie")); err != nil {
		log.Error("memberService.Apply(%d, %d) error(%v)", mid, num, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}
