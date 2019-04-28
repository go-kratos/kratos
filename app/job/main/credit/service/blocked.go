package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"go-common/app/job/main/credit/model"
	blkmdl "go-common/app/service/main/member/model/block"
	"go-common/library/log"
	xtime "go-common/library/time"
)

// InvalidJury invalid juryer.
func (s *Service) InvalidJury(c context.Context, nwMsg []byte, oldMsg []byte) (err error) {
	mr := &model.BlockedInfo{}
	if err = json.Unmarshal(nwMsg, mr); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(nwMsg), err)
		return
	}
	if err = s.dao.InvalidJury(c, model.JuryBlocked, mr.UID); err != nil {
		log.Error("s.dao.InvalidJury(%d %d) error(%v)", model.JuryBlocked, mr.UID, err)
	}
	return
}

// UnBlockAccount unblock account.
func (s *Service) UnBlockAccount(c context.Context, nwMsg []byte, oldMsg []byte) (err error) {
	nMR := &model.BlockedInfo{}
	if err = json.Unmarshal(nwMsg, nMR); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(nwMsg), err)
		return
	}
	oMR := &model.BlockedInfo{}
	if err = json.Unmarshal(oldMsg, oMR); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(oldMsg), err)
		return
	}
	if int8(oMR.Status) != model.StatusOpen {
		return
	}
	if int8(nMR.Status) != model.StatusClose {
		return
	}
	var id int64
	if id, err = s.dao.BlockedInfoID(c, nMR.UID); err != nil {
		log.Error("s.dao.BlockedInfoID(%+v) error(%v)", oMR, err)
		return
	}
	if id != nMR.ID {
		log.Warn("databus id(%d) do uid(%d) unblocked info(%d) not right!", nMR.ID, nMR.UID, id)
		return
	}
	if err = s.dao.UnBlockAccount(c, oMR); err != nil {
		log.Error("s.dao.UnBlockAccount(%+v) error(%v)", oMR, err)
	}
	return
}

// CheckBlock check user block state
func (s *Service) CheckBlock(c context.Context, mid int64) (ok bool, err error) {
	var block *blkmdl.RPCResInfo
	if block, err = s.memRPC.BlockInfo(c, &blkmdl.RPCArgInfo{MID: mid}); err != nil {
		log.Error("s.memRPC.BlockInfo(%d) error(%+v)", err)
		return
	}
	status := int8(block.BlockStatus)
	if status == model.BlockStatusOn {
		log.Warn("mid(%d) in blocked", mid)
		return
	}
	if status == model.BlockStatusForever {
		log.Warn("mid(%d) in blocked forever", mid)
		return
	}
	ok = true
	return
}

// NotifyBlockAnswer notify block answer status
func (s *Service) NotifyBlockAnswer(c context.Context, nwMsg []byte) (err error) {
	mr := &model.BlockLabourAnswerLog{}
	if err = json.Unmarshal(nwMsg, mr); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(nwMsg), err)
		return
	}
	if mr.Score < 100 {
		log.Warn("mid(%d) answer score(%d) lt 100", mr.Score, mr.Score)
		return
	}
	var ts time.Time
	if ts, err = time.ParseInLocation(model.TimeFormatSec, mr.CTime, time.Local); err != nil {
		log.Error("time.ParseInLocation(%s) error(%v)", mr.CTime, err)
		return
	}
	key := strconv.FormatInt(mr.MID, 10)
	msg := &model.LabourAnswer{MID: mr.MID, MTime: xtime.Time(ts.Unix())}
	if err = s.labourSub.Send(c, key, msg); err != nil {
		log.Error("PubLabour.Pub(%s, %+v) error (%v)", key, msg, err)
	}
	if err = s.dao.DelAnswerStateCache(c, mr.MID); err != nil {
		log.Error("DelAnswerStateCache(%d) error (%v)", mr.MID, err)
	}
	return
}
