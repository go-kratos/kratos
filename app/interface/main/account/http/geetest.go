package http

import (
	"strings"

	"go-common/app/interface/main/account/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// getChallenge get geetest params gt ,challenge
func getChallenge(c *bm.Context) {
	params := new(model.GeeCaptchaRequest)
	var (
		mid, ok = c.Get("mid")
		mobile  = strings.Contains(c.Request.UserAgent(), model.MobileUserAgentFlag)
	)
	if !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if mobile {
		params.ClientType = model.PlatH5
	} else {
		params.ClientType = model.PlatPC
	}
	params.MID = mid.(int64)
	c.JSON(geetestSvr.PreProcess(c, params))
}

func geetestValidate(c *bm.Context) {
	params := new(model.GeeCheckRequest)
	if err := c.Bind(params); err != nil {
		return
	}
	mid, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	params.MID = mid.(int64)
	if params.MID == 0 {
		c.JSON(nil, ecode.RequestErr)
	}
	if strings.Contains(c.Request.UserAgent(), model.MobileUserAgentFlag) {
		params.ClientType = model.PlatH5
	} else {
		params.ClientType = model.PlatPC
	}
	c.JSON(geetestSvr.Validate(c, params), nil)
}
