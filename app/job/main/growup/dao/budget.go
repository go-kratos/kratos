package dao

import (
	"context"
	"fmt"

	"go-common/app/job/main/growup/model"
	"go-common/library/log"
)

const (
	// get expense_daily_info total_expense
	_expenseDailyTotalSQL = "SELECT total_expense,ctype FROM expense_daily_info WHERE date = '%s'"

	// insert
	_insertDailyExpenseSQL   = "INSERT INTO expense_daily_info(day_expense, up_count, av_count, up_avg_expense, av_avg_expense, total_expense, date, ctype) VALUES (%d,%d,%d,%d,%d,%d,'%s',%d) ON DUPLICATE KEY UPDATE day_expense = VALUES(day_expense), up_count = VALUES(up_count), av_count = VALUES(av_count), up_avg_expense = values(up_avg_expense), av_avg_expense = values(av_avg_expense), total_expense = values(total_expense), date = values(date)"
	_insertMonthlyExpenseSQL = "INSERT INTO expense_monthly_info(month_expense, up_count, av_count, up_avg_expense, av_avg_expense, total_expense, date, month, ctype) VALUES (%d,%d,%d,%d,%d,%d,'%s','%s',%d) ON DUPLICATE KEY UPDATE month_expense = VALUES(month_expense), up_count = VALUES(up_count), av_count = VALUES(av_count), up_avg_expense = values(up_avg_expense), av_avg_expense = values(av_avg_expense), total_expense = values(total_expense), date = values(date)"
)

// GetTotalExpenseByDate get expense_daily_info by date
func (d *Dao) GetTotalExpenseByDate(c context.Context, date string) (expense map[int64]int64, err error) {
	expense = make(map[int64]int64)
	rows, err := d.db.Query(c, fmt.Sprintf(_expenseDailyTotalSQL, date))
	if err != nil {
		log.Error("GetTotalExpenseByDate d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var ctype, totalExpense int64
		err = rows.Scan(&totalExpense, &ctype)
		if err != nil {
			log.Error("GetTotalExpenseByDate rows.Scan error(%v)", err)
			return
		}
		expense[ctype] = totalExpense
	}
	err = rows.Err()
	return
}

// InsertDailyExpense insert expense_daily_info
func (d *Dao) InsertDailyExpense(c context.Context, e *model.BudgetExpense) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_insertDailyExpenseSQL, e.Expense, e.UpCount, e.AvCount, e.UpAvgExpense, e.AvAvgExpense, e.TotalExpense, e.Date, e.CType))
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// InsertMonthlyExpense insert expense_monthly_info
func (d *Dao) InsertMonthlyExpense(c context.Context, e *model.BudgetExpense) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_insertMonthlyExpenseSQL, e.Expense, e.UpCount, e.AvCount, e.UpAvgExpense, e.AvAvgExpense, e.TotalExpense, e.Date, e.Date.Format("2006-01"), e.CType))
	if err != nil {
		return
	}
	return res.RowsAffected()
}
