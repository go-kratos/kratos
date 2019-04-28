package service

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/job/main/credit/model"
	"go-common/library/log"
	xtime "go-common/library/time"
	"go-common/library/xstr"
)

// loadCase load grant case to list.
func (s *Service) loadCase() (err error) {
	c := context.TODO()
	if s.c.Judge.CaseLoadSwitch != model.StateCaseLoadOpen {
		log.Warn("case load switch off!")
		return
	}
	count, err := s.dao.TotalGrantCase(c)
	if err != nil {
		log.Error("s.dao.CaseGrantCount error(%v)", err)
		return
	}
	loadCount := s.c.Judge.CaseLoadMax - count
	if loadCount <= 0 {
		log.Info("load case full caseLoadMax(%d) listCount(%d)", s.c.Judge.CaseLoadMax, count)
		return
	}
	log.Info("need load count(%d) on already que len(%d) and max(%d)", loadCount, count, s.c.Judge.CaseLoadMax)
	mcases, err := s.dao.Grantcase(c, loadCount)
	if err != nil {
		log.Error("s.dao.Grantcase(%d) error(%v)", loadCount, err)
		return
	}
	if len(mcases) == 0 {
		log.Warn("granting case is zero!")
		return
	}
	var cids []int64
	for cid := range mcases {
		cids = append(cids, cid)
	}
	now := time.Now()
	stime := xtime.Time(now.Unix())
	etime := xtime.Time(now.Add(time.Duration(s.c.Judge.CaseGiveHours) * time.Hour).Unix())
	if err = s.dao.UpGrantCase(c, cids, stime, etime); err != nil {
		log.Error("s.dao.UpGrantCase(%s) error(%v)", xstr.JoinInts(cids), err)
		return
	}
	for _, v := range mcases {
		v.Stime = stime
		v.Etime = etime
	}
	if err = s.dao.SetGrantCase(c, mcases); err != nil {
		log.Error("s.dao.SetMIDCaseGrant(%s) error(%v)", xstr.JoinInts(cids), err)
	}
	log.Info("load cases(%s) to queue on start_time(%v) and end_time(%v)", xstr.JoinInts(cids), stime, etime)
	return
}

// DealCaseApplyReason deal with case apply reason.
func (s *Service) DealCaseApplyReason(c context.Context, nwMsg []byte) (err error) {
	mr := &model.CaseApplyModifyLog{}
	if err = json.Unmarshal(nwMsg, mr); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(nwMsg), err)
		return
	}
	var (
		aReasons []int64
		cas      []*model.CaseApplyModifyLog
	)
	if aReasons, err = s.dao.CaseApplyReasons(c, mr.CID); err != nil {
		log.Error("s.dao.CaseApplyReasons(%d) error(%v)", mr.CID, err)
		return
	}
	if len(aReasons) == 0 {
		log.Warn("apply reason(%s) is nil", xstr.JoinInts(aReasons))
		return
	}
	if cas, err = s.dao.CaseApplyReasonNum(c, mr.CID, aReasons); err != nil {
		log.Error("s.dao.CaseApplyReasonNum(%d,%s) error(%v)", mr.CID, xstr.JoinInts(aReasons), err)
		return
	}
	var (
		max    int
		reason int8
	)
	for _, v := range cas {
		if v.Num > max {
			max = v.Num
			reason = v.AReason
		}
	}
	standard := int(s.c.Judge.CaseVoteMax * int64(s.c.Judge.CaseVoteMaxPercent) / 100)
	if max < standard {
		log.Warn("apply reason num(%d) not enough standard(%d)", max, standard)
		return
	}
	if reason == mr.OReason {
		log.Warn("max reason(%d) eq orgin reason(%d)", reason, mr.OReason)
		return
	}
	var effect int64
	if effect, err = s.dao.UpBlockedCaseReason(c, mr.CID, reason); err != nil {
		log.Error("s.dao.UpBlockedCaseReason(%d,%d) error(%v)", mr.CID, reason, err)
		return
	}
	if effect <= 0 {
		log.Warn("update case_id(%d) and reason(%d) notIdempotent", mr.CID, reason)
		return
	}
	if model.ReasonToFreeze(reason) {
		if err = s.dao.UpBlockedCaseStatus(c, mr.CID, model.CaseStatusFreeze); err != nil {
			log.Error("s.dao.UpBlockedCaseStatus(%d,%d) error(%v)", mr.CID, model.CaseStatusFreeze, err)
			return
		}
	}
	if err = s.dao.AddBlockedCaseModifyLog(c, mr.CID, mr.AType, mr.OReason, mr.AReason); err != nil {
		log.Error("s.dao.AddBlockedCaseModifyLog(%d,%d,%d,%d) error(%v)", mr.CID, mr.AType, mr.OReason, mr.AReason, err)
	}
	return
}

