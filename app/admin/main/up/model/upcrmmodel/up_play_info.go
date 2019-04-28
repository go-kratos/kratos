package upcrmmodel

import "go-common/library/time"

var (
	//BusinessTypeVideo 1
	BusinessTypeVideo = 1
	//BusinessTypeAudio 2
	BusinessTypeAudio = 2
	//BusinessTypeArticle 3
	BusinessTypeArticle = 3
)

//UpPlayInfo  struct
type UpPlayInfo struct {
	ID                  uint32    `gorm:"column:id" json:"-"`
	Mid                 int64     `gorm:"column:mid" json:"mid"`
	BusinessType        int32     `gorm:"column:business_type" json:"-"`
	PlayCountAccumulate int64     `gorm:"column:play_count_accumulate" json:"play_count_accumulate"`
	ArticleCount        int64     `gorm:"column:article_count" json:"article_count"`
	PlayCount90Day      int64     `gorm:"column:play_count_90day" json:"play_count_90_day"`
	PlayCount30Day      int64     `gorm:"column:play_count_30day" json:"play_count_30_day"`
	PlayCount7Day       int64     `gorm:"column:play_count_7day" json:"play_count_7_day"`
	CTime               time.Time `gorm:"column:ctime" json:"-"`
	MTime               time.Time `gorm:"column:mtime" json:"-"`
}
