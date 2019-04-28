package model

import (
	xtime "go-common/library/time"
)

// VIDEO actiivty types .
const (
	VIDEO        = 1
	PICTURE      = 2
	DRAWYOO      = 3
	VIDEOLIKE    = 4
	PICTURELIKE  = 5
	DRAWYOOLIKE  = 6
	TEXT         = 7
	TEXTLIKE     = 8
	ONLINEVOTE   = 9
	QUESTION     = 10
	LOTTERY      = 11
	ARTICLE      = 12
	VIDEO2       = 13
	MUSIC        = 15
	PHONEVIDEO   = 16
	SMALLVIDEO   = 17
	RESERVATION  = 18
	MISSIONGROUP = 19
)

// SidSub def
type SidSub struct {
	Type int     `form:"type" validate:"required"`
	Lids []int64 `form:"lids,split" validate:"max=50,min=1,dive,min=1"`
}

// ListSub def
type ListSub struct {
	Page     int    `form:"page" default:"1" validate:"min=1"`
	PageSize int    `form:"pagesize" default:"15" validate:"min=1"`
	Keyword  string `form:"keyword"`
	States   []int  `form:"state,split" default:"0"`
	Types    []int  `form:"type,split" default:"0"`
	Sctime   int64  `form:"sctime"`
	Ectime   int64  `form:"ectime"`
}

// SubListRes .
type SubListRes struct {
	List []*ActSubject `json:"list"`
	Page *PageRes      `json:"page"`
}

// PageRes .
type PageRes struct {
	Num   int   `json:"num"`
	Size  int   `json:"size"`
	Total int64 `json:"total"`
}

// AddList def
type AddList struct {
	ActSubject
	Protocol  string     `form:"protocol"`
	Types     string     `form:"types"`
	Pubtime   xtime.Time `form:"pubtime" time_format:"2006-01-02 15:04:05"`
	Deltime   xtime.Time `form:"deltime" time_format:"2006-01-02 15:04:05"`
	Editime   xtime.Time `form:"editime" time_format:"2006-01-02 15:04:05"`
	Tags      string     `form:"tags"`
	Interval  int        `form:"interval"`
	Tlimit    int        `form:"tlimit"`
	Ltime     int        `form:"ltime"`
	Hot       int        `form:"hot"`
	BgmID     int64      `form:"bgm_id"`
	PasterID  int64      `form:"paster_id"`
	Oids      string     `from:"oids"`
	ScreenSet int        `form:"screen_set" default:"1"`
}

//ActSubjectProtocol def
type ActSubjectProtocol struct {
	ID        int64      `json:"id" form:"id" gorm:"column:id"`
	Sid       int64      `json:"sid" form:"sid"`
	Protocol  string     `json:"protocol" form:"protocol"`
	Mtime     xtime.Time `json:"mtime" form:"mtime" time_format:"2006-01-02 15:04:05"`
	Ctime     xtime.Time `json:"ctime" form:"ctime" time_format:"2006-01-02 15:04:05"`
	Types     string     `json:"types" form:"types"`
	Tags      string     `json:"tags" form:"tags"`
	Hot       int        `json:"hot" form:"hot"`
	Pubtime   xtime.Time `json:"pubtime" form:"pubtime" time_format:"2006-01-02 15:04:05"`
	Deltime   xtime.Time `json:"deltime" form:"deltime" time_format:"2006-01-02 15:04:05"`
	Editime   xtime.Time `json:"editime" form:"editime" time_format:"2006-01-02 15:04:05"`
	BgmID     int64      `json:"bgm_id" form:"bgm_id" gorm:"column:bgm_id"`
	PasterID  int64      `json:"paster_id" form:"paster_id" gorm:"column:paster_id"`
	Oids      string     `json:"oids" form:"oids" gorm:"column:oids"`
	ScreenSet int        `json:"screen_set" form:"screen_set" gorm:"column:screen_set"`
}

//ActTimeConfig def
type ActTimeConfig struct {
	ID       int64      `json:"id" form:"id" gorm:"column:id"`
	Sid      int64      `json:"sid" form:"sid"`
	Interval int        `json:"interval" form:"interval"`
	Ctime    xtime.Time `json:"ctime" form:"ctime" time_format:"2006-01-02 15:04:05"`
	Mtime    xtime.Time `json:"mtime" form:"mtime" time_format:"2006-01-02 15:04:05"`
	Tlimit   int        `json:"tlimit" form:"tlimit"`
	Ltime    int        `json:"ltime" form:"ltime"`
}

