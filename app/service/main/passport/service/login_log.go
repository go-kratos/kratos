package service

import (
	"context"

	"go-common/app/service/main/passport/model"
)

const (
	_maxLimit = 1000
)

var (
	_emptyLoginLogs = make([]*model.LoginLog, 0)

	_emptyLoginLogResps = make([]*model.LoginLogResp, 0)
)

// FormattedLoginLogs get the latest limit login logs with formatting IP int to IP string.
func (s *Service) FormattedLoginLogs(c context.Context, mid int64, limit int) (res []*model.LoginLogResp, err error) {
	ls, err := s.LoginLogs(c, mid, limit)
	if err != nil {
		return
	}
	if len(ls) == 0 {
		res = _emptyLoginLogResps
		return
	}
	res = make([]*model.LoginLogResp, 0)
	for _, v := range ls {
		res = append(res, model.Format(v))
	}
	return
}

// LoginLogs get the latest limit login logs.
// If the limit is less than or equal to 0, a empty result will be returned,
// else if the limit is greater than then max limit, then the limit will be set to max limit.
func (s *Service) LoginLogs(c context.Context, mid int64, limit int) (res []*model.LoginLog, err error) {
	if mid < 0 || limit <= 0 {
		res = _emptyLoginLogs
		return
	}
	if limit > _maxLimit {
		limit = _maxLimit
	}
	if res, err = s.d.LoginLogs(c, mid, limit); err != nil {
		return
	}
	if len(res) == 0 {
		res = _emptyLoginLogs
	}
	return
}
