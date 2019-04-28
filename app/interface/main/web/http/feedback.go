package http

import (
	"net/http"
	"strconv"

	"go-common/app/interface/main/web/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func feedback(c *bm.Context) {
	var (
		mid, aid, tagID int64
		content, buvid  string
		buvidCk         *http.Cookie
		midStr          interface{}
		ok              bool
		err             error
	)
	params := c.Request.Form
	content = params.Get("content")
	aidStr := params.Get("aid")
	if aidStr != "" {
		if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil {
			log.Warn("strconv.ParseInt(%s) err(%d)", aidStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	tagIDStr := params.Get("tag_id")
	if tagID, err = strconv.ParseInt(tagIDStr, 10, 64); err != nil || tagID < 1 {
		log.Warn("strconv.ParseInt(%s) error(%v)", tagIDStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if !model.CheckFeedTag(tagID) {
		log.Warn("tag_id(%d) check fail", tagID)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if buvidCk, err = c.Request.Cookie("buvid3"); err != nil {
		log.Warn("buvid3 is nil")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if buvid = buvidCk.Value; buvid == "" {
		log.Warn("buvid == nil")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if midStr, ok = c.Get("mid"); ok {
		mid = midStr.(int64)
	}
	browser := params.Get("browser")
	version := params.Get("version")
	email := params.Get("email")
	qq := params.Get("qq")
	other := params.Get("other")
	feedParams := &model.Feedback{Aid: aid, Mid: mid, TagID: tagID, Buvid: buvid, Browser: browser,
		Version: version, Content: &model.Content{Reason: content}, Email: email, QQ: qq, Other: other}
	c.JSON(nil, webSvc.Feedback(c, feedParams))
}
