package model

import (
	"go-common/library/time"
)

// BudgetDayStatistics day.
type BudgetDayStatistics struct {
	DayExpense   int       `json:"day_expense"`
	UpCount      int       `json:"up_count"`
	AvCount      int       `json:"av_count"`
	UpAvgExpense int       `json:"up_avg_expense"`
	AvAvgExpense int       `json:"av_avg_expense"`
	Date         time.Time `json:"date"`
	TotalExpense int64     `json:"total_expense"`
	ExpenseRatio string    `json:"expense_ratio"`
	DayRatio     string    `json:"day_ratio"`
}

// BudgetRatio budget ratio.
type BudgetRatio struct {
	ExpenseRatio string `json:"expense_ratio"`
	DayRatio     string `json:"day_ratio"`
	Year         int64  `json:"year"`
	Budget       int64  `json:"budget"`
}

// BudgetMonthStatistics month
type BudgetMonthStatistics struct {
	MonthExpense int64     `json:"month_expense"`
	Month        string    `json:"month"`
	Date         time.Time `json:"date"`
	UpCount      int       `json:"up_count"`
	AvCount      int       `json:"av_count"`
	UpAvgExpense int       `json:"up_avg_expense"`
	AvAvgExpense int       `json:"av_avg_expense"`
	TotalExpense int64     `json:"total_expense"`
	ExpenseRatio string    `json:"expense_ratio"`
}
