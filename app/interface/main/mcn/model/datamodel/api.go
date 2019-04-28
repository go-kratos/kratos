package datamodel

import (
	"time"
)

// GetLastDay get data daily
func GetLastDay() time.Time {
	//return time.Date(2018, 11, 18, 0, 0, 0, 0, time.Local)
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.Local)
}

// GetLastWeek get data weekly
func GetLastWeek() time.Time {
	//return time.Date(2018, 11, 18, 0, 0, 0, 0, time.Local)
	now := time.Now()
	gDate := getTuesday(now)
	if now.Before(gDate.Add(12 * time.Hour)) {
		return gDate.AddDate(0, 0, -9)
	}
	return gDate.AddDate(0, 0, -2)
}

func beginningOfDay(t time.Time) time.Time {
	d := time.Duration(-t.Hour()) * time.Hour
	return t.Truncate(time.Hour).Add(d)
}

func getTuesday(now time.Time) time.Time {
	t := beginningOfDay(now)
	weekday := int(t.Weekday())
	if weekday == int(time.Sunday) {
		weekday = int(time.Saturday) + 1
	}
	d := time.Duration(-weekday+2) * 24 * time.Hour
	return t.Truncate(time.Hour).Add(d)
}
