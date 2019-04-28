package model

import (
	"fmt"

	livemdl "go-common/app/interface/main/app-interface/model/live"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/conf/env"
)

// app-interface const
const (
	// PlatAndroid is int8 for android.
	PlatAndroid = int8(0)
	// PlatIPhone is int8 for iphone.
	PlatIPhone = int8(1)
	// PlatIPad is int8 for ipad.
	PlatIPad = int8(2)
	// PlatWPhone is int8 for wphone.
	PlatWPhone = int8(3)
	// PlatAndroidG is int8 for Android Googleplay.
	PlatAndroidG = int8(4)
	// PlatIPhoneI is int8 for Iphone Global.
	PlatIPhoneI = int8(5)
	// PlatIPadI is int8 for IPAD Global.
	PlatIPadI = int8(6)
	// PlatAndroidTV is int8 for AndroidTV Global.
	PlatAndroidTV = int8(7)
	// PlatAndroidI is int8 for Android Global.
	PlatAndroidI = int8(8)
	// PlatIpadHD is int8 for IpadHD
	PlatIpadHD = int8(9)
	// PlatAndroidB is int8 for Android Blue.
	PlatAndroidB = int8(10)
	// PlatIPhoneB is int8 for Android Blue.
	PlatIPhoneB = int8(11)

	GotoAv             = "av"
	GotoWeb            = "web"
	GotoBangumi        = "bangumi"
	GotoMovie          = "movie"
	GotoBangumiWeb     = "bangumi_web"
	GotoSp             = "sp"
	GotoLive           = "live"
	GotoGame           = "game"
	GotoAuthor         = "author"
	GotoClip           = "clip"
	GotoAlbum          = "album"
	GotoArticle        = "article"
	GotoAudio          = "audio"
	GotoSpecial        = "special"
	GotoBanner         = "banner"
	GotoSpecialS       = "special_s"
	GotoConverge       = "converge"
	GOtoRecommendWord  = "recommend_word"
	GotoPGC            = "pgc"
	GotoSuggestKeyWord = "suggest_keyword"
	GotoComic          = "comic"
	GotoChannel        = "channel"
	GotoEP             = "ep"
	GotoTwitter        = "twitter"
	GotoStar           = "star"
	GotoTicket         = "ticket"
	GotoProduct        = "product"
	GotoSpace          = "space"
	GotoSpecialerGuide = "special_guide"
	GotoDynamic        = "dynamic"
	// EnvPro is pro.
	EnvPro = "pro"
	EnvHK  = "hk"
	// EnvTest is env.
	EnvTest = "test"
	// EnvDev is env.
	EnvDev = "dev"
	// ForbidCode is forbid by law
	ForbidCode   = -110
	NoResultCode = -111

	CoverIng      = "即将上映"
	CoverPay      = "付费观看"
	CoverFree     = "免费观看"
	CoverVipFree  = "付费观看"
	CoverVipOnly  = "专享"
	CoverVipFirst = "抢先"

	Hans = "hans"
	Hant = "hant"

	// AttrNo attribute no
	AttrNo = int32(0)
	// AttrYes attribute yes
	AttrYes = int32(1)

	AttrBitArchive = uint32(0)
	AttrBitArticle = uint32(1)
	AttrBitClip    = uint32(2)
	AttrBitAlbum   = uint32(3)
	AttrBitAudio   = uint32(34)
)

// for FillURI
var (
	AvHandler = func(a *archive.Archive3) func(uri string) string {
		return func(uri string) string {
			if a == nil {
				return uri
			}
			if a.Dimension.Height != 0 || a.Dimension.Width != 0 {
				return fmt.Sprintf("%s?player_width=%d&player_height=%d&player_rotate=%d", uri, a.Dimension.Width, a.Dimension.Height, a.Dimension.Rotate)
			}
			return uri
		}
	}
	LiveHandler = func(l *livemdl.RoomInfo) func(uri string) string {
		return func(uri string) string {
			if l == nil {
				return uri
			}
			if l.BroadcastType == 0 || l.BroadcastType == 1 {
				return fmt.Sprintf("%s?broadcast_type=%d", uri, l.BroadcastType)
			}
			return uri
		}
	}
)

// IsAndroid check plat is android or ipad.
func IsAndroid(plat int8) bool {
	return plat == PlatAndroid || plat == PlatAndroidG || plat == PlatAndroidI || plat == PlatAndroidB
}

// IsIOS check plat is iphone or ipad.
func IsIOS(plat int8) bool {
	return plat == PlatIPad || plat == PlatIPhone || plat == PlatIPadI || plat == PlatIPhoneI
}

// IsIPhone check plat is iphone.
func IsIPhone(plat int8) bool {
	return plat == PlatIPhone || plat == PlatIPhoneI
}

// IsIPad check plat is pad.
func IsIPad(plat int8) bool {
	return plat == PlatIPad || plat == PlatIPadI || plat == PlatIpadHD
}

// IsIPhoneB check plat is iphone_b.
func IsIPhoneB(plat int8) bool {
	return plat == PlatIPhoneB
}

// Plat return plat by platStr or mobiApp
func Plat(mobiApp, device string) int8 {
	switch mobiApp {
	case "iphone":
		if device == "pad" {
			return PlatIPad
		}
		return PlatIPhone
	case "white":
		return PlatIPhone
	case "ipad":
		return PlatIpadHD
	case "android":
		return PlatAndroid
	case "win", "winphone":
		return PlatWPhone
	case "android_G":
		return PlatAndroidG
	case "android_i":
		return PlatAndroidI
	case "iphone_i":
		if device == "pad" {
			return PlatIPadI
		}
		return PlatIPhoneI
	case "ipad_i":
		return PlatIPadI
	case "android_tv":
		return PlatAndroidTV
	case "android_b":
		return PlatAndroidB
	case "iphone_b":
		return PlatIPhoneB
	}
	return PlatIPhone
}

