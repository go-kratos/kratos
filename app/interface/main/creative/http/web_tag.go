package http

import (
	"strconv"

	"go-common/app/interface/main/creative/model/archive"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func webTags(c *bm.Context) {
	params := c.Request.Form
	tidStr := params.Get("typeid")
	title := params.Get("title")
	filename := params.Get("filename")
	desc := params.Get("desc")
	cover := params.Get("cover")
	// check user
	midStr, ok := c.Get("mid")
	mid := midStr.(int64)
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	tid, _ := strconv.ParseInt(tidStr, 10, 16)
	if tid <= 0 {
		tid = 0
	}
	tags, _ := dataSvc.TagsWithChecked(c, mid, uint16(tid), title, filename, desc, cover, archive.TagPredictFromWeb)
	c.JSON(tags, nil)
}
