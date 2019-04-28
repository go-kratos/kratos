package service

import (
	"context"
	"sync"
	"time"

	"go-common/app/service/main/thumbup/model"
	"go-common/library/ecode"
	"go-common/library/queue/databus/report"
	"go-common/library/sync/errgroup"
)

// Stats .
func (s *Service) Stats(c context.Context, business string, originID int64, messageIDs []int64) (res map[int64]*model.Stats, err error) {
	var businessID int64
	if businessID, err = s.CheckBusinessOrigin(business, originID); err != nil {
		return
	}
	if originID != 0 {
		res, err = s.originStats(c, businessID, originID, messageIDs)
	} else {
		res, err = s.itemStats(c, businessID, originID, messageIDs)
	}
	if err != nil {
		return
	}
	for _, id := range messageIDs {
		if res[id] == nil {
			res[id] = &model.Stats{ID: id, OriginID: originID}
		}
	}
	return
}

func (s *Service) originStats(c context.Context, businessID, originID int64, messageIDs []int64) (res map[int64]*model.Stats, err error) {
	var (
		cache    bool
		addCache = true
	)
	if cache, err = s.dao.ExpireHashStatsCache(c, businessID, originID); err != nil {
		addCache = false
	}
	if cache {
		if res, err = s.dao.HashStatsCache(c, businessID, originID, messageIDs); err == nil {
			return
		}
		addCache = false
	}
	var stats map[int64]*model.Stats
	if stats, err = s.dao.OriginStats(c, businessID, originID); err != nil {
		return
	}
	res = make(map[int64]*model.Stats)
	for _, id := range messageIDs {
		res[id] = stats[id]
	}
	if addCache {
		s.cache.Do(c, func(c context.Context) {
			s.dao.AddHashStatsCacheMap(c, businessID, originID, stats)
		})
	}
	return
}

func (s *Service) itemStats(c context.Context, businessID, originID int64, messageIDs []int64) (res map[int64]*model.Stats, err error) {
	var (
		missed   []int64
		addCache = true
	)
	res, missed, err = s.dao.StatsCache(c, businessID, messageIDs)
	if err == nil && len(missed) == 0 {
		return
	} else if err != nil {
		res = nil
		missed = messageIDs
		addCache = false
	}
	var stats map[int64]*model.Stats
	if stats, err = s.dao.MessageStats(c, businessID, missed); err != nil {
		return
	}
	if res == nil {
		res = make(map[int64]*model.Stats)
	}
	for id, stat := range stats {
		res[id] = stat
	}
	if addCache {
		s.cache.Do(c, func(c context.Context) {
			s.dao.AddStatsCacheMap(c, businessID, stats)
		})
	}
	return
}

// AddStatsCache .
func (s *Service) AddStatsCache(c context.Context, businessID, originID int64, stat *model.Stats) error {
	if originID == 0 {
		return s.dao.AddStatsCache(c, businessID, stat)
	}
	return s.dao.AddHashStatsCache(c, businessID, originID, stat)
}

// StatsWithLike .
func (s *Service) StatsWithLike(c context.Context, business string, mid, originID int64, messageIDs []int64) (res map[int64]*model.StatsWithLike, err error) {
	if _, err = s.CheckBusinessOrigin(business, originID); err != nil {
		return
	}
	group, ctx := errgroup.WithContext(c)
	var stats map[int64]*model.Stats
	var likes map[int64]int8
	group.Go(func() (err error) {
		stats, err = s.Stats(ctx, business, originID, messageIDs)
		return
	})
	group.Go(func() (err error) {
		likes, _, err = s.HasLike(ctx, business, mid, messageIDs)
		return
	})
	if err = group.Wait(); err != nil {
		return
	}
	res = make(map[int64]*model.StatsWithLike)
	for id, stat := range stats {
		if stat == nil {
			stat = &model.Stats{}
		}
		var likeState int8
		if likes != nil {
			likeState = likes[id]
		}
		res[id] = &model.StatsWithLike{Stats: *stat, LikeState: likeState}
	}
	return
}

// MultiStatsWithLike multi stats
func (s *Service) MultiStatsWithLike(c context.Context, arg *model.MultiBusiness) (res map[string]map[int64]*model.StatsWithLike, err error) {
	group := errgroup.Group{}
	mutex := sync.Mutex{}
	res = make(map[string]map[int64]*model.StatsWithLike)
	for businessName, business := range arg.Businesses {
		oids := make(map[int64][]int64)
		for _, b := range business {
			oids[b.OriginID] = append(oids[b.OriginID], b.MessageID)
		}
		businessName := businessName
		for oid, ids := range oids {
			oid := oid
			ids := ids
			group.Go(func() (err error) {
				r, err := s.StatsWithLike(c, businessName, arg.Mid, oid, ids)
				if err != nil {
					return
				}
				mutex.Lock()
				if res[businessName] == nil {
					res[businessName] = r
				} else {
					for k, v := range r {
						res[businessName][k] = v
					}
				}
				mutex.Unlock()
				return
			})
		}
	}
	err = group.Wait()
	if err != nil && len(res) > 0 {
		err = nil
	}
	return
}

// UpdateCount  update count
func (s *Service) UpdateCount(c context.Context, business string, originID, messageID int64, likeChange, dislikeChange int64, ip, operator string) (err error) {
	var businessID int64
	if businessID, err = s.CheckBusinessOrigin(business, originID); err != nil {
		return
	}
	if operator == "" {
		err = ecode.RequestErr
		return
	}
	if likeChange == 0 && dislikeChange == 0 {
		return
	}
	const thumbupType = 171
	stats, err := s.Stats(c, business, originID, []int64{messageID})
	if err != nil {
		return
	}
	var likeCount, dislikeCount int64
	if stats[messageID] != nil {
		likeCount = stats[messageID].Likes
		dislikeCount = stats[messageID].Dislikes
	}
	if likeCount+likeChange < 0 {
		likeChange = -likeCount
	}
	if dislikeCount+dislikeChange < 0 {
		dislikeChange = -dislikeCount
	}
	if err = s.dao.UpdateCount(c, businessID, originID, messageID, likeChange, dislikeChange); err != nil {
		return
	}
	// add log
	_ = report.Manager(&report.ManagerInfo{
		Uname:    operator,
		Business: thumbupType,
		Type:     0,
		Ctime:    time.Now(),
		Index:    []interface{}{business, originID, messageID},
		Content: map[string]interface{}{
			"ip":             ip,
			"like_change":    likeChange,
			"dislike_change": dislikeChange,
		},
	})
	_ = s.cache.Do(c, func(c context.Context) {
		// update cache not del
		stats, err := s.Stats(c, business, originID, []int64{messageID})
		if err != nil {
			return
		}
		if stats == nil || stats[messageID] == nil {
			return
		}
		stat := stats[messageID]
		stat.Likes += likeChange
		stat.Dislikes += dislikeChange
		s.updateStatCache(c, businessID, originID, stat)
		s.dao.PubStatDatabus(c, business, 0, stat, 0)
	})
	return
}

// RawStats get stat changes
func (s *Service) RawStats(c context.Context, business string, originID, messageID int64) (res model.RawStats, err error) {
	var businessID int64
	if businessID, err = s.CheckBusinessOrigin(business, originID); err != nil {
		return
	}
	res, err = s.dao.RawStats(c, businessID, originID, messageID)
	return
}
