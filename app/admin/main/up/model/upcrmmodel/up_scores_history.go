package upcrmmodel

import (
	"fmt"
	"go-common/library/time"

	"github.com/siddontang/go-mysql/mysql"
)

var (
	// TimeFmtMysql mysql time format
	TimeFmtMysql = mysql.TimeFormat
	// TimeFmtDate with only date
	TimeFmtDate = "2006-01-02"
	// TimeFmtDateTime .
	TimeFmtDateTime = "2006-01-02 15:04:05"
)

//UpScoreHistory  struct
type UpScoreHistory struct {
	ID           uint32    `gorm:"column:id"`
	Mid          int64     `gorm:"column:mid"`
	ScoreType    int8      `gorm:"column:score_type"`
	Score        int       `gorm:"column:score"`
	GenerateDate time.Time `gorm:"column:generate_date"`
	CTime        time.Time `gorm:"column:ctime"`
	MTime        time.Time `gorm:"column:mtime"`
}

//TableName table name
func (u *UpScoreHistory) TableName() string {
	return GetUpScoreHistoryTableName(u.Mid)
}

//GetUpScoreHistoryTableName table name
func GetUpScoreHistoryTableName(mid int64) string {
	return fmt.Sprintf("up_scores_history_%02d", mid%100)
}
