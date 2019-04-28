package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"go-common/app/service/main/antispam/dao"
	"go-common/app/service/main/antispam/model"
	"go-common/app/service/main/antispam/util"

	"go-common/library/log"
)

const (
	// ThreeMonths .
	ThreeMonths = 60 * 60 * 24 * 120
)

// GetKeywordsByCond get keywords by condition.
func (s *SvcImpl) GetKeywordsByCond(ctx context.Context, cond *Condition) (ks []*model.Keyword, total int64, err error) {
	daoKs, total, err := s.KeywordDao.GetByCond(ctx, ToDaoCond(cond))
	if err == dao.ErrResourceNotExist {
		return []*model.Keyword{}, 0, nil
	}
	if err != nil {
		return nil, 0, err
	}
	ks = ToModelKeywords(daoKs)
	for _, k := range ks {
		k.SenderCounts, _ = s.antiDao.CntSendersCache(ctx, k.ID)
	}
	return ks, total, nil
}

// DeleteKeywords delete keywords and all caches related to them
func (s *SvcImpl) DeleteKeywords(ctx context.Context, ids []int64) ([]*model.Keyword, error) {
	ks, err := s.getKeywordByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	idsNeedDelete := make([]int64, 0)
	keywordsNeedDelete := make([]*model.Keyword, 0)
	for _, k := range ks {
		if k.State == model.StateDeleted {
			continue
		}
		k.State = model.StateDeleted
		idsNeedDelete = append(idsNeedDelete, k.ID)
		keywordsNeedDelete = append(keywordsNeedDelete, k)
	}
	if len(idsNeedDelete) > 0 {
		if err := s.antiDao.DelKeywordRelatedCache(ctx, keywordsNeedDelete); err != nil {
			log.Error("s.antiDao.DelKeywordRelatedCache(%+v) error(%v)", keywordsNeedDelete, err)
			return nil, err
		}
		daoKs, err := s.KeywordDao.DeleteByIDs(ctx, idsNeedDelete)
		if err != nil {
			return nil, err
		}
		return ToModelKeywords(daoKs), nil
	}
	return []*model.Keyword{}, nil
}

// GetSenderIDsByKeywordID query keyword's sender list by keyword's id
func (s *SvcImpl) GetSenderIDsByKeywordID(ctx context.Context, id int64) (*model.SenderList, error) {
	senders, err := s.antiDao.AllSendersCache(ctx, id)
	if err != nil {
		log.Error("%v", err)
		return nil, err
	}
	senderIDs := make([]int64, len(senders))
	for i, sender := range senders {
		sid, err := strconv.ParseInt(sender, 10, 64)
		if err != nil {
			return nil, err
		}
		senderIDs[i] = sid
	}
	return &model.SenderList{SenderIDs: senderIDs, Counts: len(senders)}, nil
}

// GetKeywordsByOffsetLimit query keywords by id range(offset, limit)
func (s *SvcImpl) GetKeywordsByOffsetLimit(ctx context.Context, cond *Condition) ([]*model.Keyword, error) {
	ks, err := s.KeywordDao.GetByOffsetLimit(ctx, ToDaoCond(cond))
	if err != nil {
		return nil, err
	}
	return ToModelKeywords(ks), nil
}

// GetKeywordByID .
func (s *SvcImpl) GetKeywordByID(ctx context.Context, id int64) (*model.Keyword, error) {
	k, err := s.KeywordDao.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return ToModelKeyword(k), nil
}

// OpKeyword perform update on keywords, including keyword's tag change or delete keyword.
func (s *SvcImpl) OpKeyword(ctx context.Context, id int64, newTag string) (*model.Keyword, error) {
	k, err := s.GetKeywordByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if k.State == model.StateDeleted {
		return nil, ErrResourceNotExist
	}
	if k.Tag == newTag {
		return k, nil
	}
	k.Tag = newTag
	if newTag == model.KeywordTagWhite {
		err1 := s.antiDao.DelCountRelatedCache(ctx, k)
		if err1 != nil {
			log.Error("s.antiDao.DelCountRelatedCache(%+v), error(%v)", k, err1)
			return nil, err1
		}
	}
	dk, err := s.KeywordDao.Update(ctx, ToDaoKeyword(k))
	if err != nil {
		return nil, err
	}
	return ToModelKeyword(dk), nil
}

// IncrKeywordHitCount increase keyword's hit_counts in cache,
// and persist to db only if hit_counts % 2 equals to 0
func (s *SvcImpl) IncrKeywordHitCount(ctx context.Context, k *model.Keyword) (int64, error) {
	totalCounts, err := s.antiDao.IncrTotalLimitCache(ctx, k.ID)
	if err != nil {
		log.Error("%v", err)
		return 0, err
	}
	if err := s.antiDao.TotalLimitExpire(ctx, k.ID, ThreeMonths); err != nil {
		return 0, err
	}
	if totalCounts%2 == 0 {
		k.HitCounts += 2
		s.KeywordDao.Update(ctx, ToDaoKeyword(k))
	} else {
		k.HitCounts++
	}
	return k.HitCounts, nil
}

