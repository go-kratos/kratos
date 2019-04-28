package http

import (
	"strconv"
	"time"

	"go-common/app/admin/main/reply/conf"
	"go-common/app/admin/main/reply/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func reportSearch(c *bm.Context) {
	var (
		err                     error
		typ, startTime, endTime int64
		page                    = int64(conf.Conf.Reply.PageNum)
		pageSize                = int64(conf.Conf.Reply.PageSize)
		sp                      = &model.SearchReportParams{}
	)

	params := c.Request.Form
	typeStr := params.Get("type")
	oidStr := params.Get("oid")
	uidStr := params.Get("uid")
	startTimeStr := params.Get("start_time")
	endTimeStr := params.Get("end_time")
	pageStr := params.Get("page")
	pageSizeStr := params.Get("pagesize")
	// parse params
	if typ, err = strconv.ParseInt(typeStr, 10, 64); err != nil {
		log.Warn("strconv.ParseInt(type:%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
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
	if startTimeStr != "" {
		if startTime, err = strconv.ParseInt(startTimeStr, 10, 64); err == nil {
			sp.StartTime = time.Unix(startTime, 0).Format(model.DateFormat)
		} else {
			var t time.Time
			t, err = time.Parse("2006-01-02", startTimeStr)
			if err == nil {
				sp.StartTime = t.Format(model.DateFormat)
			} else {
				sp.StartTime = startTimeStr
			}
		}
	}
	if endTimeStr != "" {
		if endTime, err = strconv.ParseInt(endTimeStr, 10, 64); err == nil {
			sp.EndTime = time.Unix(endTime, 0).Format(model.DateFormat)
		} else {
			var t time.Time
			t, err = time.Parse("2006-01-02", endTimeStr)
			if err == nil {
				sp.EndTime = t.Format(model.DateFormat)
			} else {
				sp.EndTime = endTimeStr
			}
		}
	} else if startTimeStr != "" {
		t := time.Now()
		sp.EndTime = t.Format(model.DateFormat)
	}
	if pageStr != "" {
		if page, err = strconv.ParseInt(pageStr, 10, 64); err != nil || page < 1 {
			log.Warn("strconv.ParseInt(page:%s) error(%v)", pageStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if pageSizeStr != "" {
		if pageSize, err = strconv.ParseInt(pageSizeStr, 10, 64); err != nil || pageSize < 1 || pageSize > int64(conf.Conf.Reply.PageSize) {
			log.Warn("strconv.ParseInt(page:%s) error(%v)", pageSizeStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	sp.Type = int32(typ)
	sp.Reason = params.Get("reason")
	sp.Typeids = params.Get("typeids")
	sp.Keyword = params.Get("keyword")
	sp.Nickname = params.Get("nickname")
	sp.States = params.Get("states")
	sp.Order = params.Get("order")
	sp.Sort = params.Get("sort")
	rpts, err := rpSvc.ReportSearch(c, sp, page, pageSize)
	if err != nil {
		log.Error("rpSvc.ReportSearch(%d,%d,%v) error(%v)", page, pageSize, sp, err)
		c.JSON(nil, err)
		return
	}
	res := map[string]interface{}{}
	res["data"] = rpts
	c.JSONMap(res, nil)
}

func reportDel(c *bm.Context) {
	v := new(struct {
		RpID    []int64 `form:"rpid,split" validate:"required"`
		Oid     []int64 `form:"oid,split" validate:"required"`
		Type    []int32 `form:"type,split" validate:"required"`
		AdminID int64   `form:"adid"`
		Remark  string  `form:"remark"`
		Moral   int32   `form:"moral"`
		Notify  bool    `form:"notify"`
		FTime   int64   `form:"ftime"`
		AdName  string  `form:"adname"`
		FReason int32   `form:"freason"`
		Audit   int32   `form:"audit"`
		Reason  int32   `form:"reason"`
		Content string  `form:"reason_content"`
	})
	err := c.Bind(v)
	if err != nil {
		return
	}
	if len(v.RpID) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(v.RpID) != len(v.Oid) || len(v.RpID) != len(v.Type) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// reason没传修改为-1而不是0
	params := c.Request.Form
	if params.Get("reason") == "" {
		v.Reason = -1
	}
	tMap := make(map[int32]*Compose)
	for i, tp := range v.Type {
		if c, ok := tMap[tp]; ok {
			c.Oids = append(c.Oids, v.Oid[i])
			c.RpIDs = append(c.RpIDs, v.RpID[i])
		} else {
			c = &Compose{
				Oids:  []int64{v.Oid[i]},
				RpIDs: []int64{v.RpID[i]},
			}
			tMap[tp] = c
		}
	}
	adid := v.AdminID
	if uid, ok := c.Get("uid"); ok {
		adid = uid.(int64)
	}
	adname := v.AdName
	if username, ok := c.Get("username"); ok {
		adname = username.(string)
	}
	for tp, com := range tMap {
		err = rpSvc.ReportDel(c, com.Oids, com.RpIDs, adid, v.FTime, tp, v.Audit, v.Moral, v.Reason, v.FReason, v.Notify, adname, v.Remark, v.Content)
		if err != nil {
			c.JSON(nil, err)
			return
		}
	}
	c.JSON(nil, err)
	return
}

func reportIgnore(c *bm.Context) {
	v := new(struct {
		RpID    []int64 `form:"rpid,split" validate:"required"`
		Oid     []int64 `form:"oid,split" validate:"required"`
		Type    []int32 `form:"type,split" validate:"required"`
		AdminID int64   `form:"adid"`
		Remark  string  `form:"remark"`
		Audit   int32   `form:"audit"`
	})
	err := c.Bind(v)
	if err != nil {
		return
	}
	if len(v.RpID) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(v.RpID) != len(v.Oid) || len(v.RpID) != len(v.Type) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tMap := make(map[int32]*Compose)
	for i, tp := range v.Type {
		if c, ok := tMap[tp]; ok {
			c.Oids = append(c.Oids, v.Oid[i])
			c.RpIDs = append(c.RpIDs, v.RpID[i])
		} else {
			c = &Compose{
				Oids:  []int64{v.Oid[i]},
				RpIDs: []int64{v.RpID[i]},
			}
			tMap[tp] = c
		}
	}
	adid := v.AdminID
	if uid, ok := c.Get("uid"); ok {
		adid = uid.(int64)
	}
	var adName string
	if uname, ok := c.Get("username"); ok {
		adName = uname.(string)
	}
	for tp, com := range tMap {
		err = rpSvc.ReportIgnore(c, com.Oids, com.RpIDs, adid, adName, tp, v.Audit, v.Remark, true)
	}
	c.JSON(nil, err)
	return
}

func reportRecover(c *bm.Context) {
	v := new(struct {
		RpID    []int64 `form:"rpid,split" validate:"required"`
		Oid     []int64 `form:"oid,split" validate:"required"`
		Type    []int32 `form:"type,split" validate:"required"`
		AdminID int64   `form:"adid"`
		Remark  string  `form:"remark"`
		Audit   int32   `form:"audit"`
	})
	err := c.Bind(v)
	if err != nil {
		return
	}
	if len(v.RpID) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(v.RpID) != len(v.Oid) || len(v.RpID) != len(v.Type) {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	tMap := make(map[int32]*Compose)
	for i, tp := range v.Type {
		if c, ok := tMap[tp]; ok {
			c.Oids = append(c.Oids, v.Oid[i])
			c.RpIDs = append(c.RpIDs, v.RpID[i])
		} else {
			c = &Compose{
				Oids:  []int64{v.Oid[i]},
				RpIDs: []int64{v.RpID[i]},
			}
			tMap[tp] = c
		}
	}
	adid := v.AdminID
	if uid, ok := c.Get("uid"); ok {
		adid = uid.(int64)
	}

	for tp, com := range tMap {
		err = rpSvc.ReportRecover(c, com.Oids, com.RpIDs, adid, tp, v.Audit, v.Remark)
	}
	c.JSON(nil, err)
}

func reportTransfer(c *bm.Context) {
	v := new(struct {
		RpID    []int64 `form:"rpid,split" validate:"required"`
		Oid     []int64 `form:"oid,split" validate:"required"`
		Type    []int32 `form:"type,split" validate:"required"`
		AdminID int64   `form:"adid"`
		Remark  string  `form:"remark"`
		Audit   int32   `form:"audit"`
	})
	err := c.Bind(v)
	if err != nil {
		return
	}
	if len(v.RpID) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(v.RpID) != len(v.Oid) || len(v.RpID) != len(v.Type) {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	tMap := make(map[int32]*Compose)
	for i, tp := range v.Type {
		if c, ok := tMap[tp]; ok {
			c.Oids = append(c.Oids, v.Oid[i])
			c.RpIDs = append(c.RpIDs, v.RpID[i])
		} else {
			c = &Compose{
				Oids:  []int64{v.Oid[i]},
				RpIDs: []int64{v.RpID[i]},
			}
			tMap[tp] = c
		}
	}
	adid := v.AdminID
	if uid, ok := c.Get("uid"); ok {
		adid = uid.(int64)
	}
	var adName string
	if uname, ok := c.Get("username"); ok {
		adName = uname.(string)
	}
	for tp, com := range tMap {
		err = rpSvc.ReportTransfer(c, com.Oids, com.RpIDs, adid, adName, tp, v.Audit, v.Remark)
	}
	c.JSON(nil, err)
}

func reportStateSet(c *bm.Context) {
	v := new(struct {
		RpID    []int64 `form:"rpid,split" validate:"required"`
		Oid     []int64 `form:"oid,split" validate:"required"`
		Type    []int32 `form:"type,split" validate:"required"`
		AdminID int64   `form:"adid"`
		AdName  string  `form:"adname"`
		Remark  string  `form:"remark"`
		State   int32   `form:"state"`
	})
	err := c.Bind(v)
	if err != nil {
		return
	}
	if len(v.RpID) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(v.RpID) != len(v.Oid) || len(v.RpID) != len(v.Type) {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	tMap := make(map[int32]*Compose)
	for i, tp := range v.Type {
		if c, ok := tMap[tp]; ok {
			c.Oids = append(c.Oids, v.Oid[i])
			c.RpIDs = append(c.RpIDs, v.RpID[i])
		} else {
			c = &Compose{
				Oids:  []int64{v.Oid[i]},
				RpIDs: []int64{v.RpID[i]},
			}
			tMap[tp] = c
		}
	}
	if uid, ok := c.Get("uid"); ok {
		v.AdminID = uid.(int64)
	}
	if uname, ok := c.Get("username"); ok {
		v.AdName = uname.(string)
	}
	for tp, com := range tMap {
		err = rpSvc.ReportStateSet(c, com.Oids, com.RpIDs, v.AdminID, v.AdName, tp, v.State, v.Remark, true)
	}
	c.JSON(nil, err)

}
