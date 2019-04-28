package service

import (
	"context"

	"go-common/library/log"
)

// SetLock set click lock by plat and lv
func (s *Service) SetLock(c context.Context, aid int64, plat, lock, lv int8) (err error) {
	if _, err = s.db.UpForbid(c, aid, plat, lock, lv); err != nil {
		log.Error("s.db.UpForbid(%+v) error(%v)", c, err)
		return
	}
	return
}

// SetMidForbid is
func (s *Service) SetMidForbid(c context.Context, mid int64, status int8) (err error) {
	if err = s.db.UpMidForbidStatus(c, mid, status); err != nil {
		log.Error("s.db.UpMidForbidStatus(%d, %d) error(%v)", mid, status, err)
		return
	}
	return
}
