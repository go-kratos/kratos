package http

import (
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"strconv"
)

func webFeedbacks(c *bm.Context) {
	params := c.Request.Form
	stateStr := params.Get("state")
	tagIDStr := params.Get("tag_id")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	start := params.Get("start")
	end := params.Get("end")
	platform := params.Get("platform")
	ip := metadata.String(c, metadata.RemoteIP)
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	tagID, err := strconv.ParseInt(tagIDStr, 10, 64)
	if err != nil {
		tagID = 0
	}
	pn, err := strconv.ParseInt(pnStr, 10, 64)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.ParseInt(psStr, 10, 64)
	if err != nil || pn < 1 {
		ps = 10
	}
	if platform == "" { //兼容老逻辑
		platform = "ugc"
	}
	feedbacks, count, err := fdSvc.Feedbacks(c, mid, ps, pn, tagID, stateStr, start, end, platform, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSONMap(map[string]interface{}{
		"data": feedbacks,
		"pager": map[string]int64{
			"pn":    pn,
			"ps":    ps,
			"count": count,
		},
	}, nil)
}

func webFeedbackAdd(c *bm.Context) {
	params := c.Request.Form
	tagIDStr := params.Get("tag_id")
	aid := params.Get("aid")
	title := params.Get("title")
	browser := params.Get("browser")
	content := params.Get("content")
	sessionIDStr := params.Get("session_id")
	qq := params.Get("qq")
	imgURL := params.Get("img_url")
	platform := params.Get("platform")
	ip := metadata.String(c, metadata.RemoteIP)
	tagID, err := strconv.ParseInt(tagIDStr, 10, 64)
	if err != nil {
		log.Error("tagID(%s) format error", tagIDStr)
		tagID = 0
	}
	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		log.Error("sessionID(%s) format error", sessionIDStr)
		sessionID = 0
	}
	if content == "" {
		log.Error("content empty")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if sessionID == 0 && tagID == 0 {
		log.Error("add feedback session tag empty")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if sessionID == 0 {
		// add feedback
		content = title + "#p#" + content
	}
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	if platform == "" { //兼容老逻辑
		platform = "ugc"
	}
	c.JSON(nil, fdSvc.AddFeedback(c, mid, tagID, sessionID, qq, content, aid, imgURL, browser, platform, ip))
}

func webFeedbackDetail(c *bm.Context) {
	params := c.Request.Form
	sessionIDStr := params.Get("session_id")
	ip := metadata.String(c, metadata.RemoteIP)
	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		log.Error("sessionID(%s) format error", sessionIDStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	feedbacks, err := fdSvc.Detail(c, mid, sessionID, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(feedbacks, nil)
}

func webFeedbackTags(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	tags, err := fdSvc.Tags(c, mid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(tags, nil)
}

func webFeedbackClose(c *bm.Context) {
	params := c.Request.Form
	sessionIDStr := params.Get("session_id")
	ip := metadata.String(c, metadata.RemoteIP)
	_, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		log.Error("sessionID(%s) format error", sessionIDStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, fdSvc.CloseSession(c, sessionID, ip))
}

func webFeedbackNewTags(c *bm.Context) {
	ip := metadata.String(c, metadata.RemoteIP)
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	tags, err := fdSvc.NewTags(c, mid, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(tags, nil)
}
