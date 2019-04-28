package bws

import (
	"context"

	bwsmdl "go-common/app/interface/main/activity/model/bws"
	accapi "go-common/app/service/main/account/api"
	"go-common/library/ecode"
	"go-common/library/log"
)

// Lottery get lottery account.
func (s *Service) Lottery(c context.Context, bid, loginMid, aid int64, day string) (data *bwsmdl.LotteryUser, err error) {
	var (
		mid     int64
		accData *accapi.InfoReply
	)
	if _, ok := s.lotteryMids[loginMid]; !ok {
		err = ecode.ActivityNotLotteryAdmin
		return
	}
	if _, ok := s.lotteryAids[aid]; !ok {
		err = ecode.ActivityNotLotteryAchieve
		return
	}
	if _, err = s.Achievement(c, &bwsmdl.ParamID{Bid: bid, ID: aid}); err != nil {
		return
	}
	if mid, err = s.dao.CacheLotteryMid(c, aid, day); err != nil || mid == 0 {
		err = ecode.ActivityLotteryFail
		return
	}
	log.Warn("Lottery bid(%d) loginMid(%d) aid(%d) lotteryMid(%d)", bid, loginMid, aid, mid)
	data = &bwsmdl.LotteryUser{Mid: mid}
	if accData, err = s.accClient.Info3(c, &accapi.MidReq{Mid: mid}); err != nil {
		log.Error("Lottery s.accRPC.Info3(%d) error(%v)", mid, err)
		err = nil
		return
	}
	if accData != nil && accData.Info != nil {
		data = &bwsmdl.LotteryUser{Mid: mid, Name: accData.Info.Name, Face: accData.Info.Face}
	}
	return
}

// LotteryCheck .
func (s *Service) LotteryCheck(c context.Context, mid, aid int64, day string) (data []int64, err error) {
	if !s.isAdmin(mid) {
		err = ecode.ActivityNotAdmin
		return
	}
	return s.dao.CacheLotteryMids(c, aid, day)
}
