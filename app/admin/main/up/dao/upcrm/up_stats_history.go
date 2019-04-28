package upcrm

import (
	"time"

	"go-common/app/admin/main/up/model/upcrmmodel"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

const (
	upstatshistory = "up_stats_history"
	//ISO8601DATE only date format
	ISO8601DATE = "2006-01-02"
)

//GetUpStatLastDate get last update date from db
func (d *Dao) GetUpStatLastDate(date time.Time) (lastday time.Time, err error) {
	var lasthistory upcrmmodel.UpStatsHistory
	err = d.crmdb.Table(upstatshistory).Select("generate_date").Order("generate_date desc").Limit(1).Find(&lasthistory).Error
	if err != nil {
		log.Error("get last date fail for up stat history, err=%+v", err)
		return
	}
	lastday = lasthistory.GenerateDate.Time()
	return
}

//QueryYesterday query yesterday db
func (d *Dao) QueryYesterday(date time.Time) (res []*upcrmmodel.UpStatsHistory, err error) {
	err = d.crmdb.Table(upstatshistory).Where("generate_date = ? AND type in ( ?, ?, ? )", date.Format(ISO8601DATE), upcrmmodel.ActivityType, upcrmmodel.IncrType, upcrmmodel.TotalType).Find(&res).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
		return
	}
	if err != nil {
		return nil, err
	}
	return
}

//QueryTrend query db
func (d *Dao) QueryTrend(statType int, currentDate time.Time, days int) (res []*upcrmmodel.UpStatsHistory, err error) {
	// 这种type有3种子类型，需要加起来
	if statType == upcrmmodel.ActivityType {
		days *= 3
	}
	err = d.crmdb.Table(upstatshistory).Where("type = ? AND generate_date <=?", statType, currentDate.Format(ISO8601DATE)).Order("generate_date desc").Limit(days).Find(&res).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
		return
	}
	if err != nil {
		return nil, err
	}
	return
}

//QueryDetail query db
func (d *Dao) QueryDetail(startDate time.Time, endDate time.Time) (res []*upcrmmodel.UpStatsHistory, err error) {
	err = d.crmdb.Table(upstatshistory).Where("generate_date BETWEEN ? AND ?", startDate.Format(ISO8601DATE), endDate.Format(ISO8601DATE)).Order("generate_date Desc").Find(&res).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
		return
	}
	if err != nil {
		return nil, err
	}
	return
}
