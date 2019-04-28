package http

import (
	"strconv"
	"strings"

	"go-common/app/service/main/account-recovery/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
)

// queryAccount is verify account is exist
func queryAccount(c *bm.Context) {
	params := new(model.QueryInfoReq)
	if err := c.Bind(params); err != nil {
		return
	}
	c.JSON(srv.QueryAccount(c, params))
}

// commitInfo is commit appeal info
func commitInfo(c *bm.Context) {
	uinfo := new(model.UserInfoReq)
	if err := c.Bind(uinfo); err != nil {
		return
	}
	uinfo.LoginAddrs = strings.TrimSpace(uinfo.LoginAddrs)
	uinfo.RegAddr = strings.TrimSpace(uinfo.RegAddr)
	uinfo.Unames = strings.TrimSpace(uinfo.Unames)
	uinfo.Pwds = strings.TrimSpace(uinfo.Pwds)
	uinfo.Phones = strings.TrimSpace(uinfo.Phones)
	uinfo.Emails = strings.TrimSpace(uinfo.Emails)
	uinfo.SafeAnswer = strings.TrimSpace(uinfo.SafeAnswer)
	uinfo.LinkMail = strings.TrimSpace(uinfo.LinkMail)
	uinfo.Captcha = strings.TrimSpace(uinfo.Captcha)
	uinfo.CardID = strings.TrimSpace(uinfo.CardID)

	uinfo.Business = strings.TrimSpace(uinfo.Business)
	bizMap := make(map[string]string)
	req := c.Request.Form
	uinfo.BusinessMap = bizMap
	extraArgs := model.BusinessExtraArgs(uinfo.Business)
	for _, k := range extraArgs {
		bizMap[k] = req.Get(k)
	}

	c.JSON(nil, srv.CommitInfo(c, uinfo))
}

// queryConWithBusiness query with business
func queryConWithBusiness(business string) func(*bm.Context) {
	return func(c *bm.Context) {
		params := new(model.QueryRecoveryInfoReq)
		if err := c.Bind(params); err != nil {
			return
		}
		params.Bussiness = business
		if perms, ok := c.Get(permit.CtxPermissions); ok {
			for _, p := range perms.([]string) {
				if p == "ACCOUNT_RECOVERY_ADVANCED" {
					params.IsAdvanced = true
				}
			}
		}
		req := c.Request
		if inStatus := req.Form.Get("status"); inStatus != "" {
			status, err := strconv.ParseInt(inStatus, 10, 64)
			if err != nil {
				log.Error("Invalid status: %s: %+v", inStatus, err)
				status = 0
			}
			params.Status = &status
		}
		if inGame := req.Form.Get("game"); inGame != "" {
			game, err := strconv.ParseInt(inGame, 10, 64)
			if err != nil {
				log.Error("Invalid game: %s: %+v", inGame, err)
				game = 0
			}
			params.Game = &game
		}
		if params.Size <= 0 {
			params.Size = 50
		}
		c.JSON(srv.QueryCon(c, params))
	}
}

// judge is reject or agree one operation
func judge(c *bm.Context) {
	req := new(model.JudgeReq)
	if err := c.Bind(req); err != nil {
		return
	}
	c.JSON(nil, srv.Judge(c, req))
}

// batchJudge reject or agree more operation
func batchJudge(c *bm.Context) {
	req := new(model.BatchJudgeReq)
	if err := c.Bind(req); err != nil {
		return
	}
	split := strings.Split(req.Rids, ",")
	rids := make([]int64, 0, len(split))
	for _, v := range split {
		rid, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			c.JSON(nil, err)
			return
		}
		rids = append(rids, rid)
	}
	req.RidsAry = rids
	c.JSON(nil, srv.BatchJudge(c, req))
}

// getCaptchaMail
func getCaptchaMail(c *bm.Context) {
	req := new(model.CaptchaMailReq)
	if err := c.Bind(req); err != nil {
		return
	}
	if !strings.Contains(req.LinkMail, "@") {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if req.Mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var err error
	state, err := srv.GetCaptchaMail(c, req)
	c.JSON(map[string]int64{
		"state": state,
	}, err)
}

func parseMid(midStr string) (mid int64) {
	if len(midStr) == 0 {
		return 0
	}
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		return 0
	}
	return
}

func verifyCode(c *bm.Context) {
	var err error
	arg := new(struct {
		Token string `form:"token" validate:"required"`
		Code  string `form:"code" validate:"required"` //验证码
	})
	if err = c.Bind(arg); err != nil {
		return
	}
	if err = srv.Verify(c, arg.Token, arg.Code); err != nil {
		c.JSON(nil, ecode.CreativeGeetestErr)
		return
	}
	c.JSON(nil, err)
}

func webToken(c *bm.Context) {
	c.JSON(srv.WebToken(c))
}

func compareInfo(c *bm.Context) {
	rid := parseMid(c.Request.Form.Get("rid"))
	if rid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, srv.CompareInfo(c, rid))
}

func sendMail(c *bm.Context) {
	req := new(model.SendMailReq)
	if err := c.Bind(req); err != nil {
		return
	}
	if req.RID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, srv.SendMail(c, req))
}

// gameList game list
func gameList(c *bm.Context) {
	mids := c.Request.Form.Get("mids")
	if mids == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(srv.GameList(c, mids))
}
