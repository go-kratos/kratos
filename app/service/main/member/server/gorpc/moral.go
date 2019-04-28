package gorpc

import (
	"go-common/app/service/main/member/model"
	"go-common/library/net/rpc/context"
)

// Moral get user moral.
func (r *RPC) Moral(c context.Context, arg *model.ArgMemberMid, res *model.Moral) (err error) {
	var v *model.Moral
	if v, err = r.s.Moral(c, arg.Mid); err == nil && res != nil {
		*res = *v
	}
	return
}

// MoralLog get user moral log.
func (r *RPC) MoralLog(c context.Context, arg *model.ArgMemberMid, res *[]*model.UserLog) (err error) {
	*res, err = r.s.MoralLog(c, arg.Mid)
	return
}

// AddMoral add moral.
func (r *RPC) AddMoral(c context.Context, arg *model.ArgUpdateMoral, res *struct{}) (err error) {
	return r.s.UpdateMoral(c, arg)
}

// BatchAddMoral batch add moral.
func (r *RPC) BatchAddMoral(c context.Context, arg *model.ArgUpdateMorals, res *map[int64]int64) (err error) {
	*res, err = r.s.UpdateMorals(c, arg)
	return
}
