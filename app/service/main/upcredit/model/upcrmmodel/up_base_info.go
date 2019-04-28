package upcrmmodel

import "go-common/library/time"

//UpBaseInfo db struct
type UpBaseInfo struct {
	ID                     int32
	Mid                    int32
	Name                   string
	Sex                    int8
	JoinTime               time.Time
	FirstUpTime            time.Time
	Level                  int16
	FansCount              int
	AccountState           int
	Activity               int
	ArticleCount30day      int `gorm:"article_count_30day"`
	ArticleCountAccumulate int
	VerifyType             int8
	CTime                  time.Time `gorm:"ctime"`
	MTime                  time.Time `gorm:"mtim"`
	BusinessType           int8
	CreditScore            int
	PrScore                int
	QualityScore           int
	ActiveTid              int16
	Attr                   int
}
