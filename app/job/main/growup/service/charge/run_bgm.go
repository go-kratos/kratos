package charge

import (
	"context"
	"fmt"
	"time"

	model "go-common/app/job/main/growup/model/charge"

	"go-common/library/log"
	"golang.org/x/sync/errgroup"
)

var (
	_bgmDailyCharge   = "bgm_daily_charge"
	_bgmWeeklyCharge  = "bgm_weekly_charge"
	_bgmMonthlyCharge = "bgm_monthly_charge"

	_bgmDailyStatis   = "bgm_charge_daily_statis"
	_bgmWeeklyStatis  = "bgm_charge_weekly_statis"
	_bgmMonthlyStatis = "bgm_charge_monthly_statis"
)

func (s *Service) runBgm(c context.Context, date time.Time, avBgmCharge chan []*model.AvCharge) (err error) {
	startWeeklyDate = getStartWeeklyDate(date)
	startMonthlyDate = getStartMonthlyDate(date)
	var (
		readGroup     errgroup.Group
		dailyMap      = make(map[string]*model.BgmCharge)
		sourceCh      = make(chan []*model.BgmCharge, 1000)
		dailyStatisCh = make(chan []*model.BgmCharge, 1000)
		bgmCh         = make(chan []*model.BgmCharge, 1000)
	)

	readGroup.Go(func() (err error) {
		dailyMap, err = s.bgmCharges(c, date, sourceCh, avBgmCharge)
		if err != nil {
			log.Error("s.bgmCharges error(%v)", err)
			return
		}
		log.Info("get bgm_daily_charge finished")
		return
	})

	readGroup.Go(func() (err error) {
		defer func() {
			close(bgmCh)
			close(dailyStatisCh)
		}()
		for charges := range sourceCh {
			bgmCh <- charges
			dailyStatisCh <- charges
		}
		return
	})

	var (
		weeklyMap  map[string]*model.BgmCharge
		monthlyMap map[string]*model.BgmCharge
		statisMap  map[string]*model.BgmStatis
		weeklyCh   = make(chan map[string]*model.BgmCharge, 1)
		monthlyCh  = make(chan map[string]*model.BgmCharge, 1)
	)
	readGroup.Go(func() (err error) {
		defer func() {
			close(weeklyCh)
			close(monthlyCh)
		}()
		weeklyMap, monthlyMap, statisMap, err = s.handleBgm(c, date, bgmCh)
		if err != nil {
			log.Error("s.handleBgm error(%v)", err)
			return
		}
		weeklyCh <- weeklyMap
		monthlyCh <- monthlyMap
		log.Info("handleBgm finished")
		return
	})

	var (
		dateStatis = &SectionEntries{}
		bgmDaily   = make(chan []*model.Archive, 2000)
		bgmWeekly  = make(chan []*model.Archive, 1)
		bgmMonthly = make(chan []*model.Archive, 1)
	)
	readGroup.Go(func() (err error) {
		defer close(bgmDaily)
		for bgms := range dailyStatisCh {
			bgmDaily <- transBgm2Archive(bgms)
		}
		return
	})
	readGroup.Go(func() (err error) {
		defer close(bgmWeekly)
		bgmWeekly <- transBgmMap2Archive(<-weeklyCh)
		return
	})
	readGroup.Go(func() (err error) {
		defer close(bgmMonthly)
		bgmMonthly <- transBgmMap2Archive(<-monthlyCh)
		return
	})

	readGroup.Go(func() (err error) {
		dateStatis.daily, err = s.handleDateStatis(c, bgmDaily, date, _bgmDailyStatis)
		if err != nil {
			log.Error("s.handleDateStatis(%s) error(%v)", _bgmDailyStatis, err)
		}
		return
	})
	readGroup.Go(func() (err error) {
		dateStatis.weekly, err = s.handleDateStatis(c, bgmWeekly, startWeeklyDate, _bgmWeeklyStatis)
		if err != nil {
			log.Error("s.handleDateStatis(%s) error(%v)", _bgmWeeklyStatis, err)
		}
		return
	})
	readGroup.Go(func() (err error) {
		dateStatis.monthly, err = s.handleDateStatis(c, bgmMonthly, startMonthlyDate, _bgmMonthlyStatis)
		if err != nil {
			log.Error("s.handleDateStatis(%s) error(%v)", _bgmMonthlyStatis, err)
		}
		return
	})

	if err = readGroup.Wait(); err != nil {
		log.Error("run readGroup.Wait error(%v)", err)
		return
	}

	{
		if len(dailyMap) == 0 {
			err = fmt.Errorf("Error: insert 0 bgm_daily_charge")
			return
		}
		if len(weeklyMap) == 0 {
			err = fmt.Errorf("Error: insert 0 bgm_weekly_charge")
			return
		}
		if len(monthlyMap) == 0 {
			err = fmt.Errorf("Error: insert 0 bgm_monthly_charge")
			return
		}
		if len(statisMap) == 0 {
			err = fmt.Errorf("Error: insert 0 bgm_charge_statis")
			return
		}
		if len(dateStatis.daily) == 0 {
			err = fmt.Errorf("Error: insert 0 bgm_charge_daily_statis")
			return
		}
		if len(dateStatis.weekly) == 0 {
			err = fmt.Errorf("Error: insert 0 bgm_charge_weekly_statis")
			return
		}
		if len(dateStatis.monthly) == 0 {
			err = fmt.Errorf("Error: insert 0 bgm_charge_monthly_statis")
			return
		}
	}

	// persist
	var writeGroup errgroup.Group
	// bgm_daily_charge
	writeGroup.Go(func() (err error) {
		err = s.bgmDBStore(c, _bgmDailyCharge, dailyMap)
		if err != nil {
			log.Error("s.bgmDBStore bgm_daily_charge error(%v)", err)
			return
		}
		log.Info("insert bgm_daily_charge : %d", len(dailyMap))

		// bgm_weekly_charge
		err = s.bgmDBStore(c, _bgmWeeklyCharge, weeklyMap)
		if err != nil {
			log.Error("s.bgmDBStore bgm_weekly_charge error(%v)", err)
			return
		}
		log.Info("insert bgm_weekly_charge : %d", len(weeklyMap))

		// bgm_monthly_charge
		err = s.bgmDBStore(c, _bgmMonthlyCharge, monthlyMap)
		if err != nil {
			log.Error("s.bgmDBStore bgm_monthly_charge error(%v)", err)
			return
		}
		log.Info("insert bgm_monthly_charge : %d", len(monthlyMap))

		// bgm_charge_statis
		err = s.bgmStatisDBStore(c, statisMap)
		if err != nil {
			log.Error("s.bgmStatisDBStore error(%v)", err)
			return
		}
		log.Info("insert bgm_charge_statis : %d", len(statisMap))

		// bgm_charge_daily_statis
		_, err = s.dateStatisInsert(c, dateStatis.daily, _bgmDailyStatis)
		if err != nil {
			log.Error("s.dateStatisInsert error(%v)", err)
			return
		}
		log.Info("insert bgm_charge_daily_statis : %d", len(dateStatis.daily))

		// bgm_charge_weekly_statis
		_, err = s.dateStatisInsert(c, dateStatis.weekly, _bgmWeeklyStatis)
		if err != nil {
			log.Error("s.dateStatisInsert error(%v)", err)
			return
		}
		log.Info("insert bgm_charge_weekly_statis : %d", len(dateStatis.weekly))

		// bgm_charge_monthly_statis
		_, err = s.dateStatisInsert(c, dateStatis.monthly, _bgmMonthlyStatis)
		if err != nil {
			log.Error("s.dateStatisInsert error(%v)", err)
			return
		}
		log.Info("insert bgm_charge_monthly_statis : %d", len(dateStatis.monthly))
		return
	})

	//	writeGroup.Go(func() (err error) {
	//		return
	//	})
	//
	//	writeGroup.Go(func() (err error) {
	//		return
	//	})
	//
	//	writeGroup.Go(func() (err error) {
	//		return
	//	})
	//
	//	writeGroup.Go(func() (err error) {
	//		return
	//	})
	//
	//	writeGroup.Go(func() (err error) {
	//		return
	//	})
	//
	//	writeGroup.Go(func() (err error) {
	//		return
	//	})

	if err = writeGroup.Wait(); err != nil {
		log.Error("run writeGroup.Wait error(%v)", err)
	}
	return
}
