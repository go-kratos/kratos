package upcrmdao

import (
	"go-common/app/service/main/upcredit/model/calculator"
	"go-common/app/service/main/upcredit/model/upcrmmodel"
	xtime "go-common/library/time"
	"time"
)

//InsertScoreSection insert score section
func (d *Dao) InsertScoreSection(statis calculator.OverAllStatistic, scoreType int, date time.Time) error {
	var history upcrmmodel.ScoreSectionHistory
	history.Section0 = statis.GetScore(0)
	history.Section1 = statis.GetScore(1)
	history.Section2 = statis.GetScore(2)
	history.Section3 = statis.GetScore(3)
	history.Section4 = statis.GetScore(4)
	history.Section5 = statis.GetScore(5)
	history.Section6 = statis.GetScore(6)
	history.Section7 = statis.GetScore(7)
	history.Section8 = statis.GetScore(8)
	history.Section9 = statis.GetScore(9)
	history.ScoreType = scoreType
	history.GenerateDate = xtime.Time(date.Unix())
	var now = time.Now().Unix()
	history.CTime = xtime.Time(now)
	history.MTime = xtime.Time(now)
	// insert or update
	return d.crmdb.Create(&history).Error
}
