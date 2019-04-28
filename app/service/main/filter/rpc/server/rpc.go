package server

import (
	"go-common/app/service/main/filter/conf"
	rpcmodel "go-common/app/service/main/filter/model/rpc"
	"go-common/app/service/main/filter/service"
	"go-common/library/log"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
)

// RPC represent rpc server
type RPC struct {
	s *service.Service
}

// New new rpc server.
func New(c *conf.Config, s *service.Service) (svr *rpc.Server) {
	r := &RPC{s: s}
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

// Filter use regexp and trie to filter msg.
func (r *RPC) Filter(c context.Context, a *rpcmodel.ArgFilter, res *rpcmodel.FilterRes) (err error) {
	var (
		fmsg  string
		level int8
		limit int
	)
	if fmsg, level, _, _, limit, _, err = r.s.Filter(c, a.Area, a.Message, int64(a.TypeID), a.ID, a.OID, a.MID, a.Keys, a.ReplyType); err != nil {
		log.Info("Filter arg(%+v) error(%v)", a, err)
		return
	}
	res.Result = fmsg
	res.Level = level
	res.Limit = limit
	return
}

// MFilter .
func (r *RPC) MFilter(c context.Context, arg *rpcmodel.ArgMfilter, res *map[string]*rpcmodel.FilterRes) (err error) {
	if *res, err = r.s.RPCMultiFilter(c, arg.Area, arg.Message, int64(arg.TypeID)); err != nil {
		log.Info("MFilter arg(%+v) result(%v)", res)
	}
	return
}

// FilterArea .
func (r *RPC) FilterArea(c context.Context, a *rpcmodel.ArgFilter, res *rpcmodel.FilterRes) (err error) {
	var (
		fmsg  string
		level int8
	)
	if fmsg, level, _, err = r.s.FilterArea(c, a.Area, a.Message, int64(a.TypeID)); err != nil {
		log.Info("FilterArea arg(%v) error(%v)", a)
		return
	}
	res.Result = fmsg
	res.Level = level
	return
}

// MFilterArea .
func (r *RPC) MFilterArea(c context.Context, arg *rpcmodel.ArgMfilter, res *map[string]*rpcmodel.FilterRes) (err error) {
	if *res, err = r.s.MFilterArea(c, arg.Area, arg.Message, int64(arg.TypeID)); err != nil {
		log.Info("MFilterArea arg(%v) result(%v)", arg, res)
	}
	return
}
