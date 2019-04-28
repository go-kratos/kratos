package upcrm

import (
	"time"

	"go-common/app/admin/main/up/model/upcrmmodel"
)

// QueryUpScoreHistory score type, 1质量分，2影响分，3信用分,
//				0表示所有
func (d *Dao) QueryUpScoreHistory(mid int64, scoreType []int, fromdate time.Time, todate time.Time) (result []upcrmmodel.UpScoreHistory, err error) {
	err = d.crmdb.Table(upcrmmodel.GetUpScoreHistoryTableName(mid)).Where("mid=? and score_type in (?) and generate_date >= ? and generate_date <= ?", mid, scoreType, fromdate.Format(upcrmmodel.TimeFmtDate), todate.Format(upcrmmodel.TimeFmtDate)).Find(&result).Error
	return
}

//GetLatestUpScoreDate 获取某个分数的最新记录日期，如果出错，就返回todate
func (d *Dao) GetLatestUpScoreDate(mid int64, scoreType int, todate time.Time) (date time.Time, err error) {
	date = todate
	var history upcrmmodel.UpScoreHistory
	err = d.crmdb.Table(upcrmmodel.GetUpScoreHistoryTableName(mid)).
		Select("generate_date").
		Where("mid=? and score_type = ? and generate_date <= ?", mid, scoreType, todate).
		Order("generate_date desc").
		Limit(1).
		Find(&history).
		Error
	if err != nil {
		return
	}
	date = history.GenerateDate.Time()
	return
}
