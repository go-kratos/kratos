package material

import (
	xtime "go-common/library/time"
)

// consts for workflow event

// Category model is the model for material
type Category struct {
	ID    int64      `json:"id" gorm:"column:id"`
	Name  string     `json:"name" gorm:"column:name"`
	State int8       `json:"state" gorm:"column:state"`
	Type  int64      `json:"type" gorm:"column:type"`
	Rank  int64      `json:"rank" gorm:"column:rank"`
	New   int64      `json:"new" gorm:"column:new"`
	CTime xtime.Time `json:"ctime" gorm:"column:ctime"`
	MTime xtime.Time `json:"mtime" gorm:"column:mtime"`
}

// TableName is used to identify table name in gorm
func (Category) TableName() string {
	return "material_category"
}

// CategoryPager def.
type CategoryPager struct {
	Items []*Category `json:"items"`
	Pager *Pager      `json:"pager"`
}

// WithCategory model is the model for material
type WithCategory struct {
	ID         int64 `json:"id" gorm:"column:id"`
	CategoryID int64 `json:"category_id" gorm:"column:category_id"`
	MaterialID int64 `json:"material_id" gorm:"column:material_id"`
	State      int8  `json:"state" gorm:"column:state"`
	Index      int64 `json:"index" gorm:"column:index"`
}

// TableName is used to identify table name in gorm
func (WithCategory) TableName() string {
	return "material_with_category"
}

// WithCategoryPager def.
type WithCategoryPager struct {
	TotalCount int64           `json:"total_count"`
	Pn         int             `json:"pn"`
	Ps         int             `json:"ps"`
	Items      []*WithCategory `json:"items"`
}

// CategoryParam is used to parse user request
type CategoryParam struct {
	ID    int64  `form:"id" gorm:"column:id"`
	Type  int64  `form:"type" gorm:"column:type" validate:"required"`
	UID   int64  `form:"uid" gorm:"column:uid"`
	Name  string `form:"name" gorm:"column:name" validate:"required"`
	Rank  int64  `form:"rank" gorm:"column:rank" validate:"required"`
	New   int8   `form:"new" gorm:"column:new"`
	State int8   `form:"state" gorm:"column:state"`
}

// TableName is used to identify table name in gorm
func (CategoryParam) TableName() string {
	return "material_category"
}
