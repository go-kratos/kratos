package service

import (
	"context"
	"fmt"
	"time"

	pb "go-common/app/service/main/coin/api"
	"go-common/app/service/main/coin/dao"
	coin "go-common/app/service/main/coin/model"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	_maxEXP          = 50
	_deltaEXP        = 10
	_maxArchiveSize  = 300
	_maxRandVal      = 999999
	_businessArchive = 1
)

// WebAddCoin http api old api some day need remove
func (s *Service) WebAddCoin(c context.Context, mid, upmid, maxCoin, aid, tp, multiply int64, typeid int16) (err error) {
	arg := pb.AddCoinReq{
		IP:       metadata.String(c, metadata.RemoteIP),
		Mid:      mid,
		Upmid:    upmid,
		MaxCoin:  maxCoin,
		Aid:      aid,
		Business: s.businesses[tp].Name,
		Number:   multiply,
		Typeid:   int32(typeid),
		PubTime:  0,
	}
	if tp == _businessArchive {
		err = ecode.RequestErr
		return
	}
	_, err = s.AddCoin(c, &arg)
	return
}

// AddCoin add coin to archive.
func (s *Service) AddCoin(c context.Context, arg *pb.AddCoinReq) (res *pb.AddCoinReply, err error) {
	res = &pb.AddCoinReply{}
	tp, err := s.MustCheckBusiness(arg.Business)
	if err != nil {
		return
	}
	var (
		record      *coin.Record
		now         = time.Now()
		ts          = now.Unix()
		mid         = arg.Mid
		upmid       = arg.Upmid
		maxCoin     = arg.MaxCoin
		aid         = arg.Aid
		multiply    = arg.Number
		typeid      = arg.Typeid
		pubTime     = arg.PubTime
		ip          = arg.IP
		mergeTarget = s.mergeTarget(arg.Business, arg.Aid)
	)
	if maxCoin <= 0 || maxCoin > 2 {
		maxCoin = 2
	}
	if err = s.addCoinCheck(c, mid, aid, tp, multiply, maxCoin, upmid); err != nil {
		return
	}
	if err = s.addCoin(c, tp, mid, upmid, aid, multiply, ts, ip); err != nil {
		return
	}
	if tp == _businessArchive {
		s.cache.Do(c, func(c context.Context) {
			s.sendArcMsg(c, mid, aid, tp, upmid, multiply, pubTime, ip, typeid, ts)
		})
	}
	record = &coin.Record{
		Aid:       aid,
		AvType:    tp,
		Business:  s.businesses[tp].Name,
		Up:        upmid,
		Mid:       mid,
		Multiply:  multiply,
		IPV6:      ip,
		Timestamp: ts,
	}
	if mergeTarget > 0 {
		targetCount, err1 := s.coinDao.RawItemCoin(context.Background(), mergeTarget, tp)
		if err1 == nil {
			s.cache.Do(c, func(c context.Context) {
				targetCount += multiply
				s.coinDao.AddCacheItemCoin(c, mergeTarget, targetCount, tp)
				s.coinDao.PubStat(c, mergeTarget, tp, targetCount)
			})
		}
	}
	s.job.Do(c, func(c context.Context) {
		s.coinDao.PubCoinJob(c, record.Aid, record)
	})
	var count int64
	// get count and send stat 为了尽可能实时
	if count, err = s.coinDao.RawItemCoin(c, aid, tp); err != nil {
		return
	}
	count += multiply
	s.cache.Do(c, func(c context.Context) {
		s.coinDao.AddCacheItemCoin(c, aid, count, tp)
		s.coinDao.PubStat(c, aid, tp, count)
	})
	return
}

func (s *Service) addCoin(c context.Context, tp, mid, upmid, aid, number int64, ts int64, ip string) (err error) {
	var count, upCount float64
	var tx *xsql.Tx
	if tx, err = s.coinDao.BeginTran(c); err != nil {
		return
	}
	if count, err = s.coinDao.TxUserCoin(c, tx, mid); err != nil {
		tx.Rollback()
		return
	}
	if upCount, err = s.coinDao.TxUserCoin(c, tx, upmid); err != nil {
		tx.Rollback()
		return
	}
	if err = s.coinDao.TxUpdateCoins(c, tx, mid, -float64(number)); err != nil {
		tx.Rollback()
		return
	}
	if err = s.coinDao.TxUpdateCoins(c, tx, upmid, float64(number)/10); err != nil {
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		dao.PromError("coin:dbAddCoin")
		log.Errorv(c, log.KV("log", "dbAddCoin commit"), log.KV("err", err))
		return
	}
	var to, upperTo float64
	to = Round(count - float64(number))
	upperTo = Round(upCount + float64(number)/10)
	s.cache.Do(c, func(c context.Context) {
		s.coinDao.AddCacheUserCoin(c, mid, to)
		s.coinDao.AddCacheUserCoin(c, upmid, upperTo)
		s.coinDao.IncrCoinAddedCache(c, mid, aid, tp, number)
	})
	b := s.businesses[tp]
	reason := fmt.Sprintf(b.AddCoinReason, aid, number)
	upperReason := fmt.Sprintf(b.AddCoinUpperReason, aid, number, mid)
	s.coinDao.AddLog(mid, ts, count, to, reason, ip, "", aid, coin.TypeSend)
	s.coinDao.AddLog(upmid, ts, upCount, upperTo, upperReason, ip, "", aid, coin.TypeReceive)
	return
}

