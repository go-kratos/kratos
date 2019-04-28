package archive

import (
	"context"
	"go-common/app/interface/main/creative/model/activity"
	"go-common/library/log"
)

// MissionProtocol fn
func (s *Service) MissionProtocol(c context.Context, missionID int64) (p *activity.Protocol, err error) {
	if p, err = s.act.Protocol(c, missionID); err != nil {
		log.Error("s.act.Protocol(%d) err(%v)", missionID, err)
		return
	}
	return
}

// MissionOnlineByTid fn
func (s *Service) MissionOnlineByTid(c context.Context, tid, plat int16) (res []*activity.ActWithTP, err error) {
	if res, err = s.act.MissionOnlineByTid(c, tid, plat); err != nil {
		log.Error("s.act.MissionOnlineByTid(%d,%d) err(%+v)", tid, plat, err)
		return
	}
	return
}
