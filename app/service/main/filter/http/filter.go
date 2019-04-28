package http

import (
	"strconv"
	"strings"

	"go-common/app/service/main/filter/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func filter(c *bm.Context) {
	var (
		err       error
		id        int64
		tpid      int64
		oid       int64
		mid       int64
		replyType int64
		keys      []string
		resData   = model.HTTPFilterRes{}
	)
	params := c.Request.Form
	msgStr := params.Get("msg")
	areaStr := params.Get("area")
	idStr := params.Get("id")
	tpidStr := params.Get("tpid")
	oidStr := params.Get("oid")
	midStr := params.Get("mid")
	keyStr := params.Get("key")
	replyTypeStr := params.Get("type")

	// param check
	if areaStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// default param
	if id, err = strconv.ParseInt(idStr, 10, 64); err != nil {
		id = 0
	}
	if tpid, err = strconv.ParseInt(tpidStr, 10, 64); err != nil {
		tpid = 0
	}
	if oid, err = strconv.ParseInt(oidStr, 10, 64); err != nil {
		oid = 0
	}
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		mid = 0
	}
	if replyType, err = strconv.ParseInt(replyTypeStr, 10, 64); err != nil {
		replyType = 0
	}
	if keyStr == "" {
		keys = make([]string, 0)
	} else {
		keys = append(keys, strings.Split(keyStr, "|")...)
	}
	// handle req
	resData.MSG, resData.Level, resData.TypeID, resData.Hit, resData.Limit, resData.AI, err = svc.Filter(c, areaStr, msgStr, tpid, id, oid, mid, keys, int8(replyType))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	// resp
	if resData.TypeID == nil {
		resData.TypeID = []int64{}
	}
	if resData.Hit == nil {
		resData.Hit = []string{}
	}
	c.JSON(resData, nil)
}

func mfilter(c *bm.Context) {
	var (
		err       error
		id        int64
		tpid      int64
		oid       int64
		mid       int64
		replyType int64
		keys      []string
		resData   []*model.HTTPFilterRes
	)
	params := c.Request.Form
	msgsStr := params["msg"]
	areaStr := params.Get("area")
	idStr := params.Get("id")
	oidStr := params.Get("oid")
	midStr := params.Get("mid")
	tpidStr := params.Get("tpid")
	keyStr := params.Get("key")
	replyTypeStr := params.Get("type")

	if keyStr == "" {
		keys = []string{}
	} else {
		keys = append(keys, strings.Split(keyStr, "|")...)
	}

	if areaStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	if id, err = strconv.ParseInt(idStr, 10, 64); err != nil {
		id = 0
	}
	if tpid, err = strconv.ParseInt(tpidStr, 10, 64); err != nil {
		tpid = 0
	}
	if oid, err = strconv.ParseInt(oidStr, 10, 64); err != nil {
		oid = 0
	}
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		mid = 0
	}
	if replyType, err = strconv.ParseInt(replyTypeStr, 10, 64); err != nil {
		replyType = 0
	}
	resData, err = svc.HTTPMultiFilter(c, areaStr, msgsStr, tpid, id, oid, mid, keys, int8(replyType))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if resData == nil {
		resData = []*model.HTTPFilterRes{}
	}
	c.JSON(resData, nil)
}

func areaMfilter(c *bm.Context) {
	var (
		err     error
		tpid    int64
		resData []*model.HTTPAreaFilterRes
	)
	params := c.Request.Form
	msgsStr := params["msg"]
	areaStr := params.Get("area")
	tpidStr := params.Get("tpid")

	if areaStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tpid, err = strconv.ParseInt(tpidStr, 10, 64); err != nil {
		tpid = 0
	}
	resData, err = svc.HTTPMultiAreaFilter(c, areaStr, msgsStr, tpid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if resData == nil {
		resData = []*model.HTTPAreaFilterRes{}
	}
	c.JSON(resData, nil)
}

func article(c *bm.Context) {
	var (
		err  error
		tpid int64
		data []string
	)
	params := c.Request.Form
	msgStr := params.Get("msg")
	areaStr := params.Get("area")
	tpidStr := params.Get("tpid")

	if areaStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	if tpid, err = strconv.ParseInt(tpidStr, 10, 16); err != nil {
		tpid = 0
	}
	data, err = svc.Article(c, areaStr, msgStr, tpid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if data == nil {
		data = []string{}
	}
	c.JSON(data, nil)
}

func hit(c *bm.Context) {
	var (
		err  error
		tpid int64
		data []string
	)
	params := c.Request.Form
	msgStr := params.Get("msg")
	areaStr := params.Get("area")
	tpidStr := params.Get("tpid")

	if areaStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	if tpid, err = strconv.ParseInt(tpidStr, 10, 16); err != nil {
		tpid = 0
	}
	data, err = svc.Hit(c, areaStr, msgStr, tpid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if data == nil {
		data = []string{}
	}
	c.JSON(data, nil)
}

func hitV3(c *bm.Context) {
	var (
		err   error
		tpid  int64
		data  []*model.HitRes
		level int
	)
	params := c.Request.Form
	msgStr := params.Get("msg")
	areaStr := params.Get("area")
	tpidStr := params.Get("tpid")
	levelStr := params.Get("level")

	if areaStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	if levelStr != "" {
		level, err = strconv.Atoi(levelStr)
		if err != nil {
			c.JSON(nil, err)
			return
		}
	}

	if tpid, err = strconv.ParseInt(tpidStr, 10, 16); err != nil {
		tpid = 0
	}
	data, err = svc.HitV3(c, areaStr, msgStr, tpid, int8(level))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if data == nil {
		data = []*model.HitRes{}
	}
	c.JSON(data, nil)
}
