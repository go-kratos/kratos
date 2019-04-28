package service

import (
	"context"

	"go-common/app/job/main/reply-feed/model"
	"go-common/library/log"
)

// setReplySetBatch set reply set batch.
func (s *Service) setReplySetBatch(ctx context.Context, oid int64, tp int) (err error) {
	var (
		stats []*model.ReplyStat
		rpIDs []int64
	)
	// 从DB查出满足热门评论条件的评论ID
	if rpIDs, err = s.dao.RpIDs(ctx, oid, tp); err != nil || len(rpIDs) <= 0 {
		return
	}
	// 从MC或者DB中取出reply stat
	if stats, err = s.GetStatsByID(ctx, oid, tp, rpIDs); err != nil {
		return
	}
	for _, stat := range stats {
		stat := stat
		s.statQ.Do(ctx, func(ctx context.Context) {
			s.dao.SetReplyStatMc(ctx, stat)
		})
	}
	return s.dao.SetReplySetRds(ctx, oid, tp, rpIDs)
}

// addReplySet add one rpID into redis reply set.
func (s *Service) addReplySet(ctx context.Context, oid int64, tp int, rpID int64) (err error) {
	ok, err := s.dao.ExpireReplySetRds(ctx, oid, tp)
	if err != nil {
		return
	}
	if ok {
		if err = s.dao.AddReplySetRds(ctx, oid, tp, rpID); err != nil {
			return
		}
	} else {
		if err = s.setReplySetBatch(ctx, oid, tp); err != nil {
			return
		}
	}
	return
}

func (s *Service) remSet(ctx context.Context, oid, rpID int64, tp int) (err error) {
	if err = s.dao.RemReplySetRds(ctx, oid, rpID, tp); err != nil {
		log.Error("Remove rpID from set error (%v)", err)
	}
	return
}

// func (s *Service) delSet(ctx context.Context, oid int64, tp int) (err error) {
// 	if err = s.dao.DelReplySetRds(ctx, oid, tp); err != nil {
// 		log.Error("delete reply set(oid: %d, type: %d)", oid, tp)
// 	}
// 	return
// }

// func (s *Service) delReply(ctx context.Context, oid int64, tp int) {
// 	var err error
// 	if err = s.delSet(ctx, oid, tp); err != nil {
// 		s.replyListQ.Do(ctx, func(ctx context.Context) {
// 			s.delSet(ctx, oid, tp)
// 		})
// 	}
// 	if err = s.delZSet(ctx, oid, tp); err != nil {
// 		s.replyListQ.Do(ctx, func(ctx context.Context) {
// 			s.delZSet(ctx, oid, tp)
// 		})
// 	}
// }

func (s *Service) remReply(ctx context.Context, oid int64, tp int, rpID int64) {
	var err error
	if err = s.remSet(ctx, oid, rpID, tp); err != nil {
		s.replyListQ.Do(ctx, func(ctx context.Context) {
			s.remSet(ctx, oid, rpID, tp)
		})
	}
	if err = s.remZSet(ctx, oid, tp, rpID); err != nil {
		s.replyListQ.Do(ctx, func(ctx context.Context) {
			s.remZSet(ctx, oid, tp, rpID)
		})
	}
}
