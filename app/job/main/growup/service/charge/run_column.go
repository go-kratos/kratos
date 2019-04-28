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
	_cmWeeklyCharge  = "column_weekly_charge"
	_cmMonthlyCharge = "column_monthly_charge"

	_cmDailyStatis   = "column_charge_daily_statis"
	_cmWeeklyStatis  = "column_charge_weekly_statis"
	_cmMonthlyStatis = "column_charge_monthly_statis"
)

func (s *Service) runColumn(c context.Context, date time.Time) (err error) {
	startWeeklyDate = getStartWeeklyDate(date)
	startMonthlyDate = getStartMonthlyDate(date)
	var (
		readGroup     errgroup.Group
		sourceCh      = make(chan []*model.Column, 1000)
		cmCh          = make(chan []*model.Column, 1000)
		dailyStatisCh = make(chan []*model.Column, 1000)
	)
	readGroup.Go(func() (err error) {
		err = s.columnCharges(c, date, sourceCh)
		if err != nil {
			log.Error("s.columnCharges error(%v)", err)
			return
		}
		log.Info("read column_daily_charge finished")
		return
	})

	readGroup.Go(func() (err error) {
		defer func() {
			close(cmCh)
			close(dailyStatisCh)
		}()
		for charges := range sourceCh {
			cmCh <- charges
			dailyStatisCh <- charges
		}
		return
	})

	var (
		weeklyMap  map[int64]*model.Column
		monthlyMap map[int64]*model.Column
		statisMap  map[int64]*model.ColumnStatis
		weeklyCh   = make(chan map[int64]*model.Column, 1)
		monthlyCh  = make(chan map[int64]*model.Column, 1)
	)
	readGroup.Go(func() (err error) {
		defer func() {
			close(weeklyCh)
			close(monthlyCh)
		}()
		weeklyMap, monthlyMap, statisMap, err = s.handleColumn(c, date, cmCh)
		if err != nil {
			log.Error("s.handleColumn error(%v)", err)
			return
		}
		weeklyCh <- weeklyMap
		monthlyCh <- monthlyMap
		log.Info("handleColumn finished")
		return
	})

	var (
		dateStatis = &SectionEntries{}
		cmDaily    = make(chan []*model.Archive, 2000)
		cmWeekly   = make(chan []*model.Archive, 1)
		cmMonthly  = make(chan []*model.Archive, 1)
	)
	readGroup.Go(func() (err error) {
		defer close(cmDaily)
		for cms := range dailyStatisCh {
			cmDaily <- transCm2Archive(cms)
		}
		return
	})
	readGroup.Go(func() (err error) {
		defer close(cmWeekly)
		cmWeekly <- transCmMap2Archive(<-weeklyCh)
		return
	})
	readGroup.Go(func() (err error) {
		defer close(cmMonthly)
		cmMonthly <- transCmMap2Archive(<-monthlyCh)
		return
	})

	readGroup.Go(func() (err error) {
		dateStatis.daily, err = s.handleDateStatis(c, cmDaily, date, _cmDailyStatis)
		if err != nil {
			log.Error("s.handleDateStatis(%s) error(%v)", _cmDailyStatis, err)
		}
		return
	})
	readGroup.Go(func() (err error) {
		dateStatis.weekly, err = s.handleDateStatis(c, cmWeekly, startWeeklyDate, _cmWeeklyStatis)
		if err != nil {
			log.Error("s.handleDateStatis(%s) error(%v)", _cmWeeklyStatis, err)
		}
		return
	})
	readGroup.Go(func() (err error) {
		dateStatis.monthly, err = s.handleDateStatis(c, cmMonthly, startMonthlyDate, _cmMonthlyStatis)
		if err != nil {
			log.Error("s.handleDateStatis(%s) error(%v)", _cmMonthlyStatis, err)
		}
		return
	})

	if err = readGroup.Wait(); err != nil {
		log.Error("run readGroup.Wait error(%v)", err)
		return
	}

	{
		if len(weeklyMap) == 0 {
			err = fmt.Errorf("Error: insert 0 column_weekly_charge")
			return
		}
		if len(monthlyMap) == 0 {
			err = fmt.Errorf("Error: insert 0 column_monthly_charge")
			return
		}
		if len(statisMap) == 0 {
			err = fmt.Errorf("Error: insert 0 column_charge_statis")
			return
		}
		if len(dateStatis.daily) == 0 {
			err = fmt.Errorf("Error: insert 0 column_charge_daily_statis")
			return
		}
		if len(dateStatis.weekly) == 0 {
			err = fmt.Errorf("Error: insert 0 column_charge_weekly_statis")
			return
		}
		if len(dateStatis.monthly) == 0 {
			err = fmt.Errorf("Error: insert 0 column_charge_monthly_statis")
			return
		}
	}

	// persist
	var writeGroup errgroup.Group
	// column_weekly_charge
	writeGroup.Go(func() (err error) {
		err = s.cmDBStore(c, _cmWeeklyCharge, weeklyMap)
		if err != nil {
			log.Error("s.cmDBStore column_weekly_charge error(%v)", err)
			return
		}
		log.Info("insert column_weekly_charge : %d", len(weeklyMap))

		// column_monthly_charge
		err = s.cmDBStore(c, _cmMonthlyCharge, monthlyMap)
		if err != nil {
			log.Error("s.cmDBStore column_monthly_charge error(%v)", err)
			return
		}
		log.Info("insert column_monthly_charge : %d", len(monthlyMap))

		// column_charge_statis
		err = s.cmStatisDBStore(c, statisMap)
		if err != nil {
			log.Error("s.cmStatisDBStore error(%v)", err)
			return
		}
		log.Info("insert column_charge_statis : %d", len(statisMap))

		// column_charge_daily_statis
		_, err = s.dateStatisInsert(c, dateStatis.daily, _cmDailyStatis)
		if err != nil {
			log.Error("s.dateStatisInsert error(%v)", err)
			return
		}
		log.Info("insert column_charge_daily_statis : %d", len(dateStatis.daily))

		// column_charge_weekly_statis
		_, err = s.dateStatisInsert(c, dateStatis.weekly, _cmWeeklyStatis)
		if err != nil {
			log.Error("s.dateStatisInsert error(%v)", err)
			return
		}
		log.Info("insert column_charge_weekly_statis : %d", len(dateStatis.weekly))

		// column_charge_monthly_statis
		_, err = s.dateStatisInsert(c, dateStatis.monthly, _cmMonthlyStatis)
		if err != nil {
			log.Error("s.dateStatisInsert error(%v)", err)
			return
		}
		log.Info("insert column_charge_monthly_statis : %d", len(dateStatis.monthly))
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

	if err = writeGroup.Wait(); err != nil {
		log.Error("run writeGroup.Wait error(%v)", err)
	}
	return
}
