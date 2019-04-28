package upcrmmodel

import "go-common/library/time"
import xtime "time"

// const table name
const (
	TableNameUpBaseInfo = "up_base_info"
)

//UpBaseInfo  struct
type UpBaseInfo struct {
	ID                     uint32     `gorm:"column:id"`
	Mid                    int64      `gorm:"column:mid"`
	Name                   string     `gorm:"column:name"`
	Sex                    int8       `gorm:"column:sex"`
	JoinTime               time.Time  `gorm:"column:join_time"`
	FirstUpTime            time.Time  `gorm:"column:first_up_time"`
	Level                  int16      `gorm:"column:level"`
	FansCount              int        `gorm:"column:fans_count"`
	AccountState           int8       `gorm:"column:account_state"`
	Activity               int        `gorm:"column:activity"`
	ArticleCount30day      int        `gorm:"column:article_count_30day"`
	ArticleCountAccumulate int        `gorm:"column:article_count_accumulate"`
	Birthday               xtime.Time `gorm:"column:birthday"`
	ActiveCity             string     `gorm:"column:active_city"`     // 市，存的是城市的名字
	ActiveProvince         string     `gorm:"column:active_province"` // 省，省的名字
	VerifyType             int8       `gorm:"column:verify_type"`
	BusinessType           int8       `gorm:"column:business_type"`
	CreditScore            int        `gorm:"column:credit_score"`
	PrScore                int        `gorm:"column:pr_score"`
	QualityScore           int        `gorm:"column:quality_score"`
	ActiveTid              int64      `gorm:"column:active_tid"`
	Attr                   int        `gorm:"column:attr"`
	CTime                  time.Time  `gorm:"column:ctime"`
	MTime                  time.Time  `gorm:"column:mtime"`
}

//TableName get table name
func (s *UpBaseInfo) TableName() string {
	return TableNameUpBaseInfo
}
