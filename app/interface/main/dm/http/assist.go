package http

import (
	"strconv"
	"strings"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

// assistBanned 添加up主屏蔽
func assistBanned(c *bm.Context) {
	var (
		err    error
		cid    int64
		dmids  []int64
		params = c.Request.Form
	)
	mid, _ := c.Get("mid")
	cid, err = strconv.ParseInt(params.Get("cid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if params.Get("dmids") == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if dmids, err = xstr.SplitInts(params.Get("dmids")); err != nil {
		log.Error("xstr.SplitInts(%s) error(%v)", params.Get("dmids"), err)
		return
	}
	err = dmSvc.AssistBanned(c, mid.(int64), cid, dmids)
	c.JSON(nil, err)
}

// assistBannedUpt 修改up主屏蔽状态
func assistBannedUpt(c *bm.Context) {
	var (
		err    error
		hash   string
		stat   int
		params = c.Request.Form
	)
	mid, _ := c.Get("mid")
	if hash = params.Get("hash"); hash == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if stat, err = strconv.Atoi(params.Get("stat")); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = dmSvc.AssistUptBanned(c, mid.(int64), hash, int8(stat))
	c.JSON(nil, err)

}

// assistDelete 协管删除弹幕
func assistDelete(c *bm.Context) {
	var (
		mid, _ = c.Get("mid")
		params = c.Request.Form
	)
	cid, err := strconv.ParseInt(params.Get("cid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if params.Get("dmids") == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	dmids, err := xstr.SplitInts(params.Get("dmids"))
	if err != nil || len(dmids) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = dmSvc.AssistDeleteDM(c, mid.(int64), cid, dmids)
	c.JSON(nil, err)
}

// assistBannedUsers 获取UP主屏蔽的用户列表
func assistBannedUsers(c *bm.Context) {
	var (
		err    error
		aid    int64
		hashes []string
		params = c.Request.Form
	)
	mid, _ := c.Get("mid")
	aid, err = strconv.ParseInt(params.Get("aid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	hashes, err = dmSvc.AssistBannedUsers(c, mid.(int64), aid)
	if err != nil {
		c.JSON(nil, err)
		log.Error("dmSvc.AssistBannedUsers(%v,%d) error(%v)", mid, aid, err)
		return
	}
	c.JSON(hashes, nil)
}

// AssistDelBanned2 批量撤销up主屏蔽
func AssistDelBanned2(c *bm.Context) {
	var (
		err    error
		aid    int64
		hashes []string
		params = c.Request.Form
	)
	mid, _ := c.Get("mid")
	aid, err = strconv.ParseInt(params.Get("aid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if params.Get("hashes") == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	hashes = strings.Split(params.Get("hashes"), ",")
	err = dmSvc.AssistDelBanned2(c, mid.(int64), aid, hashes)
	c.JSON(nil, err)
}
