package service

import (
	"context"
	"sync"
	"time"

	"go-common/app/job/main/reply-feed/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

// func (s *Service) delZSet(ctx context.Context, oid int64, tp int) (err error) {
// 	var names []string
// 	s.algorithmsLock.RLock()
// 	for _, algorithm := range s.algorithms {
// 		names = append(names, algorithm.Name())
// 	}
// 	s.algorithmsLock.RUnlock()
// 	if err = s.dao.DelReplyZSetRds(ctx, names, oid, tp); err != nil {
// 		log.Error("Del ZSet from redis oid(%d) type(%d) error(%v)", oid, tp, err)
// 	}
// 	return
// }

func (s *Service) remZSet(ctx context.Context, oid int64, tp int, rpID int64) (err error) {
	var (
		names []string
	)
	s.algorithmsLock.RLock()
	for _, algorithm := range s.algorithms {
		names = append(names, algorithm.Name())
	}
	s.algorithmsLock.RUnlock()
	for _, name := range names {
		if err = s.dao.RemReplyZSetRds(ctx, name, oid, tp, rpID); err != nil {
			log.Error("Remove reply (name: %s, oid: %d, type: %d, rpID: %d) from ZSet failed.", name, oid, tp, rpID)
			return
		}
	}
	return
}

func (s *Service) upsertZSet(ctx context.Context, oid int64, tp int) {
	var (
		rpIDs []int64
		rs    []*model.ReplyStat
		err   error
		ts    int64
	)
	// 获取计时器
	if ts, err = s.dao.CheckerTsRds(ctx, oid, tp); err != nil && err != redis.ErrNil {
		// 出错不刷新，如果缓存里还没有的话刷新
		return
	} else if time.Now().Unix()-ts < s.c.RefreshTime {
		// 小于CD时间不刷新
		return
	}
	// 从reply set中取rpIDs
	ok, err := s.dao.ExpireReplySetRds(ctx, oid, tp)
	if err != nil {
		return
	}
	if ok {
		// 缓存有则从redis中取
		if rpIDs, err = s.dao.ReplySetRds(ctx, oid, tp); err != nil {
			return
		}
	} else {
		// 缓存中没有从DB中取
		if rpIDs, err = s.dao.RpIDs(ctx, oid, tp); err != nil {
			return
		}
		// 异步回源
		s.taskQ.Do(ctx, func(ctx context.Context) {
			s.setReplySetBatch(ctx, oid, tp)
		})
	}
	// 从MC中获取reply stat
	if rs, err = s.GetStatsByID(ctx, oid, tp, rpIDs); err != nil {
		return
	}
	// 重新计算分值
	rsMap, err := s.recalculateScore(ctx, rs)
	if err != nil {
		return
	}
	for name, rs := range rsMap {
		name, rs := name, rs
		s.replyListQ.Do(ctx, func(ctx context.Context) {
			s.dao.SetReplyZSetRds(ctx, name, oid, tp, rs)
		})
	}
	// 更新完后更新计时器
	if err = s.dao.SetCheckerTsRds(ctx, oid, tp); err != nil {
		log.Error("set refresh checker error (%v)", err)
	}
}

// recalculateScore recalculate all e group reply list score.
func (s *Service) recalculateScore(ctx context.Context, stats []*model.ReplyStat) (rsMap map[string][]*model.ReplyScore, err error) {
	rsMap = make(map[string][]*model.ReplyScore)
	s.algorithmsLock.RLock()
	defer s.algorithmsLock.RUnlock()
	for _, algorithm := range s.algorithms {
		wg := sync.WaitGroup{}
		rs := make([]*model.ReplyScore, len(stats))
		for i := range stats {
			j := i
			wg.Add(1)
			s.calculator.JobQueue <- func() {
				rs[j] = algorithm.Score(stats[j])
				wg.Done()
			}
		}
		wg.Wait()
		rsMap[algorithm.Name()] = rs
	}
	return
}

func (s *Service) isHot(ctx context.Context, name string, oid, rpID int64, tp int) (isHot bool, err error) {
	rpIDs, err := s.dao.RangeReplyZSetRds(ctx, name, oid, tp, 0, 5)
	if err != nil || len(rpIDs) <= 0 {
		return
	}
	rs, err := s.GetStatsByID(ctx, oid, tp, rpIDs)
	if err != nil {
		return
	}
	for _, r := range rs {
		if r.RpID == rpID && r.Like >= 3 {
			isHot = true
			return
		}
	}
	return
}
