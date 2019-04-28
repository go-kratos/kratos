package http

import (
	"strconv"

	model "go-common/app/interface/main/credit/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

// requirement user status in apply jury.
func blockedUserCard(c *bm.Context) {
	var mid, _ = c.Get("mid")
	rq, err := creditSvc.BlockedUserCard(c, mid.(int64))
	if err != nil {
		log.Error("creditSvc.BlockedUserCard(%d) error(%v)", mid, err)
		c.JSON(nil, err)
	}
	c.JSON(rq, nil)
}

func blockedUserList(c *bm.Context) {
	var mid, _ = c.Get("mid")
	rq, err := creditSvc.BlockedUserList(c, mid.(int64))
	if err != nil {
		log.Error("creditSvc.BlockedUserList(%d) error(%v)", mid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(rq, nil)
}

func blockedInfo(c *bm.Context) {
	var idStr = c.Request.Form.Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt err(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	rq, err := creditSvc.BlockedInfo(c, id)
	if err != nil {
		log.Error("creditSvc.BlockedInfo(%d) error(%v)", id, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(rq, nil)
}

func blockedAppeal(c *bm.Context) {
	var (
		mid   int64
		err   error
		idStr = c.Request.Form.Get("id")
	)
	midI, ok := c.Get("mid")
	if ok {
		mid = midI.(int64)
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt err(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	rq, err := creditSvc.BlockedInfoAppeal(c, id, mid)
	if err != nil {
		log.Error("creditSvc.BlockedInfo(%d) error(%v)", id, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(rq, nil)
}

func blockedList(c *bm.Context) {
	var err error
	v := new(model.ArgBlockedList)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.PS <= 0 || v.PS > 10 {
		v.PS = 10
	}
	rq, err := creditSvc.BlockedList(c, v.OType, v.BType, v.PN, v.PS)
	if err != nil {
		log.Error("creditSvc.Blockedlist(%d,%d) error(%v)", v.OType, v.BType, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(rq, nil)
}

func announcementInfo(c *bm.Context) {
	idStr := c.Request.Form.Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt err(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	rq, err := creditSvc.AnnouncementInfo(c, id)
	if err != nil {
		log.Error("creditSvc.AnnouncementInfo(%d) error(%v)", id, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(rq, nil)
}

func announcementList(c *bm.Context) {
	var (
		params = c.Request.Form
		tpStr  = params.Get("tp")
		pnStr  = params.Get("pn")
		psStr  = params.Get("ps")
	)
	tp, err := strconv.ParseInt(tpStr, 10, 8)
	if err != nil {
		log.Error("strconv.ParseInt err(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pn, err := strconv.ParseInt(pnStr, 10, 64)
	if err != nil || pn < 1 {
		log.Error("strconv.ParseInt err(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ps, err := strconv.ParseInt(psStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt err(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ps < 0 || ps > 10 {
		ps = 10
	}
	rq, err := creditSvc.AnnouncementList(c, int8(tp), pn, ps)
	if err != nil {
		log.Error("creditSvc.AnnouncementList( tp %d) error(%v)", tp, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(rq, nil)
}

func blockedNumUser(c *bm.Context) {
	v := new(model.ArgBlockedNumUser)
	if err := c.Bind(v); err != nil {
		return
	}
	bn, err := creditSvc.BlockedNumUser(c, v.MID)
	if err != nil {
		log.Error("creditSvc.BlockedNumUser(%d) error(%v)", v.MID, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(bn, nil)
}

func batchPublishs(c *bm.Context) {
	v := new(model.ArgIDs)
	if err := c.Bind(v); err != nil {
		return
	}
	pubs, err := creditSvc.BatchPublishs(c, v.IDs)
	if err != nil {
		log.Error("creditSvc.BatchPublishs(%s) error(%v)", xstr.JoinInts(v.IDs), err)
		c.JSON(nil, err)
		return
	}
	c.JSON(pubs, nil)
}

func addBlockedInfo(c *bm.Context) {
	var err error
	v := new(model.ArgJudgeBlocked)
	if err = c.Bind(v); err != nil {
		return
	}
	err = creditSvc.AddBlockedInfo(c, v)
	if err != nil {
		log.Error("creditSvc.AddBlockedInfo(%+v) error(%v)", v, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func addBatchBlockedInfo(c *bm.Context) {
	var err error
	v := new(model.ArgJudgeBatchBlocked)
	if err = c.Bind(v); err != nil {
		return
	}
	err = creditSvc.AddBatchBlockedInfo(c, v)
	if err != nil {
		log.Error("creditSvc.AddBatchBlockedInfo(%+v) error(%v)", v, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func blkHistorys(c *bm.Context) {
	v := new(model.ArgHistory)
	if err := c.Bind(v); err != nil {
		return
	}
	rhs, err := creditSvc.BLKHistorys(c, v)
	if err != nil {
		log.Error("creditSvc.BLKHistorys(%+v) error(%v)", v, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(rhs, nil)
}

func batchBLKInfos(c *bm.Context) {
	v := new(model.ArgIDs)
	if err := c.Bind(v); err != nil {
		return
	}
	mbi, err := creditSvc.BatchBLKInfos(c, v.IDs)
	if err != nil {
		log.Error("creditSvc.BatchBLKInfos(%+v) error(%v)", v, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(mbi, nil)
}
