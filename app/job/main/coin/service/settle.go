package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go-common/app/job/main/coin/dao"
	"go-common/app/job/main/coin/model"
	"go-common/app/service/main/archive/api"
	comarcmdl "go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
)

var lockRedo sync.Mutex
var lockSettle sync.Mutex

const (
	_expPerCoin = 1
)

func (s *Service) settleproc() {
	for {
		now := time.Now()
		tmp := now.AddDate(0, 1, 0)
		// sleep one month
		time.Sleep(time.Date(tmp.Year(), tmp.Month(), 1, 0, 0, 0, 0, tmp.Location()).Sub(now))
		// settle coin
		s.Settle(0)
	}
}

// Redo redo settle.
func (s *Service) Redo(tableID int64) (err error) {
	lockRedo.Lock()
	defer lockRedo.Unlock()
	ctx := context.TODO()
	// if tableID has been seted than run test
	if tableID != 0 {
		return s.redo(ctx, tableID)
	}
	now := time.Now()
	day := now.Day()
	if day <= 25 {
		err = fmt.Errorf("redo must after every 25th of each month, today is %d ", day)
		return
	}
	period, err := s.coinDao.HitSettlePeriod(ctx, now)
	if err != nil {
		log.Error("job s.coinDao.HitCoinPeriod error(%v)", err)
		return
	}
	return s.redo(ctx, period.ID-1)
}

func (s *Service) redo(ctx context.Context, tableID int64) (err error) {
	period, err := s.coinDao.SettlePeriod(ctx, tableID)
	if err != nil {
		log.Error("s.coinDao.SettlePeriod(%d) error(%v)", tableID, err)
		return
	}
	if err = s.coinDao.ClearCoinCount(ctx, tableID, time.Now()); err != nil {
		log.Error("s.coinDao.ClearCoinCount(%d) error(%v)", tableID, err)
		return
	}
	startTime := time.Date(period.FromYear, time.Month(period.FromMonth), period.FromDay, 0, 0, 0, 0, time.Local)
	endTime := time.Date(period.ToYear, time.Month(period.ToMonth), period.ToDay, 0, 0, 0, 0, time.Local)
	var (
		coins   map[int64]int64
		aids    []int64
		aidMids map[int64]int64
		argAids *comarcmdl.ArgAids2
		arcs    map[int64]*api.Arc
		peer    = 100
	)
	for i := 0; i < dao.SHARDING; i++ {
		if coins, err = s.coinDao.TotalCoins(ctx, i, startTime, endTime); err != nil {
			log.Error("s.coinDao.TotalCoins(%d, %v, %v) error(%v)", i, startTime, endTime, err)
			return
		}
		aids = make([]int64, 0, len(coins))
		for aid := range coins {
			if aid%1000 == 1 {
				aids = append(aids, aid/1000)
			}
		}
		var (
			length = len(aids)
			cop    = length / peer
			mod    = length % peer
		)
		aidMids = make(map[int64]int64, length)
		for i := 0; i < cop; i++ {
			argAids = &comarcmdl.ArgAids2{
				Aids: aids[i*peer : peer*(i+1)],
			}
			if arcs, err = s.arcRPC.Archives3(ctx, argAids); err != nil {
				log.Error("s.arcRPC.Archives2 error(%v)", err)
				return
			}
			for aid, arc := range arcs {
				aidMids[aid] = arc.Author.Mid
			}
		}
		if mod != 0 {
			argAids = &comarcmdl.ArgAids2{
				Aids: aids[cop*peer:],
			}
			if arcs, err = s.arcRPC.Archives3(ctx, argAids); err != nil {
				log.Error("s.arcRPC.Archives2 error(%v)", err)
				return
			}
			for aid, arc := range arcs {
				aidMids[aid] = arc.Author.Mid
			}
		}
		for aid, count := range coins {
			if mid, ok := aidMids[aid/1000]; ok {
				if err = s.coinDao.UpsertSettle(ctx, tableID, mid, aid/1000, aid%1000, count, time.Now()); err != nil {
					log.Error("s.coinDao.UpdateCoinCount(%d, %d, %d) error(%v)", tableID, aid, count, err)
					return
				}
			}
		}
	}
	return
}

// Settle do settle by table.
func (s *Service) Settle(tableID int64) (err error) {
	lockSettle.Lock()
	defer lockSettle.Unlock()
	ctx := context.TODO()
	// if tableID has been seted than run test
	if tableID != 0 {
		return s.settle(ctx, tableID)
	}
	period, err := s.coinDao.HitSettlePeriod(ctx, time.Now())
	if err != nil {
		log.Error("job s.coinDao.HitCoinPeriod error(%v)", err)
		return
	}
	return s.settle(ctx, period.ID-1)
}

// job exec at every 1th of each month
func (s *Service) settle(ctx context.Context, tableID int64) (err error) {
	var (
		settles []*model.CoinSettle
	)
	period, err := s.coinDao.SettlePeriod(ctx, tableID)
	if err != nil {
		log.Error("s.coinDao.SettlePeriod(%d) error(%v)", tableID, err)
		return
	}
	var i, maxid int64
	for {
		if settles, maxid, err = s.coinDao.Every10000(ctx, tableID, i); err != nil {
			log.Error("settle: job s.coinDao.Every10000(%d) error(%v)", tableID, err)
			time.Sleep(time.Second)
			continue
		}
		if maxid == i {
			log.Info("settle: maxid %d", maxid)
			fmt.Println("maxid", maxid)
			return
		}
		for _, settle := range settles {
			if settle.State == 1 || settle.Mid == 0 {
				continue
			}
			settle.ExpTotal = settle.CoinCount*_expPerCoin - settle.ExpSub
			if settle.ExpTotal <= 0 {
				log.Errorv(ctx, log.KV("log", "settle: ExpTotal err"), log.KV("count", settle.ExpTotal))
				continue
			}
			// add exp
			var reason string
			switch settle.AvType {
			case 1:
				reason = fmt.Sprintf("%d.%d-%d.%d视频av%d投币获得奖励", period.FromMonth, period.FromDay, period.ToMonth, period.ToDay, settle.Aid)
			case 2:
				reason = fmt.Sprintf("%d.%d-%d.%d文章cv%d投币获得奖励", period.FromMonth, period.FromDay, period.ToMonth, period.ToDay, settle.Aid)
			case 3:
				reason = fmt.Sprintf("%d.%d-%d.%d音乐mv%d投币获得奖励", period.FromMonth, period.FromDay, period.ToMonth, period.ToDay, settle.Aid)
			}
			for i := 0; i < 3; i++ {
				if err = s.addExp(ctx, settle.Mid, float64(settle.ExpTotal), reason, ""); err != nil {
					time.Sleep(time.Second)
				} else {
					break
				}
			}
			if err != nil {
				log.Errorv(ctx, log.KV("log", "s.accRPC.AddExp2"), log.KV("mid", settle.Mid), log.KV("err", err))
				continue
			}
			if err = s.coinDao.UpdateSettle(ctx, tableID, settle.ID, settle.ExpTotal, time.Now()); err != nil {
				log.Error("settle: s.coinDao.UpdateState(%d, %d, %d) error(%v)", tableID, settle.ID, settle.ExpTotal, err)
				continue
			}
			log.Info("settle: aid(%d) add exp(%d) success", settle.Aid, settle.ExpTotal)
		}
		fmt.Printf("settle:%d ,len(settle) %d\n", i, len(settles))
		i = maxid
	}
}
