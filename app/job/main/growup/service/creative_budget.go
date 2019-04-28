package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/growup/model"
	"go-common/library/log"
)

// CreativeBudget creative budget
func (s *Service) CreativeBudget(c context.Context, date time.Time) (err error) {
	defer func() {
		GetTaskService().SetTaskStatus(c, TaskBudget, date.Format("2006-01-02"), err)
	}()
	err = GetTaskService().TaskReady(c, date.Format("2006-01-02"), TaskCreativeIncome)
	if err != nil {
		return
	}

	ups, err := s.GetUpIncome(c, "up_income", date.Format(_layout))
	if err != nil {
		log.Error("s.GetUpIncome(up_income) error(%v)", err)
		return
	}
	if len(ups) == 0 {
		err = fmt.Errorf("get 0 record from up_income")
		return
	}
	log.Info("daily ups(%d)", len(ups))

	monthStart := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.Local)
	monthUps, err := s.GetUpIncome(c, "up_income_monthly", monthStart.Format(_layout))
	if err != nil {
		log.Error("s.GetUpIncome(up_income_monthly) error(%v)", err)
		return
	}
	if len(monthUps) == 0 {
		err = fmt.Errorf("get 0 record from up_income_monthly")
		return
	}
	log.Info("monthly ups(%d)", len(ups))

	preTotalExpense, err := s.dao.GetTotalExpenseByDate(c, date.AddDate(0, 0, -1).Format(_layout))
	if err != nil {
		log.Error("s.dao.GetTotalExpenseByDate error(%v)", err)
		return
	}

	log.Info("CreativeBudget ready date ok(%s)", date.Format(_layout))
	avDaily, cmDaily, bgmDaily := calCreativebudget(ups, date, preTotalExpense)
	avMonthly, cmMonthly, bgmMonthly := calCreativebudget(monthUps, date, preTotalExpense)
	avMonthly.TotalExpense = avDaily.TotalExpense
	cmMonthly.TotalExpense = cmDaily.TotalExpense
	bgmMonthly.TotalExpense = bgmDaily.TotalExpense
	// insert
	if avDaily.Expense > 0 {
		_, err = s.dao.InsertDailyExpense(c, avDaily)
		if err != nil {
			log.Error("s.dao.InsertDailyExpense error(%v)", err)
			return
		}
		_, err = s.dao.InsertMonthlyExpense(c, avMonthly)
		if err != nil {
			log.Error("s.dao.InsertMonthlyExpense error(%v)", err)
			return
		}
	}
	if cmDaily.Expense > 0 {
		_, err = s.dao.InsertDailyExpense(c, cmDaily)
		if err != nil {
			log.Error("s.dao.InsertDailyExpense error(%v)", err)
			return
		}
		_, err = s.dao.InsertMonthlyExpense(c, cmMonthly)
		if err != nil {
			log.Error("s.dao.InsertMonthlyExpense error(%v)", err)
			return
		}
	}

	if bgmDaily.Expense > 0 {
		_, err = s.dao.InsertDailyExpense(c, bgmDaily)
		if err != nil {
			log.Error("s.dao.InsertDailyExpense error(%v)", err)
			return
		}
		_, err = s.dao.InsertMonthlyExpense(c, bgmMonthly)
		if err != nil {
			log.Error("s.dao.InsertMonthlyExpense error(%v)", err)
		}
	}
	return
}

func calCreativebudget(ups []*model.UpIncome, date time.Time, preTotalExpense map[int64]int64) (avBudget, cmBudget, bgmBudget *model.BudgetExpense) {
	avBudget, cmBudget, bgmBudget = &model.BudgetExpense{}, &model.BudgetExpense{}, &model.BudgetExpense{}
	var (
		avExpense, cmExpense, bgmExpense int64
		avCount, cmCount, bgmCount       int64
		upAvCount, upCmCount, upBgmCount int64
	)
	for _, up := range ups {
		if up.AvIncome > 0 {
			avCount += up.AvCount
			upAvCount++
		}
		// add up fix adjust to av
		avExpense += up.Income - up.ColumnIncome - up.BgmIncome // TODO up.AvIncome

		if up.ColumnIncome > 0 {
			cmExpense += up.ColumnIncome
			cmCount += up.ColumnCount
			upCmCount++
		}
		if up.BgmIncome > 0 {
			bgmExpense += up.BgmIncome
			bgmCount += up.BgmCount
			upBgmCount++
		}
	}
	if avCount > 0 && upAvCount > 0 {
		avBudget = &model.BudgetExpense{
			Expense:      avExpense,
			AvCount:      avCount,
			UpCount:      upAvCount,
			UpAvgExpense: avExpense / upAvCount,
			AvAvgExpense: avExpense / avCount,
			Date:         date,
			TotalExpense: preTotalExpense[0] + avExpense,
			CType:        0,
		}
	}
	if cmCount > 0 && upCmCount > 0 {
		cmBudget = &model.BudgetExpense{
			Expense:      cmExpense,
			AvCount:      cmCount,
			UpCount:      upCmCount,
			UpAvgExpense: cmExpense / upCmCount,
			AvAvgExpense: cmExpense / cmCount,
			Date:         date,
			TotalExpense: preTotalExpense[2] + cmExpense,
			CType:        2,
		}
	}
	if bgmCount > 0 && upBgmCount > 0 {
		bgmBudget = &model.BudgetExpense{
			Expense:      bgmExpense,
			AvCount:      bgmCount,
			UpCount:      upBgmCount,
			UpAvgExpense: bgmExpense / upBgmCount,
			AvAvgExpense: bgmExpense / bgmCount,
			Date:         date,
			TotalExpense: preTotalExpense[3] + bgmExpense,
			CType:        3,
		}
	}
	return
}
