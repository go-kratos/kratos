package service

import (
	"context"

	"go-common/app/service/live/wallet/conf"
	"go-common/app/service/live/wallet/dao"
	"go-common/app/service/live/wallet/model"
	"go-common/library/cache"
	"go-common/library/log"
)

// Service struct
type Service struct {
	c        *conf.Config
	dao      *dao.Dao
	runCache *cache.Cache
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:        c,
		dao:      dao.New(c),
		runCache: cache.New(1, 1024),
	}
	return s
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}

func (s *Service) pubWalletChange(c context.Context, uid int64, action string, number int64, coinType string, platform string, destCoinType string, destNumber int64) error {
	detail, err := s.dao.Detail(c, uid)
	if err != nil {
		return err
	}
	gold := detail.Gold
	if platform == "ios" {
		gold = detail.IapGold
	}
	msg := model.WalletChangeMsg{
		Uid:             uid,
		Action:          action,
		Number:          number,
		CoinType:        coinType,
		Gold:            gold,
		Silver:          detail.Silver,
		GoldRechargeCnt: detail.GoldRechargeCnt,
		GoldPayCnt:      detail.GoldPayCnt,
		SilverPayCnt:    detail.SilverPayCnt,
		Platfrom:        platform,
		DestCoinType:    destCoinType,
		DestNumber:      destNumber,
		CostBase:        detail.CostBase,
	}
	err = s.dao.Pub(c, uid, &msg)
	return err
}

func (s *Service) execByHandler(handler Handler, c context.Context, basicParam *model.BasicParam, uid int64, params ...interface{}) (v interface{}, err error) {
	ws := new(WalletService)
	ws.c = c
	ws.s = s
	ws.SetServiceHandler(handler)
	return ws.Execute(basicParam, uid, params...)
}

func (s *Service) pubWalletChangeWithDetailSnapShot(c context.Context, uid int64, action string, number int64, coinType string, platform string, destCoinType string, destNumber int64, detail *model.DetailWithSnapShot) error {
	f := func(loc string, ctx context.Context) {
		gold := detail.Gold
		if platform == "ios" {
			gold = detail.IapGold
		}
		msg := model.WalletChangeMsg{
			Uid:             uid,
			Action:          action,
			Number:          number,
			CoinType:        coinType,
			Gold:            gold,
			Silver:          detail.Silver,
			GoldRechargeCnt: detail.GoldRechargeCnt,
			GoldPayCnt:      detail.GoldPayCnt,
			SilverPayCnt:    detail.SilverPayCnt,
			Platfrom:        platform,
			DestCoinType:    destCoinType,
			DestNumber:      destNumber,
			CostBase:        detail.CostBase,
		}
		err := s.dao.Pub(ctx, uid, &msg)
		if err != nil {
			log.Error("SubError# loc:%s value:%+v", loc, msg)
		}
	}
	se := s.runCache.Save(func() {
		f("cache", context.Background())
	})
	if se != nil {
		log.Error("runCache is full")
		f("service", c)
	}
	return nil
}
