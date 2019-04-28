package http

import (
	"time"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

//moveType get archive move type stats api
func moveType(c *bm.Context) {
	req := c.Request
	param := req.Form
	stimeStr := param.Get("stime")
	etimeStr := param.Get("etime")
	if stimeStr == "" {
		stimeStr = time.Now().Format("2006-01-02") + " 00:00:00"
	}
	if etimeStr == "" {
		etimeStr = time.Now().Format("2006-01-02") + " 00:00:00"
	}
	local, _ := time.LoadLocation("Local")
	stime, err := time.ParseInLocation("2006-01-02 15:04:05", stimeStr, local)
	if stime.Unix() < 1 || err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	etime, err := time.ParseInLocation("2006-01-02 15:04:05", etimeStr, local)
	if etime.Unix() < 1 || err != nil {
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		return
	}

	c.JSON(vdaSvc.MoveType(c, stime, etime))
}

// roundFlow get archive round flow stats api
func roundFlow(c *bm.Context) {
	req := c.Request
	param := req.Form
	stimeStr := param.Get("stime")
	etimeStr := param.Get("etime")
	if stimeStr == "" {
		stimeStr = time.Now().Format("2006-01-02") + " 00:00:00"
	}
	if etimeStr == "" {
		etimeStr = time.Now().Format("2006-01-02") + " 00:00:00"
	}
	local, _ := time.LoadLocation("Local")
	stime, err := time.ParseInLocation("2006-01-02 15:04:05", stimeStr, local)
	if stime.Unix() < 1 || err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	etime, err := time.ParseInLocation("2006-01-02 15:04:05", etimeStr, local)
	if etime.Unix() < 1 || err != nil {
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		return
	}

	c.JSON(vdaSvc.RoundFlow(c, stime, etime))
}
