package http

import (
	"strconv"

	"go-common/app/interface/main/account/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// replyHistoryList
func replyHistoryList(c *bm.Context) {
	var (
		err error
		//ip        = c.RemoteIP()
		header    = c.Request.Header
		params    = c.Request.Form
		mid, _    = c.Get("mid")
		cookie    = header.Get("Cookie")
		accessKey = params.Get("access_key")
		stime     = params.Get("stime")
		etime     = params.Get("etime")
		order     = params.Get("order")
		sort      = params.Get("sort")
		pnStr     = params.Get("pn")
		psStr     = params.Get("ps")
		pn, ps    int64
	)
	if pnStr != "" {
		if pn, err = strconv.ParseInt(pnStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if psStr != "" {
		if ps, err = strconv.ParseInt(psStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	c.JSON(memberSvc.ReplyHistoryList(c, mid.(int64), stime, etime, order, sort, pn, ps, accessKey, cookie))
}

// updateSettings
func update(c *bm.Context) {
	var (
		params  = c.Request.Form
		mid, ok = c.Get("mid")
		//ip          = c.RemoteIP()
		unameStr    = params.Get("uname")
		signStr     = params.Get("usersign")
		sexStr      = params.Get("sex")
		birthdayStr = params.Get("birthday")
	)
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	settings := &model.Settings{
		Uname:    unameStr,
		Sign:     signStr,
		Sex:      sexStr,
		Birthday: birthdayStr,
	}
	log.Error("request(%v)", settings)
	c.JSON(nil, memberSvc.UpdateSettings(c, mid.(int64), settings))
}

func account(c *bm.Context) {
	var (
		//ip      = c.RemoteIP()
		mid, ok = c.Get("mid")
	)
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	c.JSON(memberSvc.SettingsInfo(c, mid.(int64)))
}

// logCoin
func logCoin(c *bm.Context) {
	var (
		//ip      = c.RemoteIP()
		mid, ok = c.Get("mid")
	)
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	c.JSON(memberSvc.LogCoin(c, mid.(int64)))
}

// coin
func coin(c *bm.Context) {
	var (
		//ip      = c.RemoteIP()
		mid, ok = c.Get("mid")
	)
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	c.JSON(memberSvc.Coin(c, mid.(int64)))
}

func logMoral(c *bm.Context) {
	var (
		//ip      = c.RemoteIP()
		mid, ok = c.Get("mid")
	)
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	c.JSON(memberSvc.LogMoral(c, mid.(int64)))
}

func logExp(c *bm.Context) {
	var (
		//ip      = c.RemoteIP()
		mid, ok = c.Get("mid")
	)
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	c.JSON(memberSvc.LogExp(c, mid.(int64)))
}

// logLogin
func logLogin(c *bm.Context) {
	var (
		//ip      = c.RemoteIP()
		mid, ok = c.Get("mid")
	)
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	c.JSON(memberSvc.LogLogin(c, mid.(int64)))
}

func reward(c *bm.Context) {
	var (
		//ip      = c.RemoteIP()
		mid, ok = c.Get("mid")
	)
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	c.JSON(memberSvc.Reward(c, mid.(int64)))
}
