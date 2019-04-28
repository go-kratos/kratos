package http

import (
	"strconv"
	"time"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// reportInfoc infoc
func reportInfoc(c *bm.Context) {
	var (
		params = c.Request.Form
		err    error
	)
	eventID := params.Get("event_id")
	eventType := params.Get("event_type")
	// header
	buvid := params.Get("buvid")
	fts := params.Get("fts")
	if _, err = strconv.ParseInt(fts, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	messageInfo := params.Get("message_info")
	if eventID == "" || eventType == "" || buvid == "" || messageInfo == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, dataflowSvr.Report(c, eventID, eventType, buvid, fts, messageInfo, time.Now()))
}
