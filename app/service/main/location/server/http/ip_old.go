package http

import (
	"strings"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// tmpInfo ip info.
func tmpInfo(c *bm.Context) {
	var ip string
	if ip = c.Request.FormValue("ip"); ip == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(svr.TmpInfo(ip))
}

// tmpInfos ip info.
func tmpInfos(c *bm.Context) {
	var ips string
	if ips = c.Request.FormValue("ip"); ips == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	zones, err := svr.TmpInfos(c, strings.Split(ips, ",")...)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if len(zones) == 1 {
		c.JSON(zones[0], nil)
	} else {
		c.JSON(zones, nil)
	}
}
