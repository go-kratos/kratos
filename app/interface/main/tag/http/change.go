package http

import (
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/tag/conf"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// up主绑定
func upBind(c *bm.Context) {
	var (
		err        error
		aid        int64
		mid        int64
		checked    []string
		regionName []string
	)
	params := c.Request.Form
	aidStr := params.Get("aid")
	midStr := params.Get("mid")
	tnamesStr := params.Get("tnames")
	regionNameStr := params.Get("region_name")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || aid < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	regionNames := strings.Split(regionNameStr, ",")
	for _, name := range regionNames {
		if name, err = svr.CheckName(name); err == nil && name != "" {
			regionName = append(regionName, name)
		}
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
	c.JSON(nil, svr.UpArcBind(c, aid, mid, checked, regionName, time.Now()))
}

// admin绑定
func adminBind(c *bm.Context) {
	var (
		err        error
		aid        int64
		amid       int64
		checked    []string
		regionName []string
	)
	params := c.Request.Form
	aidStr := params.Get("aid")
	adminMidStr := params.Get("mid")
	tnamesStr := params.Get("tnames")
	regionNameStr := params.Get("region_name")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || aid < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if amid, err = strconv.ParseInt(adminMidStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", adminMidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	regionNames := strings.Split(regionNameStr, ",")
	for _, name := range regionNames {
		if name, err = svr.CheckName(name); err == nil && name != "" {
			regionName = append(regionName, name)
		}
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
	c.JSON(nil, svr.ArcAdminBind(c, aid, amid, checked, regionName, time.Now()))
}
