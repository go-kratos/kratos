package data

import (
	"go-common/app/admin/main/up/conf"
	"go-common/app/admin/main/up/dao/data"
	"go-common/app/admin/main/up/dao/tag"
	"time"
)

//Service data service
type Service struct {
	c    *conf.Config
	data *data.Dao
	dtag *tag.Dao
}

//New get service
func New(c *conf.Config) *Service {
	s := &Service{
		c:    c,
		data: data.New(c),
		dtag: tag.New(c),
	}
	return s
}

func beginningOfDay(t time.Time) time.Time {
	d := time.Duration(-t.Hour()) * time.Hour
	return t.Truncate(time.Hour).Add(d)
}

func getTuesday(now time.Time) time.Time {
	t := beginningOfDay(now)
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	d := time.Duration(-weekday+2) * 24 * time.Hour
	return t.Truncate(time.Hour).Add(d)
}

func getSunday(now time.Time) time.Time {
	t := beginningOfDay(now)
	weekday := int(t.Weekday())
	if weekday == 0 {
		return t
	}
	d := time.Duration(7-weekday) * 24 * time.Hour
	return t.Truncate(time.Hour).Add(d)
}

func getDateLastSunday() (date time.Time) {
	t := time.Now()
	td := getTuesday(t).Add(12 * time.Hour)
	if t.Before(td) { //当前时间在本周二12点之前，则取上上周日的数据，否则取上周日的数据
		date = getSunday(t.AddDate(0, 0, -14))
	} else {
		date = getSunday(t.AddDate(0, 0, -7))
	}
	//log.Info("current time (%s) tuesday (%s) sunday (%s)", t.Format("2006-01-02 15:04:05"), td, sd)
	return
}