// List get coin added list.
func (s *Service) List(c context.Context, arg *pb.ListReq) (res *pb.ListReply, err error) {
	res = &pb.ListReply{}
	tp, err := s.MustCheckBusiness(arg.Business)
	if err != nil {
		return
	}
	res.List, err = s.coinDao.CoinList(c, arg.Mid, tp, arg.Ts-86400*30, _maxArchiveSize)
	return
}

func (s *Service) loadUserCoinAddedCache(c context.Context, mid int64) (err error) {
	var (
		addedMap map[int64]int64
	)
	if addedMap, err = s.coinDao.UserCoinsAdded(c, mid); err != nil {
		return
	}
	err = s.coinDao.SetCoinAddedsCache(c, mid, addedMap)
	return
}

// UpdateAddCoin .
func (s *Service) UpdateAddCoin(c context.Context, arg *pb.UpdateAddCoinReq) (res *pb.UpdateAddCoinReply, err error) {
	res = &pb.UpdateAddCoinReply{}
	tp, err := s.MustCheckBusiness(arg.Business)
	if err != nil {
		return
	}
	if err = s.coinDao.InsertCoinArchive(c, arg.Aid, tp, arg.Mid, arg.Timestamp, arg.Number); err != nil {
		return
	}
	if err = s.coinDao.InsertCoinMember(c, arg.Aid, tp, arg.Mid, arg.Timestamp, arg.Number, arg.Up); err != nil {
		return
	}
	if err = s.coinDao.UpdateItemCoinCount(c, arg.Aid, tp, arg.Number); err != nil {
		return
	}
	target := s.mergeTarget(arg.Business, arg.Aid)
	if target > 0 {
		if err = s.coinDao.UpdateItemCoinCount(c, target, tp, arg.Number); err != nil {
			return
		}
	}
	// user added coin for upMid of archive.
	if tp == _businessArchive {
		err = s.coinDao.UpdateCoinMemberCount(c, arg.Mid, arg.Up, arg.Number)
	}
	return
}

// ItemUserCoins get coins added of archive.
func (s *Service) ItemUserCoins(c context.Context, arg *pb.ItemUserCoinsReq) (res *pb.ItemUserCoinsReply, err error) {
	res = &pb.ItemUserCoinsReply{}
	tp, err := s.MustCheckBusiness(arg.Business)
	if err != nil {
		return
	}
	var (
		ok bool
	)
	if ok, _ = s.coinDao.ExpireCoinAdded(c, arg.Mid); ok {
		res.Number, _ = s.coinDao.CoinsAddedCache(c, arg.Mid, arg.Aid, tp)
	} else {
		res.Number, err = s.coinDao.CoinsAddedByMid(c, arg.Mid, arg.Aid, tp)
		s.cache.Do(c, func(c context.Context) {
			s.loadUserCoinAddedCache(c, arg.Mid)
		})
	}
	return
}

// AddedCoins get coin added to up.
func (s *Service) AddedCoins(c context.Context, mid, upMid int64) (a *coin.AddCoins, err error) {
	count, err := s.coinDao.AddedCoins(c, mid, upMid)
	if err != nil {
		return
	}
	a = &coin.AddCoins{Count: count}
	return
}

// UpdateItemCoins set archive coin added.
func (s *Service) UpdateItemCoins(c context.Context, aid, tp, coins int64) (err error) {
	var affect int64
	if affect, err = s.coinDao.UpdateItemCoins(c, aid, tp, coins, time.Now()); err != nil {
		log.Error("s.coinDao.UpdateItemCoins(%d, %d) error(%v)", aid, coins, err)
	}
	if affect != 1 {
		err = ecode.NothingFound
		return
	}
	log.Infov(c, log.KV("log", "UpdateItemCoins"), log.KV("aid", aid), log.KV("type", tp), log.KV("coin", coins))
	if err = s.coinDao.PubStat(c, aid, tp, coins); err != nil {
		log.Error("s.coinDao.PubStat(%d, %d) error(%v)", aid, coins, err)
	}
	return
}

// ItemCoin get creat count cache.
func (s *Service) ItemCoin(c context.Context, aid, tp int64) (count int64, err error) {
	return s.coinDao.ItemCoin(c, aid, tp)
}

// UpdateSettle coin settle.
func (s *Service) UpdateSettle(c context.Context, aid, tp, expSub int64, describe string) (err error) {
	if _, err = s.coinDao.UpdateCoinSettleBD(c, aid, tp, expSub, describe, time.Now()); err != nil {
		log.Error("s.coinDao.UpdateCoinSettleBD(%d, %d) error(%v)", aid, expSub, err)
	}
	return
}

func (s *Service) mergeTarget(business string, aid int64) int64 {
	if s.statMerge != nil && s.statMerge.Business == business && s.statMerge.Sources[aid] {
		return s.statMerge.Target
	}
	return 0
}
