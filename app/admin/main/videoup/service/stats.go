package service

import (
	"context"
	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/log"
	"time"
)

// StatsPoints get stats point
func (s *Service) StatsPoints(c context.Context, stime, etime time.Time, typeInt int8) (points []*archive.StatsPoint, err error) {
	if points, err = s.arc.StatsPoints(c, stime, etime, typeInt); err != nil {
		log.Error("s.arc.TaskTooksByHalfHour(%v,%v)", stime, etime)
		return
	}
	return
}
