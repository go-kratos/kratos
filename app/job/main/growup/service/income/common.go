package income

import (
	"time"
)

const (
	_layout    = "2006-01-02"
	_layoutSec = "2006-01-02 15:04:05"
)

var (
	startWeeklyDate  time.Time
	startMonthlyDate time.Time
)

var (
	batchSize = 2000
	_dbInsert = 1
	_dbUpdate = 2

	_video  = 0
	_column = 2
	_bgm    = 3

	_limitSize = 2000
)

func getStartWeeklyDate(date time.Time) time.Time {
	for date.Weekday() != time.Monday {
		date = date.AddDate(0, 0, -1)
	}
	return date
}

func getStartMonthlyDate(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.Local)
}
