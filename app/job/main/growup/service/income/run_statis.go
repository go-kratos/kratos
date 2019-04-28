package income

import (
	"context"
	"fmt"
	"time"

	model "go-common/app/job/main/growup/model/income"
	task "go-common/app/job/main/growup/service"

	"go-common/library/log"

	"golang.org/x/sync/errgroup"
)

// RunStatis run income statistics
func (s *Service) RunStatis(c context.Context, date time.Time) (err error) {
	var msg string
	mailReceivers := []string{"shaozhenyu@bilibili.com", "gaopeng@bilibili.com", "limengqing@bilibili.com"}
	startTime := time.Now().Unix()
	err = s.runStatis(c, date)
	if err != nil {
		msg = err.Error()
	} else {
		msg = fmt.Sprintf("%s 计算完成，耗时%ds", date.Format("2006-01-02"), time.Now().Unix()-startTime)
	}
	err = s.email.SendMail(date, msg, "创作激励每日统计%d年%d月%d日", mailReceivers...)
	if err != nil {
		log.Error("s.email.SendMail error(%v)", err)
	}
	return
}

func (s *Service) runStatis(c context.Context, date time.Time) (err error) {
	defer func() {
		task.GetTaskService().SetTaskStatus(c, task.TaskCreativeStatis, date.Format(_layout), err)
	}()
	// task status type
	err = task.GetTaskService().TaskReady(c, date.Format("2006-01-02"), task.TaskCreativeIncome)
	if err != nil {
		return
	}

	startWeeklyDate = getStartWeeklyDate(date)
	startMonthlyDate = getStartMonthlyDate(date)

	var (
		readGroup   errgroup.Group
		avSections  = SectionEntries{}
		cmSections  = SectionEntries{}
		bgmSections = SectionEntries{}

		upSection    = make([]*model.DateStatis, 0)
		upAvSection  = make([]*model.DateStatis, 0)
		upCmSection  = make([]*model.DateStatis, 0)
		upBgmSection = make([]*model.DateStatis, 0)

		upIncomeWeekly  = make(map[int64]*model.UpIncome)
		upIncomeMonthly = make(map[int64]*model.UpIncome)

		upAvStatisCh  = make(chan map[int64]*model.UpArchStatis, 1)
		upCmStatisCh  = make(chan map[int64]*model.UpArchStatis, 1)
		upBgmStatisCh = make(chan map[int64]*model.UpArchStatis, 1)

		avSourceCh  = make(chan []*model.ArchiveIncome, 1000)
		avDailyCh   = make(chan []*model.ArchiveIncome, 1000)
		avWeeklyCh  = make(chan []*model.ArchiveIncome, 1000)
		avMonthlyCh = make(chan []*model.ArchiveIncome, 1000)
		upAvCh      = make(chan []*model.ArchiveIncome, 1000)

		cmSourceCh  = make(chan []*model.ArchiveIncome, 1000)
		cmDailyCh   = make(chan []*model.ArchiveIncome, 1000)
		cmWeeklyCh  = make(chan []*model.ArchiveIncome, 1000)
		cmMonthlyCh = make(chan []*model.ArchiveIncome, 1000)
		upCmCh      = make(chan []*model.ArchiveIncome, 1000)

		bgmSourceCh  = make(chan []*model.ArchiveIncome, 1000)
		bgmDailyCh   = make(chan []*model.ArchiveIncome, 1000)
		bgmWeeklyCh  = make(chan []*model.ArchiveIncome, 1000)
		bgmMonthlyCh = make(chan []*model.ArchiveIncome, 1000)
		upBgmCh      = make(chan []*model.ArchiveIncome, 1000)

		upSourceCh = make(chan []*model.UpIncome, 1000)
		upDailyCh  = make(chan []*model.UpIncome, 1000)
		upStatisCh = make(chan []*model.UpIncome, 1000)
	)

	// get income by date to statistics
	preDate := startMonthlyDate
	if startWeeklyDate.Before(startMonthlyDate) {
		preDate = startWeeklyDate
	}
	// get av_income
	readGroup.Go(func() (err error) {
		err = s.income.dateStatisSvr.getArchiveByDate(c, avSourceCh, preDate, date, _video, _limitSize)
		if err != nil {
			log.Error("s.getArchiveByDate error(%v)", err)
		}
		return
	})

	// get column_income
	readGroup.Go(func() (err error) {
		err = s.income.dateStatisSvr.getArchiveByDate(c, cmSourceCh, preDate, date, _column, _limitSize)
		if err != nil {
			log.Error("s.getArchiveByDate error(%v)", err)
		}
		return
	})

	// get bgm_income
	readGroup.Go(func() (err error) {
		err = s.income.dateStatisSvr.getArchiveByDate(c, bgmSourceCh, preDate, date, _bgm, _limitSize)
		if err != nil {
			log.Error("s.getArchiveByDate error(%v)", err)
		}
		return
	})

	// get up_income
	readGroup.Go(func() (err error) {
		err = s.income.upIncomeSvr.getUpIncomeByDate(c, upSourceCh, date.Format(_layout), _limitSize)
		if err != nil {
			log.Error("s.getUpIncomeByDate error(%v)", err)
		}
		return
	})

	readGroup.Go(func() (err error) {
		defer func() {
			close(avDailyCh)
			close(avWeeklyCh)
			close(avMonthlyCh)
			close(upAvCh)
		}()
		for av := range avSourceCh {
			avDailyCh <- av
			avWeeklyCh <- av
			avMonthlyCh <- av
			upAvCh <- av
		}
		return
	})

	readGroup.Go(func() (err error) {
		defer func() {
			close(cmDailyCh)
			close(cmWeeklyCh)
			close(cmMonthlyCh)
			close(upCmCh)
		}()
		for cm := range cmSourceCh {
			cmDailyCh <- cm
			cmWeeklyCh <- cm
			cmMonthlyCh <- cm
			upCmCh <- cm
		}
		return
	})

	readGroup.Go(func() (err error) {
		defer func() {
			close(bgmDailyCh)
			close(bgmWeeklyCh)
			close(bgmMonthlyCh)
			close(upBgmCh)
		}()
		for bgm := range bgmSourceCh {
			bgmDailyCh <- bgm
			bgmWeeklyCh <- bgm
			bgmMonthlyCh <- bgm
			upBgmCh <- bgm
		}
		return
	})

	readGroup.Go(func() (err error) {
		defer func() {
			close(upStatisCh)
			close(upDailyCh)
		}()
		for up := range upSourceCh {
			upStatisCh <- up
			upDailyCh <- up
		}
		return
	})

	// up
	readGroup.Go(func() (err error) {
		upSection, upAvSection, upCmSection, upBgmSection, err = s.income.dateStatisSvr.handleDateUp(c, upStatisCh, date)
		if err != nil {
			log.Error("s.income.dateStatisSvr.HandleUp error(%v)", err)
		}
		return
	})

	// video
	readGroup.Go(func() (err error) {
		avSections.avDaily, err = s.income.dateStatisSvr.handleDateStatis(c, avDailyCh, date, _avIncomeDailyStatis)
		if err != nil {
			log.Error("s.income.dateStatisSvr.handleDateStatis error(%v)", err)
		}
		return
	})

	readGroup.Go(func() (err error) {
		avSections.avWeekly, err = s.income.dateStatisSvr.handleDateStatis(c, avWeeklyCh, startWeeklyDate, _avIncomeWeeklyStatis)
		if err != nil {
			log.Error("s.income.dateStatisSvr.handleDateStatis error(%v)", err)
		}
		return
	})

	readGroup.Go(func() (err error) {
		avSections.avMonthly, err = s.income.dateStatisSvr.handleDateStatis(c, avMonthlyCh, startMonthlyDate, _avIncomeMonthlyStatis)
		if err != nil {
			log.Error("s.income.dateStatisSvr.handleDateStatis error(%v)", err)
		}
		return
	})

	// column
	readGroup.Go(func() (err error) {
		cmSections.avDaily, err = s.income.dateStatisSvr.handleDateStatis(c, cmDailyCh, date, _cmIncomeDailyStatis)
		if err != nil {
			log.Error("s.income.dateStatisSvr.handleDateStatis error(%v)", err)
		}
		return
	})

	readGroup.Go(func() (err error) {
		cmSections.avWeekly, err = s.income.dateStatisSvr.handleDateStatis(c, cmWeeklyCh, startWeeklyDate, _cmIncomeWeeklyStatis)
		if err != nil {
			log.Error("s.income.dateStatisSvr.handleDateStatis error(%v)", err)
		}
		return
	})

	readGroup.Go(func() (err error) {
		cmSections.avMonthly, err = s.income.dateStatisSvr.handleDateStatis(c, cmMonthlyCh, startMonthlyDate, _cmIncomeMonthlyStatis)
		if err != nil {
			log.Error("s.income.dateStatisSvr.handleDateStatis error(%v)", err)
		}
		return
	})

	// bgm
	readGroup.Go(func() (err error) {
		bgmSections.avDaily, err = s.income.dateStatisSvr.handleDateStatis(c, bgmDailyCh, date, _bgmIncomeDailyStatis)
		if err != nil {
			log.Error("s.income.dateStatisSvr.handleDateStatis error(%v)", err)
		}
		return
	})

	readGroup.Go(func() (err error) {
		bgmSections.avWeekly, err = s.income.dateStatisSvr.handleDateStatis(c, bgmWeeklyCh, startWeeklyDate, _bgmIncomeWeeklyStatis)
		if err != nil {
			log.Error("s.income.dateStatisSvr.handleDateStatis error(%v)", err)
		}
		return
	})

	readGroup.Go(func() (err error) {
		bgmSections.avMonthly, err = s.income.dateStatisSvr.handleDateStatis(c, bgmMonthlyCh, startMonthlyDate, _bgmIncomeMonthlyStatis)
		if err != nil {
			log.Error("s.income.dateStatisSvr.handleDateStatis error(%v)", err)
		}
		return
	})

	// up_av_statis
	readGroup.Go(func() (err error) {
		err = s.income.upIncomeSvr.handleUpArchStatis(c, upAvStatisCh, upAvCh)
		if err != nil {
			log.Error("p.upIncomeSvr.handleUpArchStatis error(%v)", err)
		}
		return
	})

	// up_column_statis
	readGroup.Go(func() (err error) {
		err = s.income.upIncomeSvr.handleUpArchStatis(c, upCmStatisCh, upCmCh)
		if err != nil {
			log.Error("p.upIncomeSvr.handleUpArchStatis error(%v)", err)
		}
		return
	})

	// up_bgm_statis
	readGroup.Go(func() (err error) {
		err = s.income.upIncomeSvr.handleUpArchStatis(c, upBgmStatisCh, upBgmCh)
		if err != nil {
			log.Error("p.upIncomeSvr.handleUpArchStatis error(%v)", err)
		}
		return
	})

	// up_income_weekly up_income_monthly
	readGroup.Go(func() (err error) {
		upIncomeWeekly, upIncomeMonthly, err = s.income.upIncomeSvr.handleUpIncomeWeeklyAndMonthly(c, date, upAvStatisCh, upCmStatisCh, upBgmStatisCh, upDailyCh)
		if err != nil {
			log.Error("p.upIncomeSvr.handleUpIncomeWeeklyAndMonthly error(%v)", err)
		}
		return
	})

	if err = readGroup.Wait(); err != nil {
		log.Error("run readGroup.Wait error(%v)", err)
		return
	}

	{
		if len(avSections.avDaily) == 0 {
			err = fmt.Errorf("Error: insert 0 av_income_daily_statis")
			return
		}
		if len(avSections.avWeekly) == 0 {
			err = fmt.Errorf("Error: insert 0 av_income_weekly_statis")
			return
		}
		if len(avSections.avMonthly) == 0 {
			err = fmt.Errorf("Error: insert 0 av_income_monthly_statis")
			return
		}
		if len(cmSections.avDaily) == 0 {
			err = fmt.Errorf("Error: insert 0 cm_income_daily_statis")
			return
		}
		if len(cmSections.avWeekly) == 0 {
			err = fmt.Errorf("Error: insert 0 cm_income_weekly_statis")
			return
		}
		if len(cmSections.avMonthly) == 0 {
			err = fmt.Errorf("Error: insert 0 cm_income_monthly_statis")
			return
		}
		if len(bgmSections.avDaily) == 0 {
			err = fmt.Errorf("Error: insert 0 bgm_income_daily_statis")
			return
		}
		if len(bgmSections.avWeekly) == 0 {
			err = fmt.Errorf("Error: insert 0 bgm_income_weekly_statis")
			return
		}
		if len(bgmSections.avMonthly) == 0 {
			err = fmt.Errorf("Error: insert 0 bgm_income_monthly_statis")
			return
		}
		if len(upSection) == 0 {
			err = fmt.Errorf("Error: insert 0 up_income_daily_statis")
			return
		}
		if len(upAvSection) == 0 {
			err = fmt.Errorf("Error: insert 0 up_av_daily_statis")
			return
		}
		if len(upCmSection) == 0 {
			err = fmt.Errorf("Error: insert 0 up_column_daily_statis")
			return
		}
		if len(upBgmSection) == 0 {
			err = fmt.Errorf("Error: insert 0 up_bgm_daily_statis")
			return
		}
		if len(upIncomeWeekly) == 0 {
			err = fmt.Errorf("Error: insert 0 up_income_weekly")
			return
		}
		if len(upIncomeMonthly) == 0 {
			err = fmt.Errorf("Error: insert 0 up_income_monthly")
			return
		}
	}

	// persistent
	var writeGroup errgroup.Group
	// av_income_daily_statis
	writeGroup.Go(func() (err error) {
		_, err = s.income.dateStatisSvr.incomeDateStatisInsert(c, avSections.avDaily, _avIncomeDailyStatis)
		if err != nil {
			log.Error("s.income.dateStatisSvr.incomeDateStatisInsert(daily) error(%v)", err)
			return
		}
		log.Info("insert av_income_daily_statis : %d", len(avSections.avDaily))
		return
	})

	// av_income_weekly_statis
	writeGroup.Go(func() (err error) {
		_, err = s.income.dateStatisSvr.incomeDateStatisInsert(c, avSections.avWeekly, _avIncomeWeeklyStatis)
		if err != nil {
			log.Error("s.income.dateStatisSvr.incomeDateStatisInsert(weekly) error(%v)", err)
			return
		}
		log.Info("insert av_income_weekly_statis : %d", len(avSections.avWeekly))
		return
	})

	// av_income_monthly_statis
	writeGroup.Go(func() (err error) {
		_, err = s.income.dateStatisSvr.incomeDateStatisInsert(c, avSections.avMonthly, _avIncomeMonthlyStatis)
		if err != nil {
			log.Error("s.income.dateStatisSvr.incomeDateStatisInsert(monthly) error(%v)", err)
			return
		}
		log.Info("insert av_income_monthly_statis : %d", len(avSections.avMonthly))
		return
	})

	// column_income_daily_statis
	writeGroup.Go(func() (err error) {
		_, err = s.income.dateStatisSvr.incomeDateStatisInsert(c, cmSections.avDaily, _cmIncomeDailyStatis)
		if err != nil {
			log.Error("s.income.dateStatisSvr.incomeDateStatisInsert(daily) error(%v)", err)
			return
		}
		log.Info("insert column_income_daily_statis : %d", len(cmSections.avDaily))
		return
	})

	// column_income_weekly_statis
	writeGroup.Go(func() (err error) {
		_, err = s.income.dateStatisSvr.incomeDateStatisInsert(c, cmSections.avWeekly, _cmIncomeWeeklyStatis)
		if err != nil {
			log.Error("s.income.dateStatisSvr.incomeDateStatisInsert(weekly) error(%v)", err)
			return
		}
		log.Info("insert column_income_weekly_statis : %d", len(cmSections.avWeekly))
		return
	})

	// column_income_monthly_statis
	writeGroup.Go(func() (err error) {
		_, err = s.income.dateStatisSvr.incomeDateStatisInsert(c, cmSections.avMonthly, _cmIncomeMonthlyStatis)
		if err != nil {
			log.Error("s.income.dateStatisSvr.incomeDateStatisInsert(monthly) error(%v)", err)
			return
		}
		log.Info("insert column_income_monthly_statis : %d", len(cmSections.avMonthly))
		return
	})

	// bgm_income_daily_statis
	writeGroup.Go(func() (err error) {
		_, err = s.income.dateStatisSvr.incomeDateStatisInsert(c, bgmSections.avDaily, _bgmIncomeDailyStatis)
		if err != nil {
			log.Error("s.income.dateStatisSvr.incomeDateStatisInsert(daily) error(%v)", err)
			return
		}
		log.Info("insert bgm_income_daily_statis : %d", len(bgmSections.avDaily))
		return
	})

	// bgm_income_weekly_statis
	writeGroup.Go(func() (err error) {
		_, err = s.income.dateStatisSvr.incomeDateStatisInsert(c, bgmSections.avWeekly, _bgmIncomeWeeklyStatis)
		if err != nil {
			log.Error("s.income.dateStatisSvr.incomeDateStatisInsert(weekly) error(%v)", err)
			return
		}
		log.Info("insert bgm_income_weekly_statis : %d", len(bgmSections.avWeekly))
		return
	})

	// bgm_income_monthly_statis
	writeGroup.Go(func() (err error) {
		_, err = s.income.dateStatisSvr.incomeDateStatisInsert(c, bgmSections.avMonthly, _bgmIncomeMonthlyStatis)
		if err != nil {
			log.Error("s.income.dateStatisSvr.incomeDateStatisInsert(monthly) error(%v)", err)
			return
		}
		log.Info("insert bgm_income_monthly_statis : %d", len(bgmSections.avMonthly))
		return
	})

	// up_income_daily_statis
	writeGroup.Go(func() (err error) {
		_, err = s.income.dateStatisSvr.upIncomeDailyStatisInsert(c, upSection, "up_income_daily_statis")
		if err != nil {
			log.Error("s.upIncomeDailyStatisInsert error(%v)", err)
			return
		}
		log.Info("insert up_income_daily_statis : %d", len(upSection))
		return
	})

	// up_av_daily_statis
	writeGroup.Go(func() (err error) {
		_, err = s.income.dateStatisSvr.upIncomeDailyStatisInsert(c, upAvSection, "up_av_daily_statis")
		if err != nil {
			log.Error("s.upIncomeDailyStatisInsert error(%v)", err)
			return
		}
		log.Info("insert up_av_daily_statis : %d", len(upAvSection))
		return
	})

	// up_column_daily_statis
	writeGroup.Go(func() (err error) {
		_, err = s.income.dateStatisSvr.upIncomeDailyStatisInsert(c, upCmSection, "up_column_daily_statis")
		if err != nil {
			log.Error("s.upIncomeDailyStatisInsert error(%v)", err)
			return
		}
		log.Info("insert up_column_daily_statis : %d", len(upCmSection))
		return
	})

	// up_bgm_daily_statis
	writeGroup.Go(func() (err error) {
		_, err = s.income.dateStatisSvr.upIncomeDailyStatisInsert(c, upBgmSection, "up_bgm_daily_statis")
		if err != nil {
			log.Error("s.upIncomeDailyStatisInsert error(%v)", err)
			return
		}
		log.Info("insert up_bgm_daily_statis : %d", len(upBgmSection))
		return
	})

	// up_income_weekly
	writeGroup.Go(func() (err error) {
		err = s.income.upIncomeSvr.UpIncomeDBStoreBatch(c, _upIncomeWeekly, upIncomeWeekly)
		if err != nil {
			log.Error("s.UpIncomeDBStoreBatch up_income_weekly error(%v)", err)
			return
		}
		log.Info("insert up_income_weekly : %d", len(upIncomeWeekly))
		return
	})

	// up_income_monthly
	writeGroup.Go(func() (err error) {
		err = s.income.upIncomeSvr.UpIncomeDBStoreBatch(c, _upIncomeMonthly, upIncomeMonthly)
		if err != nil {
			log.Error("s.UpIncomeDBStoreBatch up_income_monthly error(%v)", err)
			return
		}
		log.Info("insert up_income_monthly : %d", len(upIncomeMonthly))
		return
	})

	if err = writeGroup.Wait(); err != nil {
		log.Error("run writeGroup.Wait error(%v)", err)
	}
	return
}
