package service

import (
	"context"
	"time"

	"go-common/library/log"
)

func (s *Service) roundproc() {
	defer s.wg.Done()
	for {
		if s.closed {
			return
		}
		s.roundToEnd()
		time.Sleep(5 * time.Minute)
	}
}

func (s *Service) roundToEnd() {
	var (
		err     error
		rows    int64
		now     = time.Now()
		minTime = s.delayRoundMinTime
		maxTime time.Time
		c       = context.TODO()
	)
	if s.roundDelayCache == 0 {
		log.Error("roundEnd conf is 0")
		return
	}
	if s.delayRoundMinTime.IsZero() {
		minTime = now.Add(-time.Duration(s.roundDelayCache)*24*time.Hour - time.Hour)
	}
	maxTime = now.Add(-time.Duration(s.roundDelayCache) * 24 * time.Hour)
	if rows, err = s.arc.UpDelayRound(c, minTime, maxTime); err != nil {
		return
	}
	s.delayRoundMinTime = maxTime
	log.Info("round auto change to end(99),startTime(%v),endTime(%v),affected(%d)", minTime, maxTime, rows)
}
