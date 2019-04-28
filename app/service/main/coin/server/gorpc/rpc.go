package rpc

import (
	"time"

	pb "go-common/app/service/main/coin/api"
	"go-common/app/service/main/coin/conf"
	coin "go-common/app/service/main/coin/model"
	"go-common/app/service/main/coin/service"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
)

// RPC define rpc.
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

// AddCoins add  coin to archive.
func (r *RPC) AddCoins(c context.Context, ac *coin.ArgAddCoin, res *struct{}) (err error) {
	tp, err := r.s.CheckBusiness(ac.Business)
	if err != nil {
		return
	}
	if tp > 0 {
		ac.AvType = tp
	}
	if ac.AvType == 0 {
		ac.AvType = 1
	}
	b, err := r.s.GetBusinessName(ac.AvType)
	if err != nil {
		return
	}
	arg := &pb.AddCoinReq{
		IP:       ac.RealIP,
		Mid:      ac.Mid,
		Upmid:    ac.UpMid,
		MaxCoin:  ac.MaxCoin,
		Aid:      ac.Aid,
		Business: b,
		Number:   ac.Multiply,
		Typeid:   int32(ac.TypeID),
		PubTime:  ac.PubTime,
	}
	_, err = r.s.AddCoin(c, arg)
	return
}

// ArchiveUserCoins archive coins.
func (r *RPC) ArchiveUserCoins(c context.Context, m *coin.ArgCoinInfo, res *coin.ArchiveUserCoins) (err error) {
	tp, err := r.s.CheckBusiness(m.Business)
	if err != nil {
		return
	}
	if tp > 0 {
		m.AvType = tp
	}
	if m.AvType == 0 {
		m.AvType = 1
	}
	b, err := r.s.GetBusinessName(m.AvType)
	if err != nil {
		return
	}
	arg := pb.ItemUserCoinsReq{
		Mid:      m.Mid,
		Aid:      m.Aid,
		Business: b,
	}
	var rr *pb.ItemUserCoinsReply
	if rr, err = r.s.ItemUserCoins(c, &arg); err == nil && rr != nil {
		*res = coin.ArchiveUserCoins{Multiply: rr.Number}
	}
	return
}

// UserCoins get user coins.
func (r *RPC) UserCoins(c context.Context, arg *coin.ArgCoinInfo, res *float64) (err error) {
	reply, err := r.s.UserCoins(c, &pb.UserCoinsReq{Mid: arg.Mid})
	if reply != nil {
		*res = reply.Count
	}
	return
}

// ModifyCoin modify user coin.
func (r *RPC) ModifyCoin(c context.Context, arg *coin.ArgModifyCoin, res *float64) (err error) {
	req := &pb.ModifyCoinsReq{
		Mid:       arg.Mid,
		Count:     arg.Count,
		Reason:    arg.Reason,
		IP:        arg.IP,
		Operator:  arg.Operator,
		CheckZero: int32(arg.CheckZero),
		Ts:        time.Now().Unix(),
	}
	reply, err := r.s.ModifyCoins(c, req)
	if err != nil {
		return
	}
	*res = reply.Result
	return
}

// List coin added list.
func (r *RPC) List(c context.Context, arg *coin.ArgList, res *[]*coin.List) (err error) {
	tp, err := r.s.CheckBusiness(arg.Business)
	if err != nil {
		return
	}
	if tp > 0 {
		arg.TP = tp
	}
	b, err := r.s.GetBusinessName(arg.TP)
	if err != nil {
		return
	}
	req := &pb.ListReq{
		Mid:      arg.Mid,
		Business: b,
		Ts:       time.Now().Unix(),
	}
	reply, err := r.s.List(c, req)
	if err != nil {
		return
	}
	lists := []*coin.List{}
	for _, r := range reply.List {
		lists = append(lists, &coin.List{
			Aid:      r.Aid,
			Multiply: r.Number,
			Ts:       r.Ts,
			IP:       r.IP,
		})
	}
	*res = lists
	return
}

// UserLog user log
func (r *RPC) UserLog(c context.Context, arg *coin.ArgLog, res *[]*coin.Log) (err error) {
	req := &pb.CoinsLogReq{
		Mid:       arg.Mid,
		Recent:    arg.Recent,
		Translate: arg.Translate,
	}
	reply, err := r.s.CoinsLog(c, req)
	lists := []*coin.Log{}
	for _, r := range reply.List {
		lists = append(lists, &coin.Log{
			From:      r.From,
			To:        r.To,
			IP:        r.IP,
			Desc:      r.Desc,
			TimeStamp: r.TimeStamp,
		})
	}
	*res = lists
	return
}

// AddUserCoinExp add user coin exp for job
func (r *RPC) AddUserCoinExp(c context.Context, arg *coin.ArgAddUserCoinExp, res *struct{}) (err error) {
	tp, err := r.s.CheckBusiness(arg.BusinessName)
	if err != nil {
		return
	}
	if tp > 0 {
		arg.Business = tp
	}
	b, err := r.s.GetBusinessName(arg.Business)
	if err != nil {
		return
	}
	req := &pb.AddUserCoinExpReq{
		IP:       arg.RealIP,
		Mid:      arg.Mid,
		Business: b,
		Number:   arg.Number,
	}
	_, err = r.s.AddUserCoinExp(c, req)
	return
}

// UpdateAddCoin update db after add coin for job.
func (r *RPC) UpdateAddCoin(c context.Context, arg *coin.Record, res *struct{}) (err error) {
	tp, err := r.s.CheckBusiness(arg.Business)
	if err != nil {
		return
	}
	if tp > 0 {
		arg.AvType = tp
	}
	b, err := r.s.GetBusinessName(arg.AvType)
	if err != nil {
		return
	}
	req := &pb.UpdateAddCoinReq{
		Aid:       arg.Aid,
		Mid:       arg.Mid,
		Up:        arg.Up,
		Timestamp: arg.Timestamp,
		Number:    arg.Multiply,
		Business:  b,
		IPV6:      arg.IPV6,
	}
	_, err = r.s.UpdateAddCoin(c, req)
	return
}

// TodayExp .
func (r *RPC) TodayExp(c context.Context, arg *coin.ArgMid, res *int64) (err error) {
	req := &pb.TodayExpReq{
		Mid: arg.Mid,
	}
	reply, err := r.s.TodayExp(c, req)
	if reply != nil {
		*res = reply.Exp
	}
	return
}
