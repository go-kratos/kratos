package http

import (
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// trackArchive get archive list.
func trackArchive(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	arcTrack, err := vdaSvc.TrackArchiveInfo(c, aid)
	if err != nil {
		log.Error("vdaSvc.TrackArchiveInfo (aid(%d)) error(%v)", aid, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(arcTrack, nil)
}

//trackDetail 稿件信视频息追踪详情页
func trackDetail(c *bm.Context) {
	params := c.Request.Form
	hidStr := params.Get("hid")
	hid, err := strconv.ParseInt(hidStr, 10, 64)
	if err != nil {
		log.Error("vdaSvc.trackDetail strconv.ParseInt(%s) error(%v)", hidStr, err)
		c.JSON(nil, err)
		return
	}

	history, err := vdaSvc.TrackHistoryDetail(c, hid)
	if err != nil {
		log.Error("vdaSvc.TrackHistoryDetail (hid(%d)) error(%v)", hid, err)
		c.JSON(nil, err)
		return
	}

	c.JSON(history, nil)
}

// trackVideo get video process.
func trackVideo(c *bm.Context) {
	params := c.Request.Form
	fn := params.Get("filename")
	aid, err := strconv.ParseInt(params.Get("aid"), 10, 64)
	if err != nil {
		log.Error("trackVideo strconv.ParseInt error(%v) aid(%s)", err, params.Get("aid"))
		c.JSON(nil, ecode.RequestErr)
		return
	}

	// video
	video, err := vdaSvc.TrackVideo(c, fn, aid)
	if err != nil {
		log.Error("vdaSvc.TrackVideo(%s, aid(%d)) error(%v)", fn, aid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(video, nil)
}
