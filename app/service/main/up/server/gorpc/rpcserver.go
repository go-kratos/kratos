package gorpc

import (
	"time"

	"go-common/app/service/main/up/conf"
	"go-common/app/service/main/up/model"
	"go-common/app/service/main/up/service"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
)

//RPC rpc server
type RPC struct {
	s *service.Service
}

//New create rpc
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

// Special rpc call region
func (r *RPC) Special(c context.Context, arg *model.ArgSpecial, res *[]*model.UpSpecial) (err error) {
	*res = r.s.UpsByGroup(c, arg.GroupID)
	log.Info("[rpc.Special] arg=%+v, res count=%d", arg, len(*res))
	return
}

//Info get up info
func (r *RPC) Info(c context.Context, arg *model.ArgInfo, res *model.UpInfo) (err error) {
	if arg.Mid <= 0 {
		err = ecode.RequestErr
		log.Error("[rpc.Info] error, request mid <= 0, %d", arg.Mid)
		return
	}
	isAuthor, err := r.s.Info(c, arg.Mid, uint8(arg.From))
	if err != nil {
		log.Error("[rpc.Info] error, mid=%d, from=%d, err=%s", arg.Mid, arg.From, err)
		return
	}
	*res = model.UpInfo{IsAuthor: int32(isAuthor)}
	//log.Info("[rpc.Info] mid=%d, from=%d, result=%+v", arg.Mid, arg.From, *res)
	return
}

//UpStatBase get up stat
func (r *RPC) UpStatBase(c context.Context, arg *model.ArgMidWithDate, res *model.UpBaseStat) (err error) {
	if arg.Date.IsZero() {
		// 如果没有填，则取最新的数据，如果有填，则取对应天数的数据，这里不需要做什么操作
		arg.Date = time.Now()
		// 12点更新数据，数据表为昨天日期，所以在12点以前，要读前天的表
		arg.Date = arg.Date.Add(-12*time.Hour).AddDate(0, 0, -1)
	}
	var data, e = r.s.Data.BaseUpStat(c, arg.Mid, arg.Date.Format("20060102"))
	err = e
	if err == nil {
		data.CopyTo(res)
		log.Info("[rpc.UpStatBase] arg=%+v, res=%+v", arg, res)
	} else {
		log.Error("[rpc.UpStatBase] fail arg=%+v, err=%v", arg, err)
	}
	return
}

//SetUpSwitch set up switch
func (r *RPC) SetUpSwitch(c context.Context, arg *model.ArgUpSwitch, res *model.PBSetUpSwitchRes) (err error) {
	id, err := r.s.SetSwitch(c, arg.Mid, arg.State, uint8(arg.From))
	if err == nil {
		log.Info("[rpc.SetUpSwitch] arg=%+v, res=%+v", arg, res)
	} else {
		log.Error("[rpc.SetUpSwitch] fail arg=%+v, err=%v", arg, err)
	}
	*res = model.PBSetUpSwitchRes{Id: id}
	return
}

//UpSwitch get up switch
func (r *RPC) UpSwitch(c context.Context, arg *model.ArgUpSwitch, res *model.PBUpSwitch) (err error) {
	if arg.Mid <= 0 {
		err = ecode.RequestErr
		log.Error("[rpc.UpSwitch] error, request mid <= 0, %d", arg.Mid)
		return
	}
	state, err := r.s.UpSwitchs(c, arg.Mid, uint8(arg.From))
	if err != nil {
		log.Error("[rpc.UpSwitchs] error, mid=%d, from=%d, err=%s", arg.Mid, arg.From, err)
		return
	}
	*res = model.PBUpSwitch{State: int32(state)}
	return
}