// GrantCase  push case in list.
func (s *Service) GrantCase(c context.Context, nwMsg []byte, oldMsg []byte) (err error) {
	mr := &model.Case{}
	if err = json.Unmarshal(nwMsg, mr); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(nwMsg), err)
		return
	}
	if mr.Status != model.CaseStatusGranting {
		return
	}
	stime, err := time.ParseInLocation(time.RFC3339, mr.Stime, time.Local)
	if err != nil {
		stime, err = time.ParseInLocation("2006-01-02 15:04:05", mr.Stime, time.Local)
		if err != nil {
			log.Error("time.ParseInLocation(%s) error(%v)", mr.Stime, err)
			return
		}
		err = nil
	}
	etime, err := time.ParseInLocation(time.RFC3339, mr.Etime, time.Local)
	if err != nil {
		etime, err = time.ParseInLocation("2006-01-02 15:04:05", mr.Etime, time.Local)
		if err != nil {
			log.Error("time.ParseInLocation(%s) error(%v)", mr.Etime, err)
			return
		}
		err = nil
	}
	simCase := &model.SimCase{
		ID:         mr.ID,
		Mid:        mr.Mid,
		VoteRule:   mr.Agree,
		VoteBreak:  mr.Against,
		VoteDelete: mr.VoteDelete,
		CaseType:   mr.CaseType,
		Stime:      xtime.Time(stime.Unix()),
		Etime:      xtime.Time(etime.Unix()),
	}
	mcases := make(map[int64]*model.SimCase)
	mcases[mr.ID] = simCase
	if err = s.dao.SetGrantCase(c, mcases); err != nil {
		log.Error("s.dao.SetMIDCaseGrant(%+v) error(%v)", mr, err)
	}
	return
}

// DelGrantCase  del case in list.
func (s *Service) DelGrantCase(c context.Context, nwMsg []byte, oldMsg []byte) (err error) {
	mr := &model.Case{}
	if err = json.Unmarshal(nwMsg, mr); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(nwMsg), err)
		return
	}
	if mr.Status == model.CaseStatusDealing && mr.CaseType == model.JudeCaseTypePrivate {
		return
	}
	if mr.Status == model.CaseStatusGranting || mr.Status == model.CaseStatusDealed || mr.Status == model.CaseStatusRestart {
		return
	}
	cids := []int64{mr.ID}
	// 删除冻结和停止发放中的cid
	if err = s.dao.DelGrantCase(c, cids); err != nil {
		log.Error("s.dao.SetMIDCaseGrant(%d) error(%v)", mr.ID, err)
	}
	log.Info("cid(%d) status(%d) remove hash list on start_time(%s) and end_time(%s)", mr.ID, mr.Status, mr.Stime, mr.Etime)
	return
}

// DelCaseInfoCache del case info cache.
func (s *Service) DelCaseInfoCache(c context.Context, nwMsg []byte) (err error) {
	mr := &model.Case{}
	if err = json.Unmarshal(nwMsg, mr); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(nwMsg), err)
		return
	}
	if err = s.dao.DelCaseInfoCache(c, mr.ID); err != nil {
		log.Error("s.dao.DelCaseInfoCache(%d) error(%v)", mr.ID, err)
		return
	}
	log.Info("cid(%d) del case_info cache", mr.ID)
	return
}

// DelVoteCaseCache del vote case cache.
func (s *Service) DelVoteCaseCache(c context.Context, nwMsg []byte) (err error) {
	mr := &model.BLogCaseVote{}
	if err = json.Unmarshal(nwMsg, mr); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(nwMsg), err)
		return
	}
	if err = s.dao.DelVoteCaseCache(c, mr.MID, mr.CID); err != nil {
		log.Error("s.dao.DelCaseInfoCache(%d,%d) error(%v)", mr.MID, mr.CID, err)
		return
	}
	log.Info("mid(%d) cid(%d) del vote_case cache", mr.MID, mr.CID)
	if err = s.dao.DelJuryInfoCache(c, mr.MID); err != nil {
		log.Error("s.dao.DelJuryInfoCache(%d) error(%v)", mr.MID, err)
	}
	log.Info("mid(%d) del jury cache", mr.MID, mr.CID)
	if err = s.dao.DelCaseVoteTopCache(c, mr.MID); err != nil {
		log.Error("s.dao.DelCaseVoteTopCache(%d) error(%v)", mr.MID, err)
	}
	log.Info("mid(%d) del vote top cache", mr.MID)
	return
}

// loadDealWrongCase  deal wrong case status in cache .
func (s *Service) loadDealWrongCase() (err error) {
	var (
		cids, wcids []int64
		mcids       map[int64]int8
		c           = context.TODO()
	)
	if cids, err = s.dao.GrantCases(c); err != nil {
		log.Error("s.dao.GrantCases error(%+v)", err)
		return
	}
	if len(cids) == 0 {
		log.Warn("deal wrong (granting case is zero!)")
		return
	}
	if mcids, err = s.dao.CasesStatus(c, cids); err != nil {
		log.Error("s.dao.CasesStatus(%s) error(%+v)", xstr.JoinInts(cids), err)
		return
	}
	for cid, status := range mcids {
		if status != model.CaseStatusGranting {
			wcids = append(wcids, cid)
		}
	}
	if len(wcids) == 0 {
		return
	}
	if err = s.dao.DelGrantCase(c, wcids); err != nil {
		log.Error("s.dao.DelGrantCase(%s) error(%+v)", xstr.JoinInts(wcids), err)
		return
	}
	log.Info("load del wrong cases(%s) from queue", xstr.JoinInts(wcids))
	return
}
