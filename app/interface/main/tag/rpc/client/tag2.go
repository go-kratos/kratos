package tag

import (
	"context"

	"go-common/app/interface/main/tag/model"
)

const (
	_upBind    = "RPC.UpBind"
	_adminBind = "RPC.AdminBind"
	_userAdd   = "RPC.UserAdd"
	_userDel   = "RPC.UserDel"
	_resTags   = "RPC.ResTags"
)

var (
	_noRes = &struct{}{}
)

// UpBind .
func (s *Service) UpBind(c context.Context, arg *model.ArgBind) (err error) {
	err = s.client.Call(c, _upBind, arg, _noRes)
	return
}

// AdminBind .
func (s *Service) AdminBind(c context.Context, arg *model.ArgBind) (err error) {
	err = s.client.Call(c, _adminBind, arg, _noRes)
	return
}

// UserAdd .
func (s *Service) UserAdd(c context.Context, arg *model.ArgUserAdd) (tid int64, err error) {
	err = s.client.Call(c, _userAdd, arg, &tid)
	return
}

// UserDel .
func (s *Service) UserDel(c context.Context, arg *model.ArgUserDel) (err error) {
	err = s.client.Call(c, _userDel, arg, _noRes)
	return
}

// ResTags .
func (s *Service) ResTags(c context.Context, arg *model.ArgResTags) (res map[int64][]*model.Tag, err error) {
	err = s.client.Call(c, _resTags, arg, &res)
	return
}
