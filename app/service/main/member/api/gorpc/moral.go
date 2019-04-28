package gorpc

import (
	"context"

	"go-common/app/service/main/member/model"
)

// Moral get user moral.
func (s *Service) Moral(c context.Context, arg *model.ArgMemberMid) (res *model.Moral, err error) {
	err = s.client.Call(c, _Moral, arg, &res)
	return
}

// MoralLog get user moral log.
func (s *Service) MoralLog(c context.Context, arg *model.ArgMemberMid) (res []*model.UserLog, err error) {
	err = s.client.Call(c, _MoralLog, arg, &res)
	return
}

// AddMoral add moral .
func (s *Service) AddMoral(c context.Context, arg *model.ArgUpdateMoral) (err error) {
	err = s.client.Call(c, _addMoral, arg, &_noRes)
	return
}

// BatchAddMoral get user moral log.
func (s *Service) BatchAddMoral(c context.Context, arg *model.ArgUpdateMorals) (res map[int64]int64, err error) {
	err = s.client.Call(c, _batchAddMoral, arg, &res)
	return
}
