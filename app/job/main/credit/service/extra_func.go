package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"go-common/app/job/main/credit/model"
	"go-common/library/log"
)

// UpdateVoteCount update user vote count.
func (s *Service) UpdateVoteCount(c context.Context, mr *model.Case) {
	cvs, err := s.dao.CaseVotesCID(c, mr.ID)
	if err != nil {
		log.Error("s.dao.CaseVotesCID(%d)", mr.ID)
		return
	}
	for _, cv := range cvs {
		switch {
		case mr.JudgeType == model.JudgeTypeViolate:
			if cv.Vote == model.VoteTypeDelete || cv.Vote == model.VoteTypeViolate {
				if err = s.dao.UpdateVoteRight(c, cv.MID); err != nil {
					log.Error("s.dao.UpdateVoteRight(%d)", cv.MID)
				}
			} else if cv.Vote == model.VoteTypeLegal {
				if err = s.dao.UpdateVoteTotal(c, cv.MID); err != nil {
					log.Error("s.dao.UpdateVoteTotal(%d)", cv.MID)
				}
			}
		case mr.JudgeType == model.JudgeTypeLegal:
			if cv.Vote == model.VoteTypeLegal {
				if err = s.dao.UpdateVoteRight(c, cv.MID); err != nil {
					log.Error("s.dao.UpdateVoteRight(%d)", cv.MID)
				}
			} else if cv.Vote == model.VoteTypeViolate || cv.Vote == model.VoteTypeDelete {
				if err = s.dao.UpdateVoteTotal(c, cv.MID); err != nil {
					log.Error("s.dao.UpdateVoteTotal(%d)", cv.MID)
				}
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// UpdateCache update blocked info cache.
func (s *Service) UpdateCache(c context.Context, nwMsg []byte) (err error) {
	mr := &model.BlockedInfo{}
	if err = json.Unmarshal(nwMsg, mr); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(nwMsg), err)
		return
	}
	if mr.PublishStatus == int64(model.PublishClose) || mr.Status == int64(model.StatusClose) {
		s.dao.DelBlockedInfoIdx(c, mr)
	} else {
		s.dao.AddBlockInfoIdx(c, mr)
	}
	return
}

// DelOrigin delorigin content.
func (s *Service) DelOrigin(c context.Context, cs *model.Case) {
	switch int8(cs.OriginType) {
	case model.OriginReply:
		if cs.RelationID != "" {
			args := strings.Split(cs.RelationID, "-")
			if len(args) != 3 {
				return
			}
			s.dao.DelReply(c, args[0], args[1], args[2])
		}
	case model.OriginTag:
		if cs.RelationID != "" {
			args := strings.Split(cs.RelationID, "-")
			if len(args) != 2 {
				return
			}
			s.dao.DelTag(c, args[0], args[1])
		}
	case model.OriginDM:
		if cs.RelationID != "" {
			args := strings.Split(cs.RelationID, "-")
			if len(args) != 4 {
				return
			}
			s.dao.ReportDM(c, args[2], args[1], model.DMNotifyDel)
		}
	}
}

// DeleteIdx del cache idx.
func (s *Service) DeleteIdx(c context.Context, nwMsg []byte) (err error) {
	var opinion *model.Opinion
	if err = json.Unmarshal(nwMsg, &opinion); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(nwMsg), err)
		return
	}
	s.dao.DelOpinionCache(c, opinion.Vid)
	s.dao.DelCaseIdx(c, opinion.Cid)
	s.dao.DelVoteIdx(c, opinion.Cid)
	return
}

// jugeBlockedUser
func (s *Service) jugeBlockedUser(c context.Context, mid int64, timeStr string, bType int8) (ok bool, count int64, err error) {
	var ts, nts time.Time
	if ts, err = time.ParseInLocation(model.TimeFormatSec, timeStr, time.Local); err != nil {
		log.Error("time.ParseInLocation(%s) error(%v)", timeStr, err)
		return
	}
	switch bType {
	case model.DealTimeTypeNone:
		nts = ts
	case model.DealTimeTypeDay:
		nts = ts.AddDate(0, 0, -1)
	case model.DealTimeTypeYear:
		nts = ts.AddDate(-1, 0, 0)
	}
	if count, err = s.dao.CountBlocked(c, mid, nts); err != nil {
		log.Error("s.dao.CountBlocked(%d,%s) error(%v)", mid, nts, err)
		return
	}
	if count == 0 {
		ok = true
	}
	return
}

// DelJuryInfoCache .
func (s *Service) DelJuryInfoCache(c context.Context, nwMsg []byte) (err error) {
	mr := &model.Jury{}
	if err = json.Unmarshal(nwMsg, mr); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(nwMsg), err)
		return
	}
	if err = s.dao.DelJuryInfoCache(c, mr.Mid); err != nil {
		log.Error("s.dao.DelJuryInfoCache(%d) error(%v)", mr.Mid, err)
	}
	return
}

// DelCaseVoteTopCache .
func (s *Service) DelCaseVoteTopCache(c context.Context, nwMsg []byte) (err error) {
	mr := &model.CaseVote{}
	if err = json.Unmarshal(nwMsg, mr); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(nwMsg), err)
		return
	}
	if err = s.dao.DelCaseVoteTopCache(c, mr.MID); err != nil {
		log.Error("s.dao.DelCaseVoteTopCache(%d) error(%v)", mr.MID, err)
	}
	return
}
