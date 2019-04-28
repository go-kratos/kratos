package service

import (
	"context"

	"go-common/app/service/main/antispam/conf"
	"go-common/app/service/main/antispam/model"
	"go-common/app/service/main/antispam/util"
	"go-common/library/log"
)

const (
	// SuspOK .
	SuspOK = "ok"
	// SuspHitBlack .
	SuspHitBlack = "hit_black"
	// SuspHitRestrict .
	SuspHitRestrict = "hit_restrict"
	// SuspExceedAllowedCounts .
	SuspExceedAllowedCounts = "exceed_allowed_counts"
)

// Filter detects spam info based on different area rule
func (s *SvcImpl) Filter(ctx context.Context, ugc UserGeneratedContent) (*model.SuspiciousResp, error) {
	key, val, err := s.TrieMgr.Get(ugc.GetArea(), ugc.GetContent())
	if err != nil && err != ErrTrieNotFound {
		log.Error("%v", err)
		return nil, err
	}
	resp := &model.SuspiciousResp{
		Area:      ugc.GetArea(),
		Content:   key,
		LimitType: SuspOK,
	}
	if err == ErrTrieNotFound {
		s.pushToChan(ugc)
		return resp, nil
	}

	updateCountFn := func() {
		k, err1 := s.GetKeywordByID(ctx, val.KeywordID)
		if err1 != nil || k.State == model.StateDeleted {
			log.Error("%v", err1)
			return
		}
		if _, err = s.IncrKeywordHitCount(ctx, k); err != nil {
			log.Warn("incr keyword(id:%d) fail, error(%v)", val.KeywordID, err)
		}
		if ugc.GetSenderID() > 0 {
			if err = s.persistSenderIDs(ctx, val.KeywordID, ugc.GetSenderID()); err != nil {
				log.Warn("persistSenderIDs(sender_id: %d) fail, error(%v)", ugc.GetSenderID(), err)
			}
		}
		log.Info("before running autoWhite on keyword(%+v), limitInfo(%+v), autoWhiteConf(%+v)", k, val, conf.Conf.AutoWhite)
		if val.LimitType != model.LimitTypeBlack && k.HitCounts > conf.Conf.AutoWhite.KeywordHitCounts {
			senderCounts, _ := s.antiDao.CntSendersCache(ctx, k.ID)
			if senderCounts > conf.Conf.AutoWhite.NumOfSenders {
				senderList, err1 := s.GetSenderIDsByKeywordID(ctx, k.ID)
				if err1 != nil {
					return
				}
				if util.StdDeviation(util.Normallization(senderList.SenderIDs)) > conf.Conf.AutoWhite.Derivation {
					log.Warn("start running autoWhite on keyword(%+v), senderList(%v)", k, senderList.SenderIDs)
					if _, err = s.OpKeyword(ctx, k.ID, model.KeywordTagWhite); err != nil {
						log.Warn("auto white fail %+v", k)
						return
					}
				}
			}
		}
	}
	s.AddTask(updateCountFn)

	if val.LimitType == model.LimitTypeBlack {
		resp.LimitType = SuspHitBlack
		return resp, nil
	}
	r, err := s.GetAggregateRuleByAreaAndLimitType(ctx, ugc.GetArea(), val.LimitType)
	if err != nil {
		return nil, err
	}
	counts, err := s.antiDao.GlobalLocalLimitCache(ctx, val.KeywordID, ugc.GetOID())
	if err != nil {
		log.Error("GlobalLocalLimitCache(%d,%d) error(%v)", val.KeywordID, ugc.GetOID(), err)
		return nil, err
	}
	globalCounts, localCounts := counts[0], counts[1]
	if globalCounts < r.GlobalAllowedCounts && localCounts < r.LocalAllowedCounts {
		if ret, _ := s.antiDao.IncrGlobalLimitCache(ctx, val.KeywordID); ret == 1 {
			s.antiDao.GlobalLimitExpire(ctx, val.KeywordID, r.GlobalDurationSec)
		}
		if ret, _ := s.antiDao.IncrLocalLimitCache(ctx, val.KeywordID, ugc.GetOID()); ret == 1 {
			s.antiDao.LocalLimitExpire(ctx, val.KeywordID, ugc.GetOID(), r.LocalDurationSec)
		}
		resp.LimitType = SuspOK
		if val.LimitType == model.LimitTypeRestrict {
			resp.LimitType = SuspHitRestrict
		}
		return resp, nil
	}
	resp.LimitType = SuspExceedAllowedCounts
	return resp, nil
}
