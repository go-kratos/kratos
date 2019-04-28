package http

import (
	model "go-common/app/interface/main/credit/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// requirement user status in apply jury.
func requirement(c *bm.Context) {
	mid, _ := c.Get("mid")
	rq, err := creditSvc.Requirement(c, mid.(int64))
	if err != nil {
		log.Error("creditSvc.Requirement(%d) error(%v)", mid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(rq, nil)
}

// apply user apply jury.
func apply(c *bm.Context) {
	mid, _ := c.Get("mid")
	err := creditSvc.Apply(c, mid.(int64))
	if err != nil {
		log.Error("creditSvc.Apply(%d) error(%v)", mid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// jury jury user info.
func jury(c *bm.Context) {
	mid, _ := c.Get("mid")
	ui, err := creditSvc.Jury(c, mid.(int64))
	if err != nil {
		log.Error("creditSvc.Jury(%d) error(%v)", mid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(ui, nil)
}

// caseObtain jury user obtain case list.
func caseObtain(c *bm.Context) {
	mid, _ := c.Get("mid")
	v := new(model.ArgCid)
	if err := c.Bind(v); err != nil {
		return
	}
	id, err := creditSvc.CaseObtain(c, mid.(int64), v.Cid)
	if err != nil {
		log.Error("creditSvc.CaseObtain(%d) error(%v)", mid, err)
		c.JSON(nil, err)
		return
	}
	type reid struct {
		CID int64 `json:"id"`
	}
	var data reid
	data.CID = id
	c.JSON(data, nil)
}

// caseObtainByID jury user obtain case list.
func caseObtainByID(c *bm.Context) {
	mid, _ := c.Get("mid")
	v := new(model.ArgCid)
	if err := c.Bind(v); err != nil {
		return
	}
	if v.Cid == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err := creditSvc.CaseObtainByID(c, mid.(int64), v.Cid)
	if err != nil {
		log.Error("creditSvc.CaseObtain(%d) error(%v)", mid, err)
		c.JSON(nil, err)
		return
	}
	type reid struct {
		CID int64 `json:"id"`
	}
	var data reid
	data.CID = v.Cid
	c.JSON(data, nil)
}

// vote jury user vote case.
func vote(c *bm.Context) {
	mid, _ := c.Get("mid")
	v := new(model.ArgVote)
	if err := c.Bind(v); err != nil {
		return
	}
	if err := creditSvc.Vote(c, mid.(int64), v.Cid, v.Attr, v.Vote, v.AType, v.AReason, v.Content, v.Likes, v.Hates); err != nil {
		log.Error("creditSvc.Vote(%d,%+v) error(%v)", mid, v, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// voteInfo jury user vote info.
func voteInfo(c *bm.Context) {
	mid, _ := c.Get("mid")
	v := new(model.ArgCid)
	if err := c.Bind(v); err != nil {
		return
	}
	vi, err := creditSvc.VoteInfo(c, mid.(int64), v.Cid)
	if err != nil {
		log.Error("creditSvc.VoteInfo(%d %d) error(%v)", mid, v.Cid, err)
		c.JSON(nil, err)
		return
	}
	var data interface{}
	if vi == nil || vi.MID == 0 {
		data = &struct{}{}
	} else {
		data = vi
	}
	c.JSON(data, nil)
}

// caseInfo jury get case info.
func caseInfo(c *bm.Context) {
	v := new(model.ArgCid)
	if err := c.Bind(v); err != nil {
		return
	}
	ci, err := creditSvc.CaseInfo(c, v.Cid)
	if err != nil {
		log.Error("creditSvc.CaseInfo(%d) error(%v)", v.Cid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(ci, nil)
}

// juryCase jury user case info contain vote.
func juryCase(c *bm.Context) {
	mid, _ := c.Get("mid")
	v := new(model.ArgCid)
	if err := c.Bind(v); err != nil {
		return
	}
	jc, err := creditSvc.JuryCase(c, mid.(int64), v.Cid)
	if err != nil {
		log.Error("creditSvc.JuryCase(%d %d) error(%v)", mid, v.Cid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(jc, nil)
}

// spJuryCase get specific jury case info.
func spJuryCase(c *bm.Context) {
	var mid int64
	iMid, ok := c.Get("mid")
	if ok {
		mid = iMid.(int64)
	}
	v := new(model.ArgCid)
	if err := c.Bind(v); err != nil {
		return
	}
	jc, err := creditSvc.SpJuryCase(c, mid, v.Cid)
	if err != nil {
		log.Error("creditSvc.JuryCase(%d %d) error(%v)", mid, v.Cid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(jc, nil)
}

// caseList user case list.
func caseList(c *bm.Context) {
	mid, _ := c.Get("mid")
	v := new(model.ArgPage)
	if err := c.Bind(v); err != nil {
		return
	}
	cl, err := creditSvc.CaseList(c, mid.(int64), v.PN, v.PS)
	if err != nil {
		log.Error("creditSvc.CaseList(%d,%d,%d) error(%v)", mid, v.PN, v.PS, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(cl, nil)
}

// notice get jury notice.
func notice(c *bm.Context) {
	n, err := creditSvc.Notice(c)
	if err != nil {
		log.Error("creditSvc.notice error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(n, nil)
}

// reasonList get reason list.
func reasonList(c *bm.Context) {
	n, err := creditSvc.ReasonList(c)
	if err != nil {
		log.Error("creditSvc.ReasonList error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(n, nil)
}

// kpiList get kpi list.
func kpiList(c *bm.Context) {
	mid, _ := c.Get("mid")
	n, err := creditSvc.KPIList(c, mid.(int64))
	if err != nil {
		log.Error("creditSvc.KpiList error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(n, nil)
}

// voteOpinion get vote opinion.
func voteOpinion(c *bm.Context) {
	v := new(model.ArgOpinion)
	if err := c.Bind(v); err != nil {
		return
	}
	ops, count, err := creditSvc.VoteOpinion(c, v.Cid, v.PN, v.PS, v.Otype)
	if err != nil {
		log.Error("creditSvc.VoteOpinion error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(&model.OpinionRes{
		Count:   count,
		Opinion: ops,
	}, nil)
}

// caseOpinion get case opinion.
func caseOpinion(c *bm.Context) {
	v := new(model.ArgOpinion)
	if err := c.Bind(v); err != nil {
		return
	}
	ops, count, err := creditSvc.CaseOpinion(c, v.Cid, v.PN, v.PS)
	if err != nil {
		log.Error("creditSvc.CaseOpinion error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(&model.OpinionRes{
		Count:   count,
		Opinion: ops,
	}, nil)
}

// batchBLKCases get case info.
func batchBLKCases(c *bm.Context) {
	v := new(model.ArgIDs)
	if err := c.Bind(v); err != nil {
		return
	}
	cases, err := creditSvc.BatchBLKCases(c, v.IDs)
	if err != nil {
		log.Error("creditSvc.BatchBLKCases(%+v) error(%+v)", v.IDs, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(cases, nil)
}
