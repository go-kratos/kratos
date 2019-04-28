package http

import (
	"encoding/json"
	"strconv"

	"go-common/app/interface/main/credit/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func delOpinion(c *bm.Context) {
	var (
		err     error
		cid     int64
		opid    int64
		params  = c.Request.Form
		cidStr  = params.Get("cid")
		opidStr = params.Get("opid")
	)
	if cid, err = strconv.ParseInt(cidStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt err(err) %v", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if opid, err = strconv.ParseInt(opidStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt err(err) %v", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = creditSvc.DelOpinion(c, cid, opid); err != nil {
		log.Error("creditSvc.DelOpinion(%d,%d) error(err) %v", cid, opid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func addBlockedCase(c *bm.Context) {
	var (
		err    error
		params = c.Request.Form
		data   = params.Get("data")
		bc     = make([]*model.ArgJudgeCase, 0)
	)
	if err = json.Unmarshal([]byte(data), &bc); err != nil {
		log.Error("json.Unmarshal(%s),err:%v.", data, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = creditSvc.AddBlockedCases(c, bc); err != nil {
		log.Error("creditSvc.AddBlockedCases error(err) %v", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func batchJuryInfos(c *bm.Context) {
	v := new(model.ArgMIDs)
	if err := c.Bind(v); err != nil {
		return
	}
	mbj, err := creditSvc.JuryInfos(c, v.MIDs)
	if err != nil {
		log.Error("creditSvc.JuryInfos(%+v) error(%v)", v, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(mbj, nil)
}
