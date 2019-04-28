package model

import (
	xtime "go-common/library/time"
)

// resource const
const (
	IconTypeFix     = 1
	IconTypeRandom  = 2
	IconTypeBangumi = 3

	NoCategory = 0
	IsCategory = 1

	AsgTypePic   = int8(0)
	AsgTypeVideo = int8(1)
	// pgc mobile
	AsgTypeURL     = int8(2)
	AsgTypeBangumi = int8(3)
	AsgTypeLive    = int8(4)
	AsgTypeGame    = int8(5)
	AsgTypeAv      = int8(6)
)

// IconTypes icon_type
var IconTypes = map[int]string{
	IconTypeFix:     "fix",
	IconTypeRandom:  "random",
	IconTypeBangumi: "bangumi",
}

// Rule resource_assignmen rule
type Rule struct {
	Cover int32  `json:"is_cover"`
	Style int32  `json:"style"`
	Label string `json:"label"`
	Intro string `json:"intro"`
}

// Resource struct
type Resource struct {
	ID          int           `json:"id"`
	Platform    int           `json:"platform"`
	Name        string        `json:"name"`
	Parent      int           `json:"parent"`
	State       int           `json:"-"`
	Counter     int           `json:"counter"`
	Position    int           `json:"position"`
	Rule        string        `json:"rule"`
	Size        string        `json:"size"`
	Previce     string        `json:"preview"`
	Desc        string        `json:"description"`
	Mark        string        `json:"mark"`
	Assignments []*Assignment `json:"assignments"`
	CTime       xtime.Time    `json:"ctime"`
	MTime       xtime.Time    `json:"mtime"`
	Level       int64         `json:"level"`
	Type        int           `json:"type"`
	IsAd        int           `json:"is_ad"`
}

// Assignment struct
type Assignment struct {
	ID             int        `json:"id"`
	AsgID          int        `json:"-"`
	Name           string     `json:"name"`
	ContractID     string     `json:"contract_id"`
	ResID          int        `json:"resource_id"`
	Pic            string     `json:"pic"`
	LitPic         string     `json:"litpic"`
	URL            string     `json:"url"`
	Rule           string     `json:"rule"`
	Weight         int        `json:"weight"`
	Agency         string     `json:"agency"`
	Price          float32    `json:"price"`
	State          int        `json:"state"`
	Atype          int8       `json:"atype"`
	Username       string     `json:"username"`
	PlayerCategory int8       `json:"player_category"`
	ApplyGroupID   int        `json:"-"`
	STime          xtime.Time `json:"stime"`
	ETime          xtime.Time `json:"etime"`
	CTime          xtime.Time `json:"ctime"`
	MTime          xtime.Time `json:"mtime"`
}

// IndexIcon struct
type IndexIcon struct {
	ID       int64      `json:"id"`
	Type     int        `json:"type"`
	Title    string     `json:"title"`
	State    int        `json:"state"`
	Links    []string   `json:"links"`
	Icon     string     `json:"icon"`
	Weight   int        `json:"weight"`
	UserName string     `json:"-"`
	StTime   xtime.Time `json:"sttime"`
	EndTime  xtime.Time `json:"endtime"`
	DelTime  xtime.Time `json:"deltime"`
	CTime    xtime.Time `json:"ctime"`
	MTime    xtime.Time `json:"mtime"`
}

// PlayerIcon struct
type PlayerIcon struct {
	URL1  string     `json:"url1"`
	Hash1 string     `json:"hash1"`
	URL2  string     `json:"url2"`
	Hash2 string     `json:"hash2"`
	CTime xtime.Time `json:"ctime"`
}

// ResWarnInfo for email
type ResWarnInfo struct {
	AID            int64
	URL            string
	AssignmentID   int
	AssignmentName string
	ResourceName   string
	ResourceID     int
	MaterialID     int
	UserName       string
	STime          xtime.Time `json:"stime"`
	ETime          xtime.Time `json:"etime"`
	ApplyGroupID   int
}

// Cmtbox live danmaku box
type Cmtbox struct {
	ID          int64      `json:"id"`
	LoadCID     int64      `json:"load_cid"`
	Server      string     `json:"server"`
	Port        string     `json:"port"`
	SizeFactor  string     `json:"size_factor"`
	SpeedFactor string     `json:"speed_factor"`
	MaxOnscreen string     `json:"max_onscreen"`
	Style       string     `json:"style"`
	StyleParam  string     `json:"style_param"`
	TopMargin   string     `json:"top_margin"`
	State       string     `json:"state"`
	CTime       xtime.Time `json:"ctime"`
	MTime       xtime.Time `json:"mtime"`
}
