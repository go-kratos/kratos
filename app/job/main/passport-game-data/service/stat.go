package service

import (
	"context"
	"time"

	"go-common/app/job/main/passport-game-data/model"
)

// CompareProcStat get actual compare proc stat.
func (s *Service) CompareProcStat(c context.Context) *model.ProcStat {
	var c2lStat *model.CompareProcStat
	if s.c.Compare.Cloud2Local.On {
		cc := s.c2lC
		c2lStat = &model.CompareProcStat{
			StartTime:     s.c.Compare.Cloud2Local.StartTime,
			EndTime:       s.c.Compare.Cloud2Local.EndTime,
			StepDuration:  model.JsonDuration(cc.StepDuration),
			LoopDuration:  model.JsonDuration(cc.LoopDuration),
			DelayDuration: model.JsonDuration(cc.DelayDuration),

			BatchSize:           cc.BatchSize,
			BatchMissRetryCount: cc.BatchMissRetryCount,

			Debug: cc.Debug,
			Fix:   cc.Fix,

			CurrentRangeStart:        model.JSONTime(cc.st),
			CurrentRangeEnd:          model.JSONTime(cc.ed),
			CurrentRangeRecordsCount: cc.rangeCount,
			TotalRangeRecordsCount:   cc.totalCount,
			DiffCount:                cc.diffCount,

			Sleeping: cc.sleeping,
		}
		if cc.sleeping {
			c2lStat.SleepSeconds = cc.sleepingSeconds
			c2lStat.SleepFrom = time.Unix(cc.sleepFromTs, 0).Format(_timeFormat)
			c2lStat.SleepRemainSeconds = cc.sleepingSeconds - (time.Now().Unix() - cc.sleepFromTs)
		}
	}

	var l2cStat *model.CompareProcStat
	if s.c.Compare.Local2Cloud.On {
		cc := s.l2cC
		l2cStat = &model.CompareProcStat{
			StartTime:     s.c.Compare.Local2Cloud.StartTime,
			EndTime:       s.c.Compare.Local2Cloud.EndTime,
			StepDuration:  model.JsonDuration(cc.StepDuration),
			LoopDuration:  model.JsonDuration(cc.LoopDuration),
			DelayDuration: model.JsonDuration(cc.DelayDuration),

			BatchSize:           cc.BatchSize,
			BatchMissRetryCount: cc.BatchMissRetryCount,

			Debug: cc.Debug,
			Fix:   cc.Fix,

			CurrentRangeStart:        model.JSONTime(cc.st),
			CurrentRangeEnd:          model.JSONTime(cc.ed),
			CurrentRangeRecordsCount: cc.rangeCount,
			TotalRangeRecordsCount:   cc.totalCount,
			DiffCount:                cc.diffCount,

			Sleeping: cc.sleeping,
		}
		if cc.sleeping {
			l2cStat.SleepSeconds = cc.sleepingSeconds
			l2cStat.SleepFrom = time.Unix(cc.sleepFromTs, 0).Format(_timeFormat)
			l2cStat.SleepRemainSeconds = cc.sleepingSeconds - (time.Now().Unix() - cc.sleepFromTs)
		}
	}

	return &model.ProcStat{
		Cloud2Local: c2lStat,
		Local2Cloud: l2cStat,
	}
}
