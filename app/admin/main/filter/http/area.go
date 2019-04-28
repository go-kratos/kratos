package http

import (
	"strconv"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func areaGroupList(c *bm.Context) {
	var (
		params = c.Request.Form
		err    error

		pnStr  = params.Get("pn")
		psStr  = params.Get("ps")
		pn, ps int
	)
	if pn, err = strconv.Atoi(pnStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ps, err = strconv.Atoi(psStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ps <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	total, list, err := svc.AreaGroupList(c, ps, pn)
	var data = map[string]interface{}{
		"total": total,
		"list":  list,
	}
	c.JSON(data, err)
}

func areaGroupAdd(c *bm.Context) {
	var (
		params = c.Request.Form
		err    error

		groupName = params.Get("name")
		adidStr   = params.Get("adid")
		adName    = params.Get("ad_name")

		adid int
	)
	if groupName == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if adName == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if adid, err = strconv.Atoi(adidStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if adid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svc.AddAreaGroup(c, groupName, adid, adName))
}

func areaList(c *bm.Context) {
	var (
		params = c.Request.Form
		err    error

		groupIDStr      = params.Get("group_id")
		pnStr           = params.Get("pn")
		psStr           = params.Get("ps")
		pn, ps, groupID int
	)
	if pn, err = strconv.Atoi(pnStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ps, err = strconv.Atoi(psStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ps <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if groupID, err = strconv.Atoi(groupIDStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if groupID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	total, list, err := svc.AreaList(c, groupID, pn, ps)
	var data = map[string]interface{}{
		"total": total,
		"list":  list,
	}
	c.JSON(data, err)
}

func areaAdd(c *bm.Context) {
	var (
		err    error
		params = c.Request.Form

		areaName      = params.Get("name")
		areaShowName  = params.Get("show_name")
		commonFlagStr = params.Get("common_flag")
		groupIDStr    = params.Get("group_id")
		adIDStr       = params.Get("adid")
		adName        = params.Get("ad_name")

		groupID, adID int
		commonFlag    bool
	)
	if areaName == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if areaShowName == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if adName == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if commonFlagStr == "1" {
		commonFlag = true
	} else if commonFlagStr == "0" {
		commonFlag = false
	} else {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if groupID, err = strconv.Atoi(groupIDStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if groupID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if adID, err = strconv.Atoi(adIDStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svc.AddArea(c, groupID, areaName, areaShowName, commonFlag, adID, adName))
}

func areaEdit(c *bm.Context) {
	var (
		err    error
		params = c.Request.Form

		areaIDStr     = params.Get("id")
		commonFlagStr = params.Get("common_flag")
		adIDStr       = params.Get("adid")
		adName        = params.Get("ad_name")
		comment       = params.Get("comment")

		areaID, adID int
		commonFlag   bool
	)
	if comment == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if commonFlagStr == "1" {
		commonFlag = true
	} else if commonFlagStr == "0" {
		commonFlag = false
	} else {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if areaID, err = strconv.Atoi(areaIDStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if areaID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if adID, err = strconv.Atoi(adIDStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svc.EditArea(c, areaID, commonFlag, adID, adName, comment))
}

func areaLog(c *bm.Context) {
	var (
		params = c.Request.Form
		err    error

		idStr = params.Get("id")
		id    int
	)
	if id, err = strconv.Atoi(idStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if id <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(svc.AreaLog(c, id))
}
