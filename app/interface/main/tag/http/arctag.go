package http

import (
	"strconv"
	"time"

	"go-common/app/interface/main/tag/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

var (
	reasonMap = map[string]int8{
		"内容不相关":  1,
		"敏感信息":   2,
		"恶意攻击":   3,
		"剧透内容":   4,
		"恶意删除":   5,
		"大量违规操作": 6,
		"无":      7,
	}
)

// add archive tag for user
func addArcTagForOuter(c *bm.Context) {
	var (
		err     error
		aid     int64
		succTid int64
		params  = c.Request.Form
	)
	name := params.Get("tag_name")
	aidStr := params.Get("aid")
	mid, _ := c.Get("mid")
	if name, err = svr.CheckName(name); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || aid < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if succTid, err = svr.Add(c, aid, mid.(int64), name, time.Now()); err != nil {
		log.Error("tagSvr.AddArcForOuter(%d, %d, %s) error(%v)", aid, mid, name, err)
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 1)
	data["tid"] = succTid
	c.JSONMap(data, nil)
}

// delete archive tag for user
func delArcTagForOuter(c *bm.Context) {
	var (
		aid int64
		tid int64
		err error
	)
	mid, _ := c.Get("mid")
	params := c.Request.Form
	tidStr := params.Get("tag_id")
	aidStr := params.Get("aid")
	if tid, err = strconv.ParseInt(tidStr, 10, 64); err != nil || tid < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", tidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || aid < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = svr.Del(c, aid, mid.(int64), tid, time.Now()); err != nil {
		log.Error("arcTagSvr.Del(%d, %d, %d) error(%v)", aid, mid, tid, err)
	}
	c.JSON(nil, err)
}

func arcTags(c *bm.Context) {
	var (
		aid int64
		mid int64
		ts  []*model.Tag
		err error
	)
	params := c.Request.Form
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	aidStr := params.Get("aid")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || aid < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ts, err = svr.ArcTags(c, aid, mid); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(ts, nil)
}

// todo chang to service
func multiArcTags(c *bm.Context) {
	var (
		err  error
		mid  int64
		aids []int64
	)
	params := c.Request.Form
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	aidStr := params.Get("aids")
	if aids, err = xstr.SplitInts(aidStr); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(svr.MutiArcTags(c, mid, aids))
}

func reportArcTag(c *bm.Context) {
	var (
		tid    int64
		aid    int64
		reason int8
		err    error
		ok     bool
	)
	params := c.Request.Form
	aidStr := params.Get("aid")
	tidStr := params.Get("tag_id")
	reasonStr := params.Get("reason")
	mid, _ := c.Get("mid")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || aid < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tid, err = strconv.ParseInt(tidStr, 10, 64); err != nil || tid < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", tidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if reason, ok = reasonMap[reasonStr]; !ok {
		log.Error("reason not in reason map")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = svr.AddReport(c, aid, tid, mid.(int64), reason, time.Now()); err != nil {
		log.Error("svr.AddReport(%d, %d, %d, %d) error(%v)", tid, aid, mid, reason, err)
	}
	c.JSON(nil, err)
}

func logReport(c *bm.Context) {
	var (
		err    error
		ok     bool
		aid    int64
		logID  int64
		reason int8
	)
	params := c.Request.Form
	aidStr := params.Get("aid")
	logIDStr := params.Get("log_id")
	reasonStr := params.Get("reason")
	mid, _ := c.Get("mid")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || aid < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if logID, err = strconv.ParseInt(logIDStr, 10, 64); err != nil || logID < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", logIDStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if reason, ok = reasonMap[reasonStr]; !ok {
		log.Error("reason not in reason map")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = svr.LogReport(c, aid, logID, mid.(int64), reason, time.Now()); err != nil {
		log.Error("svr.LogReport(%d, %d, %d, %d) error(%v)", aid, logID, mid, reason, err)
	}
	c.JSON(nil, err)
}

func arcTagLog(c *bm.Context) {
	var (
		aid, mid int64
		pn       int
		ps       int
		logs     []*model.ArcTagLog
		err      error
	)
	params := c.Request.Form
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	aidStr := params.Get("aid")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || aid < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn, err = strconv.Atoi(pnStr); err != nil || pn < 1 {
		pn = 1
	}
	if ps, err = strconv.Atoi(psStr); err != nil || ps < 1 {
		ps = 20
	}
	if logs, err = svr.Logs(c, aid, mid, pn, ps); err != nil {
		log.Error("arcTagSvr.Logs(%d, %d, %d, %d) error(%v)", aid, mid, pn, ps, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(logs, nil)
}
