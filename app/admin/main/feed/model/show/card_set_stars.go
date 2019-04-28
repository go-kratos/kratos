package show

import (
	"go-common/app/admin/main/feed/model/common"
	xtime "go-common/library/time"
)

//PopularStars channel tab
type PopularStars struct {
	ID        int64      `json:"id"`
	Type      string     `json:"type"`
	Value     string     `json:"value"`
	Title     string     `json:"title"`
	LongTitle string     `json:"longtitle"`
	Content   string     `json:"content"`
	Deleted   int        `json:"deleted"`
	Person    string     `json:"person"`
	Source    int        `json:"source"`
	Status    int        `json:"status"`
	Mtime     xtime.Time `json:"mtime"`
}

//PopularStarsPager .
type PopularStarsPager struct {
	Item []*PopularStars `json:"item"`
	Page common.Page     `json:"page"`
}

// TableName .
func (a PopularStars) TableName() string {
	return "card_set"
}

/*
---------------------------
 struct param
---------------------------
*/

//PopularStarsAP popular stars add param
type PopularStarsAP struct {
	Type      string `form:"type" validate:"required"`
	Value     string `form:"value" validate:"required"`
	Title     string `form:"title" validate:"required"`
	LongTitle string `form:"longtitle" validate:"required"`
	Content   string `form:"content" validate:"required"`
	UID       int64  `form:"person" gorm:"column:uid"`
	Person    string `form:"person"`
	Source    int
	Status    int
}

//PopularStarsAIAP popular stars ai add param
type PopularStarsAIAP struct {
	Mid  int64   `form:"mid"`
	Aids []int64 `form:"aids"`
}

//AiValue ai insert value
type AiValue struct {
	ID int64 `json:"id"`
}

//PopularStarsUP channel tab update param
type PopularStarsUP struct {
	ID        int64  `form:"id" validate:"required"`
	Type      string `form:"type" validate:"required"`
	Value     string `form:"value" validate:"required"`
	Title     string `form:"title" validate:"required"`
	LongTitle string `form:"longtitle"`
	Content   string `form:"content" validate:"required"`
	Status    int    `form:"status"`
}

//PopularStarsLP channel tab list param
type PopularStarsLP struct {
	ID        int    `form:"id"`
	Title     string `form:"title"`
	LongTitle string `form:"longtitle"`
	Person    string `form:"person"`
	Source    int    `form:"source" default:"-1"`
	Status    int    `form:"status"`
	Ps        int    `form:"ps" default:"20"` // 分页大小
	Pn        int    `form:"pn" default:"1"`  // 第几个分页
}

// TableName .
func (a PopularStarsAP) TableName() string {
	return "card_set"
}

// TableName .
func (a PopularStarsUP) TableName() string {
	return "card_set"
}