// IsOverseas is overseas
func IsOverseas(plat int8) bool {
	return plat == PlatAndroidI || plat == PlatIPhoneI || plat == PlatIPadI
}

// FillURI deal app schema.
func FillURI(gt, param string, f func(uri string) string) (uri string) {
	switch gt {
	case GotoAv, "":
		uri = "bilibili://video/" + param
	case GotoLive:
		uri = "bilibili://live/" + param
	case GotoBangumi:
		uri = "bilibili://bangumi/season/" + param
	case GotoBangumiWeb:
		uri = "http://bangumi.bilibili.com/anime/" + param
	case GotoGame:
		uri = "bilibili://game_center/detail?id=" + param + "&sourceType=adPut"
	case GotoSp:
		uri = "bilibili://splist/" + param
	case GotoAuthor:
		uri = "bilibili://author/" + param
	case GotoClip:
		uri = "bilibili://clip/" + param
	case GotoAlbum:
		uri = "bilibili://album/" + param
	case GotoArticle:
		uri = "bilibili://article/" + param
	case GotoWeb:
		uri = param
	case GotoPGC:
		uri = "https://www.bilibili.com/bangumi/play/ss" + param
	case GotoChannel:
		uri = "bilibili://pegasus/channel/" + param + "/"
	case GotoEP:
		uri = "https://www.bilibili.com/bangumi/play/ep" + param
	case GotoTwitter:
		uri = "bilibili://pictureshow/detail/" + param
	case GotoSpace:
		uri = "bilibili://space/" + param
	case GotoDynamic:
		uri = "bilibili://following/detail/" + param
	}
	if f != nil {
		uri = f(uri)
	}
	return
}

// StatusMark cover status mark
func StatusMark(status int) string {
	if status == 0 {
		return CoverIng
	} else if status == 1 {
		return CoverPay
	} else if status == 2 {
		return CoverFree
	} else if status == 3 {
		return CoverVipFree
	} else if status == 4 {
		return CoverVipOnly
	} else if status == 5 {
		return CoverVipFirst
	}
	return ""
}

// AttrVal get attribute value
func AttrVal(attr int32, bit uint32) (v int32) {
	v = (attr >> bit) & int32(1)
	return
}

// AttrSet set attribute value
func AttrSet(attr int32, v int32, bit uint32) int32 {
	return attr&(^(1 << bit)) | (v << bit)
}

// Direction define
type Direction int

// app-interface const
const (
	Upward   Direction = 1
	Downward Direction = 2
)

// Cursor struct
type Cursor struct {
	Current   int64
	Direction Direction
	Size      int
}

// Latest judge cursor Current
func (c *Cursor) Latest() bool {
	return c.Current == 0
}

// MoveUpward judge cursor Direction
func (c *Cursor) MoveUpward() bool {
	return c.Direction == Upward
}

// MoveDownward judge cursor Direction
func (c *Cursor) MoveDownward() bool {
	return c.Direction == Downward
}

// NewCursor judge cuser
func NewCursor(maxID int64, minID int64, size int) (cuser *Cursor, err error) {
	if maxID < 0 || minID < 0 {
		err = fmt.Errorf("either max_id(%d) or min_id(%d) < 0", maxID, minID)
		return
	}
	if (minID * maxID) != 0 {
		err = fmt.Errorf("both max_id(%d) and max_id(%d) > 0", maxID, minID)
		return
	}
	if minID == 0 && maxID == 0 {
		cuser = &Cursor{Current: 0, Direction: Downward, Size: size}
	} else if maxID > 0 {
		cuser = &Cursor{Current: maxID, Direction: Downward, Size: size}
	} else {
		cuser = &Cursor{Current: minID, Direction: Upward, Size: size}
	}
	return
}

// InvalidBuild invalid build
func InvalidBuild(srcBuild, cfgBuild int, cfgCond string) bool {
	if cfgBuild != 0 && cfgCond != "" {
		switch cfgCond {
		case "gt":
			if cfgBuild >= srcBuild {
				return true
			}
		case "lt":
			if cfgBuild <= srcBuild {
				return true
			}
		case "eq":
			if cfgBuild != srcBuild {
				return true
			}
		case "ne":
			if cfgBuild == srcBuild {
				return true
			}
		}
	}
	return false
}

// env sh001 run
func EnvRun() (res bool) {
	var _zone = "sh001"
	if env.Zone == _zone {
		return true
	}
	return false
}

// FormMediaType media type
func FormMediaType(mediaType int) (mediaName string) {
	switch mediaType {
	case 1:
		mediaName = "番剧"
	case 2:
		mediaName = "电影"
	case 3:
		mediaName = "纪录片"
	case 4:
		mediaName = "国创"
	case 5:
		mediaName = "电视剧"
	case 6:
		mediaName = "漫画"
	case 7:
		mediaName = "综艺"
	case 123:
		mediaName = "电视剧"
	case 124:
		mediaName = "电视剧"
	case 125:
		mediaName = "纪录片"
	case 126:
		mediaName = "电影"
	case 127:
		mediaName = "动漫"
	}
	return
}

// ReasonStyle reason style
type ReasonStyle struct {
	Text             string `json:"text,omitempty"`
	TextColor        string `json:"text_color,omitempty"`
	TextColorNight   string `json:"text_color_night,omitempty"`
	BgColor          string `json:"bg_color,omitempty"`
	BgColorNight     string `json:"bg_color_night,omitempty"`
	BorderColor      string `json:"border_color,omitempty"`
	BorderColorNight string `json:"border_color_night,omitempty"`
	BgStyle          int8   `json:"bg_style,omitempty"`
}
