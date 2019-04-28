package model

import (
	"go-common/app/service/main/archive/api"
	"go-common/library/time"
)

const (
	// EnvPro is pro.
	EnvPro = "pro"
	// EnvTest is env.
	EnvTest = "test"
	// EnvDev is env.
	EnvDev = "dev"
)

const (
	// TnameMaxLen tag name length .
	TnameMaxLen = 32
	// MaxTagNum max tag num.
	MaxTagNum = 50

	// MaxTopicSortNum .
	MaxTopicSortNum = 10
	// MaxChannelSortNum .
	MaxChannelSortNum = 400

	// TagStateNormal tag
	TagStateNormal = 0
	// TagStateDel .
	TagStateDel = 1
	// TagStateHide .
	TagStateHide = 2

	// RoleUp user role.
	RoleUp = 0
	// RoleUser user role.
	RoleUser = 1
	// RoleAdmin user role.
	RoleAdmin = 2

	// TagAdd tag operation type
	TagAdd = int8(1)
	// TagDel tag operation type
	TagDel = int8(2)

	// TotalScore .
	TotalScore = int8(100)

	// UserBannedNone .
	UserBannedNone = int32(0)
)

// Report .
type Report struct {
	ID         int64     `json:"id"`
	Aid        int64     `json:"aid"`
	Tid        int64     `json:"tag_id"`
	OpMid      int64     `json:"opmid"`
	Action     int8      `json:"type"`
	ParentID   int64     `json:"-"`
	PartID     int16     `json:"part_id"`
	Reason     int8      `json:"reason"`
	IsDelMoral int8      `json:"is_del_moral"`
	State      int8      `json:"state"`
	CTime      time.Time `json:"_"`
	MTime      time.Time `json:"-"`
	RptMid     int64     `json:"rptmid"`
	IsFirst    int64     `json:"is_first"`
}

// Detail .
type Detail struct {
	Info    *Tag          `json:"info"`
	Similar []*SimilarTag `json:"similar"`
	News    struct {
		Count    int        `json:"count"`
		Archives []*api.Arc `json:"archives"`
	} `json:"news"`
}

// UploadTag .
type UploadTag struct {
	Rid        int64  `json:"rid"`
	Tid        int64  `json:"tid"`
	Tname      string `json:"tname"`
	Rank       int64  `json:"rank"`
	IsBusiness int8   `json:"-"`
}

// Filter .
type Filter struct {
	Level int    `json:"level"`
	Msg   string `json:"msg"`
}

// Synonym .
type Synonym struct {
	Parent int64   `json:"parent"`
	Childs []int64 `json:"childs"`
}

// TagInfo TagInfo.
type TagInfo struct {
	ID           int64     `json:"id"`
	Type         int32     `json:"type"`
	Name         string    `json:"name"`
	Cover        string    `json:"cover"`
	HeadCover    string    `json:"head_cover"`
	Content      string    `json:"content"`
	ShortContent string    `json:"short_content"`
	Verify       int32     `json:"-"`
	Attr         int32     `json:"-"`
	Attention    int32     `json:"attention"`
	State        int32     `json:"-"`
	Bind         int64     `json:"bind,omitempty"`
	Sub          int64     `json:"sub,omitempty"`
	Activity     int32     `json:"activity"`
	INTShield    int32     `json:"int_shield"` // International Shield. 国际版是否屏蔽
	CTime        time.Time `json:"ctime"`
	MTime        time.Time `json:"-"`
}

// TagTop  web-interface tag top struct, include tag info, and similar tags.
type TagTop struct {
	Tag      *Tag          `json:"tag"`
	Similars []*SimilarTag `json:"similars"`
}
