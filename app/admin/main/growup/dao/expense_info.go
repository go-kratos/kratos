package dao

import (
	"context"
	"fmt"
	"time"

	"go-common/app/admin/main/growup/model"
	"go-common/library/log"
)

const (
	// select
	_getAllDayExpenseSQL     = "SELECT day_expense, up_count, av_count, up_avg_expense, av_avg_expense, total_expense, date FROM expense_daily_info WHERE date >= ? AND ctype = ? ORDER BY date DESC limit ?,?"
	_getAllMonthExpenseSQL   = "SELECT month_expense, up_count, av_count, up_avg_expense, av_avg_expense, total_expense, date, month FROM expense_monthly_info WHERE month <= ? AND month >=? AND ctype = ? ORDER BY month DESC LIMIT ?,?"
	_getDayTotalExpenseSQL   = "SELECT total_expense FROM expense_daily_info WHERE date = ? AND ctype = ?"
	_getLatelyExpenseDateSQL = "SELECT date FROM expense_%s_info WHERE ctype = ? ORDER BY date DESC LIMIT 1"

	// count(*)
	_expenseDayCountSQL   = "SELECT count(*) FROM expense_daily_info WHERE date >= ? AND ctype = ?"
	_expenseMonthCountSQL = "SELECT count(*) FROM expense_monthly_info WHERE month <= ? AND month >= ? AND ctype = ?"
)

// GetDayExpenseCount get expense_daily_info count
func (d *Dao) GetDayExpenseCount(c context.Context, beginDate time.Time, ctype int) (total int, err error) {
	err = d.rddb.QueryRow(c, _expenseDayCountSQL, beginDate, ctype).Scan(&total)
	return
}

// GetAllDayExpenseInfo get year all day expense.
func (d *Dao) GetAllDayExpenseInfo(c context.Context, beginDate time.Time, ctype, from, limit int) (infos []*model.BudgetDayStatistics, err error) {
	rows, err := d.rddb.Query(c, _getAllDayExpenseSQL, beginDate, ctype, from, limit)
	if err != nil {
		log.Error("dao.GetAllDayExpenseInfo query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := &model.BudgetDayStatistics{}
		if err = rows.Scan(&a.DayExpense, &a.UpCount, &a.AvCount, &a.UpAvgExpense, &a.AvAvgExpense, &a.TotalExpense, &a.Date); err != nil {
			log.Error("dao.GetAllDayExpenseInfo scan error(%v)", err)
			return
		}
		infos = append(infos, a)
	}
	err = rows.Err()
	return
}

// GetDayTotalExpenseInfo get one day total_expense.
func (d *Dao) GetDayTotalExpenseInfo(c context.Context, date time.Time, ctype int) (totalExpense int64, err error) {
	err = d.rddb.QueryRow(c, _getDayTotalExpenseSQL, date, ctype).Scan(&totalExpense)
	return
}

// GetMonthExpenseCount get expense month count
func (d *Dao) GetMonthExpenseCount(c context.Context, month, beginMonth string, ctype int) (total int, err error) {
	err = d.rddb.QueryRow(c, _expenseMonthCountSQL, month, beginMonth, ctype).Scan(&total)
	return
}

// GetAllMonthExpenseInfo get all month expense.
func (d *Dao) GetAllMonthExpenseInfo(c context.Context, month, beginMonth string, ctype, from, limit int) (infos []*model.BudgetMonthStatistics, err error) {
	rows, err := d.rddb.Query(c, _getAllMonthExpenseSQL, month, beginMonth, ctype, from, limit)
	if err != nil {
		log.Error("dao.GetAllMonthExpenseInfo query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := &model.BudgetMonthStatistics{}
		if err = rows.Scan(&a.MonthExpense, &a.UpCount, &a.AvCount, &a.UpAvgExpense, &a.AvAvgExpense, &a.TotalExpense, &a.Date, &a.Month); err != nil {
			log.Error("dao.GetAllMonthExpenseInfo scan error(%v)", err)
			return
		}
		infos = append(infos, a)
	}
	err = rows.Err()
	return
}

// GetLatelyExpenseDate get lately date.
func (d *Dao) GetLatelyExpenseDate(c context.Context, table string, ctype int) (date time.Time, err error) {
	err = d.rddb.QueryRow(c, fmt.Sprintf(_getLatelyExpenseDateSQL, table), ctype).Scan(&date)
	return
}
