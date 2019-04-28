package service

import (
	"context"

	"go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
)

// UserNoticeState .
func (s *Service) UserNoticeState(c context.Context, mid int64) (res model.NoticeState, err error) {
	state, err := s.dao.NoticeState(c, mid)
	if err != nil {
		return
	}
	res = model.NewNoticeState(state)
	return
}

// UpdateUserNoticeState .
func (s *Service) UpdateUserNoticeState(c context.Context, mid int64, typ string) (err error) {
	state, err := s.UserNoticeState(c, mid)
	if err != nil {
		return
	}
	if _, ok := state[typ]; !ok {
		err = ecode.RequestErr
		return
	}
	state[typ] = true
	err = s.dao.UpdateNoticeState(c, mid, state.ToInt64())
	return
}
