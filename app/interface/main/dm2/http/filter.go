package http

import (
	"encoding/json"
	"strconv"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func upFilters(c *bm.Context) {
	mid, _ := c.Get("mid")
	data, err := dmSvc.UpFilters(c, mid.(int64))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

// 前端协管/up主禁言用户接口
func addUpFilterID(c *bm.Context) {
	var (
		p       = c.Request.Form
		fltList = make([]*model.UpFilter, 0)
		fltMap  = make(map[string]string)
	)
	mid, _ := c.Get("mid")
	fType, err := strconv.ParseInt(p.Get("type"), 10, 64)
	if err != nil || fType != int64(model.FilterTypeID) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	oid, err := strconv.ParseInt(p.Get("oid"), 10, 64)
	if err != nil || oid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = json.Unmarshal([]byte(p.Get("filters")), &fltList); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	for _, v := range fltList {
		fltMap[v.Filter] = v.Comment
	}
	if err = dmSvc.AddUpFilterID(c, mid.(int64), oid, fltMap); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// 创作中心up主添加屏蔽词
func editUpFilters(c *bm.Context) {
	var (
		p       = c.Request.Form
		filters = make([]*model.UpFilter, 0)
		fType   int64
		fltMap  = make(map[string]string)
	)
	mid, _ := c.Get("mid")
	active, err := strconv.ParseInt(p.Get("active"), 10, 64)
	if err != nil || (int8(active) != model.FilterActive && int8(active) != model.FilterUnActive) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = json.Unmarshal([]byte(p.Get("filters")), &filters); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if fType, err = strconv.ParseInt(p.Get("type"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	for _, filter := range filters {
		fltMap[filter.Filter] = filter.Comment
	}
	switch int8(active) {
	case model.FilterActive:
		err = dmSvc.AddUpFilters(c, mid.(int64), int8(fType), fltMap)
	case model.FilterUnActive:
		flts := make([]string, 0)
		for f := range fltMap {
			flts = append(flts, f)
		}
		if _, err = dmSvc.EditUpFilters(c, mid.(int64), int8(fType), model.FilterUnActive, flts); err != nil {
			break
		}
	default:
		err = ecode.RequestErr
	}
	c.JSON(nil, err)
}
