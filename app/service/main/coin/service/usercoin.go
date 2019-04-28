package service

import (
	"context"
	"regexp"

	pb "go-common/app/service/main/coin/api"
	"go-common/app/service/main/coin/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// ModifyCoins modify user coins.
func (s *Service) ModifyCoins(c context.Context, arg *pb.ModifyCoinsReq) (res *pb.ModifyCoinsReply, err error) {
	res = &pb.ModifyCoinsReply{}
	var (
		count float64
	)
	log.Infov(c, log.KV("log", "ModifyCoins"), log.KV("mid", arg.Mid), log.KV("count", arg.Count), log.KV("reason", arg.Reason), log.KV("ip", arg.IP), log.KV("operator", arg.Operator))
	count, err = s.coinDao.RawUserCoin(c, arg.Mid)
	if err != nil {
		return
	}
	res.Result = Round(count + arg.Count)
	// checkZere为0才返回硬币不足
	if arg.Count < 0 && res.Result < 0 && arg.CheckZero != 1 {
		err = ecode.LackOfCoins
		return
	}
	if err = s.coinDao.UpdateCoin(c, arg.Mid, arg.Count); err != nil {
		return
	}
	s.coinDao.AddCacheUserCoin(c, arg.Mid, res.Result)
	if arg.Operator == "nolog" {
		return
	}
	s.coinDao.AddLog(arg.Mid, arg.Ts, count, res.Result, arg.Reason, arg.IP, arg.Operator, 0, model.TypeNone)
	return
}

// UserCoins get user coins.
func (s *Service) UserCoins(c context.Context, arg *pb.UserCoinsReq) (res *pb.UserCoinsReply, err error) {
	res = &pb.UserCoinsReply{}
	res.Count, err = s.coinDao.UserCoin(c, arg.Mid)
	return
}

// CoinsLog coins log
func (s *Service) CoinsLog(c context.Context, arg *pb.CoinsLogReq) (res *pb.CoinsLogReply, err error) {
	res = &pb.CoinsLogReply{}
	res.List, err = s.coinDao.CoinLog(c, arg.Mid)
	if res.List == nil {
		res.List = []*pb.ModelLog{}
	}
	if !arg.Translate {
		return
	}
	for _, l := range res.List {
		if l.Desc != "" {
			l.Desc = translateLog(l.Desc)
		}
	}
	return
}

func translateLog(l string) (res string) {
	var match bool
	// 依赖顺序
	m := [][2]string{
		{`cv Rating for (?P<var1>[0-9]+) : ([0-9]+) from ([0-9]+)`, `专栏 cv$var1 收到打赏`},
		{`cv Rating for (?P<var1>[0-9]+)`, `给专栏 cv$var1 打赏`},
		{`mv Rating for (?P<var1>[0-9]+) : ([0-9]+) from ([0-9]+)`, `音频 mv$var1 收到打赏`},
		{`Rating for (?P<var1>[0-9]+) : ([0-9]+) from ([0-9]+)`, `视频 av$var1 收到打赏`},
		{`cv Rating for (?P<var1>[0-9]+)`, `给专栏 $var1 打赏`},
		{`mv Rating for (?P<var1>[0-9]+)`, `给音乐 $var1 打赏`},
		{`Rating for (?P<var1>[0-9]+)`, `给视频 av$var1 打赏`},
		{`ASS:DOWNLOAD SHARE:([0-9]+)`, `其他用户下载弹幕分享积分`},
		{`通过审核（AID:(?P<var1>[0-9]+)）`, `投稿 av$var1 通过审核`},
		{`删除投稿 AID:(?P<var1>[0-9]+)`, `删除投稿 av$var1`},
		{`删除收藏 AID:(?P<var1>[0-9]+)`, `投稿 av$var1 低于100人收藏，取消过100收藏的奖励积分`},
		{`删除已审核投稿 AID:(?P<var1>[0-9]+)`, `删除投稿 av$var1`},
		{`取消审核状态（AID:(?P<var1>[0-9]+)）`, `投稿被退回 av$var1`},
		{`UPDATE:NICK:`, `修改昵称`},
		{`^管理:(?P<var1>.+) Operator : (.+)$`, `$matches[1]`},
		{`stow : (?P<var1>[0-9]+)`, `投稿 av$var1 超过100人收藏`},
		{`Buy stow limits (?P<var1>[0-9]+)`, `购买收藏上限 $var1 个`},
		{`BuyRank::([0-9]+)`, `购买标识`},
		{`Click Ads b`, `支持广告，支持网站发展`},
		{`Activity Award`, `活动奖励`},
		{`兑换`, `礼品兑换`},
		{`2015萌战活动`, `投票资格`},
	}
	for _, r := range m {
		if match, res = regexpReplace(r[0], l, r[1]); match {
			return
		}
	}
	return l
}

func regexpReplace(reg, src, temp string) (match bool, res string) {
	result := []byte{}
	pattern := regexp.MustCompile(reg)
	for _, submatches := range pattern.FindAllStringSubmatchIndex(src, -1) {
		result = pattern.ExpandString(result, temp, src, submatches)
		match = true
	}
	return match, string(result)
}
