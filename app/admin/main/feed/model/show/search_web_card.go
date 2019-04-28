package show

import (
	"go-common/app/admin/main/feed/model/common"
	"go-common/library/time"
)

//SearchWebCard web card
type SearchWebCard struct {
	ID      int64     `form:"id" gorm:"column:id" json:"id"`
	Type    int64     `form:"type" gorm:"column:type" json:"type"`
	Title   string    `form:"title" gorm:"column:title" json:"title"`
	Desc    string    `form:"desc" gorm:"column:desc" json:"desc"`
	Cover   string    `form:"cover" gorm:"column:cover" json:"cover"`
	ReType  int64     `form:"re_type" gorm:"column:re_type" json:"re_type"`
	ReValue string    `form:"re_value" gorm:"column:re_value" json:"re_value"`
	Corner  string    `form:"corner" gorm:"column:corner" json:"corner"`
	Person  string    `form:"person" gorm:"column:person" json:"person"`
	Deleted int64     `form:"deleted" gorm:"column:deleted" json:"deleted"`
	Ctime   time.Time `form:"string" gorm:"column:ctime" json:"ctime"`
	Mtime   time.Time `form:"string" gorm:"column:mtime" json:"mtime"`
}

//SearchWebCardPager .
type SearchWebCardPager struct {
	Item []*SearchWebCard `json:"item"`
	Page common.Page      `json:"page"`
}

// TableName .
func (a SearchWebCard) TableName() string {
	return "search_web_card"
}

/*
---------------------------
 struct param
---------------------------
*/

//SearchWebCardAP web card add param
type SearchWebCardAP struct {
	Type    int64  `form:"type" gorm:"column:type" json:"type"`
	Title   string `form:"title" gorm:"column:title" json:"title"`
	Desc    string `form:"desc" gorm:"column:desc" json:"desc"`
	Cover   string `form:"cover" gorm:"column:cover" json:"cover"`
	ReType  int64  `form:"re_type" gorm:"column:re_type" json:"re_type"`
	ReValue string `form:"re_value" gorm:"column:re_value" json:"re_value"`
	Corner  string `form:"corner" gorm:"column:corner" json:"corner"`
	Person  string `form:"person" gorm:"column:person" json:"person"`
}

//SearchWebCardUP web card update param
type SearchWebCardUP struct {
	ID      int64  `form:"id" gorm:"column:id" json:"id"`
	Type    int64  `form:"type" gorm:"column:type" json:"type"`
	Title   string `form:"title" gorm:"column:title" json:"title"`
	Desc    string `form:"desc" gorm:"column:desc" json:"desc"`
	Cover   string `form:"cover" gorm:"column:cover" json:"cover"`
	ReType  int64  `form:"re_type" gorm:"column:re_type" json:"re_type"`
	ReValue string `form:"re_value" gorm:"column:re_value" json:"re_value"`
	Corner  string `form:"corner" gorm:"column:corner" json:"corner"`
}

//SearchWebCardLP search web card list param
type SearchWebCardLP struct {
	ID     int    `form:"id"`
	Person string `form:"person"`
	Title  string `form:"title"`
	Ps     int    `form:"ps" default:"20"` // 分页大小
	Pn     int    `form:"pn" default:"1"`  // 第几个分页
	STime  string `form:"stime"`
	ETime  string `form:"etime"`
}

// TableName .
func (a SearchWebCardAP) TableName() string {
	return "search_web_card"
}

// TableName .
func (a SearchWebCardUP) TableName() string {
	return "search_web_card"
}
