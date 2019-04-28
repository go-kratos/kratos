package http

import (
	"strconv"

	"go-common/app/interface/main/dm2/model"
	"go-common/app/interface/main/dm2/model/oplog"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

// editState edit dm state.
// 0:正常1:删除 2:保护 3:取消保护
func editState(c *bm.Context) {
	var (
		p     = c.Request.Form
		dmids = make([]int64, 0)
	)
	mid, _ := c.Get("mid")
	oid, err := strconv.ParseInt(p.Get("oid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tp, err := strconv.ParseInt(p.Get("type"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	state, err := strconv.ParseInt(p.Get("state"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ids, err := xstr.SplitInts(p.Get("dmids"))
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	for _, dmid := range ids {
		if dmid == 0 {
			log.Warn("dmid is zero")
			continue
		}
		dmids = append(dmids, dmid)
	}
	if len(dmids) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	switch state {
	case 0, 1:
		err = dmSvc.EditDMState(c, int32(tp), mid.(int64), oid, int32(state), dmids, oplog.SourceUp, oplog.OperatorUp)
		if err != nil {
			log.Error("dmSvc.EditDMState(mid:%d, oid:%d, state:%d, dmids:%v) error(%v)", mid.(int64), oid, state, dmids, err)
			c.JSON(nil, err)
			return
		}
		err = dmSvc.UptSearchDMState(c, dmids, oid, int32(state), int32(tp))
	case 2, 3:
		var (
			affectIds []int64
		)
		attr := model.AttrYes
		if state == 3 {
			attr = model.AttrNo
		}
		affectIds, err = dmSvc.EditDMAttr(c, int32(tp), mid.(int64), oid, model.AttrProtect, attr, dmids, oplog.SourceUp, oplog.OperatorUp)
		if err != nil {
			log.Error("dmSvc.EditDMAttr(mid:%d, oid:%d, attr:%d, dmids:%v) error(%v)", mid.(int64), oid, attr, dmids, err)
			c.JSON(nil, err)
			return
		}
		if len(affectIds) > 0 {
			err = dmSvc.UptSearchDMAttr(c, affectIds, oid, attr, int32(tp))
		}
	default:
		err = ecode.RequestErr
	}
	c.JSON(nil, err)
}

func editPool(c *bm.Context) {
	p := c.Request.Form
	mid, _ := c.Get("mid")
	oid, err := strconv.ParseInt(p.Get("oid"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tp, err := strconv.ParseInt(p.Get("type"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pool, err := strconv.ParseInt(p.Get("pool"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	dmids, err := xstr.SplitInts(p.Get("dmids"))
	if err != nil || len(dmids) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = dmSvc.EditDMPool(c, int32(tp), mid.(int64), oid, int32(pool), dmids, oplog.SourceUp, oplog.OperatorUp)
	if err != nil {
		log.Error("dmSvc.EditDMStat(oid:%d dmids:%v) error(%v)", oid, dmids, err)
		c.JSON(nil, err)
		return
	}
	err = dmSvc.UptSearchDMPool(c, dmids, oid, int32(pool), int32(tp))
	c.JSON(nil, err)
}
