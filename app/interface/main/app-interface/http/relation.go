package http

import (
	"strconv"

	model "go-common/app/interface/main/app-interface/model/relation"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// followings get user's following list.
func followings(c *bm.Context) {
	var (
		mid, vmid int64
		pn, ps    int
		version   uint64
		self      bool
		err       error
	)
	params := c.Request.Form
	midInter, ok := c.Get("mid")
	if ok {
		mid = midInter.(int64)
	}
	versionStr := params.Get("re_version")
	order := params.Get("order")
	if vmid, err = strconv.ParseInt(params.Get("vmid"), 10, 64); err != nil || vmid < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	self = mid == vmid
	if pn, _ = strconv.Atoi(params.Get("pn")); pn < 1 {
		pn = 1
	}
	if !self && pn > 5 {
		c.JSON(nil, ecode.RelFollowingGuestLimit)
		return
	}
	if ps, _ = strconv.Atoi(params.Get("ps")); ps < 1 || ps > 50 {
		ps = 50
	}
	if versionStr != "" {
		if version, err = strconv.ParseUint(versionStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if order != "asc" {
		order = "desc"
	}
	followings, crc32v, total, err := relSvr.Followings(c, vmid, mid, pn, ps, version, order)
	c.JSON(struct {
		List      []*model.Following `json:"list"`
		ReVersion uint32             `json:"re_version"`
		Total     int                `json:"total"`
	}{followings, crc32v, total}, err)
}

func tag(c *bm.Context) {
	var (
		mid, tid int64
		pn, ps   int
		err      error
	)
	params := c.Request.Form
	midInter, _ := c.Get("mid")
	mid = midInter.(int64)
	if tid, err = strconv.ParseInt(params.Get("tagid"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn, _ = strconv.Atoi(params.Get("pn")); pn < 1 {
		pn = 1
	}
	if ps, _ = strconv.Atoi(params.Get("ps")); ps < 1 || ps > 50 {
		ps = 50
	}
	c.JSON(relSvr.Tag(c, mid, tid, pn, ps))
}
