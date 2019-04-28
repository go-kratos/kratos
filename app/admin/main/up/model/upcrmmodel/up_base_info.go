package upcrmmodel

import "go-common/library/time"
import xtime "time"

//UpBaseInfo  struct
type UpBaseInfo struct {
	ID                     uint32 `gorm:"column:id"`
	Mid                    int64
	Name                   string
	Sex                    int8
	JoinTime               time.Time
	FirstUpTime            time.Time
	Level                  int16
	FansCount              int
	AccountState           int8
	Activity               int
	ArticleCount30day      int `gorm:"column:article_count_30day"`
	ArticleCountAccumulate int
	Birthday               xtime.Time
	ActiveCity             string // 市，存的是城市的名字
	ActiveProvince         string // 省，省的名字
	VerifyType             int8
	BusinessType           int8
	CreditScore            int
	PrScore                int
	QualityScore           int
	ActiveTid              int64
	Attr                   int
	CTime                  time.Time `gorm:"column:ctime"`
	MTime                  time.Time `gorm:"column:mtime"`
}
