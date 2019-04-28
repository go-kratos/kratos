package upcrmdao

import (
	"fmt"
	"github.com/siddontang/go-mysql/mysql"
	"go-common/app/service/main/upcredit/model/upcrmmodel"
	"go-common/library/log"
	xtime "go-common/library/time"
	"strings"
	"time"
)

const (
	//TimeFmtMysql mysql time format
	TimeFmtMysql = mysql.TimeFormat
	//TimeFmtDate with only date
	TimeFmtDate = "2006-01-02"
)

//UpQualityInfo struct
type UpQualityInfo struct {
	Mid          int64  `json:"mid"`
	QualityValue int    `json:"quality_value"`
	PrValue      int    `json:"pr_value"`
	Cdate        string `json:"cdate"` // 产生时间 "2018-01-01"
}

// AsPrScore copy to db struct
func (u *UpQualityInfo) AsPrScore() (history *upcrmmodel.UpScoreHistory) {
	if u == nil {
		return &upcrmmodel.UpScoreHistory{}
	}
	history = &upcrmmodel.UpScoreHistory{
		Mid:       u.Mid,
		ScoreType: upcrmmodel.ScoreTypePr,
		Score:     u.PrValue,
	}
	var date, _ = time.Parse(TimeFmtDate, u.Cdate)
	history.GenerateDate = xtime.Time(date.Unix())
	return
}

// AsQualityScore copy to db struct
func (u *UpQualityInfo) AsQualityScore() (history *upcrmmodel.UpScoreHistory) {
	if u == nil {
		return &upcrmmodel.UpScoreHistory{}
	}
	history = &upcrmmodel.UpScoreHistory{
		Mid:       u.Mid,
		ScoreType: upcrmmodel.ScoreTypeQuality,
		Score:     u.QualityValue,
	}
	var date, _ = time.Parse(TimeFmtDate, u.Cdate)
	history.GenerateDate = xtime.Time(date.Unix())
	return
}

//UpdateCreditScore update score
func (d *Dao) UpdateCreditScore(score int, mid int64) (affectRow int64, err error) {
	var db = d.crmdb.Model(upcrmmodel.UpBaseInfo{}).Where("mid = ? and business_type = 1", mid).Update("credit_score", score)
	return db.RowsAffected, db.Error
}

//UpdateQualityAndPrScore update score
func (d *Dao) UpdateQualityAndPrScore(prScore int, qualityScore int, mid int64) (affectRow int64, err error) {
	var db = d.crmdb.Model(upcrmmodel.UpBaseInfo{}).Where("mid = ? and business_type = 1", mid).Update(map[string]int{"pr_score": prScore, "quality_score": qualityScore})
	return db.RowsAffected, db.Error
}

//InsertScoreHistory insert into score history
func (d *Dao) InsertScoreHistory(info *UpQualityInfo) (affectRow int64, err error) {
	var qualityScoreSt = info.AsQualityScore()
	err = d.crmdb.Save(qualityScoreSt).Error
	if err != nil {
		log.Error("insert quality score error, err=%+v", err)
	}
	var prScore = info.AsPrScore()
	err = d.crmdb.Save(prScore).Error
	if err != nil {
		log.Error("insert pr score error, err=%+v", err)
	}
	return
}

//InsertBatchScoreHistory insert batch sql
func (d *Dao) InsertBatchScoreHistory(infoList []*UpQualityInfo, tablenum int) (affectRow int64, err error) {
	var batchSQL = fmt.Sprintf("insert into up_scores_history_%02d (mid, score_type, score, generate_date) values ", tablenum)
	var valueString []string
	var valueArgs []interface{}
	for _, info := range infoList {
		valueString = append(valueString, "(?,?,?,?),(?,?,?,?)")
		valueArgs = append(valueArgs, info.Mid, upcrmmodel.ScoreTypePr, info.PrValue, info.Cdate)
		valueArgs = append(valueArgs, info.Mid, upcrmmodel.ScoreTypeQuality, info.QualityValue, info.Cdate)
	}
	var db = d.crmdb.Exec(batchSQL+strings.Join(valueString, ","), valueArgs...)
	affectRow = db.RowsAffected
	err = db.Error
	return
}
