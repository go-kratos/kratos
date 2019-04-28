package income

import (
	"go-common/library/time"
)

// AvChargeStatis av charge statistics
type AvChargeStatis struct {
	ID          int64
	AvID        int64
	MID         int64
	TagID       int64
	IsOriginal  int
	UploadTime  time.Time
	TotalCharge int64
	IsDeleted   int
	CTime       time.Time
	MTime       time.Time
	DBState     int // 1-insert 2-update
}
