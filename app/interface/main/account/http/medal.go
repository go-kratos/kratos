package http

import (
	"strconv"

	"go-common/app/service/main/usersuit/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func medalUserInfo(c *bm.Context) {
	var (
		err    error
		mid, _ = c.Get("mid")
		//ip     = c.RemoteIP()
		data *model.MedalUserInfo
	)
	if data, err = memberSvc.MedalUserInfo(c, mid.(int64)); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func medalInstall(c *bm.Context) {
	var (
		err              error
		nid, isActivated int64
		mid, _           = c.Get("mid")
	)
	nidStr := c.Request.Form.Get("nid")
	if nid, err = strconv.ParseInt(nidStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	isActivatedStr := c.Request.Form.Get("isActivated")
	if isActivated, err = strconv.ParseInt(isActivatedStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = memberSvc.MedalInstall(c, mid.(int64), nid, int8(isActivated))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func medalPopup(c *bm.Context) {
	var (
		err    error
		mid, _ = c.Get("mid")
		data   *model.MedalPopup
	)
	if data, err = memberSvc.MedalPopup(c, mid.(int64)); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func medalMyInfo(c *bm.Context) {
	var (
		mid, _ = c.Get("mid")
	)
	data, err := memberSvc.MedalMyInfo(c, mid.(int64))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func medalAllInfo(c *bm.Context) {
	var (
		mid, _ = c.Get("mid")
	)
	data, err := memberSvc.MedalAllInfo(c, mid.(int64))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}
