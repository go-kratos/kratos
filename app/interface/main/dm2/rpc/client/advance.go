package client

import (
	"context"

	"go-common/app/interface/main/dm2/model"
)

const (
	_buyAdvance    = "RPC.BuyAdvance"
	_advanceState  = "RPC.AdvanceState"
	_advances      = "RPC.Advances"
	_passAdvance   = "RPC.PassAdvance"
	_denyAdvance   = "RPC.DenyAdvance"
	_cancelAdvance = "RPC.CancelAdvance"
)

// BuyAdvance 购买高级弹幕
func (s *Service) BuyAdvance(c context.Context, arg *model.ArgAdvance) (err error) {
	err = s.client.Call(c, _buyAdvance, arg, &_noArg)
	return
}

// AdvanceState 高级弹幕状态
func (s *Service) AdvanceState(c context.Context, arg *model.ArgAdvance) (res *model.AdvState, err error) {
	err = s.client.Call(c, _advanceState, arg, &res)
	return
}

// Advances 高级弹幕申请列表
func (s *Service) Advances(c context.Context, arg *model.ArgMid) (res []*model.Advance, err error) {
	err = s.client.Call(c, _advances, arg, &res)
	return
}

// PassAdvance 通过高级弹幕申请
func (s *Service) PassAdvance(c context.Context, arg *model.ArgUpAdvance) (err error) {
	err = s.client.Call(c, _passAdvance, arg, _noArg)
	return
}

// DenyAdvance 拒绝高级弹幕申请
func (s *Service) DenyAdvance(c context.Context, arg *model.ArgUpAdvance) (err error) {
	err = s.client.Call(c, _denyAdvance, arg, _noArg)
	return
}

// CancelAdvance 取消高级弹幕申请
func (s *Service) CancelAdvance(c context.Context, arg *model.ArgUpAdvance) (err error) {
	err = s.client.Call(c, _cancelAdvance, arg, _noArg)
	return
}
