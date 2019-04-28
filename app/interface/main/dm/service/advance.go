package service

import (
	"context"

	dm2Mdl "go-common/app/interface/main/dm2/model"
	"go-common/library/log"
)

// BuyAdvance 购买高级弹幕
func (s *Service) BuyAdvance(c context.Context, mid, cid int64, mode string) (err error) {
	arg := &dm2Mdl.ArgAdvance{
		Mid:  mid,
		Cid:  cid,
		Mode: mode,
	}
	if err = s.dmRPC.BuyAdvance(c, arg); err != nil {
		log.Error("dmRPC.BuyAdvance(%v) error(%v)")
	}
	return
}

// AdvanceState 高级弹幕状态
func (s *Service) AdvanceState(c context.Context, mid, cid int64, mode string) (state *dm2Mdl.AdvState, err error) {
	arg := &dm2Mdl.ArgAdvance{
		Mid:  mid,
		Cid:  cid,
		Mode: mode,
	}
	if state, err = s.dmRPC.AdvanceState(c, arg); err != nil {
		log.Error("dmRPC.AdvanceState(%v) error(%v)", arg, err)
	}
	return
}

// Advances 高级弹幕申请列表
func (s *Service) Advances(c context.Context, mid int64) (res []*dm2Mdl.Advance, err error) {
	arg := &dm2Mdl.ArgMid{
		Mid: mid,
	}
	if res, err = s.dmRPC.Advances(c, arg); err != nil {
		log.Error("dmRPC.Advances(%v) error(%v)", arg, err)
	}
	return
}

// PassAdvance 通过高级弹幕申请
func (s *Service) PassAdvance(c context.Context, mid, id int64) (err error) {
	arg := &dm2Mdl.ArgUpAdvance{
		Mid: mid,
		ID:  id,
	}
	if err = s.dmRPC.PassAdvance(c, arg); err != nil {
		log.Error("dmRPC.PassAdvance(%v) error(%v)", arg, err)
	}
	return
}

// DenyAdvance 拒绝高级弹幕申请
func (s *Service) DenyAdvance(c context.Context, mid, id int64) (err error) {
	arg := &dm2Mdl.ArgUpAdvance{
		Mid: mid,
		ID:  id,
	}
	if err = s.dmRPC.DenyAdvance(c, arg); err != nil {
		log.Error("dmRPC.DenyAdvance(%v) error(%v)", arg, err)
	}
	return
}

// CancelAdvance 取消高级弹幕申请
func (s *Service) CancelAdvance(c context.Context, mid, id int64) (err error) {
	arg := &dm2Mdl.ArgUpAdvance{
		Mid: mid,
		ID:  id,
	}
	if err = s.dmRPC.CancelAdvance(c, arg); err != nil {
		log.Error("dmRPC.CancelAdvance(%v) error(%v)", arg, err)
	}
	return
}
