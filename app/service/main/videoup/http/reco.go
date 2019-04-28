package http

import (
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

// Recos fn
func Recos(c *bm.Context) {
	params := c.Request.Form
	aidStr := params.Get("aid")
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	recos, err := vdpSvc.Recos(c, aid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(recos, nil)
}

// RecoUpdate fn
func RecoUpdate(c *bm.Context) {
	params := c.Request.Form
	recoIdsStr := params.Get("ids")
	aidStr := params.Get("aid")
	midStr := params.Get("mid")
	mid, _ := strconv.ParseInt(midStr, 10, 64)
	if mid <= 0 {
		log.Error("http.archivesByMid  mid(%d) <=0 ", mid)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	recoIds := []int64{}
	if len(recoIdsStr) > 0 {
		recoIds, err = xstr.SplitInts(recoIdsStr)
	}
	if err != nil {
		log.Error("idsStr splitInts error(%v) | recoIdsStr(%s)", err, recoIdsStr)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(recoIds) > 8 {
		log.Error("length of idsStr has over Max 8, error(%v) | recoIdsStr(%s)", err, recoIdsStr)
		c.JSON(nil, ecode.CreativeRecommendOverMax)
		return
	}
	err = vdpSvc.RecoUpdate(c, mid, aid, recoIds)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}
