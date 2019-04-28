package service

import (
	"context"

	"go-common/app/job/main/reply-feed/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

func (s *Service) updateStat(ctx context.Context, rpID int64, stat *model.ReplyStat) {
	s.statQ.Do(ctx, func(ctx context.Context) {
		s.dao.SetReplyStatMc(ctx, stat)
	})
}

// GetStatByID 从缓存获取单条评论stat获取不到则从DB获取
func (s *Service) GetStatByID(ctx context.Context, oid int64, tp int, rpID int64) (stat *model.ReplyStat, err error) {
	stats, err := s.GetStatsByID(ctx, oid, tp, []int64{rpID})
	if err != nil {
		return
	}
	if len(stats) > 0 {
		stat = stats[0]
	} else {
		err = ecode.ReplyNotExist
		log.Error("reply not exists rpID %d", rpID)
	}
	return
}

// GetStatsByID 从缓存获取多条评论stat获取不到则从DB获取
func (s *Service) GetStatsByID(ctx context.Context, oid int64, tp int, rpIDs []int64) (rs []*model.ReplyStat, err error) {
	var (
		rsMap   map[int64]*model.ReplyStat
		missIDs []int64
		missed  map[int64]*model.ReplyStat
	)
	if rsMap, missIDs, err = s.dao.ReplyStatsMc(ctx, rpIDs); err != nil {
		return
	}
	for _, stat := range rsMap {
		rs = append(rs, stat)
	}
	if len(missIDs) > 0 {
		// miss从DB查
		if missed, err = s.getStatsByIDDB(ctx, oid, tp, missIDs); err != nil {
			rs = nil
			return
		}
		for _, stat := range missed {
			stat := stat
			rs = append(rs, stat)
			s.statQ.Do(ctx, func(ctx context.Context) {
				s.dao.SetReplyStatMc(ctx, stat)
			})
		}
	}
	return
}

// getStatsByIDDB 从数据库获取热门评论stats，这里不需要一致性，所以跨表查再聚合
func (s *Service) getStatsByIDDB(ctx context.Context, oid int64, tp int, rpIDs []int64) (rs map[int64]*model.ReplyStat, err error) {
	if len(rpIDs) == 0 {
		return
	}
	replyMap, err := s.dao.ReplyLHRCStatsByID(ctx, oid, rpIDs)
	if err != nil {
		return
	}
	reportMap, err := s.dao.ReportStatsByID(ctx, oid, rpIDs)
	if err != nil {
		return
	}
	ctime, err := s.dao.SubjectStats(ctx, oid, tp)
	if err != nil {
		return
	}
	for rpID := range replyMap {
		r, ok := reportMap[rpID]
		if ok && r != nil {
			replyMap[rpID].Report = r.Report
		}
		replyMap[rpID].SubjectTime = ctime
	}
	rs = replyMap
	return
}

// GetStatsDB 从数据库获取热门评论stats，这里不需要一致性，所以跨表查再聚合
func (s *Service) GetStatsDB(ctx context.Context, oid int64, tp int) (rs []*model.ReplyStat, err error) {
	replyMap, err := s.dao.ReplyLHRCStats(ctx, oid, tp)
	if err != nil {
		return
	}
	var RpIDs []int64
	for rpID := range replyMap {
		RpIDs = append(RpIDs, rpID)
	}
	reportMap, err := s.dao.ReportStatsByID(ctx, oid, RpIDs)
	if err != nil {
		return
	}
	ctime, err := s.dao.SubjectStats(ctx, oid, tp)
	if err != nil {
		return
	}
	for rpID := range replyMap {
		r, ok := reportMap[rpID]
		if ok && r != nil {
			replyMap[rpID].Report = r.Report
		}
		replyMap[rpID].SubjectTime = ctime
		rs = append(rs, replyMap[rpID])
	}
	return
}
