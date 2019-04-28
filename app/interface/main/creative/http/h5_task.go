package http

import (
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"

	bm "go-common/library/net/http/blademaster"
	"strings"
)

func h5TaskBind(c *bm.Context) {
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	mid, _ := midI.(int64)

	// check white list
	if task := whiteSvc.TaskWhiteList(mid); task != 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	id, err := newcomerSvc.TaskBind(c, mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"id": id,
	}, nil)
}

func h5TaskList(c *bm.Context) {
	params := c.Request.Form
	// check user
	midStr, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	mid := midStr.(int64)

	// check white list
	if task := whiteSvc.TaskWhiteList(mid); task != 1 {
		log.Warn("h5TaskList whiteSvc.TaskWhiteList mid(%d)", mid)
		c.JSON(nil, ecode.RequestErr)
		return
	}

	from := params.Get("from")
	if !strings.EqualFold(from, "ios") && !strings.EqualFold(from, "android") {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	tks, err := newcomerSvc.H5TaskList(c, mid, from)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(tks, nil)
}

func h5RewardReceive(c *bm.Context) {
	params := c.Request.Form
	ridStr := params.Get("reward_id")
	rewardTypeStr := params.Get("reward_type")

	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	mid, _ := midI.(int64)

	ip := metadata.String(c, metadata.RemoteIP)
	var (
		err        error
		rewardID   int64
		rewardType int
	)

	rewardID, err = toInt64(ridStr)
	if err != nil || rewardID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	rewardType, err = toInt(rewardTypeStr)
	if err != nil || rewardType < 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	msg, err := newcomerSvc.RewardReceive(c, mid, rewardID, int8(rewardType), ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"msg": msg,
	}, nil)
}

func h5RewardActivate(c *bm.Context) {
	params := c.Request.Form
	idStr := params.Get("id")
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	mid, _ := midI.(int64)

	id, err := toInt64(idStr)
	if err != nil || id <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	ip := metadata.String(c, metadata.RemoteIP)
	row, err := newcomerSvc.RewardActivate(c, mid, id, ip)
	if err != nil {
		c.JSON(nil, err)
		return
	}

	c.JSON(map[string]interface{}{
		"row": row,
	}, nil)
}

func h5RewardReceiveList(c *bm.Context) {
	// check user
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	mid, _ := midI.(int64)

	recs, err := newcomerSvc.RewardReceives(c, mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(recs, nil)
}
