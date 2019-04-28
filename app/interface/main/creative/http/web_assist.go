package http

//assist 创作中心协管相关

import (
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"strconv"
	"time"
)

func webAssists(c *bm.Context) {
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	ip := metadata.String(c, metadata.RemoteIP)
	assists, err := assistSvc.Assists(c, mid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(assists, nil)
}

func webAssistLogs(c *bm.Context) {
	req := c.Request
	params := req.Form
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	ip := metadata.String(c, metadata.RemoteIP)
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	pn, err := strconv.ParseInt(pnStr, 10, 64)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.ParseInt(psStr, 10, 64)
	if err != nil || ps <= 10 {
		ps = 10
	}
	assistMidStr := params.Get("assist_mid")
	assistMid, err := strconv.ParseInt(assistMidStr, 10, 64)
	if err != nil {
		assistMid = 0
	}

	stimeStr := params.Get("stime")
	stime, err := strconv.ParseInt(stimeStr, 10, 64)
	if err != nil || stime <= 0 {
		stime = time.Now().Add(-time.Hour * 72).Unix()
	}
	etimeStr := params.Get("etime")
	etime, err := strconv.ParseInt(etimeStr, 10, 64)
	if err != nil || etime <= 0 {
		etime = time.Now().Unix()
	}
	assistLogs, pager, err := assistSvc.AssistLogs(c, mid, assistMid, pn, ps, stime, etime, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSONMap(map[string]interface{}{
		"pager": pager,
		"data":  assistLogs,
	}, nil)
}

func webAssistAdd(c *bm.Context) {
	req := c.Request
	params := req.Form
	ck := c.Request.Header.Get("cookie")
	ak := params.Get("access_key")
	mainStr := params.Get("main")
	liveStr := params.Get("live")
	assistMidStr := params.Get("assist_mid")
	main := 1
	live := 0
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	var (
		err       error
		assistMid int64
		m, l      int
	)
	assistMid, err = strconv.ParseInt(assistMidStr, 10, 64)
	if err != nil || assistMid == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if mainStr != "" {
		m, err = strconv.Atoi(mainStr)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		main = m
	}
	if liveStr != "" {
		l, err = strconv.Atoi(liveStr)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		live = l
	}
	ip := metadata.String(c, metadata.RemoteIP)
	if err = assistSvc.AddAssist(c, mid, assistMid, int8(main), int8(live), ip, ak, ck); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func webAssistSet(c *bm.Context) {
	req := c.Request
	params := req.Form
	ck := c.Request.Header.Get("cookie")
	ak := params.Get("access_key")
	mainStr := params.Get("main")
	liveStr := params.Get("live")
	midI, ok := c.Get("mid")
	main := 1
	live := 0
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	if mainStr != "" {
		m, err := strconv.Atoi(mainStr)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		main = m
	}
	if liveStr != "" {
		l, err := strconv.Atoi(liveStr)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		live = l
	}
	var (
		err       error
		assistMid int64
	)
	assistMidStr := params.Get("assist_mid")
	assistMid, err = strconv.ParseInt(assistMidStr, 10, 64)
	if err != nil || assistMid == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mid, _ := midI.(int64)
	ip := metadata.String(c, metadata.RemoteIP)
	if err = assistSvc.SetAssist(c, mid, assistMid, int8(main), int8(live), ip, ak, ck); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func webAssistDel(c *bm.Context) {
	req := c.Request
	params := req.Form
	ck := c.Request.Header.Get("cookie")
	ak := params.Get("access_key")
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	assistMidStr := params.Get("assist_mid")
	assistMid, err := strconv.ParseInt(assistMidStr, 10, 64)
	if err != nil || assistMid == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mid, _ := midI.(int64)
	ip := metadata.String(c, metadata.RemoteIP)
	if err := assistSvc.DelAssist(c, mid, assistMid, ip, ak, ck); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func webAssistLogRevoc(c *bm.Context) {
	req := c.Request
	params := req.Form
	ck := c.Request.Header.Get("cookie")
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	logIDStr := params.Get("log_id")
	assistMidStr := params.Get("assist_mid")
	logID, _ := strconv.ParseInt(logIDStr, 10, 64)
	assistMid, _ := strconv.ParseInt(assistMidStr, 10, 64)
	if assistMid < 1 || logID < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mid, _ := midI.(int64)
	ip := metadata.String(c, metadata.RemoteIP)
	if err := assistSvc.RevocAssistLog(c, mid, assistMid, logID, ck, ip); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func webAssistStatus(c *bm.Context) {
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	ip := metadata.String(c, metadata.RemoteIP)
	status, err := assistSvc.LiveStatus(c, mid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]int8{
		"live": status,
	}, nil)
}
