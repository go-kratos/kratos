package music

import (
	xtime "go-common/library/time"
)

// consts for workflow event

// Category model is the model for music
type Category struct {
	ID          int64      `json:"id" gorm:"column:id"`
	Pid         int64      `json:"pid" gorm:"column:pid"`
	Name        string     `json:"name" gorm:"column:name"`
	Index       int64      `json:"index" gorm:"column:index"`
	CameraIndex int64      `json:"camera_index" gorm:"column:camera_index"`
	State       int8       `json:"state" gorm:"column:state"`
	CTime       xtime.Time `json:"ctime" gorm:"column:ctime"`
	MTime       xtime.Time `json:"mtime" gorm:"column:mtime"`
}

// TableName is used to identify table name in gorm
func (Category) TableName() string {
	return "music_category"
}

// CategoryPager def.
type CategoryPager struct {
	Items []*Category `json:"items"`
	Pager *Pager      `json:"pager"`
}

// Pager Pager def.
type Pager struct {
	Num   int   `json:"num"`
	Size  int   `json:"size"`
	Total int64 `json:"total"`
}

// SidNotify model is the model for music
type SidNotify struct {
	Sid      int64 `json:"sid"`
	MidFirst bool  `json:"mid_first"`
	SidFirst bool  `json:"sid_first"`
}

// WithCategory model is the model for music
type WithCategory struct {
	ID    int64      `json:"id" gorm:"column:id"`
	Sid   int64      `json:"sid" gorm:"column:sid"`
	Tid   int64      `json:"tid" gorm:"column:tid"`
	State int8       `json:"state" gorm:"column:state"`
	Index int64      `json:"index" gorm:"column:index"`
	CTime xtime.Time `json:"ctime" gorm:"column:ctime"`
	MTime xtime.Time `json:"mtime" gorm:"column:mtime"`
}

// TableName is used to identify table name in gorm
func (WithCategory) TableName() string {
	return "music_with_category"
}

// WithCategoryPager def.
type WithCategoryPager struct {
	TotalCount int64           `json:"total_count"`
	Pn         int             `json:"pn"`
	Ps         int             `json:"ps"`
	Items      []*WithCategory `json:"items"`
}

//CategoryList .
type CategoryList struct {
	ID           int64      `json:"id" gorm:"column:id"`
	Sid          int64      `json:"sid" gorm:"column:sid"`
	Name         string     `json:"name" gorm:"column:name"`
	FrontName    string     `json:"frontname" gorm:"column:frontname"`
	Musicians    string     `json:"musicians" gorm:"column:musicians"`
	Cooperate    int8       `json:"cooperate" gorm:"column:cooperate"`
	Mid          int64      `json:"mid" gorm:"column:mid"`
	Tid          int64      `json:"tid" gorm:"-"`
	Pid          int64      `json:"pid" gorm:"-"`
	Cover        string     `json:"cover" gorm:"column:cover"`
	MaterialName string     `json:"material_name" gorm:"-"`
	CategoryName string     `json:"category_name" gorm:"-"`
	MusicState   string     `json:"music_state" gorm:"-"`
	Stat         string     `json:"stat" gorm:"column:stat"`
	Categorys    string     `json:"categorys" gorm:"column:categorys"`
	Playurl      string     `json:"playurl" gorm:"column:playurl"`
	State        int8       `json:"state" gorm:"column:state"`
	Index        int        `json:"index" gorm:"column:index"`
	Duration     int32      `json:"duration" gorm:"column:duration"`
	Filesize     int32      `json:"filesize" gorm:"column:filesize"`
	PubTime      xtime.Time `json:"pubtime" gorm:"column:pubtime"`
	SyncTime     xtime.Time `json:"synctime" gorm:"column:synctime"`
	Tags         string     `json:"tags" gorm:"-"`
	Timeline     string     `json:"timeline" gorm:"-"`
	CTime        xtime.Time `json:"ctime" gorm:"column:ctime"`
	MTime        xtime.Time `json:"mtime" gorm:"column:mtime"`
}

// TableName is used to identify table name in gorm
func (CategoryList) TableName() string {
	return "music_with_category"
}

// CategoryListPager def.
type CategoryListPager struct {
	Items []*CategoryList `json:"items"`
	Pager *Pager          `json:"pager"`
}
