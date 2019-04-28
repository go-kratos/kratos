package service

import (
	"context"
	"encoding/json"

	model "go-common/app/job/main/reply/model/reply"
	"go-common/library/log"
	xtime "go-common/library/time"
)

func (s *Service) actionUp(c context.Context, msg *consumerMsg) {
	var d struct {
		Op     string     `json:"op"`
		Action uint32     `json:"action"`
		Oid    int64      `json:"oid"`
		Tp     int64      `json:"tp"`
		RpID   int64      `json:"rpid"`
		MTime  xtime.Time `json:"mtime"`
	}
	if err := json.Unmarshal([]byte(msg.Data), &d); err != nil {
		log.Error("json.Unmarshal() error(%v)", err)
		return
	}
	if d.Oid == 0 || d.RpID == 0 {
		log.Error("The structure of action(%s) from rpCh was wrong", msg.Data)
		return
	}
	rp, err := s.getReplyCache(c, d.Oid, d.RpID)
	if err != nil {
		log.Error("s.getReply failed , oid(%d), RpID(%d) err(%v)", d.Oid, d.RpID, err)
		return
	}
	if rp == nil {
		log.Error("reply is nil oid(%d) RpID(%d)", d.Oid, d.RpID)
		return
	}
	var state int8
	switch {
	case d.Op == "show":
		if rp.State != model.ReplyStateHidden {
			log.Warn("reply state(%d) is not hidden", rp.State)
			return
		}
		state = model.ReplyStateNormal
	case d.Op == "hide":
		if rp.State != model.ReplyStateNormal && rp.State != model.ReplyStateGarbage && rp.State != model.ReplyStateFiltered && rp.State != model.ReplyStateFolded {
			log.Warn("reply state(%d) is not normal", rp.State)
			return
		}
		state = model.ReplyStateHidden
	case d.Op == "top_add":
		if err := s.topAdd(c, rp, d.MTime, d.Action, model.SubAttrUpperTop); err != nil {
			log.Error("s.topAdd(oid %d) err(%v)", d.Oid, err)
		}
		return
	}
	if rows, err := s.dao.Reply.UpState(c, d.Oid, d.RpID, state, d.MTime.Time()); err != nil || rows < 1 {
		log.Error("dao.Reply.Update(%d, %d, %d), error(%v) and rows==0", d.Oid, d.RpID, state, err)
	} else {
		//	if rp, err := s.dao.Mc.GetReply(nil, d.RpID); err == nil && rp != nil {
		rp.State = state
		if err = s.dao.Mc.AddReply(c, rp); err != nil {
			log.Error("s.dao.Mc.AddReply(%d, %d, %d), error(%v) and rows==0", d.Oid, d.RpID, state, err)
		}
		//	}
	}
	//	callSearchUp(rp.Oid, rp.RpID, rp.Type, false)
}
