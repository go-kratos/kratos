package service

import (
	"context"
	"time"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/log"
)

// TaskTooksByHalfHour get task books by ctime
func (s *Service) TaskTooksByHalfHour(c context.Context, stime, etime time.Time) (tooks []*archive.TaskTook, err error) {
	if tooks, err = s.dao.TaskTooksByHalfHour(c, stime, etime); err != nil {
		log.Error("s.dao.TaskTooksByHalfHour(%v,%v)", stime, etime)
		return
	}
	return
}
