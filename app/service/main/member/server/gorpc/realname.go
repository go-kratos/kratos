package gorpc

import (
	"go-common/app/service/main/member/model"
	"go-common/library/net/rpc/context"
)

// RealnameStatus is
func (r *RPC) RealnameStatus(c context.Context, arg *model.ArgMemberMid, res *model.RealnameStatus) (err error) {
	var v model.RealnameStatus
	if v, err = r.s.RealnameStatus(c, arg.Mid); err == nil && res != nil {
		*res = v
	}
	return
}

// RealnameApplyStatus is
func (r *RPC) RealnameApplyStatus(c context.Context, arg *model.ArgMemberMid, res *model.RealnameApplyStatusInfo) (err error) {
	var v *model.RealnameApplyStatusInfo
	if v, err = r.s.RealnameApplyStatus(c, arg.Mid); err == nil && v != nil {
		*res = *v
	}
	return
}

// RealnameTelCapture is
func (r *RPC) RealnameTelCapture(c context.Context, arg *model.ArgMemberMid, res *struct{}) (err error) {
	_, err = r.s.RealnameTelCapture(c, arg.Mid)
	return
}

// RealnameApply is
func (r *RPC) RealnameApply(c context.Context, arg *model.ArgRealnameApply, res *struct{}) (err error) {
	err = r.s.RealnameApply(c, arg.MID, arg.CaptureCode, arg.Realname, arg.CardType, arg.CardCode, arg.Country, arg.HandIMGToken, arg.FrontIMGToken, arg.BackIMGToken)
	return
}

// RealnameAlipayApply commit a alipay realname apply
func (r *RPC) RealnameAlipayApply(c context.Context, arg *model.ArgRealnameAlipayApply, res *struct{}) (err error) {
	err = r.s.RealnameAlipayApply(c, arg.MID, arg.CaptureCode, arg.Realname, arg.CardCode, arg.IMGToken, arg.Bizno)
	return
}

// RealnameAlipayConfirm confirm a alipay realname apply
func (r *RPC) RealnameAlipayConfirm(c context.Context, arg *model.ArgRealnameAlipayConfirm, res *struct{}) (err error) {
	err = r.s.RealnameAlipayConfirm(c, arg.MID, arg.Pass, arg.Reason)
	return
}

// RealnameAlipayBizno get alipay realname certify bizno by mid
func (r *RPC) RealnameAlipayBizno(c context.Context, arg *model.ArgMemberMid, res *model.RealnameAlipayInfo) (err error) {
	var bizno string
	if bizno, err = r.s.RealnameAlipayBizno(c, arg.Mid); err == nil {
		(*res).Bizno = bizno
	}
	return
}

// RealnameDetail detail about realname by mid
func (r *RPC) RealnameDetail(ctx context.Context, arg *model.ArgMemberMid, res *model.RealnameDetail) error {
	detail, err := r.s.RealnameDetail(ctx, arg.Mid)
	if err != nil {
		return err
	}
	*res = *detail
	return nil
}
