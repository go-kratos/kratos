package model

const (
	// TypeUser .
	TypeUser = int32(0) // 普通tag // TypeUser tag type
	// TypeUpper .
	TypeUpper = int32(1) // up主tag
	// TypeOfficailCategory .
	TypeOfficailCategory = int32(2) // 官方-分类tag
	// TypeOfficailContent .
	TypeOfficailContent = int32(3) // 官方-内容tag
	// TypeOfficailActivity .
	TypeOfficailActivity = int32(4) // 官方-活动tag

	// TagStateNormal .
	TagStateNormal = int32(0) // tag state
	// TagStateDelete .
	TagStateDelete = int32(1)
	// TagStateHide .
	TagStateHide = int32(2)

	// AttrNo .
	AttrNo = int32(0) // attr
	// AttrYes .
	AttrYes = int32(1)

	// SpamActionAdd .
	SpamActionAdd = int32(1) // spam
	// SpamActionDel .
	SpamActionDel = int32(2)

	// TnameMaxLen .
	TnameMaxLen = 32
	// MaxSubNum MaxSubNum.
	MaxSubNum = 400
	// UserBannedNone .
	UserBannedNone = int32(0)

	// ChannelMaxGroups channel max groups num.
	ChannelMaxGroups = int32(8)
)

// Detail .
type Detail struct {
	Info    *Tag          `json:"info"`
	Similar []*TagSimilar `json:"similar"`
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

// HotTags .
type HotTags struct {
	Rid  int64     `json:"rid"`
	Tags []*HotTag `json:"tags"`
}

// HotTag .
type HotTag struct {
	Rid       int64  `json:"-"`
	Tid       int64  `json:"tid"`
	Tname     string `json:"tname"`
	HighLight int64  `json:"highlight"`
	IsAtten   int8   `json:"is_atten"`
}

// UploadTag .
type UploadTag struct {
	Rid        int64  `json:"rid"`
	Tid        int64  `json:"tid"`
	Tname      string `json:"tname"`
	Rank       int64  `json:"rank"`
	IsBusiness int8   `json:"-"`
}
