package server

import (
	"go-common/app/service/main/assist/model/assist"

	"go-common/library/net/rpc/context"
)

func (r *RPC) AssistLogAdd(c context.Context, arg *assist.ArgAssistLogAdd, res *struct{}) (err error) {
	err = r.s.AddLog(c, arg.Mid, arg.AssistMid, arg.Type, arg.Action, arg.SubjectID, arg.ObjectID, arg.Detail)
	return
}

func (r *RPC) AssistLogCancel(c context.Context, arg *assist.ArgAssistLog, res *struct{}) (err error) {
	err = r.s.CancelLog(c, arg.Mid, arg.AssistMid, arg.LogID)
	return
}

func (r *RPC) AssistLogs(c context.Context, arg *assist.ArgAssistLogs, res *[]*assist.Log) (err error) {
	*res, err = r.s.Logs(c, arg.Mid, arg.AssistMid, arg.Stime, arg.Etime, arg.Pn, arg.Ps)
	return
}

func (r *RPC) AssistLogInfo(c context.Context, arg *assist.ArgAssistLog, res *assist.Log) (err error) {
	var info *assist.Log
	if info, err = r.s.LogInfo(c, arg.LogID, arg.Mid, arg.AssistMid); err == nil && info != nil {
		*res = *info
	}
	return
}
