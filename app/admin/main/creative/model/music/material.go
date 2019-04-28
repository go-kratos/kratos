package music

import (
	xtime "go-common/library/time"
)

// consts for workflow event

// Material model is the model for music
type Material struct {
	ID    int64      `json:"id" gorm:"column:id"`
	Pid   int64      `json:"pid" gorm:"column:pid"`
	Name  string     `json:"name" gorm:"column:name"`
	Index int64      `json:"index" gorm:"column:index"`
	State int8       `json:"state" gorm:"column:state"`
	CTime xtime.Time `json:"ctime" gorm:"column:ctime"`
	MTime xtime.Time `json:"mtime" gorm:"column:mtime"`
}

// TableName is used to identify table name in gorm
func (Material) TableName() string {
	return "music_material"
}

// MaterialPager def.
type MaterialPager struct {
	Items []*Material `json:"items"`
	Pager *Pager      `json:"pager"`
}

// MaterialMixParent model is the model for music
type MaterialMixParent struct {
	Material
	PName string `json:"p_name" gorm:"column:p_name"`
}

// TableName is used to identify table name in gorm
func (MaterialMixParent) TableName() string {
	return "music_material"
}

// MaterialMixParentPager def.
type MaterialMixParentPager struct {
	Items []*MaterialMixParent `json:"items"`
	Pager *Pager               `json:"pager"`
}

// WithMaterial model is the model for music
type WithMaterial struct {
	ID    int64      `json:"id" gorm:"column:id"`
	Sid   int64      `json:"sid" gorm:"column:sid"`
	Tid   int64      `json:"tid" gorm:"column:tid"`
	State int8       `json:"state" gorm:"column:state"`
	Index int64      `json:"index" gorm:"column:index"`
	CTime xtime.Time `json:"ctime" gorm:"column:ctime"`
	MTime xtime.Time `json:"mtime" gorm:"column:mtime"`
}

// TableName is used to identify table name in gorm
func (WithMaterial) TableName() string {
	return "music_with_material"
}

// WithMaterialPager def.
type WithMaterialPager struct {
	Pager *Pager          `json:"pager"`
	Items []*WithMaterial `json:"items"`
}
