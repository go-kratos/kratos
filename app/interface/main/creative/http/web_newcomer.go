package http

import (
	"go-common/app/interface/main/creative/model/newcomer"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/xstr"
)

func webTaskList(c *bm.Context) {
	params := c.Request.Form
	tyStr := params.Get("type")

	// check user
	midStr, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	mid := midStr.(int64)

	// check white list
	if task := whiteSvc.TaskWhiteList(mid); task != 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	ty, err := toInt(tyStr)
	if err != nil {
		c.JSON(nil, err)
		return
	}

	tks, err := newcomerSvc.TaskList(c, mid, int8(ty))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(tks, nil)
}

func webRewardReceive(c *bm.Context) {
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

func webRewardActivate(c *bm.Context) {
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

func webRewardReceiveList(c *bm.Context) {
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

func webTaskBind(c *bm.Context) {
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

func growAccountStateInternal(c *bm.Context) {
	params := c.Request.Form
	midsStr := params.Get("mids")
	var (
		err  error
		mids []int64
	)
	mids, err = xstr.SplitInts(midsStr)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", midsStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	for _, mid := range mids {
		pubSvc.TaskPub(mid, newcomer.MsgForGrowAccount, newcomer.MsgFinishedCount)
	}
	c.JSON(nil, nil)
}

//webTaskMakeup to compensation update task status
func webTaskMakeup(c *bm.Context) {
	params := c.Request.Form
	// check user
	midStr, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	mid := midStr.(int64)
	if !dataSvc.IsWhite(mid) {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}

	tmid := params.Get("tmid")
	tid, err := toInt64(tmid)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// check white list
	if task := whiteSvc.TaskWhiteList(tid); task != 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	if err := newcomerSvc.TaskMakeup(c, tid); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON("ok", nil)
}

// taskPubList to apply task list
func taskPubList(c *bm.Context) {
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	mid, _ := midI.(int64)

	data, err := newcomerSvc.TaskPubList(c, mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}
