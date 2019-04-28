package service

import (
	"context"

	pb "go-common/app/service/main/coin/api"
	"go-common/app/service/main/coin/dao"
	memmdl "go-common/app/service/main/member/api"
	"go-common/library/log"
)

// AddUserCoinExp .
func (s *Service) AddUserCoinExp(c context.Context, arg *pb.AddUserCoinExpReq) (res *pb.AddUserCoinExpReply, err error) {
	res = &pb.AddUserCoinExpReply{}
	tp, err := s.MustCheckBusiness(arg.Business)
	if err != nil {
		return
	}
	var (
		todayExp, exp int64
	)
	if todayExp, err = s.coinDao.Exp(c, arg.Mid); err != nil {
		return
	}
	if todayExp < _maxEXP {
		exp = arg.Number * _deltaEXP
		if todayExp += arg.Number * _deltaEXP; todayExp > _maxEXP {
			todayExp = _maxEXP
			exp = _deltaEXP
		}
		reason := s.businesses[tp].AddExpReason
		if err = s.addExp(c, arg.Mid, float64(exp), reason, arg.IP); err != nil {
			return
		}
		if err = s.coinDao.SetTodayExpCache(c, arg.Mid, todayExp); err != nil {
			return
		}
	}
	return
}

func (s *Service) addExp(c context.Context, mid int64, count float64, reason, ip string) (err error) {
	argExp := &memmdl.AddExpReq{
		Mid:     mid,
		Count:   count,
		Operate: "coin",
		Reason:  reason,
		Ip:      ip,
	}
	if _, err = s.memRPC.UpdateExp(c, argExp); err != nil {
		log.Errorv(c, log.KV("log", "s.coinDao.IncrExp()"), log.KV("mid", mid), log.KV("err", err), log.KV("reason", reason), log.KV("count", count))
		dao.PromError("exp:addExp")
		return
	}
	log.Infov(c, log.KV("log", "add exp"), log.KV("mid", mid), log.KV("count", count), log.KV("reason", reason), log.KV("err", err))
	return
}

// TodayExp get today coin added exp.
func (s *Service) TodayExp(c context.Context, arg *pb.TodayExpReq) (res *pb.TodayExpReply, err error) {
	res = &pb.TodayExpReply{}
	res.Exp, err = s.coinDao.Exp(c, arg.Mid)
	return
}
