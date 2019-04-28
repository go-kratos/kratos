package http

import (
	"strconv"
	"strings"

	"go-common/app/admin/main/filter/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

func whiteAddArea(c *bm.Context) {
	var (
		err   error
		mode  int64
		adid  int64
		areas []string
		tps   []int64
	)
	params := c.Request.Form
	contentStr := params.Get("filter")
	modeStr := params.Get("mode")
	areaStr := params.Get("area")
	tpStr := params.Get("tpid")
	adidStr := params.Get("adid")
	nameStr := params.Get("name")
	commentStr := params.Get("comment")
	if modeStr != "" {
		if mode, err = strconv.ParseInt(modeStr, 10, 8); err != nil {
			log.Error("strconv.ParseInt() err(%v)", err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if areaStr == "" {
		areas = append(areas, _baseArea)
	} else {
		areas = strings.Split(areaStr, ",")
	}
	if tpStr != "" {
		if tps, err = xstr.SplitInts(tpStr); err != nil {
			log.Error("xstr.SplintInts(tpStr %s) err(%v)", tpStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	} else {
		tps = []int64{0}
	}
	if contentStr == "" || len([]rune(contentStr)) > 64 {
		log.Error("contentStr == nil or contentStr > 60 contentStr(%s)", contentStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if adid, err = strconv.ParseInt(adidStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", adidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if nameStr == "" {
		log.Error("nameStr == nil", nameStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len([]rune(commentStr)) > 50 || len([]rune(commentStr)) < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svc.AddAreaWhite(c, contentStr, int8(mode), areas, tps, adid, nameStr, commentStr))
}

func whiteDel(c *bm.Context) {
	var (
		err       error
		adid      int64
		contentID int64
	)

	params := c.Request.Form
	contentStr := params.Get("filter_id")
	adidStr := params.Get("adid")
	nameStr := params.Get("name")
	reasonStr := params.Get("reason")
	if contentID, err = strconv.ParseInt(contentStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", contentStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if adid, err = strconv.ParseInt(adidStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", adidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if nameStr == "" {
		log.Error("nameStr == nil", nameStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len([]rune(reasonStr)) > 50 || len([]rune(reasonStr)) < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svc.DeleteWhite(c, contentID, adid, nameStr, reasonStr))
}

func whiteSearch(c *bm.Context) {
	var (
		err    error
		total  int64
		pn, ps int64
		rs     []*model.WhiteInfo
	)

	params := c.Request.Form
	contentStr := params.Get("filter")
	areaStr := params.Get("area")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	if pnStr != "" {
		if pn, err = strconv.ParseInt(pnStr, 10, 64); err != nil || pn < 1 {
			log.Error("strconv.ParseInt(%s) error(%v)", pnStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	} else {
		pn = 1
	}
	if psStr != "" {
		if ps, err = strconv.ParseInt(psStr, 10, 64); err != nil || ps < 0 {
			log.Error("strconv.ParseInt(%s) error(%v)", psStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	} else {
		ps = 20
	}
	if areaStr == "" {
		areaStr = _baseArea
	}
	if rs, total, err = svc.SearchWhite(c, contentStr, areaStr, pn, ps); err != nil {
		c.JSON(nil, err)
		return
	}
	var data = map[string]interface{}{
		"rules": rs,
		"total": total,
		"pn":    pn,
		"ps":    ps,
	}
	c.JSON(data, nil)
}

func whiteEditInfo(c *bm.Context) {
	var (
		err       error
		contentID int64
	)
	params := c.Request.Form
	contentStr := params.Get("filter_id")
	if contentID, err = strconv.ParseInt(contentStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", contentStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(svc.WhiteInfo(c, contentID))
}

func whiteEdit(c *bm.Context) {
	var (
		err   error
		mode  int64
		adid  int64
		areas []string
		tps   []int64
	)

	params := c.Request.Form
	contentStr := params.Get("filter")
	areaStr := params.Get("area")
	tpStr := params.Get("tpid")
	modeStr := params.Get("mode")
	reasonStr := params.Get("reason")
	commentStr := params.Get("comment")
	adidStr := params.Get("adid")
	nameStr := params.Get("name")
	if modeStr != "" {
		if mode, err = strconv.ParseInt(modeStr, 10, 8); err != nil {
			log.Error("strconv.ParseInt(%s) err(%v)", modeStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if areaStr == "" {
		areas = append(areas, _baseArea)
	} else {
		areas = strings.Split(areaStr, ",")
	}
	if tpStr != "" {
		if tps, err = xstr.SplitInts(tpStr); err != nil {
			log.Error("xstr.SplintInts(tpStr %s) err(%v)", tpStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	} else {
		tps = []int64{0}
	}
	if contentStr == "" {
		log.Error("contentStr == nil")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if adid, err = strconv.ParseInt(adidStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", adidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if nameStr == "" {
		log.Error("nameStr == nil", nameStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len([]rune(reasonStr)) > 50 || len([]rune(reasonStr)) < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len([]rune(commentStr)) > 50 || len([]rune(commentStr)) < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svc.EditWhite(c, contentStr, int8(mode), areas, tps, adid, nameStr, commentStr, reasonStr))
}

func whiteEditLog(c *bm.Context) {
	var (
		err       error
		contentID int64
	)

	params := c.Request.Form
	contentStr := params.Get("filter_id")
	if contentID, err = strconv.ParseInt(contentStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", contentStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(svc.WhiteEditLog(c, contentID))
}
