package music

import (
	xtime "go-common/library/time"
)

// consts for workflow event
const (
	MusicDelete = -100
	MusicOpen   = 0
)

// Music model is the model for music
type Music struct {
	ID           int64      `json:"id" gorm:"column:id"`
	Sid          int64      `json:"sid" gorm:"column:sid"`
	Name         string     `json:"name" gorm:"column:name"`
	FrontName    string     `json:"frontname" gorm:"column:frontname"`
	Musicians    string     `json:"musicians" gorm:"column:musicians"`
	Cooperate    int8       `json:"cooperate" gorm:"column:cooperate"`
	Mid          int64      `json:"mid" gorm:"column:mid"`
	Tid          int64      `json:"tid" gorm:"-"`
	Tags         string     `json:"tags" gorm:"tags"`
	Timeline     string     `json:"timeline" gorm:"timeline"`
	Rid          int64      `json:"rid" gorm:"-"`
	Pid          int64      `json:"pid" gorm:"-"`
	Cover        string     `json:"cover" gorm:"column:cover"`
	MaterialName string     `json:"material_name" gorm:"-"`
	CategoryName string     `json:"category_name" gorm:"-"`
	Stat         string     `json:"stat" gorm:"column:stat"`
	Categorys    string     `json:"categorys" gorm:"column:categorys"`
	Playurl      string     `json:"playurl" gorm:"column:playurl"`
	State        int8       `json:"state" gorm:"column:state"`
	Duration     int32      `json:"duration" gorm:"column:duration"`
	Filesize     int32      `json:"filesize" gorm:"column:filesize"`
	PubTime      xtime.Time `json:"pubtime" gorm:"column:pubtime"`
	SyncTime     xtime.Time `json:"synctime" gorm:"column:synctime"`
	CTime        xtime.Time `json:"ctime" gorm:"column:ctime"`
	MTime        xtime.Time `json:"mtime" gorm:"column:mtime"`
}

// TableName is used to identify table name in gorm
func (Music) TableName() string {
	return "music"
}

// ResultPager def.
type ResultPager struct {
	Items []*Music `json:"items"`
	Pager *Pager   `json:"pager"`
}

// Param is used to parse user request
type Param struct {
	ID        int64      `form:"id" gorm:"column:id"`
	Sid       int64      `form:"sid" validate:"required"`
	Name      string     `form:"name" validate:"required"`
	Musicians string     `form:"musicians"`
	Mid       int64      `form:"mid" validate:"required"`
	Cover     string     `form:"cover" validate:"required"`
	Stat      string     `form:"stat" `
	Categorys string     `form:"categorys" `
	Playurl   string     `form:"playurl" `
	State     int8       `form:"state"`
	Duration  int32      `form:"duration" `
	Filesize  int32      `form:"filesize" `
	UID       int64      `form:"uid" `
	PubTime   xtime.Time `form:"pubtime"`
	SyncTime  xtime.Time `form:"synctime"`
	Tags      string     `form:"tags"`
	Timeline  string     `form:"timeline"`
}

// TableName is used to identify table name in gorm
func (Param) TableName() string {
	return "music"
}

// CategoryParam is used to parse user request
type CategoryParam struct {
	ID          int64  `form:"id" gorm:"column:id"`
	Pid         int64  `form:"pid" gorm:"column:pid"`
	UID         int64  `form:"uid" gorm:"column:uid"`
	Name        string `form:"name" gorm:"column:name" validate:"required"`
	Index       int64  `form:"index" gorm:"column:index" validate:"required"`
	CameraIndex int64  `form:"camera_index" gorm:"column:camera_index" validate:"required"`
	State       int8   `form:"state" gorm:"column:state"`
}

// TableName is used to identify table name in gorm
func (CategoryParam) TableName() string {
	return "music_category"
}

// MaterialParam is used to parse user request
type MaterialParam struct {
	ID    int64  `form:"id" gorm:"column:id"`
	Pid   int64  `form:"pid" gorm:"column:pid"`
	UID   int64  `form:"uid" gorm:"column:uid"`
	Name  string `form:"name" gorm:"column:name" validate:"required"`
	Index int64  `form:"index" gorm:"column:index"`
	State int8   `form:"state" gorm:"column:state"`
}

// TableName is used to identify table name in gorm
func (MaterialParam) TableName() string {
	return "music_material"
}

// WithMaterialParam is used to parse user request
type WithMaterialParam struct {
	ID    int64 `form:"id" gorm:"column:id"`
	UID   int64 `form:"uid" gorm:"column:uid"`
	Sid   int64 `form:"sid" gorm:"column:sid" validate:"required,min=1"`
	Tid   int64 `form:"tid" gorm:"column:tid" validate:"required,min=1"`
	Index int64 `form:"index" gorm:"column:index"`
	State int8  `form:"state" gorm:"column:state"`
}

// IndexParam is used to parse user request
type IndexParam struct {
	ID          int64 `form:"id"  validate:"required"`
	Index       int64 `form:"index" validate:"required"`
	SwitchID    int64 `form:"switch_id"  validate:"required"`
	SwitchIndex int64 `form:"switch_index" validate:"required"`
	UID         int64 `form:"uid"`
}

// TableName is used to identify table name in gorm
func (WithMaterialParam) TableName() string {
	return "music_with_material"
}

// BatchMusicWithMaterialParam is used to parse user request
type BatchMusicWithMaterialParam struct {
	UID     int64  `form:"uid" gorm:"column:uid"`
	SidList string `form:"sid_list"  validate:"required"`
	Tid     int64  `form:"tid" gorm:"column:tid" validate:"required,min=1"`
	Index   int64  `form:"index" gorm:"column:index"`
	State   int8   `form:"state" gorm:"column:state"`
}

// TableName is used to identify table name in gorm
func (BatchMusicWithMaterialParam) TableName() string {
	return "music_with_material"
}

// WithCategoryParam is used to parse user request
type WithCategoryParam struct {
	ID    int64 `form:"id" gorm:"column:id"`
	UID   int64 `form:"uid" gorm:"column:uid"`
	Sid   int64 `form:"sid" gorm:"column:sid" validate:"required,min=1"`
	Tid   int64 `form:"tid" gorm:"column:tid" validate:"required,min=1"`
	Index int64 `form:"index" gorm:"column:index" default:"1"`
	State int8  `form:"state" gorm:"column:state"`
}

// BatchMusicWithCategoryParam is used to parse user request
type BatchMusicWithCategoryParam struct {
	UID      int64  `form:"uid" gorm:"column:uid"`
	SidList  string `form:"sid_list"  validate:"required"`
	SendList string `form:"send_list"`
	Tid      int64  `form:"tid" gorm:"column:tid" validate:"required,min=1"`
	Index    int64  `form:"index" gorm:"column:index" default:"1"`
	State    int8   `form:"state" gorm:"column:state"`
}

// TableName is used to identify table name in gorm
func (BatchMusicWithCategoryParam) TableName() string {
	return "music_with_category"
}

// TableName is used to identify table name in gorm
func (WithCategoryParam) TableName() string {
	return "music_with_category"
}

// LogParam is used to parse user request
type LogParam struct {
	ID     int64  `json:"id"`
	UID    int64  `json:"uid"`
	UName  string `json:"uname"`
	Action string `json:"action"`
	Name   string `json:"name"`
}
