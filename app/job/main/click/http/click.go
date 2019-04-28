package http

import (
	"strconv"

	"go-common/app/job/main/click/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func click(c *bm.Context) {
	var (
		aid, click                    int64
		aidStr, platformStr, clickStr string
		err                           error
	)
	params := c.Request.Form
	if clickStr = params.Get("click"); clickStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if click, err = strconv.ParseInt(clickStr, 10, 64); err != nil || click < 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if aidStr = params.Get("aid"); aidStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if platformStr = params.Get("platform"); platformStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if platformStr != model.TypeForAndroid &&
		platformStr != model.TypeForH5 &&
		platformStr != model.TypeForIOS &&
		platformStr != model.TypeForOutside &&
		platformStr != model.TypeForWeb &&
		platformStr != model.TypeForAndroidTv {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = srv.SetSpecial(c, aid, click, platformStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func lock(c *bm.Context) {
	var (
		aid, plat, lv, lock                 int64
		aidStr, platformStr, lvStr, lockStr string
		err                                 error
	)
	params := c.Request.Form
	if aidStr = params.Get("aid"); aidStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if platformStr = params.Get("platform"); platformStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if lvStr = params.Get("lv"); lvStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if lockStr = params.Get("lock"); lockStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if plat, err = strconv.ParseInt(platformStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if lv, err = strconv.ParseInt(lvStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if lock, err = strconv.ParseInt(lockStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = srv.SetLock(c, aid, int8(plat), int8(lock), int8(lv)); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func lockMid(c *bm.Context) {
	var (
		mid       int64
		status    int64
		statusStr string
		err       error
	)
	params := c.Request.Form
	mid, err = strconv.ParseInt(params.Get("mid"), 10, 64)
	if mid <= 0 || err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if statusStr = params.Get("status"); statusStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if status, err = strconv.ParseInt(statusStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if status != 0 && status != 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = srv.SetMidForbid(c, mid, int8(status))
	c.JSON(nil, err)
}
