package service

import (
	"context"
	"errors"
	"fmt"

	"go-common/app/service/main/reply-feed/model"
	"go-common/library/log"
)

/*
SlotsMapping
*/

// SlotStatsManager SlotStatsManager
func (s *Service) SlotStatsManager(ctx context.Context) ([]*model.SlotsStat, error) {
	return s.dao.SlotsStatManager(ctx)
}

// NewEGroup new experimental group.
func (s *Service) NewEGroup(ctx context.Context, name, algorithm, weight string, percent int) (err error) {
	var (
		usedSlots []int64
		slots     []int64
	)
	if usedSlots, _, _, err = s.dao.SlotsStatByName(ctx, name); err != nil {
		return
	}
	if len(usedSlots) > 0 {
		return errors.New("duplicate name")
	}
	if slots, err = s.dao.IdleSlots(ctx, percent); err != nil {
		return
	}
	return s.dao.UpdateSlotsStat(ctx, name, algorithm, weight, slots, model.StateInactive)
}

// ResizeSlots resize slots
func (s *Service) ResizeSlots(ctx context.Context, name string, percent int) (err error) {
	var (
		count             int
		slots             []int64
		idelSlots         []int64
		algorithm, weight string
	)
	if slots, algorithm, weight, err = s.dao.SlotsStatByName(ctx, name); err != nil {
		return
	}
	if percent > len(slots) {
		if count, err = s.dao.CountIdleSlot(ctx); err != nil {
			return
		}
		if percent > count {
			err = errors.New("out of slot")
			return
		}
		if idelSlots, err = s.dao.IdleSlots(ctx, percent-len(slots)); err != nil {
			return
		}
		if err = s.dao.UpdateSlotsStat(ctx, name, algorithm, weight, idelSlots, model.StateActive); err != nil {
			return
		}
	} else {
		if err = s.dao.UpdateSlotsStat(ctx, model.DefaultSlotName, model.DefaultAlgorithm, model.DefaultWeight, slots[percent:], model.StateActive); err != nil {
			return
		}
	}
	return
}

// EditSlotsStat edit a test set weight.
func (s *Service) EditSlotsStat(ctx context.Context, name, algorithm, weight string, slots []int64) (err error) {
	if err = s.dao.UpdateSlotsStat(ctx, name, algorithm, weight, slots, model.StateInactive); err != nil {
		log.Error("Edit SlotsMapping Failed, Error (%v)", err)
	}
	return
}

// ModifyState modify test set state, activate or inactivate by name.
func (s *Service) ModifyState(ctx context.Context, name string, state int) (err error) {
	return s.dao.ModifyState(ctx, name, state)
}

// ResetEGroup ResetEGroup
func (s *Service) ResetEGroup(ctx context.Context, name string) (err error) {
	var slots []int64
	if slots, _, _, err = s.dao.SlotsStatByName(ctx, name); err != nil {
		return
	}
	if err = s.dao.UpdateSlotsStat(ctx, model.DefaultSlotName, model.DefaultAlgorithm, model.DefaultWeight, slots, model.StateActive); err != nil {
		return
	}
	return
}

/*
TODO(Statistics)
Statistics
*/

// StatisticsByDate GetStatisticsByDate
func (s *Service) StatisticsByDate(ctx context.Context, req *model.SSReq) (res map[string]map[int]*model.StatisticsStat, err error) {
	res = make(map[string]map[int]*model.StatisticsStat)
	stats, err := s.dao.StatisticsByDate(ctx, req.DateFrom, req.DateEnd)
	if err != nil {
		return
	}
	groupedStatistics := stats.GroupByName()
	for name, ss := range groupedStatistics {
		dateStatistics := make(map[int]*model.StatisticsStat)
		for _, s := range ss {
			if _, ok := dateStatistics[s.Date]; ok {
				dateStatistics[s.Date] = dateStatistics[s.Date].MergeByDate(s)
			} else {
				dateStatistics[s.Date] = s
			}
		}
		res[name] = dateStatistics
	}
	return
}

// StatisticsByHour GetStatisticsByHour
func (s *Service) StatisticsByHour(ctx context.Context, req *model.SSReq) (res map[string]map[string]*model.StatisticsStat, err error) {
	res = make(map[string]map[string]*model.StatisticsStat)
	stats, err := s.dao.StatisticsByDate(ctx, req.DateFrom, req.DateEnd)
	if err != nil {
		return
	}
	groupedStatistics := stats.GroupByName()
	for name, ss := range groupedStatistics {
		hourStatistics := make(map[string]*model.StatisticsStat)
		for _, s := range ss {
			if s.Hour < 10 {
				hourStatistics[fmt.Sprintf("%d-0%d", s.Date, s.Hour)] = s
			} else {
				hourStatistics[fmt.Sprintf("%d-%d", s.Date, s.Hour)] = s
			}
		}
		res[name] = hourStatistics
	}
	return
}
