package model

import (
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
)

// const .
const (
	// TypeUnknown unkown
	TypeUnknown = int8(-1)
	// TypeOffline offline
	TypeOffline = int8(0)
	// TypeBangumi bangumi
	TypeBangumi = int8(1)
	// TypeMovie movie
	TypeMovie = int8(2)
	// TypeUGC UGC
	TypeUGC = int8(3)
	// TypePGC PGC
	TypePGC = int8(4)
	// TypeArticle Article
	TypeArticle = int8(5)
	// TypeLive Live
	TypeLive = int8(6)
	// TypeCorpus corpus
	TypeCorpus = int8(7)
	// TypeComic comic
	TypeComic = int8(8)

	// SubTypeOffline archive subtype
	SubTypeOffline = int8(1)
	// SubTypeBangumi bangumi
	SubTypeBangumi = int8(1)
	// SubTypeFilm film
	SubTypeFilm = int8(2)
	// SubTypeDoc documentary
	SubTypeDoc = int8(3)
	// SubTypeNation nation
	SubTypeNation = int8(4)
	// SubTypeTV TV
	SubTypeTV = int8(5)

	// DeviceUnknown unknown
	DeviceUnknown = int8(0)
	// DeviceIphone iphoneTV
	DeviceIphone = int8(1)
	// DevicePC PC
	DevicePC = int8(2)
	// DeviceAndroid android
	DeviceAndroid = int8(3)
	// DeviceAndroidTV android TV
	DeviceAndroidTV = int8(33)
	// DeviceIpad ipad
	DeviceIpad = int8(4)
	// DeviceWP8 WP8
	DeviceWP8 = int8(5)
	// DeviceUWP UWP
	DeviceUWP = int8(6)

	// ShadowUnknown unknown
	ShadowUnknown = int64(-1)
	// ShadowOff off
	ShadowOff = int64(0)
	// ShadowOn on
	ShadowOn = int64(1)

	// ProComplete progress complete
	ProComplete = int64(-1)

	// PlatformAndroid platform android.
	PlatformAndroid string = "android"
	// PlatformIOS platform ios.
	PlatformIOS string = "ios"

	// DevicePad device pad.
	DevicePad string = "pad"
	// MobileAppAndroidTV mobile app android tv.
	MobileAppAndroidTV string = "android_tv_yst"

	HistoryLog      = 171
	HistoryClear    = "history_clear"
	HistoryClearTyp = "history_clear_%s"
	ToviewClear     = "toview_clear"
)

var businesses = map[string]int8{
	"pgc":          TypePGC,
	"article":      TypeArticle,
	"archive":      TypeUGC,
	"live":         TypeLive,
	"article-list": TypeCorpus,
	"comic":        TypeComic,
}

var businessIDs = map[int8]string{
	TypeOffline: "archive",
	TypeMovie:   "pgc",
	TypeBangumi: "pgc",
	TypePGC:     "pgc",
	TypeArticle: "article",
	TypeUGC:     "archive",
	TypeLive:    "live",
	TypeCorpus:  "article-list",
	TypeComic:   "comic",
}

// BusinessByTP .
func BusinessByTP(b int8) string {
	return businessIDs[b]
}

// CheckBusiness .
func CheckBusiness(bs string) (tp int8, err error) {
	if bs == "" {
		return
	}
	tp, ok := businesses[bs]
	if !ok {
		err = ecode.AppDenied
	}
	return
}

// MustCheckBusiness .
func MustCheckBusiness(bs string) (tp int8, err error) {
	if bs == "" {
		err = ecode.RequestErr
		return
	}
	tp, ok := businesses[bs]
	if !ok {
		err = ecode.AppDenied
	}
	return
}

// Merge report merge in history.
type Merge struct {
	Mid int64 `json:"mid"`
	Now int64 `json:"now"`
}

// Video video history.
type Video struct {
	*archive.Archive3
	Favorite    bool           `json:"favorite"` // video favorite
	TP          int8           `json:"type"`     // video type
	STP         int8           `json:"sub_type"` // video type
	DT          int8           `json:"device"`   // device type
	Page        *archive.Page3 `json:"page,omitempty"`
	Count       int            `json:"count,omitempty"`
	BangumiInfo *Bangumi       `json:"bangumi,omitempty"`
	Progress    int64          `json:"progress"`
	ViewAt      int64          `json:"view_at"`
}

// Season season.
type Season struct {
	ID            int64  `json:"season_id"`
	Title         string `json:"title"`
	SeasonStatus  int    `json:"season_status"`
	IsFinish      int    `json:"is_finish"`
	TotalCount    int32  `json:"total_count"`
	NewestEpid    int64  `json:"newest_ep_id"`
	NewestEpIndex string `json:"newest_ep_index"`
	SeasonType    int    `json:"season_type,omitempty"`
	Mode          int    `json:"mode,omitempty"`
}

// Bangumi bangumi.
type Bangumi struct {
	Epid          int64   `json:"ep_id"`
	Title         string  `json:"title"`
	LongTitle     string  `json:"long_title"`
	EpisodeStatus int     `json:"episode_status"`
	Follow        int     `json:"follow"`
	Cover         string  `json:"cover"`
	Season        *Season `json:"season"`
}

// BangumiSeason season.
type BangumiSeason struct {
	ID       int64 `json:"season_id"`
	Epid     int64 `json:"episode_id"`
	EpidType int64 `json:"season_type"`
}