// ExpireKeyword clean the keywords which satify the following condition:
// 1. hit_counts <= 3
// 2. scan range from "one month ago - 5 day" to "one month ago"
func (s *SvcImpl) ExpireKeyword(ctx context.Context, dbLimit int64) error {
	until := time.Now().AddDate(0, -1, 0)
	start := until.AddDate(0, 0, -5)
	cond := &Condition{
		State: model.StateDefault,
		Tags: []string{
			model.KeywordTagDefaultLimit,
			model.KeywordTagRestrictLimit,
			model.KeywordTagWhite,
		},
		StartTime: &start,
		EndTime:   &until,
		HitCounts: "3",
		Pagination: &util.Pagination{
			CurPage: 1,
			PerPage: dbLimit,
		},
	}
	ks, err := s.KeywordDao.GetRubbish(ctx, ToDaoCond(cond))
	if err != nil {
		return err
	}
	needExpireIDs := make([]int64, 0)
	for _, k := range ks {
		needExpireIDs = append(needExpireIDs, k.ID)
	}
	_, err = s.DeleteKeywords(ctx, needExpireIDs)
	return err
}

func (s *SvcImpl) persistSenderIDs(ctx context.Context, keywordID, senderID int64) error {
	totalCounts, err := s.antiDao.ZaddSendersCache(ctx, keywordID, time.Now().UnixNano(), senderID)
	if err != nil {
		return err
	}
	if totalCounts <= s.Option.MaxSenderNum {
		return nil
	}
	extraCounts := totalCounts - s.Option.MaxSenderNum
	senderIDs, err := s.antiDao.SendersCache(ctx, keywordID, 0, extraCounts)
	if err != nil {
		log.Error("s.antiDao.SendersCache(%d,%d,%d)%v", keywordID, 0, extraCounts, err)
		return err
	}
	if len(senderIDs) != int(extraCounts) {
		log.Warn("got wrong number of senderIDs:keywordID(%d), want senderIDs(%v), length(%d), got(%d)",
			keywordID, senderIDs, len(senderIDs), extraCounts)
	}
	for _, sid := range senderIDs {
		ret, err := s.antiDao.ZremSendersCache(ctx, keywordID, sid)
		if err != nil {
			log.Error("%v", err)
			return err
		}
		if ret != 1 {
			err = errors.New("fail to remove senederID from senderID list")
			log.Error("%v", err)
			return err
		}
	}
	return nil
}

// PersistKeyword persist catched keyword
func (s *SvcImpl) PersistKeyword(ctx context.Context, catchedKeyword *model.Keyword) (*model.Keyword, error) {
	keyword, err := s.getKeywordByAreaAndContent(ctx, catchedKeyword.Area, catchedKeyword.Content)
	if err != nil {
		insertedKeyword, err := s.insertKeyword(ctx, catchedKeyword)
		if err != nil {
			return nil, err
		}
		return insertedKeyword, nil
	}
	if keyword.State == model.StateDeleted {
		// the keyword was deleted before,
		// now it's hit again, restore to init state
		keyword.State = model.StateDefault
		keyword.OriginContent = catchedKeyword.OriginContent
		keyword.RegexpName = catchedKeyword.RegexpName
		keyword.CTime = catchedKeyword.CTime
		keyword.Tag = catchedKeyword.Tag
		keyword.HitCounts = 1

		k, err := s.KeywordDao.Update(ctx, ToDaoKeyword(keyword))
		if err != nil {
			return nil, err
		}
		return ToModelKeyword(k), nil
	}
	s.IncrKeywordHitCount(ctx, keyword)
	return keyword, nil
}

func (s *SvcImpl) insertKeyword(ctx context.Context, k *model.Keyword) (*model.Keyword, error) {
	res, err := s.KeywordDao.Insert(ctx, ToDaoKeyword(k))
	if err != nil {
		return nil, err
	}
	return ToModelKeyword(res), nil
}

func (s *SvcImpl) getKeywordByAreaAndContent(ctx context.Context,
	area, content string) (*model.Keyword, error) {
	k, err := s.KeywordDao.GetByAreaAndContent(ctx,
		ToDaoCond(&Condition{
			Area:     area,
			Contents: []string{content},
		}))
	if err != nil {
		return nil, err
	}
	return ToModelKeyword(k), nil
}

func (s *SvcImpl) getKeywordByIDs(ctx context.Context,
	ids []int64) ([]*model.Keyword, error) {
	ks, err := s.KeywordDao.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	return ToModelKeywords(ks), err
}
