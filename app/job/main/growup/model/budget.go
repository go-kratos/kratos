package model

import (
	"time"
)

// BudgetExpense budget expense
type BudgetExpense struct {
	Expense      int64
	UpCount      int64
	AvCount      int64
	UpAvgExpense int64
	AvAvgExpense int64
	Date         time.Time
	TotalExpense int64
	CType        int64
}
