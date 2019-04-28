package service

import (
	"context"
	"time"

	"go-common/app/job/main/reply-feed/model"
	"go-common/library/log"
)

func hourNow() int {
	return time.Now().Hour()
}

func lastHour() int {
	hour := hourNow()
	if hour == 0 {
		return 23
	}
	return hour - 1
}

// AddUV ...
func (s *Service) addUV(ctx context.Context, value *model.StatsMsg, isHot bool) {
	var action string
	switch value.Action {
	case model.DatabusActionLike:
		action = model.StatisticActionLike
	case model.DatabusActionHate:
		action = model.StatisticActionHate
	case model.DatabusActionReport:
		action = model.StatisticActionReport
	case model.DatabusActionReply:
		if value.Reply.IsRoot() {
			action = model.StatisticActionRootReply
		} else {
			action = model.StatisticActionChildReply
		}
	}
	if action == "" || value.Mid == 0 {
		return
	}
	s.uvQ.Do(ctx, func(ctx context.Context) {
		if isHot {
			s.dao.AddUV(ctx, action, hourNow(), int(value.Sharding()), value.Mid, model.StatisticKindHot)
		}
		s.dao.AddUV(ctx, action, hourNow(), int(value.Sharding()), value.Mid, model.StatisticKindTotal)
	})
}

// uvStatistics ...
func (s *Service) uvStatistics(ctx context.Context, slots []int, stat *model.StatisticsStat) {
	var (
		keys     []string
		lastHour = lastHour()
		x, y, z  = len(model.StatisticKinds), len(model.StatisticActions), len(slots)
		idxMap   = make([][][]int, x)
		idx      int
	)
	for i, kind := range model.StatisticKinds {
		idxMap[i] = make([][]int, y)
		for j, action := range model.StatisticActions {
			idxMap[i][j] = make([]int, z)
			for k, slot := range slots {
				keys = append(keys, s.dao.KeyUV(action, lastHour, slot, kind))
				idxMap[i][j][k] = idx
				idx++
			}
		}
	}
	counts, err := s.dao.CountUV(ctx, keys)
	if err != nil || len(counts) != len(keys) {
		return
	}
	for i, kind := range model.StatisticKinds {
		for j, action := range model.StatisticActions {
			for k := range slots {
				count := counts[idxMap[i][j][k]]
				switch {
				case kind == model.StatisticKindTotal:
					switch action {
					case model.StatisticActionRootReply:
						stat.TotalRootUV += count
					case model.StatisticActionChildReply:
						stat.TotalChildUV += count
					case model.StatisticActionLike:
						stat.TotalLikeUV += count
					case model.StatisticActionHate:
						stat.TotalHateUV += count
					case model.StatisticActionReport:
						stat.TotalReportUV += count
					}
				case kind == model.StatisticKindHot:
					switch action {
					case model.StatisticActionChildReply:
						stat.HotChildUV += count
					case model.StatisticActionLike:
						stat.HotLikeUV += count
					case model.StatisticActionHate:
						stat.HotHateUV += count
					case model.StatisticActionReport:
						stat.HotReportUV += count
					}
				}
			}
		}
	}
}

// persistStatistics persist statistics
func (s *Service) persistStatistics() {
	ctx := context.Background()
	statisticsMap := make(map[string]*model.StatisticsStat)
	nameMapping := make(map[string][]int)
	s.statisticsLock.RLock()
	for slot, stat := range s.statisticsStats {
		nameMapping[stat.Name] = append(nameMapping[stat.Name], slot)
		s, ok := statisticsMap[stat.Name]
		if ok {
			statisticsMap[stat.Name] = s.Merge(&stat)
		} else {
			statisticsMap[stat.Name] = &stat
		}
	}
	s.statisticsLock.RUnlock()
	now := time.Now()
	year, month, day := now.Date()
	date := year*10000 + int(month)*100 + day
	hour := now.Hour()
	for name, stat := range statisticsMap {
		s.uvStatistics(ctx, nameMapping[name], stat)
		err := s.dao.UpsertStatistics(ctx, name, date, hour, stat)
		var (
			retryTimes    = 0
			maxRetryTimes = 5
		)
		for err != nil && retryTimes < maxRetryTimes {
			time.Sleep(s.bc.Backoff(retryTimes))
			err = s.dao.UpsertStatistics(ctx, name, date, hour, stat)
			retryTimes++
		}
		if retryTimes >= maxRetryTimes {
			log.Error("upsert statistics error retry reached max retry times.")
		}
	}
	for i := range s.statisticsStats {
		reset(&s.statisticsStats[i])
	}
	log.Warn("reset statistics at %v", now)
}

func reset(stat *model.StatisticsStat) {
	stat.HotChildReply = 0
	stat.HotHate = 0
	stat.HotLike = 0
	stat.HotReport = 0
	stat.TotalChildReply = 0
	stat.TotalRootReply = 0
	stat.TotalReport = 0
	stat.TotalLike = 0
	stat.TotalHate = 0
}
