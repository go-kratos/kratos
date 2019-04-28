package http

import (
	"strconv"

	model "go-common/app/interface/main/credit/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

func getQs(c *bm.Context) {
	mid, _ := c.Get("mid")
	qs, err := creditSvc.GetQs(c, mid.(int64))
	if err != nil {
		log.Error("creditSvc.getQs(%d) error(%v)", mid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(qs, nil)
}

func commitQs(c *bm.Context) {
	var (
		mid, _    = c.Get("mid")
		params    = c.Request.Form
		idStr     = params.Get("id")
		buvidStr  = params.Get("buvid")
		anStr     = params.Get("ans")
		labourAns = &model.LabourAns{}
		refer     = c.Request.Referer()
		ua        = c.Request.UserAgent()
	)
	ids, err := xstr.SplitInts(idStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	answer, err := xstr.SplitInts(anStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	labourAns.ID = ids
	labourAns.Answer = answer
	commitRes, err := creditSvc.CommitQs(c, mid.(int64), refer, ua, buvidStr, labourAns)
	if err != nil {
		log.Error("creditSvc.commitQs(%d) error(%v)", mid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(commitRes, nil)
}

func addQs(c *bm.Context) {
	var (
		params    = c.Request.Form
		question  = params.Get("question")
		ansStr    = params.Get("ans")
		avIDStr   = params.Get("av_id")
		statusStr = params.Get("status")
		sourceStr = params.Get("source")
	)
	ans, err := strconv.ParseInt(ansStr, 10, 64)
	if err != nil || ans < 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	avID, err := strconv.ParseInt(avIDStr, 10, 64)
	if err != nil || avID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	status, err := strconv.ParseInt(statusStr, 10, 64)
	if err != nil || status <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	source, err := strconv.ParseInt(sourceStr, 10, 64)
	if err != nil || source < 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var qs = &model.LabourQs{Question: question, Ans: ans, AvID: avID, Status: status, Source: source}
	if err = creditSvc.AddQs(c, qs); err != nil {
		log.Error("creditSvc.AddQs(%+v) error(%v)", qs, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func setQs(c *bm.Context) {
	var (
		params    = c.Request.Form
		id, _     = strconv.ParseInt(params.Get("id"), 10, 64)
		ans, _    = strconv.ParseInt(params.Get("ans"), 10, 64)
		status, _ = strconv.ParseInt(params.Get("status"), 10, 64)
	)
	if err := creditSvc.SetQs(c, id, ans, status); err != nil {
		log.Error("creditSvc.SetQs(id:%d ans:%d status:%d) error(%v)", id, ans, status, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func delQs(c *bm.Context) {
	var err error
	v := new(model.ArgDElQS)
	if err = c.Bind(v); err != nil {
		return
	}
	if err = creditSvc.DelQs(c, v.ID, v.IsDel); err != nil {
		log.Error("creditSvc.delQs(id:%d isDel:%d error(%v)", v.ID, v.IsDel, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func isAnswered(c *bm.Context) {
	v := new(model.ArgHistory)
	if err := c.Bind(v); err != nil {
		return
	}
	anState, err := creditSvc.IsAnswered(c, v.MID, v.STime)
	if err != nil {
		log.Error("creditSvc.isanswered(mid:%d mtime:%d error(%v)", v.MID, v.STime, err)
		c.JSON(nil, err)
		return
	}
	type data struct {
		Status int8 `json:"status"`
	}
	var rs data
	rs.Status = anState
	c.JSON(rs, nil)
}
