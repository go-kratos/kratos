package tag

import (
	"go-common/library/time"
)

// TagInfo tag_info
type TagInfo struct {
	ID              int64
	TagName         string
	CategoryID      int64
	BusinessID      int
	AdjustType      int
	Ratio           int
	IsCommon        int
	ActivityID      int64
	UploadStartTime time.Time
	UploadEndTime   time.Time
	StartAt         time.Time
	EndAt           time.Time
}

// AvTagRatio av tag ratio
type AvTagRatio struct {
	AvID        int64
	MID         int64
	TagID       int64
	AdjustType  int
	Ratio       int
	Income      int64
	BaseIncome  int64
	TotalIncome int64
	TaxMoney    int64
	Date        string
}

// up_tag_income up tag income info
type UpTagIncome struct {
	ID    int64
	MID   int64
	TagID int64
}
