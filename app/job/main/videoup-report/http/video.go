package http

import (
	"time"

	"go-common/app/job/main/videoup-report/model/archive"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// reportsByType request report by type
func reportsByType(c *bm.Context, t int8) {
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

	c.JSON(vdaSvc.VideoReports(c, t, stime, etime))
}

// videoAudit get video audit reports
func videoAudit(c *bm.Context) {
	reportsByType(c, archive.ReportTypeVideoAudit)
}

// videoXcode get video xcode reports
func videoXcode(c *bm.Context) {
	reportsByType(c, archive.ReportTypeXcode)
}
