package http

import (
	"strconv"
	"time"

	"go-common/app/admin/main/reply/conf"
	"go-common/app/admin/main/reply/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

func monitorStats(c *bm.Context) {
	var (
		err       error
		mode      int64
		startTime int64
		endTime   int64
		page      = int64(conf.Conf.Reply.PageNum)
		pageSize  = int64(conf.Conf.Reply.PageSize)
	)
	params := c.Request.Form
	modeStr := params.Get("mode")
	sortStr := params.Get("sort")
	pageStr := params.Get("page")
	pageSizeStr := params.Get("pagesize")
	adminIDsStr := params.Get("adminids")
	startTimeStr := params.Get("start_time")
	endTimeStr := params.Get("end_time")
	orderStr := params.Get("order_time")
	if mode, err = strconv.ParseInt(modeStr, 10, 64); err != nil {
		log.Warn("strconv.ParseInt(mode:%s) error(%v)", modeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pageStr != "" {
		if page, err = strconv.ParseInt(pageStr, 10, 64); err != nil {
			log.Warn("strconv.ParseInt(page:%s) error(%v)", pageStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if pageSizeStr != "" {
		if pageSize, err = strconv.ParseInt(pageSizeStr, 10, 64); err != nil {
			log.Warn("strconv.ParseInt(pageSize:%s) error(%v)", pageSizeStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if startTimeStr != "" {
		if startTime, err = strconv.ParseInt(startTimeStr, 10, 64); err == nil {
			startTimeStr = time.Unix(startTime, 0).Format(model.DateSimpleFormat)
		}
	}
	if endTimeStr != "" {
		if endTime, err = strconv.ParseInt(endTimeStr, 10, 64); err == nil {
			endTimeStr = time.Unix(endTime, 0).Format(model.DateSimpleFormat)
		}
	} else {
		t := time.Now()
		endTimeStr = t.Format(model.DateSimpleFormat)

	}
	stats, err := rpSvc.MonitorStats(c, mode, page, pageSize, adminIDsStr, sortStr, orderStr, startTimeStr, endTimeStr)
	if err != nil {
		log.Error("svc.MonitorStats(%d,%d,%d,%s,%s,%s,%s,%s) error(%v)", mode, page, pageSize, adminIDsStr, sortStr, orderStr, startTimeStr, endTimeStr, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(stats, nil)
}

func monitorSearch(c *bm.Context) {
	var (
		err       error
		mode, typ int64
		page      = int64(conf.Conf.Reply.PageNum)
		pageSize  = int64(conf.Conf.Reply.PageSize)
		sp        = &model.SearchMonitorParams{}
	)

	params := c.Request.Form
	modeStr := params.Get("mode")
	typeStr := params.Get("type")
	oidStr := params.Get("oid")
	uidStr := params.Get("uid")
	pageStr := params.Get("page")
	pageSizeStr := params.Get("pagesize")
	if typ, err = strconv.ParseInt(typeStr, 10, 64); err != nil {
		log.Warn("strconv.ParseInt(type:%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if modeStr != "" {
		if mode, err = strconv.ParseInt(modeStr, 10, 64); err != nil {
			log.Warn("strconv.ParseInt(mode:%s) error(%v)", typeStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if oidStr != "" {
		if sp.Oid, err = strconv.ParseInt(oidStr, 10, 64); err != nil {
			log.Warn("strconv.ParseInt(oid:%s) error(%v)", oidStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if uidStr != "" {
		if sp.UID, err = strconv.ParseInt(uidStr, 10, 64); err != nil {
			log.Warn("strconv.ParseInt(uid:%s) error(%v)", uidStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	sp.Mode = int8(mode)
	sp.Type = int8(typ)
	sp.Keyword = params.Get("keyword")
	sp.NickName = params.Get("nickname")
	sp.Sort = params.Get("sort")
	sp.Order = params.Get("order")
	if pageStr != "" {
		if page, err = strconv.ParseInt(pageStr, 10, 64); err != nil || page < 1 {
			log.Warn("strconv.ParseInt(page:%s) error(%v)", pageStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if pageSizeStr != "" {
		if pageSize, err = strconv.ParseInt(pageSizeStr, 10, 64); err != nil || pageSize < 1 || pageSize > int64(conf.Conf.Reply.PageSize) {
			log.Warn("strconv.ParseInt(pagesize:%s) error(%v)", pageSizeStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	rpts, err := rpSvc.MonitorSearch(c, sp, page, pageSize)
	if err != nil {
		log.Error("svc.ReportSearch(%d,%d,%v) error(%v)", page, pageSize, sp, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(rpts, nil)
}

func monitorState(c *bm.Context) {
	var (
		err     error
		typ     int64
		state   int64
		adminID int64
		oids    []int64
	)

	params := c.Request.Form
	typeStr := params.Get("type")
	oidStr := params.Get("oid")
	stateStr := params.Get("state")
	adminIDStr := params.Get("adid")
	remark := params.Get("remark")
	if oids, err = xstr.SplitInts(oidStr); err != nil {
		log.Warn("strconv.ParseInt(oid:%s) error(%v)", oidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if typ, err = strconv.ParseInt(typeStr, 10, 64); err != nil {
		log.Warn("strconv.ParseInt(type:%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if state, err = strconv.ParseInt(stateStr, 10, 64); err != nil {
		log.Warn("strconv.ParseInt(state:%s) error(%v)", stateStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if adminIDStr != "" {
		if adminID, err = strconv.ParseInt(adminIDStr, 10, 64); err != nil {
			log.Warn("strconv.ParseInt(admin:%s) error(%v)", adminIDStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if uid, ok := c.Get("uid"); ok {
		adminID = uid.(int64)
	}
	var adName string
	if uname, ok := c.Get("username"); ok {
		adName = uname.(string)
	}
	for _, oid := range oids {
		if err = rpSvc.UpMonitorState(c, adminID, adName, oid, int32(typ), int32(state), remark); err != nil {
			log.Error("svc.MonitorState(%d,%d,%d,%d,%s) error(%v)", oid, typ, state, adminID, remark, err)
			c.JSON(nil, err)
			return
		}
	}
	c.JSON(nil, nil)
}

// monitorLog get monitor log.
func monitorLog(c *bm.Context) {
	v := new(struct {
		Oid       int64  `form:"oid" `
		Type      int32  `form:"type" `
		StartTime string `form:"start_time"`
		EndTime   string `form:"end_time"`
		Page      int64  `form:"pn"`
		PageSize  int64  `form:"ps"`
		Mid       int64  `form:"mid"`
		Order     string `form:"order"`
		Sort      string `form:"sort"`
	})
	var err error
	err = c.Bind(v)
	if err != nil {
		return
	}

	sp := model.LogSearchParam{
		Oid:       v.Oid,
		Type:      v.Type,
		Mid:       v.Mid,
		CtimeFrom: v.StartTime,
		CtimeTo:   v.EndTime,
		Pn:        v.Page,
		Ps:        v.PageSize,
		Order:     v.Order,
		Sort:      v.Sort,
	}
	data, err := rpSvc.MointorLog(c, sp)
	res := map[string]interface{}{}
	res["data"] = data

	c.JSONMap(res, err)
	return
}
