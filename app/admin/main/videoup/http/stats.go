package http

import (
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"strconv"
	"time"
)

// statsPoints get stats points data
func statsPoints(c *bm.Context) {
	req := c.Request
	params := req.Form
	stimeStr := params.Get("stime")
	etimeStr := params.Get("etime")
	typeStr := params.Get("type")
	if stimeStr == "" {
		stimeStr = time.Now().Format("2006-01-02") + " 00:00:00"
	}
	if etimeStr == "" {
		etimeStr = time.Now().Format("2006-01-02 15:04:05")
	}
	if typeStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	local, _ := time.LoadLocation("Local")
	stime, err := time.ParseInLocation("2006-01-02 15:04:05", stimeStr, local)
	if stime.Unix() < 1 || err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	etime, err := time.ParseInLocation("2006-01-02 15:04:05", etimeStr, local)
	if etime.Unix() < 1 || err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	typeInt, err := strconv.Atoi(typeStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	points, err := vdaSvc.StatsPoints(c, stime, etime, int8(typeInt))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(points, nil)
}
