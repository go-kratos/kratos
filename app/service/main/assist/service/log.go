package service

import (
	"context"
	"time"

	"go-common/app/service/main/assist/model/assist"
	"go-common/library/ecode"
	"go-common/library/log"
)

const typeLive = 3

// Logs get assist logs.
func (s *Service) Logs(c context.Context, mid, assistMid int64, stime, etime time.Time, pn, ps int) (logs []*assist.Log, err error) {
	if assistMid > 0 && stime.IsZero() && etime.IsZero() {
		logs, err = s.ass.LogsByAssist(c, mid, assistMid, pn, ps)
	} else if assistMid == 0 && !stime.IsZero() && !etime.IsZero() {
		logs, err = s.ass.LogsByCtime(c, mid, stime, etime, pn, ps)
	} else if assistMid > 0 && !stime.IsZero() && !etime.IsZero() {
		logs, err = s.ass.LogsByAssistCtime(c, mid, assistMid, stime, etime, pn, ps)
	} else {
		logs, err = s.ass.Logs(c, mid, pn, ps)
	}
	if err != nil {
		log.Error("s.ass.Logs(%d,%d,%v,%v,%d,%d) error(%v)", mid, assistMid, stime, etime, pn, ps, err)
		return
	}
	return
}

// LogCnt count by mid, assistMid, etime,stime
func (s *Service) LogCnt(c context.Context, mid, assistMid int64, stime, etime time.Time) (cnt int64, err error) {
	if assistMid > 0 && stime.IsZero() && etime.IsZero() {
		cnt, err = s.ass.LogCntAssist(c, mid, assistMid)
	} else if assistMid == 0 && !stime.IsZero() && !etime.IsZero() {
		cnt, err = s.ass.LogCntCtime(c, mid, stime, etime)
	} else if assistMid > 0 && !stime.IsZero() && !etime.IsZero() {
		cnt, err = s.ass.LogCntAssistCtime(c, mid, assistMid, stime, etime)
	} else {
		cnt, err = s.ass.LogCnt(c, mid)
	}
	if err != nil {
		log.Error("s.ass.LogCnt(%d,%d,%v,%v,%d,%d) error(%v)", mid, assistMid, stime, etime, err)
		return
	}
	return
}

// AddLog add assist log.
func (s *Service) AddLog(c context.Context, mid, assistMid, tp, act, subID int64, objIDStr string, detail string) (err error) {
	if mid == assistMid {
		log.Info("s.ass.AddLog mid eq assistMid (%d,%d,%d,%d,%s,%d,%s)", mid, assistMid, tp, act, subID, objIDStr, detail)
		return nil
	}
	//get assist info and check field, just except live type
	if tp != typeLive {
		ar := &assist.AssistRes{Assist: 0, Allow: 0}
		if ar, err = s.Assist(c, mid, assistMid, tp); err != nil {
			log.Error("s.Assist(%d,%d,%d,%d,%s,%d,%s) error(%v)", mid, assistMid, tp, act, subID, objIDStr, detail, err)
			return
		}
		if ar.Assist == 0 {
			err = ecode.AssistNotExist
			return
		}
	}
	//check type and act
	if _, ok := assist.TypeEnum[tp]; !ok {
		err = ecode.AssistForbidType
		return
	}
	if _, ok := assist.ActEnum[act]; !ok {
		err = ecode.AssistForbidAction
		return
	}
	if _, err = s.ass.AddLog(c, mid, assistMid, tp, act, subID, objIDStr, detail); err != nil {
		log.Error("s.ass.AddLog(%d,%d,%d,%d,%s,%d,%s) error(%v)", mid, assistMid, tp, act, subID, objIDStr, detail, err)
		return
	}
	if err = s.ass.IncrLogCount(c, mid, assistMid, tp); err != nil {
		log.Error("s.ass.IncrCount(%d,%d,%d,%d,%s,%d,%s) error(%v)", mid, assistMid, tp, act, subID, objIDStr, detail, err)
		return
	}
	return
}

// CancelLog cancel this asssist action.
func (s *Service) CancelLog(c context.Context, mid, assistMid, logID int64) (err error) {
	var logInfo *assist.Log
	logInfo, err = s.ass.LogInfo(c, logID, mid, assistMid)
	if err != nil {
		err = ecode.AssistLogNotExist
		log.Error("s.ass.LogInfo(%d,%d,%d) error(%v)", logID, mid, assistMid, err)
		return
	}
	if logInfo.State == 1 {
		err = ecode.AssistAlreadyCancel
		return
	}
	rows, err := s.ass.CancelLog(c, logID, mid, assistMid)
	if err != nil {
		log.Error("s.ass.CancelLog(%d,%d,%d,%v) error(%v)", logID, mid, assistMid, rows, err)
		return
	}
	return
}

// LogInfo get one log info.
func (s *Service) LogInfo(c context.Context, id, mid, assistMid int64) (logInfo *assist.Log, err error) {
	logInfo, err = s.ass.LogInfo(c, id, mid, assistMid)
	if err != nil {
		log.Error("s.ass.LogInfo(%d,%d,%d) error(%v)", id, mid, assistMid, err)
		return
	}
	return
}

// LogObj get one log info.
func (s *Service) LogObj(c context.Context, mid, objID, tp, act int64) (logInfo *assist.Log, err error) {
	logInfo, err = s.ass.LogObj(c, mid, objID, tp, act)
	if err != nil {
		log.Error("s.ass.LogInfo(%d,%d,%d,%d) error(%v)", mid, objID, tp, act, err)
		return
	}
	return
}
