package tag

import (
	"go-common/library/time"
)

// AvIncome av income
type AvIncome struct {
	ID          int64
	AvID        int64
	MID         int64
	Income      int64
	TotalIncome int64
	TaxMoney    int64
	Date        time.Time
}

// UpIncome up income
type UpIncome struct {
	ID                int64
	MID               int64
	Income            int64
	BaseIncome        int64
	TotalIncome       int64
	TaxMoney          int64
	AvIncome          int64
	AvBaseIncome      int64
	AvTotalIncome     int64
	AvTax             int64
	ColumnIncome      int64
	ColumnBaseIncome  int64
	ColumnTotalIncome int64
	ColumnTax         int64
	BgmIncome         int64
	BgmBaseIncome     int64
	BgmTotalIncome    int64
	BgmTax            int64
	Date              time.Time
}

// AvCharge av_charge
type AvCharge struct {
	ID         int64
	AvID       int64
	MID        int64
	CategoryID int64
	IncCharge  int64
	ActivityID int64
	UploadTime time.Time
	TagID      int64
	IsDeleted  int
}

// ArchiveIncome av income
type ArchiveIncome struct {
	ID          int64
	AID         int64
	MID         int64
	Income      int64
	BaseIncome  int64
	TotalIncome int64
	TaxMoney    int64
	Date        time.Time
}

// ArchiveCharge av column bgm
type ArchiveCharge struct {
	ID         int64
	AID        int64
	MID        int64
	CategoryID int64
	IncCharge  int64
	ActivityID int64
	UploadTime time.Time
	TagID      int64
	IsDeleted  int
}
