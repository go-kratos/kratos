package gorpc

import (
	"context"

	"go-common/app/service/main/member/model"
)

const (
	_RealnameStatus        = "RPC.RealnameStatus"
	_RealnameApplyStatus   = "RPC.RealnameApplyStatus"
	_RealnameTelCapture    = "RPC.RealnameTelCapture"
	_RealnameApply         = "RPC.RealnameApply"
	_RealnameDetail        = "RPC.RealnameDetail"
	_RealnameAlipayApply   = "RPC.RealnameAlipayApply"
	_RealnameAlipayConfirm = "RPC.RealnameAlipayConfirm"
	_RealnameAlipayBizno   = "RPC.RealnameAlipayBizno"
)

// RealnameStatus get realname current status by mid
func (s *Service) RealnameStatus(c context.Context, arg *model.ArgMemberMid) (res *model.RealnameStatus, err error) {
	err = s.client.Call(c, _RealnameStatus, arg, &res)
	return
}

// RealnameApplyStatus get user realname apply status
func (s *Service) RealnameApplyStatus(c context.Context, arg *model.ArgMemberMid) (res *model.RealnameApplyStatusInfo, err error) {
	err = s.client.Call(c, _RealnameApplyStatus, arg, &res)
	return
}

// RealnameTelCapture get user telphone capture
func (s *Service) RealnameTelCapture(c context.Context, arg *model.ArgMemberMid) (err error) {
	err = s.client.Call(c, _RealnameTelCapture, arg, &_noRes)
	return
}

// RealnameApply put a realname apply
func (s *Service) RealnameApply(c context.Context, arg *model.ArgRealnameApply) (err error) {
	err = s.client.Call(c, _RealnameApply, arg, &_noRes)
	return
}

// RealnameAlipayApply put a alipay realname apply
func (s *Service) RealnameAlipayApply(c context.Context, arg *model.ArgRealnameAlipayApply) (err error) {
	err = s.client.Call(c, _RealnameAlipayApply, arg, &_noRes)
	return
}

// RealnameAlipayConfirm confirm a alipay realname apply positivly
func (s *Service) RealnameAlipayConfirm(c context.Context, arg *model.ArgRealnameAlipayConfirm) (err error) {
	err = s.client.Call(c, _RealnameAlipayConfirm, arg, &_noRes)
	return
}

// RealnameAlipayBizno get alipay realname certify bizno by mid
func (s *Service) RealnameAlipayBizno(c context.Context, arg *model.ArgMemberMid) (res *model.RealnameAlipayInfo, err error) {
	err = s.client.Call(c, _RealnameAlipayBizno, arg, &res)
	return
}

// RealnameDetail is
func (s *Service) RealnameDetail(c context.Context, arg *model.ArgMemberMid) (res *model.RealnameDetail, err error) {
	res = new(model.RealnameDetail)
	err = s.client.Call(c, _RealnameDetail, arg, res)
	return
}
