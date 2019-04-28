package income

import "go-common/library/time"

// UpCharge up charge
type UpCharge struct {
	MID         int64
	AvCount     int64
	IncCharge   int64
	TotalCharge int64
	Date        time.Time
	IsDeleted   int
	CTime       time.Time
	MTime       time.Time
}
