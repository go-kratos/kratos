package http

import (
	"time"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func memberStats(c *bm.Context) {
	c.JSON(srv.GetStats(c))
}

func taskTooks(c *bm.Context) {
	req := c.Request
	params := req.Form
	stimeStr := params.Get("stime")
	etimeStr := params.Get("etime")
	if stimeStr == "" {
		stimeStr = time.Now().Format("2006-01-02") + " 00:00:00"
	}
	if etimeStr == "" {
		etimeStr = time.Now().Format("2006-01-02 15:04:05")
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
	tooks, err := srv.TaskTooksByHalfHour(c, stime, etime)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(tooks, nil)
}