// ActSubject def.
type ActSubject struct {
	ID         int64      `json:"id,omitempty" form:"id" gorm:"column:id"`
	Oid        int64      `json:"oid,omitempty" form:"oid"`
	Type       int        `json:"type,omitempty" form:"type"`
	State      int        `json:"state,omitempty" form:"state"`
	Level      int        `json:"level,omitempty" form:"level"`
	Flag       int64      `json:"flag,omitempty" form:"flag"`
	Rank       int64      `json:"rank,omitempty" form:"rank"`
	Stime      xtime.Time `json:"stime,omitempty" form:"stime" time_format:"2006-01-02 15:04:05"`
	Etime      xtime.Time `json:"etime,omitempty" form:"etime" time_format:"2006-01-02 15:04:05"`
	Ctime      xtime.Time `json:"ctime,omitempty" form:"ctime" time_format:"2006-01-02 15:04:05"`
	Mtime      xtime.Time `json:"mtime,omitempty" form:"mtime" time_format:"2006-01-02 15:04:05"`
	Lstime     xtime.Time `json:"lstime,omitempty" form:"lstime" time_format:"2006-01-02 15:04:05"`
	Letime     xtime.Time `json:"letime,omitempty" form:"letime" time_format:"2006-01-02 15:04:05"`
	Uetime     xtime.Time `json:"uetime,omitempty" form:"uetime" time_format:"2006-01-02 15:04:05"`
	Ustime     xtime.Time `json:"ustime,omitempty" form:"ustime" time_format:"2006-01-02 15:04:05"`
	Name       string     `json:"name,omitempty" form:"name"`
	Author     string     `json:"author,omitempty" form:"author"`
	ActURL     string     `json:"act_url,omitempty" form:"act_url"`
	Cover      string     `json:"cover,omitempty" form:"cover"`
	Dic        string     `json:"dic,omitempty" form:"dic"`
	H5Cover    string     `json:"h5_cover,omitempty" form:"h5_cover"`
	LikeLimit  int        `json:"like_limit" form:"like_limit"`
	AndroidURL string     `json:"android_url"`
	IosURL     string     `json:"ios_url"`
}

// ActSubjectResult .
type ActSubjectResult struct {
	*ActSubject
	Aids []int64 `json:"aids,omitempty"`
}

// Like def.
type Like struct {
	ID       int64       `json:"id" form:"id" gorm:"column:id"`
	Sid      int64       `json:"sid" form:"sid"`
	Type     int         `json:"type" form:"type"`
	Mid      int64       `json:"mid" form:"mid"`
	Wid      int64       `json:"wid" form:"wid"`
	State    int         `json:"state" form:"state"`
	StickTop int         `json:"stick_top" form:"stick_top"`
	Ctime    xtime.Time  `json:"ctime" form:"ctime" time_format:"2006-01-02 15:04:05"`
	Mtime    xtime.Time  `json:"mtime" form:"mtime" time_format:"2006-01-02 15:04:05"`
	Object   interface{} `json:"object,omiempty" gorm:"-"`
	Like     int64       `json:"like,omiempty" gorm:"-"`
}

//LikeAction def
type LikeAction struct {
	ID     int64      `form:"id" gorm:"column:id"`
	Lid    int64      `form:"lid"`
	Mid    int64      `form:"mid"`
	Action int64      `form:"action"`
	Ctime  xtime.Time `form:"ctime" time_format:"2006-01-02 15:04:05"`
	Mtime  xtime.Time `form:"mtime" time_format:"2006-01-02 15:04:05"`
	Sid    int64      `form:"sid"`
	IP     int64      `form:"ip" gorm:"column:ip"`
}

// TableName LikeAction def
func (LikeAction) TableName() string {
	return "like_action"
}

// TableName ActMatchs def.
func (ActSubject) TableName() string {
	return "act_subject"
}

// TableName Likes def
func (Like) TableName() string {
	return "likes"
}

// TableName ActSubjectProtocol def
func (ActSubjectProtocol) TableName() string {
	return "act_subject_protocol"
}

// TableName ActTimeConfig def
func (ActTimeConfig) TableName() string {
	return "act_time_config"
}
