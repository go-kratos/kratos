package upcrmmodel

import (
	"fmt"
	"go-common/library/time"
)

const (
	//ScoreTypeQuality 1
	ScoreTypeQuality = 1 // 质量分
	//ScoreTypePr 2
	ScoreTypePr = 2 // 影响力分
	//ScoreTypeCredit 3
	ScoreTypeCredit = 3 // 信用分

	//UpScoreHistoryTableCount table count
	UpScoreHistoryTableCount = 100
)

//UpScoreHistory db struct
type UpScoreHistory struct {
	ID           uint      `gorm:"primary_key" json:"-"`
	Mid          int64     `json:"mid"`
	ScoreType    int       `json:"-"`
	Score        int       ` json:"score"`
	GenerateDate time.Time `json:"date"`
	CTime        time.Time `gorm:"column:ctime" json:"-"`
	MTime        time.Time `gorm:"column:mtime" json:"-"`
}

func getTableNameUpScoreHistory(mid int64) string {
	return fmt.Sprintf("up_scores_history_%02d", mid%UpScoreHistoryTableCount)
}

//TableName table name
func (u *UpScoreHistory) TableName() string {
	return getTableNameUpScoreHistory(u.Mid)
}
