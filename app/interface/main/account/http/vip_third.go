package http

import (
	"net/http"

	"go-common/app/interface/main/account/model"
	idtv1 "go-common/app/service/main/identify/api/grpc"
	vipmol "go-common/app/service/main/vip/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/metadata"
)

//
// vip 第三方[ele]接入gateway
//

// openID
func openIDByOAuth2Code(c *bm.Context) {
	var err error
	a := new(model.ArgAuthCode)
	if err = c.Bind(a); err != nil {
		return
	}
	a.IP = metadata.String(c, metadata.RemoteIP)
	a.APPID = vipmol.EleAppID
	c.JSON(vipSvc.OpenIDByAuthCode(c, a))
}

func openBindByOutOpenID(c *bm.Context) {
	var err error
	a := new(model.ArgBind)
	if err = c.Bind(a); err != nil {
		return
	}
	a.AppID = vipmol.EleAppID
	c.JSON(nil, vipSvc.OpenBindByOutOpenID(c, a))
}

func userInfoByOpenID(c *bm.Context) {
	var err error
	a := new(model.ArgUserInfoByOpenID)
	if err = c.Bind(a); err != nil {
		return
	}
	a.AppID = vipmol.EleAppID
	c.JSON(vipSvc.UserInfoByOpenID(c, a))
}

func bilibiliVipGrant(c *bm.Context) {
	var err error
	a := new(model.ArgBilibiliVipGrant)
	if err = c.Bind(a); err != nil {
		return
	}
	a.AppID = vipmol.EleAppID
	c.JSON(nil, vipSvc.BilibiliVipGrant(c, a))
}

func bilibiliPrizeGrant(c *bm.Context) {
	var err error
	a := new(model.ArgBilibiliPrizeGrant)
	if err = c.Bind(a); err != nil {
		return
	}
	a.AppID = vipmol.EleAppID
	c.JSON(vipSvc.BilibiliPrizeGrant(c, a))
}

func openAuthCallBack(c *bm.Context) {
	var err error
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	a := new(model.ArgOpenAuthCallBack)
	if err = c.Bind(a); err != nil {
		return
	}
	// verify csrf.
	verifyState(c, authn, a.State)
	a.AppID = vipmol.EleAppID
	a.Mid = midI.(int64)
	c.Redirect(http.StatusFound, vipSvc.OpenAuthCallBack(c, a))
}

func eleOAuthURL(c *bm.Context) {
	var (
		state string
		err   error
	)
	if state, err = csrf(c, authn); err != nil {
		return
	}
	c.JSON(vipSvc.ElemeOAuthURI(c, state), nil)
}

func verifyState(ctx *bm.Context, a *auth.Auth, state string) (err error) {
	var csrfStr string
	if csrfStr, err = csrf(ctx, a); err != nil {
		return
	}
	if csrfStr != state {
		return ecode.CsrfNotMatchErr
	}
	return
}

func csrf(ctx *bm.Context, a *auth.Auth) (string, error) {
	req := ctx.Request
	cookie := req.Header.Get("Cookie")
	reply, err := a.GetCookieInfo(ctx, &idtv1.GetCookieInfoReq{Cookie: cookie})
	if err != nil {
		return "", err
	}
	if !reply.IsLogin {
		return "", ecode.NoLogin
	}
	return reply.Csrf, nil
}
