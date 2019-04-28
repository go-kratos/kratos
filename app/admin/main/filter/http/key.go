package http

import (
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/filter/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func keyAdd(c *bm.Context) {
	var (
		err   error
		mode  int64
		level int64
		stime int64
		etime int64
		adid  int64
		areas []string
	)
	params := c.Request.Form
	areaStr := params.Get("area")
	keyStr := params.Get("key")
	ruleStr := params.Get("rule")
	nameStr := params.Get("name")
	modeStr := params.Get("mode")
	levelStr := params.Get("level")
	stimeStr := params.Get("stime")
	etimeStr := params.Get("etime")
	adidStr := params.Get("adid")
	commentStr := params.Get("comment")
	if areaStr == "" || keyStr == "" || ruleStr == "" || nameStr == "" || levelStr == "" {
		log.Error("strconv.ParseInt(%s,%s,%s,%s) err(%v)", areaStr, keyStr, ruleStr, nameStr, levelStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len([]rune(commentStr)) > 128 || len([]rune(ruleStr)) > 255 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	areas = strings.Split(areaStr, ",")
	if mode, err = strconv.ParseInt(modeStr, 10, 8); err != nil {
		log.Error("strconv.ParseInt(%s) err(%v)", modeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if level, err = strconv.ParseInt(levelStr, 10, 8); err != nil {
		log.Error("strconv.ParseInt(%s) err(%v)", levelStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if adid, err = strconv.ParseInt(adidStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) err(%v)", adidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if stimeStr != "" {
		stime, err = strconv.ParseInt(stimeStr, 10, 64)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	} else {
		stime = time.Now().Unix()
	}
	if etimeStr != "" {
		etime, err = strconv.ParseInt(etimeStr, 10, 64)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	} else {
		etime = stime + 3600*24*365*10 // 10年
	}
	c.JSON(nil, svc.AddKey(c, areas, keyStr, ruleStr, commentStr, nameStr, int8(mode), int8(level), adid, stime, etime))
}

func keyDelFid(c *bm.Context) {
	var (
		err  error
		fid  int64
		adid int64
	)
	params := c.Request.Form
	keyStr := params.Get("key")
	nameStr := params.Get("name")
	fidStr := params.Get("fid")
	adidStr := params.Get("adid")
	commentStr := params.Get("comment")
	reasonStr := params.Get("reason")
	if keyStr == "" || nameStr == "" || reasonStr == "" || len([]rune(reasonStr)) > 50 {
		log.Error("strconv.ParseInt(%s,%s,%s,%s) err(%v)", keyStr, commentStr, nameStr, reasonStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if fid, err = strconv.ParseInt(fidStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", fidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if adid, err = strconv.ParseInt(adidStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", adidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svc.DelKeyFid(c, keyStr, fid, adid, commentStr, nameStr, reasonStr))
}

func keyEditInfo(c *bm.Context) {
	var (
		err error
		fid int64
		fil *model.KeyInfo
	)
	params := c.Request.Form
	keyStr := params.Get("key")
	fidStr := params.Get("fid")
	if keyStr == "" || fidStr == "" {
		log.Error("strconv.ParseInt(%s,%s) error(%v)", keyStr, fidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if fid, err = strconv.ParseInt(fidStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", fidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if fil, err = svc.EditInfo(c, keyStr, fid); err != nil {
		c.JSON(nil, err)
		return
	}
	var data = map[string]interface{}{
		"rule": fil,
	}
	c.JSON(data, nil)
}

func keyEdit(c *bm.Context) {
	var (
		err   error
		fid   int64
		mode  int64
		level int64
		stime int64
		etime int64
		adid  int64
		areas []string
	)
	params := c.Request.Form
	keyStr := params.Get("key")
	areaStr := params.Get("area")
	ruleStr := params.Get("rule")
	nameStr := params.Get("name")
	fidStr := params.Get("fid")
	modeStr := params.Get("mode")
	levelStr := params.Get("level")
	stimeStr := params.Get("stime")
	etimeStr := params.Get("etime")
	adidStr := params.Get("adid")
	commentStr := params.Get("comment")
	reasonStr := params.Get("reason")
	if areaStr == "" || keyStr == "" || ruleStr == "" || nameStr == "" || levelStr == "" || reasonStr == "" {
		log.Error("strconv.ParseInt(%s,%s,%s,%s,%s,%s) err(%v)", keyStr, areaStr, commentStr, nameStr,
			levelStr, reasonStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len([]rune(commentStr)) > 128 || len([]rune(ruleStr)) > 255 || len([]rune(reasonStr)) > 50 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if fid, err = strconv.ParseInt(fidStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", fidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	areas = strings.Split(areaStr, ",")
	if mode, err = strconv.ParseInt(modeStr, 10, 8); err != nil {
		log.Error("strconv.ParseInt(%s) err(%v)", modeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}

	if level, err = strconv.ParseInt(levelStr, 10, 8); err != nil {
		log.Error("strconv.ParseInt(%s) err(%v)", levelStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if adid, err = strconv.ParseInt(adidStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) err(%v)", adidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if stimeStr != "" {
		stime, err = strconv.ParseInt(stimeStr, 10, 64)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	} else {
		stime = time.Now().Unix()
	}
	if etimeStr != "" {
		etime, err = strconv.ParseInt(etimeStr, 10, 64)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	} else {
		etime = stime + 3600*24*365*10 // 10年
	}
	c.JSON(nil, svc.EditKey(c, keyStr, fid, areas, int8(mode), ruleStr, int8(level), stime, etime, adid, nameStr, commentStr, reasonStr))
}

func keySearch(c *bm.Context) {
	var (
		err   error
		pn    int64
		ps    int64
		total int64
		state int64
		rs    []*model.KeyInfo
	)
	params := c.Request.Form
	keyStr := params.Get("key")
	commentStr := params.Get("comment")
	psStr := params.Get("ps")
	pnStr := params.Get("pn")
	stateStr := params.Get("state")
	if pn, err = strconv.ParseInt(pnStr, 10, 64); err != nil || pn < 1 {
		pn = 1
	}
	if ps, err = strconv.ParseInt(psStr, 10, 64); err != nil || ps < 0 || ps > 20 {
		ps = 20
	}
	if state, err = strconv.ParseInt(stateStr, 10, 64); err != nil {
		state = 0
	}
	if total, rs, err = svc.SearchKey(c, keyStr, commentStr, pn, ps, int8(state)); err != nil {
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

func keyLog(c *bm.Context) {
	var (
		err error
	)
	params := c.Request.Form
	keyStr := params.Get("key")
	if keyStr == "" {
		log.Error("strconv.ParseInt(%s) err(%v)", keyStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(svc.FkLog(c, keyStr))
}
