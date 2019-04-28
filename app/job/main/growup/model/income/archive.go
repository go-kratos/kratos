package income

import (
	"go-common/library/time"
)

// ArchiveIncome include av income and column income
type ArchiveIncome struct {
	ID          int64
	AID         int64
	MID         int64
	TagID       int64
	IsOriginal  int
	UploadTime  time.Time
	Income      int64
	TaxMoney    int64
	TotalIncome int64
	Date        time.Time
}
