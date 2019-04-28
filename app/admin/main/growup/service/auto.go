package service

import (
	"context"
	"fmt"
	"time"

	"go-common/library/log"
)

// AutoDismiss update account state to 6
func (s *Service) AutoDismiss(c context.Context, operator string, typ int, mid int64, reason string) (err error) {
	ups, err := s.dao.UpsVideoInfo(c, fmt.Sprintf("mid = %d", mid))
	if err != nil {
		log.Error("s.dao.UpsVideoInfo error(%v)", err)
		return
	}
	if len(ups) <= 0 {
		return
	}
	up := ups[0]
	return s.Dismiss(c, operator, typ, up.AccountState, mid, reason)
}

// AutoForbid update account state to 7 and add a n days CD
func (s *Service) AutoForbid(c context.Context, operator string, typ int, mid int64, reason string, days, second int) (err error) {
	ups, err := s.dao.UpsVideoInfo(c, fmt.Sprintf("mid = %d", mid))
	if err != nil {
		log.Error("s.dao.UpsVideoInfo error(%v)", err)
		return
	}
	if len(ups) <= 0 {
		return
	}
	up := ups[0]
	switch up.AccountState {
	case 3:
		return s.Forbid(c, operator, typ, 3, mid, reason, days, second)
	case 7:
		return s.Forbid(c, operator, typ, 7, mid, reason, days, int(int64(up.ExpiredIn)+int64(second)-time.Now().Unix()))
	}
	return
}
