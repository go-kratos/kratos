package service

import (
	"context"

	"go-common/app/service/main/push/model"
	"go-common/library/log"
	xtime "go-common/library/time"
)

func (s *Service) setChCounter(taskID string, v int64) {
	s.pmu.Lock()
	s.chCounter[taskID] += v
	s.pmu.Unlock()
}

func (s *Service) chCounterVal(taskID string) int64 {
	s.pmu.RLock()
	defer s.pmu.RUnlock()
	return s.chCounter[taskID]
}

const (
	_pgStatus = iota
	_pgBeginTime
	_pgEndTime
	_pgPushTime
	_pgMidTotal
	_pgMidValid
	_pgMidMissed
	_pgMidMissedSuccess
	_pgMidMissedFailed
	_pgTokenTotal
	_pgTokenSuccess
	_pgTokenValid
	_pgTokenFailed
	_pgTokenDelay
	_pgRetryTimes
)

func (s *Service) updateProgressproc() {
	defer s.waiter.Done()
	for {
		f, ok := <-s.progressCh
		if !ok {
			log.Info("updateProgressproc exit")
			return
		}
		f()
	}
}

// AddMidProgress .
func (s *Service) AddMidProgress(ctx context.Context, task string, midTotal, midValid int64) error {
	p := &model.Progress{MidTotal: midTotal, MidValid: midValid}
	return s.dao.UpdateTaskProgress(ctx, task, p)
}

func (s *Service) setProgress(taskID string, typ int, v int64) {
	s.ppmu.Lock()
	defer s.ppmu.Unlock()
	if s.progress[taskID] == nil {
		s.progress[taskID] = &model.Progress{Brands: make(map[int]int64)}
	}
	s.setBaseProgress(taskID, typ, v)
}

func (s *Service) setBaseProgress(taskID string, typ int, v int64) {
	switch typ {
	case _pgStatus:
		s.progress[taskID].Status = int8(v)
	case _pgBeginTime:
		s.progress[taskID].BeginTime = xtime.Time(v)
	case _pgEndTime:
		s.progress[taskID].EndTime = xtime.Time(v)
	case _pgPushTime:
		s.progress[taskID].PushTime = xtime.Time(v)
	case _pgRetryTimes:
		s.progress[taskID].RetryTimes += v
	case _pgMidTotal:
		s.progress[taskID].MidTotal += v
	case _pgMidValid:
		s.progress[taskID].MidValid += v
	case _pgMidMissed:
		s.progress[taskID].MidMissed += v
	case _pgMidMissedSuccess:
		s.progress[taskID].MidMissedSuccess += v
	case _pgMidMissedFailed:
		s.progress[taskID].MidMissedFailed += v
	case _pgTokenTotal:
		s.progress[taskID].TokenTotal += v
	case _pgTokenSuccess:
		s.progress[taskID].TokenSuccess += v
	case _pgTokenValid:
		s.progress[taskID].TokenValid += v
	case _pgTokenFailed:
		s.progress[taskID].TokenFailed += v
	case _pgTokenDelay:
		s.progress[taskID].TokenDelay += v
	}
}

func (s *Service) setBrandProgress(taskID string, brand int, v int64) {
	s.ppmu.Lock()
	defer s.ppmu.Unlock()
	if s.progress[taskID] == nil {
		s.progress[taskID] = &model.Progress{Brands: make(map[int]int64)}
	}
	s.progress[taskID].Brands[brand] += v
}

func (s *Service) fetchProgress(taskID string) *model.Progress {
	s.ppmu.RLock()
	defer s.ppmu.RUnlock()
	return s.progress[taskID]
}
