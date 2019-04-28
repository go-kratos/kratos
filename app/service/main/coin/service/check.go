package service

import (
	"context"
	"fmt"

	pb "go-common/app/service/main/coin/api"
	"go-common/app/service/main/coin/dao"
	mml "go-common/app/service/main/member/api"
	mml2 "go-common/app/service/main/member/model"
	"go-common/app/service/main/member/model/block"
	bml "go-common/app/service/main/member/model/block"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup.v2"
)

func boolToInt32(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

func blockStatusToSilence(status bml.BlockStatus) int32 {
	return boolToInt32(status == bml.BlockStatusForever || status == bml.BlockStatusLimit)
}

func identificationStatus(realNameStatus mml2.RealnameStatus) int32 {
	return boolToInt32(realNameStatus == mml2.RealnameStatusTrue)
}

// User .
func (s *Service) checkUser(c context.Context, mid int64) (err error) {
	eg := errgroup.WithCancel(c)
	var rank, identification int32
	eg.Go(func(c context.Context) error {
		var mb *mml.MemberInfoReply
		var e error
		if mb, e = s.memRPC.Member(c, &mml.MemberMidReq{Mid: mid}); (e != nil) || (mb == nil) {
			log.Error("d.mRPC.Member(%d) err(%v) (%+v)", mid, e, mb)
			return nil
		}
		if mb.LevelInfo.Cur < 1 {
			return ecode.UserLevelLow
		}
		rank = int32(mb.BaseInfo.Rank)
		return nil
	})
	eg.Go(func(c context.Context) error {
		mb, e := s.memRPC.BlockInfo(c, &mml.MemberMidReq{Mid: mid})
		if (e != nil) || (mb == nil) {
			log.Error("d.mRPC.Member(%d) err(%v) (%+v)", mid, e, mb)
			return nil
		}
		if blockStatusToSilence(block.BlockStatus(uint8(mb.BlockStatus))) == 1 {
			return ecode.UserDisabled
		}
		return nil
	})
	eg.Go(func(c context.Context) error {
		mb, e := s.memRPC.RealnameStatus(c, &mml.MemberMidReq{Mid: mid})
		if (e != nil) || (mb == nil) {
			log.Error("d.mRPC.Member(%d) err(%v) (%+v)", mid, e, mb)
			return nil
		}
		identification = identificationStatus(mml2.RealnameStatus(mb.RealnameStatus))
		return nil
	})
	eg.Go(func(c context.Context) error {
		mb, e := s.memRPC.Moral(c, &mml.MemberMidReq{Mid: mid})
		if (e != nil) || (mb == nil) {
			log.Error("d.mRPC.Member(%d) err(%v) (%+v)", mid, e, mb)
			return nil
		}
		if (int32(mb.Moral) / 100) < 60 {
			return ecode.LackOfScores
		}
		return nil
	})
	eg.Go(func(c context.Context) error {
		pass, e := s.coinDao.PassportDetail(c, mid)
		if e != nil {
			log.Error("d.OldMyInfo(%d) err %v", mid, e)
			return nil
		}
		if !pass.BindEmail && !pass.BindTel {
			return ecode.UserInactive
		}
		if !pass.BindTel {
			return ecode.MobileNoVerfiy
		}
		return nil
	})
	if err = eg.Wait(); err != nil {
		dao.PromError("check:user")
		return
	}
	if identification == 0 && rank == 5000 {
		err = ecode.UserNoMember
	}
	return
}

// addCoinCheck check whether user can add coin
func (s *Service) addCoinCheck(c context.Context, mid, aid, tp, multiply, maxCoin, upmid int64) (err error) {
	var (
		added int64
		exist bool
	)
	if _, ok := s.businesses[tp]; !ok {
		err = ecode.RequestErr
		return
	}
	if upmid == mid {
		log.Errorv(c, log.KV("log", "user can not add coin to self archive"), log.KV("mid", mid))
		err = ecode.CoinCannotAddToSelf
		return
	}
	if multiply > maxCoin {
		log.Errorv(c, log.KV("log", fmt.Sprintf("multiply(%d) can not bigger than maxCoin(%d)", multiply, maxCoin)), log.KV("mid", mid))
		err = ecode.CoinIllegaMultiply
		return
	}
	if err = s.checkUser(c, mid); err != nil {
		log.Errorv(c, log.KV("log", "checkUser error"), log.KV("mid", mid), log.KV("err", err))
		return
	}
	if exist, err = s.coinDao.ExpireCoinAdded(c, mid); err == nil && exist {
		added, _ = s.coinDao.CoinsAddedCache(c, mid, aid, tp)
	}
	if !exist || (added == 0) {
		if added, err = s.coinDao.CoinsAddedByMid(c, mid, aid, tp); err != nil {
			return
		}
		s.cache.Do(c, func(c context.Context) {
			s.coinDao.SetCoinAddedCache(c, mid, aid, tp, added)
			if !exist {
				s.loadUserCoinAddedCache(c, mid)
			}
		})
	}
	if added+multiply > maxCoin {
		log.Errorv(c, log.KV("log", "add too much coins"), log.KV("mid", mid), log.KV("err", err))
		err = ecode.CoinOverMax
		return
	}
	var coins *pb.UserCoinsReply
	if coins, err = s.UserCoins(c, &pb.UserCoinsReq{Mid: mid}); err != nil {
		return
	}
	if coins.Count < (float64)(multiply) {
		log.Errorv(c, log.KV("log", "have not enough money"), log.KV("mid", mid), log.KV("coins", coins.Count))
		err = ecode.LackOfCoins
	}
	return
}
