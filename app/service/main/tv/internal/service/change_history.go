package service

import (
	"context"

	"go-common/app/service/main/tv/internal/model"
	xtime "go-common/library/time"
)

func (s *Service) ChangeHistory(c context.Context, hid int32) (ch *model.UserChangeHistory, err error) {
	return s.dao.UserChangeHistoryByID(c, hid)
}

func (s *Service) ChangeHistorys(c context.Context, mid int64, from, to, pn, ps int32) (chs []*model.UserChangeHistory, total int, err error) {
	if from == 0 || to == 0 {
		return s.dao.UserChangeHistorysByMid(c, mid, pn, ps)
	}
	return s.dao.UserChangeHistorysByMidAndCtime(c, mid, xtime.Time(from), xtime.Time(to), pn, ps)
}
