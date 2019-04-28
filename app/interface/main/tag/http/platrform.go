package http

import (
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/tag/conf"
	"go-common/app/interface/main/tag/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

func platformUpBind(c *bm.Context) {
	var (
		err     error
		oid     int64
		mid     int64
		checked []string
		typeID  int64
	)
	params := c.Request.Form
	aidStr := params.Get("oid")
	midStr := params.Get("mid")
	tnamesStr := params.Get("names")
	typeStr := params.Get("type")
	if oid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || oid < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if typeID, err = strconv.ParseInt(typeStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tNames := strings.Split(tnamesStr, ",")
	for _, name := range tNames {
		if name, err = svr.CheckName(name); err == nil && name != "" {
			checked = append(checked, name)
		}
	}
	if len(checked) > conf.Conf.Tag.ArcTagMaxNum {
		checked = checked[0:conf.Conf.Tag.ArcTagMaxNum]
	}
	c.JSON(nil, svr.UpResBind(c, oid, mid, checked, int8(typeID), time.Now()))
}

func platformAdminBind(c *bm.Context) {
	var (
		err     error
		oid     int64
		mid     int64
		checked []string
		typeID  int64
	)
	params := c.Request.Form
	aidStr := params.Get("oid")
	midStr := params.Get("mid")
	tnamesStr := params.Get("names")
	typeStr := params.Get("type")
	if oid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || oid < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if typeID, err = strconv.ParseInt(typeStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tNames := strings.Split(tnamesStr, ",")
	for _, name := range tNames {
		if name, err = svr.CheckName(name); err == nil && name != "" {
			checked = append(checked, name)
		}
	}
	if len(checked) > conf.Conf.Tag.ArcTagMaxNum {
		checked = checked[0:conf.Conf.Tag.ArcTagMaxNum]
	}
	c.JSON(nil, svr.ResAdminBind(c, oid, mid, checked, int8(typeID), time.Now()))
}

func platformListTag(c *bm.Context) {
	var (
		err    error
		typeID int64
		mid    int64
		oids   []int64
		tm     map[int64][]*model.Tag
	)
	params := c.Request.Form
	oidStr := params.Get("oid")
	typeStr := params.Get("type")
	midStr := params.Get("mid")
	if oids, err = xstr.SplitInts(oidStr); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", oidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if typeID, err = strconv.ParseInt(typeStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", typeStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if midStr != "" {
		if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if tm, err = svr.ResTags(c, oids, mid, int8(typeID)); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(tm, nil)
}
