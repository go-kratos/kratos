package server

import (
	"go-common/app/service/main/assist/model/assist"

	"go-common/library/net/rpc/context"
)

func (r *RPC) Assists(c context.Context, arg *assist.ArgAssists, res *[]*assist.Assist) (err error) {
	*res, err = r.s.Assists(c, arg.Mid)
	return
}

func (r *RPC) AssistIDs(c context.Context, arg *assist.ArgAssists, res *[]int64) (err error) {
	*res, err = r.s.AssistIDs(c, arg.Mid)
	return
}

func (r *RPC) Assist(c context.Context, arg *assist.ArgAssist, res *assist.AssistRes) (err error) {
	var info *assist.AssistRes
	if info, err = r.s.Assist(c, arg.Mid, arg.AssistMid, arg.Type); err == nil && info != nil {
		*res = *info
	}
	return
}

func (r *RPC) AddAssist(c context.Context, arg *assist.ArgAssist, res *struct{}) (err error) {
	err = r.s.AddAssist(c, arg.Mid, arg.AssistMid)
	return
}

func (r *RPC) DelAssist(c context.Context, arg *assist.ArgAssist, res *struct{}) (err error) {
	err = r.s.DelAssist(c, arg.Mid, arg.AssistMid)
	return
}

func (r *RPC) AssistUps(c context.Context, arg *assist.ArgAssistUps, res *assist.AssistUpsPager) (err error) {
	var info *assist.AssistUpsPager
	if info, err = r.s.AssistUps(c, arg.AssistMid, arg.Pn, arg.Ps); err == nil && info != nil {
		*res = *info
	}
	return
}

// AssistExit notice: reuse arg *assist.ArgAssist, except Type field
func (r *RPC) AssistExit(c context.Context, arg *assist.ArgAssist, res *struct{}) (err error) {
	err = r.s.Exit(c, arg.Mid, arg.AssistMid)
	return
}
