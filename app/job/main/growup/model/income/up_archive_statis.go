package income

import (
	"go-common/library/time"
)

// UpArchStatis up archive statistics
type UpArchStatis struct {
	MID         int64
	WeeklyDate  time.Time
	WeeklyAIDs  string
	MonthlyDate time.Time
	MonthlyAIDs string
}

// UpAvStatis up av statistics
type UpAvStatis struct {
	ID           int64
	MID          int64
	WeeklyDate   time.Time
	WeeklyAvIDs  string
	MonthlyDate  time.Time
	MonthlyAvIDs string
	IsDeleted    int
	CTime        time.Time
	MTime        time.Time
	DBState      int
}
