package http

import (
	"strconv"
	"time"

	"go-common/app/interface/main/tag/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

func addSub(c *bm.Context) {
	var (
		err     error
		ok      bool
		tids    []int64
		addTids []int64
		tMap    map[int64]struct{}
	)
	mid, _ := c.Get("mid")
	params := c.Request.Form
	tidsStr := params.Get("tag_id")
	if tids, err = xstr.SplitInts(tidsStr); err != nil || len(tids) <= 0 {
		log.Error("tids is nil")
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// filter the same tid
	tMap = make(map[int64]struct{})
	for _, tid := range tids {
		if tid < 1 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if _, ok = tMap[tid]; !ok {
			tMap[tid] = struct{}{}
			addTids = append(addTids, tid)
		}
	}
	if err = svr.AddSub(c, mid.(int64), addTids, time.Now()); err != nil {
		log.Error("tagSubSvr.AddSub(%d, %d) error(%v)", mid, addTids, err)
	}
	c.JSON(nil, err)
}

func cancelSub(c *bm.Context) {
	var (
		err    error
		tid    int64
		params = c.Request.Form
	)
	mid, _ := c.Get("mid")
	tidStr := params.Get("tag_id")
	if tid, err = strconv.ParseInt(tidStr, 10, 64); err != nil || tid < 1 {
		log.Error("strconv.ParseInt(%s) error(%v)", tidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = svr.CancelSub(c, tid, mid.(int64), time.Now()); err != nil {
		log.Error("tagSubSvr.CancelSub(%d, %d) error(%v)", tid, mid, err)
	}
	c.JSON(nil, err)
}

func subTags(c *bm.Context) {
	var (
		err    error
		pn     int
		ps     int
		total  int
		order  int
		mid    int64
		vmid   int64
		ts     []*model.Tag
		params = c.Request.Form
	)
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	vmidStr := params.Get("vmid")
	orderStr := params.Get("order")
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if pn, err = strconv.Atoi(pnStr); err != nil || pn < 1 {
		pn = 1
	}
	if ps, err = strconv.Atoi(psStr); err != nil || ps < 1 || ps > model.SubTagMaxNum {
		ps = model.SubTagMaxNum
	}
	if vmid, err = strconv.ParseInt(vmidStr, 10, 64); err != nil || vmid < 0 {
		vmid = 0
	}
	if orderStr != "" {
		if order, err = strconv.Atoi(orderStr); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	} else {
		order = -1
	}
	// order -1:desc 1:esc
	if order == 1 || order == -1 {
		if ts, total, err = svr.SubTags(c, mid, vmid, pn, ps, order); err != nil {
			c.JSON(nil, err)
			return
		}
	} else {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(ts) == 0 {
		ts = []*model.Tag{}
	}
	data := make(map[string]interface{}, 2)
	data["total"] = total
	data["data"] = ts
	c.JSONMap(data, nil)
}

var _emptySubArc = []*model.SubArcs{}

// subArcs get new arcs of subscribed tag .
func subArcs(c *bm.Context) {
	var (
		err    error
		mid    int64
		vmid   int64
		as     []*model.SubArcs
		params = c.Request.Form
	)
	vmidStr := params.Get("vmid")
	// check params
	if vmid, err = strconv.ParseInt(vmidStr, 10, 64); err != nil || vmid < 0 {
		vmid = 0
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if mid == 0 && vmid == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// service
	if as, err = svr.SubArcs(c, mid, vmid); err != nil {
		c.JSON(nil, err)
		return
	}
	if len(as) == 0 {
		as = _emptySubArc
	}
	c.JSON(as, nil)
}

func customSortTags(c *bm.Context) {
	var (
		err           error
		mid, vmid     int64
		tp, order     int
		ps, pn, total int
		params        = c.Request.Form
	)
	tpStr := params.Get("type")
	orderStr := params.Get("order")
	vmidStr := params.Get("vmid")
	psStr := params.Get("ps")
	pnStr := params.Get("pn")
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	vmid, _ = strconv.ParseInt(vmidStr, 10, 64)
	if vmid > 0 {
		mid = vmid
	}
	if mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tp, err = strconv.Atoi(tpStr); err != nil {
		log.Error("strconv.Atoi(%s) error(%v)", tpStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn, err = strconv.Atoi(pnStr); err != nil || pn < 1 {
		pn = 1
	}
	if ps, err = strconv.Atoi(psStr); err != nil || ps < 1 || ps > model.SubTagMaxNum {
		ps = model.SubTagMaxNum
	}
	order, _ = strconv.Atoi(orderStr)
	if order != model.SortOrderASC {
		order = model.SortOrderDESC
	}
	cst, sst, total, err := svr.CustomSubTags(c, mid, order, tp, ps, pn)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 3)
	data["page"] = map[string]int{
		"num":   pn,
		"size":  ps,
		"total": total,
	}
	data["custom"] = cst
	data["standard"] = sst
	c.JSON(data, nil)
}

func upCustomSortTags(c *bm.Context) {
	var (
		err    error
		mid    int64
		tp     int
		tids   []int64
		params = c.Request.Form
	)
	tidsStr := params.Get("tids")
	tpStr := params.Get("type")
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tidsStr != "" {
		if tids, err = xstr.SplitInts(tidsStr); err != nil {
			log.Error("xstr.SplitInts(%s) error(%v)", tidsStr, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if tp, err = strconv.Atoi(tpStr); err != nil {
		log.Error(" strconv.Atoi(%s) error(%v)", tpStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(tids) > model.MaxTopicSortNum {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svr.UpCustomSubTags(c, mid, tids, tp))
}
