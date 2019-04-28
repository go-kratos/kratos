package service

import (
	"context"
	"time"

	"go-common/app/job/main/spy/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// UpdateStatData update spy stat data.
func (s *Service) UpdateStatData(c context.Context, m *model.SpyStatMessage) (err error) {
	//TODO check event resaon
	if s.allEventName[m.EventName] == 0 {
		log.Error("event name not found %+v", err)
		err = ecode.SpyEventNotExist
		return
	}
	stat := &model.Statistics{
		TargetMid: m.TargetMid,
		TargetID:  m.TargetID,
		EventID:   s.allEventName[m.EventName],
		State:     model.WaiteCheck,
		Quantity:  m.Quantity,
		Ctime:     time.Now(),
	}
	if stat.TargetID != 0 {
		_, ok := s.activityEvents[m.EventName]
		if ok {
			stat.Type = model.ActivityType
		} else {
			stat.Type = model.ArchiveType
		}
	}
	// add stat
	if model.ResetStat == m.Type {
		if _, err = s.dao.AddStatistics(c, stat); err != nil {
			log.Error("%+v", err)
			return
		}
	} else {
		if _, err = s.dao.AddIncrStatistics(c, stat); err != nil {
			log.Error("%+v", err)
			return
		}
	}
	return
}
