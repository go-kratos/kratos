package gorpc

import (
	"go-common/app/interface/main/history/conf"
	"go-common/app/interface/main/history/model"
	"go-common/app/interface/main/history/service"
	"go-common/library/ecode"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
)

// RPC represent rpc server
type RPC struct {
	svc *service.Service
}

// New init rpc.
func New(c *conf.Config, s *service.Service) (svr *rpc.Server) {
	r := &RPC{svc: s}
	svr = rpc.NewServer(c.RPCServer)
	if err := svr.Register(r); err != nil {
		panic(err)
	}
	return
}

// Ping check connection success.
func (r *RPC) Ping(c context.Context, arg *struct{}, res *struct{}) (err error) {
	return
}

// Progress report a user hisotry.
func (r *RPC) Progress(c context.Context, arg *model.ArgPro, res *map[int64]*model.History) (err error) {
	*res, err = r.svc.Progress(c, arg.Mid, arg.Aids)
	return
}

// Position report a user hisotry.
func (r *RPC) Position(c context.Context, arg *model.ArgPos, res *model.History) (err error) {
	var v *model.History
	tp, err := model.CheckBusiness(arg.Business)
	if err != nil {
		return
	}
	if tp > 0 {
		arg.TP = tp
	}
	v, err = r.svc.Position(c, arg.Mid, arg.Aid, arg.TP)
	if err == nil {
		*res = *v
	}
	return
}

// Add (c context.Context, mid, src, rtime int64, ip string, h *model.History) (err error) .
func (r *RPC) Add(c context.Context, arg *model.ArgHistory, res *struct{}) (err error) {
	if arg.History != nil {
		var tp int8
		if tp, err = model.CheckBusiness(arg.History.Business); err != nil {
			return
		} else if tp > 0 {
			arg.History.TP = tp
		}
	}
	err = r.svc.AddHistory(c, arg.Mid, arg.Realtime, arg.History)
	return
}

// Delete delete histories
func (r *RPC) Delete(c context.Context, arg *model.ArgDelete, res *struct{}) (err error) {
	if len(arg.Resources) == 0 {
		err = ecode.RequestErr
		return
	}
	histories := make([]*model.History, 0, len(arg.Resources))
	for _, history := range arg.Resources {
		var tp int8
		if tp, err = model.MustCheckBusiness(history.Business); err != nil {
			return
		}
		histories = append(histories, &model.History{
			Aid: history.Oid,
			TP:  tp,
		})
	}
	err = r.svc.Delete(c, arg.Mid, histories)
	return
}

// History  return all history .
func (r *RPC) History(c context.Context, arg *model.ArgHistories, res *[]*model.Resource) (err error) {
	tp, err := model.CheckBusiness(arg.Business)
	if err != nil {
		return
	}
	if tp > 0 {
		arg.TP = tp
	}
	*res, err = r.svc.Histories(c, arg.Mid, arg.TP, arg.Pn, arg.Ps)
	return
}

// HistoryCursor  return all history .
func (r *RPC) HistoryCursor(c context.Context, arg *model.ArgCursor, res *[]*model.Resource) (err error) {
	tp, err := model.CheckBusiness(arg.Business)
	if err != nil {
		return
	}
	if tp > 0 {
		arg.TP = tp
	}
	var tps []int8
	for _, b := range arg.Businesses {
		tp, err = model.CheckBusiness(b)
		if err != nil {
			return
		}
		tps = append(tps, tp)
	}
	*res, err = r.svc.HistoryCursor(c, arg.Mid, arg.Max, arg.ViewAt, arg.Ps, arg.TP, tps, arg.RealIP)
	return
}

// Clear clear history
func (r *RPC) Clear(c context.Context, arg *model.ArgClear, res *struct{}) (err error) {
	var tps []int8
	for _, b := range arg.Businesses {
		var tp int8
		if tp, err = model.MustCheckBusiness(b); err != nil {
			return
		}
		tps = append(tps, tp)
	}
	err = r.svc.ClearHistory(c, arg.Mid, tps)
	return
}
