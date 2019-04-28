package http

import (
	"strconv"

	pushmdl "go-common/app/service/main/push/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func setSettingInternal(c *bm.Context) {
	var (
		params  = c.Request.Form
		midStr  = params.Get("mid")
		typeStr = params.Get("type")
		valStr  = params.Get("value")
	)
	mid, _ := strconv.ParseInt(midStr, 10, 64)
	if mid <= 0 {
		log.Warn("mid(%s) is wrong", midStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	typ, _ := strconv.Atoi(typeStr)
	if _, ok := pushmdl.Settings[typ]; !ok {
		log.Warn("type(%s) is wrong", typeStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	val, _ := strconv.Atoi(valStr)
	if val != pushmdl.SwitchOn && val != pushmdl.SwitchOff {
		log.Warn("value(%s) is wrong", valStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, pushSrv.SetSetting(c, mid, typ, val))
}
