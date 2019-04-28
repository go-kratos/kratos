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
	_avChargeDailyStatis   = "av_charge_daily_statis"
	_avChargeWeeklyStatis  = "av_charge_weekly_statis"
	_avChargeMonthlyStatis = "av_charge_monthly_statis"

	_avWeeklyCharge  = "av_weekly_charge"
	_avMonthlyCharge = "av_monthly_charge"
)

func (s *Service) runVideo(c context.Context, date time.Time, avBgmCharge chan []*model.AvCharge) (err error) {
	startWeeklyDate = getStartWeeklyDate(date)
	startMonthlyDate = getStartMonthlyDate(date)
	var (
		readGroup     errgroup.Group
		sourceCh      = make(chan []*model.AvCharge, 1000)
		avChargeCh    = make(chan []*model.AvCharge, 1000)
		upChargeCh    = make(chan []*model.AvCharge, 1000)
		dailyStatisCh = make(chan []*model.AvCharge, 1000)
	)
	readGroup.Go(func() (err error) {
		err = s.avDailyCharges(c, date, sourceCh)
		if err != nil {
			log.Error("s.avCharge.AvCharges error(%v)", err)
			return
		}
		log.Info("read av_daily_charge finished")
		return
	})

	readGroup.Go(func() (err error) {
		defer func() {
			close(avChargeCh)
			close(upChargeCh)
			close(dailyStatisCh)
			close(avBgmCharge)
		}()
		for charges := range sourceCh {
			avChargeCh <- charges
			upChargeCh <- charges
			dailyStatisCh <- charges
			avBgmCharge <- charges
		}
		return
	})

	var (
		weeklyChargeMap  map[int64]*model.AvCharge
		monthlyChargeMap map[int64]*model.AvCharge
		chargeStatisMap  map[int64]*model.AvChargeStatis
		weeklyCh         = make(chan map[int64]*model.AvCharge, 1)
		monthlyCh        = make(chan map[int64]*model.AvCharge, 1)
	)
	readGroup.Go(func() (err error) {
		defer func() {
			close(weeklyCh)
			close(monthlyCh)
		}()
		weeklyChargeMap, monthlyChargeMap, chargeStatisMap, err = s.handleAvCharge(c, date, avChargeCh)
		if err != nil {
			log.Error("s.handleAvCharge error(%v)", err)
			return
		}
		weeklyCh <- weeklyChargeMap
		monthlyCh <- monthlyChargeMap
		log.Info("handleAvCharge finished")
		return
	})

	var (
		dateStatis = &SectionEntries{}
		avDaily    = make(chan []*model.Archive, 2000)
		avWeekly   = make(chan []*model.Archive, 1)
		avMonthly  = make(chan []*model.Archive, 1)
	)
	readGroup.Go(func() (err error) {
		defer close(avDaily)
		for avs := range dailyStatisCh {
			avDaily <- transAv2Archive(avs)
		}
		return
	})
	readGroup.Go(func() (err error) {
		defer close(avWeekly)
		avWeekly <- transAvMap2Archive(<-weeklyCh)
		return
	})
	readGroup.Go(func() (err error) {
		defer close(avMonthly)
		avMonthly <- transAvMap2Archive(<-monthlyCh)
		return
	})

	readGroup.Go(func() (err error) {
		dateStatis.daily, err = s.handleDateStatis(c, avDaily, date, _avChargeDailyStatis)
		if err != nil {
			log.Error("s.handleDateStatis(%s) error(%v)", _avChargeDailyStatis, err)
		}
		return
	})
	readGroup.Go(func() (err error) {
		dateStatis.weekly, err = s.handleDateStatis(c, avWeekly, startWeeklyDate, _avChargeWeeklyStatis)
		if err != nil {
			log.Error("s.handleDateStatis(%s) error(%v)", _avChargeWeeklyStatis, err)
		}
		return
	})
	readGroup.Go(func() (err error) {
		dateStatis.monthly, err = s.handleDateStatis(c, avMonthly, startMonthlyDate, _avChargeMonthlyStatis)
		if err != nil {
			log.Error("s.handleDateStatis(%s) error(%v)", _avChargeMonthlyStatis, err)
		}
		return
	})

	var (
		daily   map[int64]*model.UpCharge
		weekly  map[int64]*model.UpCharge
		monthly map[int64]*model.UpCharge
	)
	readGroup.Go(func() (err error) {
		daily, weekly, monthly, err = s.calUpCharge(c, date, upChargeCh)
		if err != nil {
			log.Error("s.calUpCharge error(%v)", err)
			return
		}
		log.Info("s.calUpCharge finished")
		return
	})

	if err = readGroup.Wait(); err != nil {
		log.Error("run readGroup.Wait error(%v)", err)
		return
	}

	{
		if len(weeklyChargeMap) == 0 {
			err = fmt.Errorf("Error: insert 0 av_weekly_charge")
			return
		}
		if len(monthlyChargeMap) == 0 {
			err = fmt.Errorf("Error: insert 0 av_monthly_charge")
			return
		}
		if len(chargeStatisMap) == 0 {
			err = fmt.Errorf("Error: insert 0 av_charge_statis")
			return
		}
		if len(daily) == 0 {
			err = fmt.Errorf("Error: insert 0 up_daily_charge")
			return
		}
		if len(weekly) == 0 {
			err = fmt.Errorf("Error: insert 0 up_weekly_charge")
			return
		}
		if len(monthly) == 0 {
			err = fmt.Errorf("Error: insert 0 up_monthly_charge")
			return
		}
		if len(dateStatis.daily) == 0 {
			err = fmt.Errorf("Error: insert 0 av_charge_daily_statis")
			return
		}
		if len(dateStatis.weekly) == 0 {
			err = fmt.Errorf("Error: insert 0 av_charge_weekly_statis")
			return
		}
		if len(dateStatis.monthly) == 0 {
			err = fmt.Errorf("Error: insert 0 av_charge_monthly_statis")
			return
		}
	}

	// persist
	var writeGroup errgroup.Group
	// av_weekly_charge
	writeGroup.Go(func() (err error) {
		err = s.AvChargeDBStore(c, _avWeeklyCharge, weeklyChargeMap)
		if err != nil {
			log.Error("s.AvChargeDBStore av_weekly_charge error(%v)", err)
			return
		}
		log.Info("insert av_weekly_charge : %d", len(weeklyChargeMap))
		return
	})

	// av_monthly_charge
	writeGroup.Go(func() (err error) {
		err = s.AvChargeDBStore(c, _avMonthlyCharge, monthlyChargeMap)
		if err != nil {
			log.Error("s.AvChargeDBStore av_monthly_charge error(%v)", err)
			return
		}
		log.Info("insert av_monthly_charge : %d", len(monthlyChargeMap))

		// av_charge_statis
		err = s.AvChargeStatisDBStore(c, chargeStatisMap)
		if err != nil {
			log.Error("s.AvChargeStatisDBStore error(%v)", err)
			return
		}
		log.Info("insert av_charge_statis : %d", len(chargeStatisMap))

		// up_daily_charge
		err = s.BatchInsertUpCharge(c, "up_daily_charge", daily)
		if err != nil {
			log.Error("s.BatchInsertUpDailyCharge error(%v)", err)
			return
		}
		log.Info("insert up_daily_charge : %d", len(daily))

		// up_weekly_charge
		err = s.BatchInsertUpCharge(c, "up_weekly_charge", weekly)
		if err != nil {
			log.Error("s.BatchInsertUpWeeklyCharge error(%v)", err)
			return
		}
		log.Info("insert up_weekly_charge : %d", len(weekly))

		// up_monthly_charge
		err = s.BatchInsertUpCharge(c, "up_monthly_charge", monthly)
		if err != nil {
			log.Error("s.BatchInsertUpMonthlyCharge error(%v)", err)
			return
		}
		log.Info("insert up_monthly_charge : %d", len(monthly))

		// av_charge_daily_statis
		_, err = s.dateStatisInsert(c, dateStatis.daily, _avChargeDailyStatis)
		if err != nil {
			log.Error("s.dateStatisInsert error(%v)", err)
			return
		}
		log.Info("insert av_charge_daily_statis : %d", len(dateStatis.daily))

		// av_charge_weekly_statis
		_, err = s.dateStatisInsert(c, dateStatis.weekly, _avChargeWeeklyStatis)
		if err != nil {
			log.Error("s.dateStatisInsert error(%v)", err)
			return
		}
		log.Info("insert av_charge_weekly_statis : %d", len(dateStatis.weekly))

		// av_charge_monthly_statis
		_, err = s.dateStatisInsert(c, dateStatis.monthly, _avChargeMonthlyStatis)
		if err != nil {
			log.Error("s.dateStatisInsert error(%v)", err)
			return
		}
		log.Info("insert av_charge_monthly_statis : %d", len(dateStatis.monthly))
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
	//
	//	writeGroup.Go(func() (err error) {
	//		return
	//	})

	if err = writeGroup.Wait(); err != nil {
		log.Error("run writeGroup.Wait error(%v)", err)
	}
	return
}
