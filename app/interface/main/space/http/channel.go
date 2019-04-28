package http

import (
	"strconv"

	"go-common/app/interface/main/space/conf"
	"go-common/app/interface/main/space/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

func channel(c *bm.Context) {
	var (
		vmid, mid, cid int64
		isGuest        bool
		err            error
	)
	params := c.Request.Form
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	vmidStr := params.Get("mid")
	cidStr := params.Get("cid")
	guestStr := params.Get("guest")
	if vmid, err = strconv.ParseInt(vmidStr, 10, 64); err != nil || (vmid <= 0 && mid <= 0) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if cid, err = strconv.ParseInt(cidStr, 10, 64); err != nil || cid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if guestStr != "" {
		if isGuest, err = strconv.ParseBool(guestStr); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if !isGuest && vmid > 0 && mid != vmid {
		mid = vmid
	}
	c.JSON(spcSvc.Channel(c, mid, cid))
}

func channelIndex(c *bm.Context) {
	var (
		vmid, mid int64
		isGuest   bool
		err       error
	)
	params := c.Request.Form
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	vmidStr := params.Get("mid")
	guestStr := params.Get("guest")
	if vmid, err = strconv.ParseInt(vmidStr, 10, 64); err != nil || (vmid <= 0 && mid <= 0) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if guestStr != "" {
		if isGuest, err = strconv.ParseBool(guestStr); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if !isGuest && vmid > 0 && mid != vmid {
		isGuest = true
		mid = vmid
	}
	c.JSON(spcSvc.ChannelIndex(c, mid, isGuest))
}

func channelList(c *bm.Context) {
	var (
		vmid, mid int64
		channels  []*model.Channel
		isGuest   bool
		err       error
	)
	params := c.Request.Form
	vmidStr := params.Get("mid")
	guestStr := params.Get("guest")
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if vmid, err = strconv.ParseInt(vmidStr, 10, 64); err != nil || (vmid <= 0 && mid <= 0) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if guestStr != "" {
		if isGuest, err = strconv.ParseBool(guestStr); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if !isGuest && vmid > 0 && mid != vmid {
		isGuest = true
		mid = vmid
	}
	if channels, err = spcSvc.ChannelList(c, mid, isGuest); err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 2)
	data["count"] = len(channels)
	data["list"] = channels
	c.JSON(data, nil)
}

func addChannel(c *bm.Context) {
	var (
		mid, cid int64
		err      error
	)
	params := c.Request.Form
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	name := params.Get("name")
	intro := params.Get("intro")
	if name == "" || len([]rune(name)) > conf.Conf.Rule.MaxChNameLen {
		c.JSON(nil, ecode.ChNameToLong)
		return
	}
	if intro != "" && len([]rune(intro)) > conf.Conf.Rule.MaxChIntroLen {
		c.JSON(nil, ecode.ChIntroToLong)
		return
	}
	if cid, err = spcSvc.AddChannel(c, mid, name, intro); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(struct {
		Cid int64 `json:"cid"`
	}{Cid: cid}, nil)
}

func editChannel(c *bm.Context) {
	var (
		mid, cid int64
		err      error
	)
	params := c.Request.Form
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	cidStr := params.Get("cid")
	name := params.Get("name")
	intro := params.Get("intro")
	if cid, err = strconv.ParseInt(cidStr, 10, 64); err != nil || cid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if name == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len([]rune(name)) > conf.Conf.Rule.MaxChNameLen {
		c.JSON(nil, ecode.ChNameToLong)
		return
	}
	if intro != "" && len([]rune(intro)) > conf.Conf.Rule.MaxChIntroLen {
		c.JSON(nil, ecode.ChIntroToLong)
		return
	}
	c.JSON(nil, spcSvc.EditChannel(c, mid, cid, name, intro))
}

func delChannel(c *bm.Context) {
	var (
		mid, cid int64
		err      error
	)
	params := c.Request.Form
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	cidStr := params.Get("cid")
	if cid, err = strconv.ParseInt(cidStr, 10, 64); err != nil || cid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, spcSvc.DelChannel(c, mid, cid))
}

func channelVideo(c *bm.Context) {
	var (
		vmid, mid, cid int64
		pn, ps         int
		isGuest, order bool
		channelDetail  *model.ChannelDetail
		err            error
	)
	params := c.Request.Form
	vmidStr := params.Get("mid")
	cidStr := params.Get("cid")
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	guestStr := params.Get("guest")
	orderStr := params.Get("order")
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if vmid, err = strconv.ParseInt(vmidStr, 10, 64); err != nil || (vmid <= 0 && mid <= 0) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if cid, err = strconv.ParseInt(cidStr, 10, 64); err != nil || cid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn, err = strconv.Atoi(pnStr); err != nil || pn < 1 {
		pn = 1
	}
	if ps, err = strconv.Atoi(psStr); err != nil || ps < 1 || ps > conf.Conf.Rule.MaxChArcsPs {
		ps = conf.Conf.Rule.MaxChArcsPs
	}
	if guestStr != "" {
		if isGuest, err = strconv.ParseBool(guestStr); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if !isGuest && vmid > 0 && mid != vmid {
		isGuest = true
		mid = vmid
	}
	if orderStr != "" {
		if order, err = strconv.ParseBool(orderStr); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if channelDetail, err = spcSvc.ChannelVideos(c, mid, cid, pn, ps, isGuest, order); err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 2)
	page := map[string]int{
		"num":   pn,
		"size":  ps,
		"count": channelDetail.Count,
	}
	data["page"] = page
	data["list"] = channelDetail
	c.JSON(data, nil)
}

func addChannelVideo(c *bm.Context) {
	var (
		mid, cid int64
		aids     []int64
		err      error
	)
	params := c.Request.Form
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	cidStr := params.Get("cid")
	if cid, err = strconv.ParseInt(cidStr, 10, 64); err != nil || cid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aidStr := params.Get("aids")
	if aidStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if aids, err = xstr.SplitInts(aidStr); err != nil || len(aids) == 0 || len(aids) > conf.Conf.Rule.MaxChArcAddLimit {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aidMap := make(map[int64]int64, len(aids))
	for _, aid := range aids {
		aidMap[aid] = aid
	}
	if len(aidMap) < len(aids) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(spcSvc.AddChannelArc(c, mid, cid, aids))
}

func delChannelVideo(c *bm.Context) {
	var (
		mid, cid, aid int64
		err           error
	)
	params := c.Request.Form
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	cidStr := params.Get("cid")
	if cid, err = strconv.ParseInt(cidStr, 10, 64); err != nil || cid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aidStr := params.Get("aid")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, spcSvc.DelChannelArc(c, mid, cid, aid))
}

func sortChannelVideo(c *bm.Context) {
	var (
		mid, cid, aid int64
		orderNum      int
		err           error
	)
	params := c.Request.Form
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	cidStr := params.Get("cid")
	aidStr := params.Get("aid")
	toStr := params.Get("to")
	if cid, err = strconv.ParseInt(cidStr, 10, 64); err != nil || cid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if orderNum, err = strconv.Atoi(toStr); err != nil || orderNum < 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, spcSvc.SortChannelArc(c, mid, cid, aid, orderNum))
}

func checkChannelVideo(c *bm.Context) {
	var (
		mid, cid int64
		err      error
	)
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	cidStr := c.Request.Form.Get("cid")
	if cid, err = strconv.ParseInt(cidStr, 10, 64); err != nil || cid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, spcSvc.CheckChannelVideo(c, mid, cid))
}
