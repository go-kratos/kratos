package client

import (
	"context"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/net/rpc"
)

const (
	_subjectInfos   = "RPC.SubjectInfos"
	_editDMState    = "RPC.EditDMState"
	_editDMPool     = "RPC.EditDMPool"
	_editDMAttr     = "RPC.EditDMAttr"
	_addUserFilters = "RPC.AddUserFilters"
	_userFilters    = "RPC.UserFilters"
	_delUserFilters = "RPC.DelUserFilters"
	_addUpFilters   = "RPC.AddUpFilters"
	_upFilters      = "RPC.UpFilters"
	_banUsers       = "RPC.BanUsers"
	_cancelBanUsers = "RPC.CancelBanUsers"
	_editUpFilters  = "RPC.EditUpFilters"
	_addGblFilter   = "RPC.AddGlobalFilter"
	_globalFilters  = "RPC.GlobalFilters"
	_delGlbFilters  = "RPC.DelGlobalFilters"
)

const (
	_appid = "community.service.dm"
)

var (
	_noArg = &struct{}{}
)

// Service dm rpc client.
type Service struct {
	client *rpc.Client2
}

// New new a dm rpc client.
func New(c *rpc.ClientConfig) (s *Service) {
	s = &Service{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return
}

// SubjectInfos multi get dm subject info by type and oids.
func (s *Service) SubjectInfos(c context.Context, arg *model.ArgOids) (res map[int64]*model.SubjectInfo, err error) {
	err = s.client.Call(c, _subjectInfos, arg, &res)
	return
}

// EditDMState update dm state.
// 0：正常、1：删除10：用户删除、11：举报脚本删除
func (s *Service) EditDMState(c context.Context, arg *model.ArgEditDMState) (err error) {
	err = s.client.Call(c, _editDMState, arg, _noArg)
	return
}

// EditDMAttr update dm attr.
func (s *Service) EditDMAttr(c context.Context, arg *model.ArgEditDMAttr) (err error) {
	err = s.client.Call(c, _editDMAttr, arg, _noArg)
	return
}

// EditDMPool update dm pool.
// 0:普通弹幕池、1:字幕弹幕池
func (s *Service) EditDMPool(c context.Context, arg *model.ArgEditDMPool) (err error) {
	err = s.client.Call(c, _editDMPool, arg, _noArg)
	return
}

// AddUserFilters add user filter.
func (s *Service) AddUserFilters(c context.Context, arg *model.ArgAddUserFilters) (res []*model.UserFilter, err error) {
	err = s.client.Call(c, _addUserFilters, arg, &res)
	return
}

// UserFilters multi get user filters.
func (s *Service) UserFilters(c context.Context, arg *model.ArgMid) (res []*model.UserFilter, err error) {
	err = s.client.Call(c, _userFilters, arg, &res)
	return
}

// DelUserFilters delete user filters by filter id.
func (s *Service) DelUserFilters(c context.Context, arg *model.ArgDelUserFilters) (affect int64, err error) {
	err = s.client.Call(c, _delUserFilters, arg, &affect)
	return
}

// AddUpFilters add up filters.
func (s *Service) AddUpFilters(c context.Context, arg *model.ArgAddUpFilters) (err error) {
	err = s.client.Call(c, _addUpFilters, arg, &_noArg)
	return
}

// UpFilters multi get up filters.
func (s *Service) UpFilters(c context.Context, arg *model.ArgUpFilters) (res []*model.UpFilter, err error) {
	err = s.client.Call(c, _upFilters, arg, &res)
	return
}

// BanUsers ban user by upper or assist.
func (s *Service) BanUsers(c context.Context, arg *model.ArgBanUsers) (err error) {
	err = s.client.Call(c, _banUsers, arg, &_noArg)
	return
}

// CancelBanUsers cancel users by upper or assiat.
func (s *Service) CancelBanUsers(c context.Context, arg *model.ArgCancelBanUsers) (err error) {
	err = s.client.Call(c, _cancelBanUsers, arg, &_noArg)
	return
}

// EditUpFilters edit upper filters.
func (s *Service) EditUpFilters(c context.Context, arg *model.ArgEditUpFilters) (affect int64, err error) {
	err = s.client.Call(c, _editUpFilters, arg, &affect)
	return
}

// AddGlobalFilter add global filters.
func (s *Service) AddGlobalFilter(c context.Context, arg *model.ArgAddGlobalFilter) (res *model.GlobalFilter, err error) {
	err = s.client.Call(c, _addGblFilter, arg, &res)
	return
}

// GlobalFilters multi get global filters.
func (s *Service) GlobalFilters(c context.Context, arg *model.ArgGlobalFilters) (res []*model.GlobalFilter, err error) {
	err = s.client.Call(c, _globalFilters, arg, &res)
	return
}

// DelGlobalFilters delete global filter.
func (s *Service) DelGlobalFilters(c context.Context, arg *model.ArgDelGlobalFilters) (affect int64, err error) {
	err = s.client.Call(c, _delGlbFilters, arg, &affect)
	return
}
