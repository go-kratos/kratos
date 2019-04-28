package http

import (
	"strconv"

	pushmdl "go-common/app/service/main/push/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var settingBizMap = map[int32]string{
	pushmdl.UserSettingArchive: "archive",
	pushmdl.UserSettingLive:    "live",
}

func setSetting(c *bm.Context) {
	var (
		params    = c.Request.Form
		midStr, _ = c.Get("mid")
	)
	mid, _ := midStr.(int64)
	if mid <= 0 {
		log.Warn("mid(%s) is wrong", midStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	typ, _ := strconv.Atoi(params.Get("type"))
	if _, ok := pushmdl.Settings[typ]; !ok {
		log.Warn("type(%s) is wrong", params.Get("type"))
		c.JSON(nil, ecode.RequestErr)
		return
	}
	val, _ := strconv.Atoi(params.Get("value"))
	if val != pushmdl.SwitchOn && val != pushmdl.SwitchOff {
		log.Warn("value(%s) is wrong", params.Get("value"))
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, pushSrv.SetSetting(c, mid, typ, val))
}

func setting(c *bm.Context) {
	midStr, _ := c.Get("mid")
	mid, _ := midStr.(int64)
	if mid <= 0 {
		log.Warn("mid(%s) is wrong", midStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	setting, err := pushSrv.Setting(c, mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	res := make(map[string]int32)
	for t, v := range setting {
		if _, ok := settingBizMap[t]; !ok {
			continue
		}
		res[settingBizMap[t]] = v
	}
	c.JSON(res, nil)
}
