package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/app/admin/main/growup/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

// BudgetDayStatistics budget day statistics.
func (s *Service) BudgetDayStatistics(c context.Context, ctype, from, limit int) (total int, infos []*model.BudgetDayStatistics, err error) {
	infos = make([]*model.BudgetDayStatistics, 0)
	latelyDate, err := s.dao.GetLatelyExpenseDate(c, "daily", ctype)
	if err == sql.ErrNoRows {
		err = nil
		return
	}
	if err != nil {
		log.Error("s.BudgetDayGraph dao.GetLatelyExpenseDate error(%v)", err)
		return
	}
	beginDate := time.Date(latelyDate.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	total, err = s.dao.GetDayExpenseCount(c, beginDate, ctype)
	if err != nil {
		log.Error("s.dao.GetDayExpenseCount error(%v)", err)
		return
	}
	if total == 0 {
		return
	}

	infos, err = s.dao.GetAllDayExpenseInfo(c, beginDate, ctype, from, limit)
	if err != nil {
		log.Error("s.dao.GetAllDayExpenseInfo error(%v)", err)
		return
	}
	var annualBudget, dayBudget int64
	switch ctype {
	case _video:
		annualBudget, dayBudget = s.conf.Budget.Video.AnnualBudget, s.conf.Budget.Video.DayBudget
	case _column:
		annualBudget, dayBudget = s.conf.Budget.Column.AnnualBudget, s.conf.Budget.Column.DayBudget
	case _bgm:
		annualBudget, dayBudget = s.conf.Budget.Bgm.AnnualBudget, s.conf.Budget.Bgm.DayBudget
	}
	for _, info := range infos {
		info.ExpenseRatio = strconv.FormatFloat(float64(info.TotalExpense)/float64(annualBudget), 'f', 2, 32)
		info.DayRatio = strconv.FormatFloat(float64(info.DayExpense)/float64(dayBudget), 'f', 2, 32)
	}
	return
}

// BudgetDayGraph get day graph.
func (s *Service) BudgetDayGraph(c context.Context, ctype int) (ratioInfo *model.BudgetRatio, err error) {
	ratioInfo = new(model.BudgetRatio)
	switch ctype {
	case _video:
		ratioInfo.Year, ratioInfo.Budget = s.conf.Budget.Video.Year, s.conf.Budget.Video.AnnualBudget
	case _column:
		ratioInfo.Year, ratioInfo.Budget = s.conf.Budget.Column.Year, s.conf.Budget.Column.AnnualBudget
	case _bgm:
		ratioInfo.Year, ratioInfo.Budget = s.conf.Budget.Bgm.Year, s.conf.Budget.Bgm.AnnualBudget
	}
	latelyDate, err := s.dao.GetLatelyExpenseDate(c, "daily", ctype)
	if err == sql.ErrNoRows {
		err = nil
		return
	}
	if err != nil {
		log.Error("s.BudgetDayGraph dao.GetLatelyExpenseDate error(%v)", err)
		return
	}
	t := time.Date(latelyDate.Year(), latelyDate.Month(), latelyDate.Day(), 0, 0, 0, 0, time.Local)
	totalExpense, err := s.dao.GetDayTotalExpenseInfo(c, t, ctype)
	if err != nil {
		log.Error("s.BudgetDayGraph dao.GetDayTotalExpenseInfo error(%v)", err)
		return
	}
	ratioInfo.ExpenseRatio = strconv.FormatFloat(float64(totalExpense)/float64(ratioInfo.Budget), 'f', 2, 32)
	ratioInfo.DayRatio = strconv.FormatFloat(float64(getGoneDays(latelyDate)*100)/float64(getYearDays(latelyDate.Year())), 'f', 2, 32)
	return
}

// BudgetMonthStatistics budget month statistics
func (s *Service) BudgetMonthStatistics(c context.Context, ctype, from, limit int) (total int, infos []*model.BudgetMonthStatistics, err error) {
	latelyDate, err := s.dao.GetLatelyExpenseDate(c, "monthly", ctype)
	if err == sql.ErrNoRows {
		err = nil
		return
	}
	if err != nil {
		log.Error("s.BudgetMonthStatistics dao.GetLatelyExpenseDate error(%v)", err)
		return
	}
	var month, beginMonth string
	if int(latelyDate.Month()) < 10 {
		month = fmt.Sprintf("%d-0%d", latelyDate.Year(), latelyDate.Month())
	} else {
		month = fmt.Sprintf("%d-%d", latelyDate.Year(), latelyDate.Month())
	}
	beginMonth = fmt.Sprintf("%d-01", latelyDate.Year())
	total, err = s.dao.GetMonthExpenseCount(c, month, beginMonth, ctype)
	if err != nil {
		log.Error("s.dao.GetMonthExpenseCount error(%v)", err)
		return
	}
	infos, err = s.dao.GetAllMonthExpenseInfo(c, month, beginMonth, ctype, from, limit)
	if err != nil {
		log.Error("s.BudgetMonthStatistics dao.GetAllMonthExpenseInfo error(%v)", err)
		return
	}
	if len(infos) <= 0 {
		infos = make([]*model.BudgetMonthStatistics, 0)
	}
	var dayBudget int64
	switch ctype {
	case _video:
		dayBudget = s.conf.Budget.Video.DayBudget
	case _column:
		dayBudget = s.conf.Budget.Column.DayBudget
	case _bgm:
		dayBudget = s.conf.Budget.Bgm.DayBudget
	}
	for _, info := range infos {
		info.ExpenseRatio = strconv.FormatFloat(float64(info.MonthExpense)/float64(int(dayBudget)*getMonthDays(time.Unix(int64(info.Date), 0))), 'f', 2, 32)
		date := time.Unix(int64(info.Date), 0)
		info.Month = strconv.Itoa(date.Year())
		info.Month += "-"
		if int(date.Month()) < 10 {
			info.Month += "0"
		}
		info.Month += strconv.Itoa(int(date.Month()))
	}
	return
}

func getMonthDays(date time.Time) (count int) {
	begin := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.Local)
	end := time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, time.Local)
	return int(end.Sub(begin).Hours() / 24)
}

func getYearDays(year int) (count int) {
	yearBegin := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	yearEnd := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.Local)
	return int(yearEnd.Sub(yearBegin).Hours() / 24)
}

func getGoneDays(date time.Time) (count int) {
	yearBegin := time.Date(date.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	return int(date.Sub(yearBegin).Hours()/24) + 1
}
