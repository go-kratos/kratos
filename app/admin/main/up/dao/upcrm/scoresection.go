package upcrm

import (
	"time"

	"go-common/app/admin/main/up/model/upcrmmodel"
	"go-common/library/ecode"
)

const (
	//ScoreTypeQuality 质量分
	ScoreTypeQuality = 1
	//ScoreTypePr 影响力
	ScoreTypePr = 2
	//ScoreTypeCredit 信用分
	ScoreTypeCredit = 3

	//ScoreSectionTableName table name
	ScoreSectionTableName = "score_section_history"
)

//ScoreQueryHistory get history
func (d *Dao) ScoreQueryHistory(scoreType int, date time.Time) (result upcrmmodel.ScoreSectionHistory, err error) {
	err = d.crmdb.Model(&result).
		Where("score_type=? AND generate_date = ?", scoreType, date.Format("2006-01-02")).
		Find(&result).Error
	if err == ecode.NothingFound {
		err = nil
	}
	return
}

//GetLastHistory get last update date
func (d *Dao) GetLastHistory(scoreType int) (lastHistoryDate time.Time, err error) {
	var model upcrmmodel.ScoreSectionHistory
	err = d.crmdb.Table(ScoreSectionTableName).Select("generate_date").Where("score_type=?", scoreType).Order("generate_date desc").Limit(1).Find(&model).Error
	if err != nil {
		return
	}
	lastHistoryDate = model.GenerateDate.Time()
	return
}
