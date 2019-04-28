package upcrm

import (
	"time"

	"go-common/app/admin/main/up/model/upcrmmodel"
)

//QueryUpRank query up rank
func (d *Dao) QueryUpRank(rankType int, date time.Time) (result []upcrmmodel.UpRank, err error) {
	err = d.crmdb.Model(&upcrmmodel.UpRank{}).Where("type=? and generate_date=?", rankType, date.Format(upcrmmodel.TimeFmtDate)).Find(&result).Error
	return
}

//QueryUpRankAll query up rank all
func (d *Dao) QueryUpRankAll(date time.Time) (result []upcrmmodel.UpRank, err error) {
	err = d.crmdb.Model(&upcrmmodel.UpRank{}).Where("generate_date=?", date.Format(upcrmmodel.TimeFmtDate)).Find(&result).Error
	return
}

//GetUpRankLatestDate get last generate date
func (d *Dao) GetUpRankLatestDate() (date time.Time, err error) {
	var rankInfo = upcrmmodel.UpRank{}
	err = d.crmdb.Model(&rankInfo).Select("generate_date").Order("generate_date desc").Limit(1).Find(&rankInfo).Error
	if err == nil {
		date, err = time.Parse(time.RFC3339, rankInfo.GenerateDate)
	}
	return
}
