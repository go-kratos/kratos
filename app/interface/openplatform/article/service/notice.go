package service

import (
	"context"
	"time"

	artmdl "go-common/app/interface/openplatform/article/model"
)

const (
	_platAll            = 0
	_platAndroid        = 1
	_platIOS            = 2
	_equal              = 0
	_greaterThanOrEqual = 1
	_lessThanOrEqual    = 2
)

func (s *Service) loadNoticeproc() {
	for {
		if notices, err := s.dao.Notices(context.TODO(), time.Now()); err == nil {
			s.notices = notices
		}
		time.Sleep(time.Minute)
	}
}

// Notice get notice
func (s *Service) Notice(platform string, build int) (res *artmdl.Notice) {
	var plat int
	if platform == "android" {
		plat = _platAndroid
	}
	if platform == "ios" {
		plat = _platIOS
	}
	for _, notice := range s.notices {
		var ok bool
		if (notice.Plat == _platAll) || (notice.Plat == plat) {
			switch notice.Condition {
			case _equal:
				ok = build == notice.Build
			case _greaterThanOrEqual:
				ok = build >= notice.Build
			case _lessThanOrEqual:
				ok = build <= notice.Build
			}
		}
		if ok {
			notice.Content = notice.Title
			return notice
		}
	}
	return nil
}
