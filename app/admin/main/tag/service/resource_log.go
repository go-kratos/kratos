package service

import (
	"context"

	"go-common/app/admin/main/tag/model"
	"go-common/library/ecode"
)

var (
	_emtpyResTagLog = make([]*model.ResTagLog, 0)
)

// ResourceLogs resource log.
func (s *Service) ResourceLogs(c context.Context, oid int64, tp, role, action, pn, ps int32) (res []*model.ResTagLog, total int64, err error) {
	var (
		start = (pn - 1) * ps
		end   = ps
	)
	if total, err = s.dao.ResTagLogCount(c, oid, tp, role, action); err != nil || total <= 0 {
		return _emtpyResTagLog, 0, err
	}
	res, _ = s.dao.ResourceLogs(c, oid, tp, role, action, start, end)
	if len(res) == 0 {
		res = _emtpyResTagLog
		total = 0
	}
	return
}

// UpdateResLogState UpdateResLogState.
func (s *Service) UpdateResLogState(c context.Context, id, oid int64, tp, state int32) (err error) {
	affect, err := s.dao.UpdateResLogState(c, id, oid, tp, state)
	if err != nil {
		return
	}
	if affect == 0 {
		err = ecode.TagOperateFail
	}
	return
}
